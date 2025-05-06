package models

import (
	"time"

	"gorm.io/gorm"
)

// Session представляет сессию пользователя
type Session struct {
	gorm.Model
	UserID    uint      `gorm:"not null"`                   // ID пользователя
	Token     string    `gorm:"uniqueIndex;not null"`       // Уникальный токен сессии
	ExpiresAt time.Time `gorm:"not null"`                   // Время истечения сессии
	UserAgent string    `gorm:"not null"`                   // User-Agent браузера
	IP        string    `gorm:"not null"`                   // IP адрес пользователя
	IsActive  bool      `gorm:"default:true;not null"`      // Активна ли сессия
	User      User      `gorm:"foreignKey:UserID;not null"` // Связь с пользователем
}

// IsExpired проверяет, истекла ли сессия
func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// TableName указывает имя таблицы в базе данных
func (s *Session) TableName() string {
	return "sessions"
}
