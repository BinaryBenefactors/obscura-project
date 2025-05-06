package repository

import (
	"context"

	"obscura.app/backend/internal/domain/models"
)

// UserRepository определяет методы для работы с пользователями
type UserRepository interface {
	// Create создает нового пользователя
	Create(ctx context.Context, user *models.User) error

	// GetByID получает пользователя по ID
	GetByID(ctx context.Context, id uint) (*models.User, error)

	// GetByEmail получает пользователя по email
	GetByEmail(ctx context.Context, email string) (*models.User, error)

	// Update обновляет информацию о пользователе
	Update(ctx context.Context, user *models.User) error

	// Delete удаляет пользователя
	Delete(ctx context.Context, id uint) error

	// List возвращает список пользователей с пагинацией
	List(ctx context.Context, offset, limit int) ([]models.User, error)
}
