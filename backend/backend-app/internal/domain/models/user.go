package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User представляет пользователя в системе
type User struct {
	gorm.Model            // Встраивает ID, CreatedAt, UpdatedAt, DeletedAt
	Email      string     `gorm:"uniqueIndex;not null"` // Уникальный email
	Password   string     `gorm:"not null"`             // Хешированный пароль
	Name       string     `gorm:"not null"`             // Имя пользователя
	Role       string     `gorm:"default:'user'"`       // Роль пользователя (user/admin)
	Active     bool       `gorm:"default:true"`         // Активен ли пользователь
	LastLogin  *time.Time // Время последнего входа
}

// BeforeCreate хук, который выполняется перед созданием пользователя
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// Хешируем пароль перед сохранением
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// BeforeUpdate хук, который выполняется перед обновлением пользователя
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	// Если пароль был изменен, хешируем его
	if tx.Statement.Changed("Password") {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	return nil
}

// CheckPassword проверяет, соответствует ли пароль хешу
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// TableName указывает имя таблицы в базе данных
func (u *User) TableName() string {
	return "users"
}
