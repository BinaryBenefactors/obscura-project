package repository

import (
	"context"

	"gorm.io/gorm"
	"obscura.app/backend/internal/domain/models"
)

// GormUserRepository реализация UserRepository с использованием GORM
type GormUserRepository struct {
	db *gorm.DB
}

// NewGormUserRepository создает новый репозиторий пользователей
func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

// Create создает нового пользователя
func (r *GormUserRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// GetByID получает пользователя по ID
func (r *GormUserRepository) GetByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail получает пользователя по email
func (r *GormUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update обновляет информацию о пользователе
func (r *GormUserRepository) Update(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// Delete удаляет пользователя
func (r *GormUserRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, id).Error
}

// List возвращает список пользователей с пагинацией
func (r *GormUserRepository) List(ctx context.Context, offset, limit int) ([]models.User, error) {
	var users []models.User
	err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
