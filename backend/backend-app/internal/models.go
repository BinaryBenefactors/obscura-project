package internal

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User модель пользователя
// @Description User account information
type User struct {
	ID        uint      `json:"id" gorm:"primarykey" example:"1"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null" example:"user@example.com"`
	Password  string    `json:"-" gorm:"not null"`
	Name      string    `json:"name" gorm:"not null" example:"John Doe"`
	CreatedAt time.Time `json:"created_at" example:"2024-01-15T09:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2024-01-15T09:00:00Z"`
}

// File модель загруженного файла  
// @Description Uploaded file information
type File struct {
	ID           string    `json:"id" gorm:"primarykey" example:"550e8400-e29b-41d4-a716-446655440000"`
	UserID       uint      `json:"user_id" gorm:"not null" example:"1"`
	OriginalName string    `json:"original_name" gorm:"not null" example:"photo.jpg"`
	FileName     string    `json:"file_name" gorm:"not null" example:"550e8400-e29b-41d4-a716-446655440000.jpg"`
	FileSize     int64     `json:"file_size" gorm:"not null" example:"1048576"`
	MimeType     string    `json:"mime_type" gorm:"not null" example:"image/jpeg"`
	Status       string    `json:"status" gorm:"default:'uploaded'" example:"uploaded" enums:"uploaded,processing,completed,failed"`
	UploadedAt   time.Time `json:"uploaded_at" example:"2024-01-15T09:00:00Z"`
	User         User      `json:"-" gorm:"foreignKey:UserID"`
}

// UserStats статистика пользователя
// @Description User statistics and usage information
type UserStats struct {
	TotalFiles        int                `json:"total_files" example:"25"`
	TotalSize         int64              `json:"total_size" example:"52428800"`
	TotalSizeMB       float64            `json:"total_size_mb" example:"50.0"`
	UploadedToday     int                `json:"uploaded_today" example:"3"`
	UploadedThisWeek  int                `json:"uploaded_this_week" example:"8"`
	UploadedThisMonth int                `json:"uploaded_this_month" example:"15"`
	FilesByStatus     map[string]int     `json:"files_by_status" example:"uploaded:20,completed:5"`
	FilesByType       map[string]int     `json:"files_by_type" example:"image:15,video:10"`
	RecentFiles       []File             `json:"recent_files"`
}

// RegisterRequest запрос регистрации
// @Description User registration request
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required,min=6,max=128" example:"securePassword123"`
	Name     string `json:"name" binding:"required,min=2,max=100" example:"John Doe"`
}

// LoginRequest запрос авторизации
// @Description User login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required" example:"securePassword123"`
}

// UpdateProfileRequest запрос обновления профиля
// @Description Profile update request (all fields are optional)
type UpdateProfileRequest struct {
	Email    string `json:"email,omitempty" binding:"omitempty,email" example:"newemail@example.com"`
	Name     string `json:"name,omitempty" binding:"omitempty,min=2,max=100" example:"Jane Doe"`
	Password string `json:"password,omitempty" binding:"omitempty,min=6,max=128" example:"newSecurePassword456"`
}

// AuthResponse ответ авторизации
// @Description Authentication response with JWT token
type AuthResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User  User   `json:"user"`
}

// ErrorResponse ошибка API
// @Description Error response
type ErrorResponse struct {
	Error string `json:"error" example:"Invalid credentials"`
}

// SuccessResponse успешный ответ
// @Description Success response with optional data
type SuccessResponse struct {
	Message string      `json:"message" example:"Operation completed successfully"`
	Data    interface{} `json:"data,omitempty"`
}

// ValidationError ошибка валидации
// @Description Validation error details
type ValidationError struct {
	Field   string `json:"field" example:"email"`
	Message string `json:"message" example:"Email is required"`
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
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

// Статусы файлов
const (
	StatusUploaded   = "uploaded"
	StatusProcessing = "processing" 
	StatusCompleted  = "completed"
	StatusFailed     = "failed"
)
