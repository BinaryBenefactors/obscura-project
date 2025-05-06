package repository

import (
	"context"

	"obscura.app/backend/internal/domain/models"
)

// SessionRepository определяет методы для работы с сессиями
type SessionRepository interface {
	// Create создает новую сессию
	Create(ctx context.Context, session *models.Session) error

	// GetByToken получает сессию по токену
	GetByToken(ctx context.Context, token string) (*models.Session, error)

	// GetByUserID получает все активные сессии пользователя
	GetByUserID(ctx context.Context, userID uint) ([]models.Session, error)

	// Delete удаляет сессию
	Delete(ctx context.Context, id uint) error

	// DeleteExpired удаляет все истекшие сессии
	DeleteExpired(ctx context.Context) error

	// DeleteByUserID удаляет все сессии пользователя
	DeleteByUserID(ctx context.Context, userID uint) error
}
