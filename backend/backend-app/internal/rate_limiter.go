package internal

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// RateLimiter структура для ограничения количества запросов
type RateLimiter struct {
	mu      sync.RWMutex
	clients map[string]*ClientInfo
	limit   int           // Максимум запросов
	window  time.Duration // Временное окно
}

// ClientInfo информация о клиенте
type ClientInfo struct {
	Count     int
	FirstSeen time.Time
	LastSeen  time.Time
}

// NewRateLimiter создает новый rate limiter
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		clients: make(map[string]*ClientInfo),
		limit:   limit,
		window:  window,
	}

	// Запускаем очистку старых записей каждые 10 минут
	go rl.cleanup()

	return rl
}

// generateClientID создает уникальный ID клиента на основе IP, User-Agent и других headers
func (rl *RateLimiter) generateClientID(r *http.Request) string {
	ip := getClientIP(r)
	userAgent := r.Header.Get("User-Agent")
	acceptLang := r.Header.Get("Accept-Language")

	// Создаем отпечаток клиента
	fingerprint := fmt.Sprintf("%s|%s|%s", ip, userAgent, acceptLang)

	// Хешируем для компактности
	hash := md5.Sum([]byte(fingerprint))
	return fmt.Sprintf("%x", hash)
}

// getClientIP извлекает реальный IP клиента
func getClientIP(r *http.Request) string {
	// Проверяем заголовки прокси
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		return forwarded
	}
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		return realIP
	}
	return r.RemoteAddr
}

// IsAllowed проверяет, разрешен ли запрос для данного клиента
func (rl *RateLimiter) IsAllowed(r *http.Request) (bool, int, time.Duration) {
	clientID := rl.generateClientID(r)
	now := time.Now()

	rl.mu.Lock()
	defer rl.mu.Unlock()

	client, exists := rl.clients[clientID]
	if !exists {
		// Новый клиент
		rl.clients[clientID] = &ClientInfo{
			Count:     1,
			FirstSeen: now,
			LastSeen:  now,
		}
		return true, 1, 0
	}

	// Если прошло больше времени чем окно, сбрасываем счетчик
	if now.Sub(client.FirstSeen) >= rl.window {
		client.Count = 1
		client.FirstSeen = now
		client.LastSeen = now
		return true, 1, 0
	}

	// Обновляем время последнего обращения
	client.LastSeen = now

	// Проверяем лимит
	if client.Count >= rl.limit {
		// Вычисляем время до сброса
		resetTime := client.FirstSeen.Add(rl.window)
		waitTime := time.Until(resetTime)
		return false, client.Count, waitTime
	}

	// Увеличиваем счетчик
	client.Count++
	return true, client.Count, 0
}

// cleanup очищает старые записи клиентов
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()

		for clientID, client := range rl.clients {
			// Удаляем клиентов, которые не обращались больше 2 часов
			if now.Sub(client.LastSeen) > 2*time.Hour {
				delete(rl.clients, clientID)
			}
		}

		rl.mu.Unlock()
	}
}

// GetStats возвращает статистику rate limiter'а
func (rl *RateLimiter) GetStats() map[string]interface{} {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	return map[string]interface{}{
		"total_clients": len(rl.clients),
		"limit":         rl.limit,
		"window_hours":  rl.window.Hours(),
	}
}
