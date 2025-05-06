package repository

import (
	"context"
	"time"

	"gorm.io/gorm"
	"obscura.app/backend/internal/domain/models"
)

// GormSessionRepository реализация SessionRepository с использованием GORM
type GormSessionRepository struct {
	db *gorm.DB
}

// NewGormSessionRepository создает новый репозиторий сессий
func NewGormSessionRepository(db *gorm.DB) *GormSessionRepository {
	return &GormSessionRepository{db: db}
}

// Create создает новую сессию
func (r *GormSessionRepository) Create(ctx context.Context, session *models.Session) error {
	return r.db.WithContext(ctx).Create(session).Error
}

// GetByToken получает сессию по токену
func (r *GormSessionRepository) GetByToken(ctx context.Context, token string) (*models.Session, error) {
	var session models.Session
	err := r.db.WithContext(ctx).
		Where("token = ? AND is_active = ? AND expires_at > ?", token, true, time.Now()).
		First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// GetByUserID получает все активные сессии пользователя
func (r *GormSessionRepository) GetByUserID(ctx context.Context, userID uint) ([]models.Session, error) {
	var sessions []models.Session
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_active = ? AND expires_at > ?", userID, true, time.Now()).
		Find(&sessions).Error
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

// Delete удаляет сессию
func (r *GormSessionRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Session{}, id).Error
}

// DeleteExpired удаляет все истекшие сессии
func (r *GormSessionRepository) DeleteExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("expires_at <= ? OR is_active = ?", time.Now(), false).
		Delete(&models.Session{}).Error
}

// DeleteByUserID удаляет все сессии пользователя
func (r *GormSessionRepository) DeleteByUserID(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&models.Session{}).Error
}
