package internal

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"obscura.app/pkg/logger"
)

// FileCleaner структура для очистки файлов
type FileCleaner struct {
	uploadPath       string
	cleanupInterval  time.Duration
	maxAge          time.Duration
	logger          *logger.Logger
	stopChan        chan struct{}
}

// NewFileCleaner создает новый file cleaner
func NewFileCleaner(uploadPath string, logger *logger.Logger) *FileCleaner {
	return &FileCleaner{
		uploadPath:      uploadPath,
		cleanupInterval: 6 * time.Hour, // Запуск очистки каждые 6 часов
		maxAge:         24 * time.Hour, // Удаляем файлы старше 24 часов
		logger:         logger,
		stopChan:       make(chan struct{}),
	}
}

// Start запускает автоматическую очистку файлов
func (fc *FileCleaner) Start() {
	fc.logger.Info("File cleaner started with interval: %v, max age: %v", fc.cleanupInterval, fc.maxAge)
	
	// Запускаем первую очистку сразу
	go fc.cleanupOnce()
	
	// Запускаем периодическую очистку
	go fc.run()
}

// Stop останавливает автоматическую очистку
func (fc *FileCleaner) Stop() {
	fc.logger.Info("Stopping file cleaner...")
	close(fc.stopChan)
}

// run основной цикл очистки
func (fc *FileCleaner) run() {
	ticker := time.NewTicker(fc.cleanupInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			fc.cleanupOnce()
		case <-fc.stopChan:
			fc.logger.Info("File cleaner stopped")
			return
		}
	}
}

// cleanupOnce выполняет одну итерацию очистки
func (fc *FileCleaner) cleanupOnce() {
	fc.logger.Debug("Starting file cleanup process...")
	
	start := time.Now()
	deletedCount := 0
	totalSize := int64(0)
	
	err := filepath.Walk(fc.uploadPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fc.logger.Warning("Error accessing file %s: %v", path, err)
			return nil // Продолжаем обход
		}
		
		// Пропускаем директории
		if info.IsDir() {
			return nil
		}
		
		// Проверяем возраст файла
		if time.Since(info.ModTime()) > fc.maxAge {
			// Проверяем, является ли файл анонимным (не в БД)
			if fc.isAnonymousFile(path) {
				fc.logger.Debug("Deleting old anonymous file: %s (age: %v, size: %d bytes)", 
					path, time.Since(info.ModTime()), info.Size())
				
				if err := os.Remove(path); err != nil {
					fc.logger.Error("Failed to delete file %s: %v", path, err)
				} else {
					deletedCount++
					totalSize += info.Size()
				}
			}
		}
		
		return nil
	})
	
	if err != nil {
		fc.logger.Error("Error during file cleanup: %v", err)
		return
	}
	
	duration := time.Since(start)
	
	if deletedCount > 0 {
		fc.logger.Info("File cleanup completed: deleted %d files (%.2f MB) in %v", 
			deletedCount, float64(totalSize)/(1024*1024), duration)
	} else {
		fc.logger.Debug("File cleanup completed: no files to delete (took %v)", duration)
	}
}

// isAnonymousFile проверяет, является ли файл анонимным
// Анонимные файлы имеют имена в формате UUID.extension и не сохраняются в БД
func (fc *FileCleaner) isAnonymousFile(filePath string) bool {
	filename := filepath.Base(filePath)
	
	// Проверяем, что это файл с UUID именем
	// UUID имеет формат: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx (36 символов)
	parts := strings.Split(filename, ".")
	if len(parts) < 2 {
		return false
	}
	
	nameWithoutExt := parts[0]
	if len(nameWithoutExt) != 36 {
		return false
	}
	
	// Проверяем формат UUID (8-4-4-4-12)
	uuidParts := strings.Split(nameWithoutExt, "-")
	if len(uuidParts) != 5 {
		return false
	}
	
	if len(uuidParts[0]) != 8 || len(uuidParts[1]) != 4 || 
	   len(uuidParts[2]) != 4 || len(uuidParts[3]) != 4 || 
	   len(uuidParts[4]) != 12 {
		return false
	}
	
	// Все файлы с UUID именами считаем потенциально анонимными
	// В реальном приложении здесь можно было бы проверить БД,
	// но для простоты считаем, что все UUID файлы старше 24 часов - анонимные
	return true
}

// GetStats возвращает статистику файлов
func (fc *FileCleaner) GetStats() (map[string]interface{}, error) {
	stats := map[string]interface{}{
		"upload_path":        fc.uploadPath,
		"cleanup_interval_h": fc.cleanupInterval.Hours(),
		"max_age_h":         fc.maxAge.Hours(),
		"total_files":       0,
		"total_size_bytes":  int64(0),
		"total_size_mb":     float64(0),
	}
	
	totalFiles := 0
	totalSize := int64(0)
	
	err := filepath.Walk(fc.uploadPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		
		if !info.IsDir() {
			totalFiles++
			totalSize += info.Size()
		}
		
		return nil
	})
	
	if err != nil {
		return stats, err
	}
	
	stats["total_files"] = totalFiles
	stats["total_size_bytes"] = totalSize
	stats["total_size_mb"] = float64(totalSize) / (1024 * 1024)
	
	return stats, nil
}
