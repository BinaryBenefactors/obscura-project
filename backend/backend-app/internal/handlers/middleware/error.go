package middleware

import (
	"encoding/json"
	"net/http"

	"obscura.app/backend/internal/domain/models"
	"obscura.app/backend/pkg/logger"
)

func ErrorHandler(logger *logger.Logger, next func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := next(w, r)
		if err != nil {
			logger.Error("Request handling error: %v", err)

			response := models.ErrorResponse{
				Error:   "Internal Server Error",
				Details: err.Error(),
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			if err := json.NewEncoder(w).Encode(response); err != nil {
				logger.Error("Failed to encode error response: %v", err)
			}
		}
	}
}

func SendError(w http.ResponseWriter, message string, statusCode int) error {
	response := models.ErrorResponse{
		Error: message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(response)
}
