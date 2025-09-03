package internal

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"obscura.app/pkg/logger"
)

type Database struct {
	DB     *gorm.DB
	logger *logger.Logger
}

func NewDatabase(cfg *Config, log *logger.Logger) (*Database, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Автомиграция
	err = db.AutoMigrate(&User{}, &File{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Info("Database connected and migrated successfully")
	return &Database{DB: db, logger: log}, nil
}

// Методы для работы с пользователями
func (d *Database) CreateUser(user *User) error {
	if err := user.HashPassword(); err != nil {
		return err
	}
	return d.DB.Create(user).Error
}

func (d *Database) GetUserByEmail(email string) (*User, error) {
	var user User
	err := d.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (d *Database) GetUserByID(id uint) (*User, error) {
	var user User
	err := d.DB.First(&user, id).Error
	return &user, err
}

func (d *Database) UpdateUser(user *User) error {
	return d.DB.Save(user).Error
}

func (d *Database) DeleteUser(id uint) error {
	return d.DB.Delete(&User{}, id).Error
}

// Методы для работы с файлами
func (d *Database) CreateFile(file *File) error {
	if err := d.DB.Create(file).Error; err != nil {
		return err
	}
	
	return d.UpdateUserStats(file.UserID, 1, 0, 0, file.FileSize)
}

func (d *Database) GetFileByID(id string) (*File, error) {
	var file File
	err := d.DB.First(&file, "id = ?", id).Error
	return &file, err
}

func (d *Database) GetUserFiles(userID uint) ([]File, error) {
	var files []File
	err := d.DB.Where("user_id = ?", userID).Order("uploaded_at DESC").Find(&files).Error
	return files, err
}

func (d *Database) UpdateFileStatus(id string, status string) error {
	updates := map[string]interface{}{
		"status": status,
	}

	// Если статус "processing", обновляем время начала обработки
	if status == StatusProcessing {
		updates["processed_at"] = time.Now()
	}

	return d.DB.Model(&File{}).Where("id = ?", id).Updates(updates).Error
}

func (d *Database) UpdateFileProcessing(id string, processedName string, processedSize int64, status string, errorMessage string) error {
	// Сначала получаем текущий статус файла
	file, err := d.GetFileByID(id)
	if err != nil {
		return err
	}
	
	updates := map[string]interface{}{
		"status":       status,
		"processed_at": time.Now(),
	}

	if processedName != "" {
		updates["processed_name"] = processedName
	}

	if processedSize > 0 {
		updates["processed_size"] = processedSize
	}

	if errorMessage != "" {
		updates["error_message"] = errorMessage
	} else {
		updates["error_message"] = ""
	}

	// Обновляем файл
	if err := d.DB.Model(&File{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return err
	}
	
	// Обновляем статистику пользователя при изменении статуса обработки
	if file.Status != status {
		var processedDelta, failedDelta int
		
		if status == StatusCompleted {
			processedDelta = 1
			if file.Status == StatusFailed {
				failedDelta = -1
			}
		} else if status == StatusFailed {
			failedDelta = 1
			if file.Status == StatusCompleted {
				processedDelta = -1
			}
		}
		
		if processedDelta != 0 || failedDelta != 0 {
			return d.UpdateUserStats(file.UserID, 0, processedDelta, failedDelta, 0)
		}
	}
	
	return nil
}

func (d *Database) GetFilesByStatus(status string) ([]File, error) {
	var files []File
	err := d.DB.Where("status = ?", status).Find(&files).Error
	return files, err
}

func (d *Database) GetPendingFiles() ([]File, error) {
	var files []File
	err := d.DB.Where("status = ?", StatusUploaded).
		Order("uploaded_at ASC").
		Find(&files).Error
	return files, err
}

func (d *Database) DeleteFile(id string) error {
	return d.DB.Delete(&File{}, "id = ?", id).Error
}


// Получение статистики пользователя
func (d *Database) GetUserStats(userID uint) (*User, error) {
	user, err := d.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	
	return user, nil
}

func (d *Database) UpdateUserStats(userID uint, filesDelta int, processedDelta int, failedDelta int, sizeDelta int64) error {
	return d.DB.Model(&User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"total_files":        gorm.Expr("total_files + ?", filesDelta),
		"total_processed":    gorm.Expr("total_processed + ?", processedDelta),
		"total_failed":       gorm.Expr("total_failed + ?", failedDelta),
		"total_size":         gorm.Expr("total_size + ?", sizeDelta),
		"last_stats_update":  time.Now(),
	}).Error
}
