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
	TotalFiles        int       `json:"total_files" gorm:"default:0" example:"50"`
	TotalProcessed    int       `json:"total_processed" gorm:"default:0" example:"45"`
	TotalFailed       int       `json:"total_failed" gorm:"default:0" example:"5"`
	TotalSize         int64     `json:"total_size" gorm:"default:0" example:"524288000"`
	LastStatsUpdate   time.Time `json:"last_stats_update" example:"2024-01-15T09:00:00Z"`
	CreatedAt time.Time `json:"created_at" example:"2025-01-15T09:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2025-01-15T09:00:00Z"`
}

// File модель загруженного файла
// @Description Uploaded file information with processing status
type File struct {
	ID            string    `json:"id" gorm:"primarykey" example:"550e8400-e29b-41d4-a716-446655440000"`
	UserID        uint      `json:"user_id" gorm:"not null" example:"1"`
	OriginalName  string    `json:"original_name" gorm:"not null" example:"photo.jpg"`
	FileName      string    `json:"file_name" gorm:"not null" example:"550e8400-e29b-41d4-a716-446655440000.jpg"`
	ProcessedName string    `json:"processed_name,omitempty" gorm:"" example:"550e8400-e29b-41d4-a716-446655440000_processed.jpg"`
	FileSize      int64     `json:"file_size" gorm:"not null" example:"1048576"`
	ProcessedSize int64     `json:"processed_size,omitempty" gorm:"" example:"1048576"`
	MimeType      string    `json:"mime_type" gorm:"not null" example:"image/jpeg"`
	Status        string    `json:"status" gorm:"default:'uploaded'" example:"uploaded" enums:"uploaded,processing,completed,failed"`
	ErrorMessage  string    `json:"error_message,omitempty" gorm:"" example:"Processing failed: invalid format"`
	UploadedAt    time.Time `json:"uploaded_at" example:"2025-01-15T09:00:00Z"`
	ProcessedAt   time.Time `json:"processed_at,omitempty" example:"2025-01-15T09:05:00Z"`
	User          User      `json:"-" gorm:"foreignKey:UserID"`
}

// ProcessingRequest запрос на обработку файла ML-сервисом
// @Description ML processing request
type ProcessingRequest struct {
	FileID   string            `json:"file_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	FilePath string            `json:"file_path" example:"/uploads/550e8400-e29b-41d4-a716-446655440000.jpg"`
	MimeType string            `json:"mime_type" example:"image/jpeg"`
	Options  ProcessingOptions `json:"options"`
}

// ProcessingOptions опции обработки файла
// @Description File processing options
type ProcessingOptions struct {
	BlurType    string   `json:"blur_type" example:"gaussian" enums:"gaussian,motion,pixelate"`
	Intensity   int      `json:"intensity" example:"5"`
	ObjectTypes []string `json:"object_types" example:"face,person,car"`
}

// ProcessingResponse ответ ML-сервиса
// @Description ML processing response
type ProcessingResponse struct {
	Success        bool     `json:"success" example:"true"`
	ProcessedPath  string   `json:"processed_path,omitempty" example:"/uploads/550e8400-e29b-41d4-a716-446655440000_processed.jpg"`
	ProcessedSize  int64    `json:"processed_size,omitempty" example:"1048576"`
	ObjectsFound   []string `json:"objects_found,omitempty" example:"face,person"`
	ProcessingTime int      `json:"processing_time_ms,omitempty" example:"2500"`
	ErrorMessage   string   `json:"error_message,omitempty" example:"Failed to detect objects"`
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

// IsProcessed проверяет, обработан ли файл
func (f *File) IsProcessed() bool {
	return f.Status == StatusCompleted && f.ProcessedName != ""
}

// CanBeProcessed проверяет, может ли файл быть обработан
func (f *File) CanBeProcessed() bool {
	return f.Status == StatusUploaded || f.Status == StatusFailed
}

// Статусы файлов
const (
	StatusUploaded   = "uploaded"
	StatusProcessing = "processing"
	StatusCompleted  = "completed"
	StatusFailed     = "failed"
)
