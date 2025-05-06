package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type contextKey string

const (
	UserIDKey contextKey = "userID"
)

// OptionalAuthMiddleware позволяет анонимный доступ, но добавляет userID если есть токен
func OptionalAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userID string

		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			token := strings.TrimPrefix(authHeader, "Bearer ")
			if token == "demo_token" { // Временная заглушка
				userID = "user123"
			}
		}

		if userID == "" {
			sessionID := r.Header.Get("X-Session-ID")
			if sessionID == "" {
				sessionID = generateSessionID()
			}
			userID = "anon_" + sessionID
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// AuthMiddleware требует аутентификации
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Authorization required"})
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token != "demo_token" { // Временная заглушка
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid token"})
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, "user123")
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func generateSessionID() string {
	return "session_" + strings.ReplaceAll(time.Now().String(), " ", "_")
}
