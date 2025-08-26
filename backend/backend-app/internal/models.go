package internal

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User модель пользователя
type User struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null"`
	Password  string    `json:"-" gorm:"not null"`
	Name      string    `json:"name" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// HashPassword хеширует пароль перед сохранением
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword проверяет соответствие пароля
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// File модель загруженного файла
type File struct {
	ID           string    `json:"id" gorm:"primarykey"`
	UserID       uint      `json:"user_id" gorm:"not null"`
	OriginalName string    `json:"original_name" gorm:"not null"`
	FileName     string    `json:"file_name" gorm:"not null"`
	FileSize     int64     `json:"file_size" gorm:"not null"`
	MimeType     string    `json:"mime_type" gorm:"not null"`
	Status       string    `json:"status" gorm:"default:'uploaded'"`
	UploadedAt   time.Time `json:"uploaded_at"`
	User         User      `json:"-" gorm:"foreignKey:UserID"` // Убираем из JSON
}

// Статусы файлов
const (
	StatusUploaded   = "uploaded"
	StatusProcessing = "processing"
	StatusCompleted  = "completed"
	StatusFailed     = "failed"
)

// UserStats статистика пользователя
type UserStats struct {
	TotalFiles      int     `json:"total_files"`
	TotalSize       int64   `json:"total_size"`
	TotalSizeMB     float64 `json:"total_size_mb"`
	UploadedToday   int     `json:"uploaded_today"`
	UploadedThisWeek int    `json:"uploaded_this_week"`
	UploadedThisMonth int   `json:"uploaded_this_month"`
	FilesByStatus   map[string]int `json:"files_by_status"`
	FilesByType     map[string]int `json:"files_by_type"`
	RecentFiles     []File  `json:"recent_files"`
}

// Структуры для API запросов/ответов
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateProfileRequest struct {
	Email    string `json:"email,omitempty"`
	Name     string `json:"name,omitempty"`
	Password string `json:"password,omitempty"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
