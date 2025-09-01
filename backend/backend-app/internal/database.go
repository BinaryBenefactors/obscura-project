package internal

import (
	"fmt"
	"strings"
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
	return d.DB.Create(file).Error
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
		updates["error_message"] = "" // Очищаем предыдущие ошибки при успехе
	}
	
	return d.DB.Model(&File{}).Where("id = ?", id).Updates(updates).Error
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

// Получение статистики пользователя с учетом обработки
func (d *Database) GetUserStats(userID uint) (*UserStats, error) {
	var stats UserStats
	
	// Общее количество файлов и их размер
	var totalSize int64
	var totalFiles int64
	
	err := d.DB.Model(&File{}).
		Where("user_id = ?", userID).
		Select("COUNT(*) as count, COALESCE(SUM(file_size), 0) as total").
		Row().
		Scan(&totalFiles, &totalSize)
	if err != nil {
		return nil, fmt.Errorf("failed to get total stats: %w", err)
	}
	
	stats.TotalFiles = int(totalFiles)
	stats.TotalSize = totalSize
	stats.TotalSizeMB = float64(totalSize) / (1024 * 1024)
	
	// Статистика по статусам
	var processedCount, processingCount, failedCount int64
	
	d.DB.Model(&File{}).Where("user_id = ? AND status = ?", userID, StatusCompleted).Count(&processedCount)
	d.DB.Model(&File{}).Where("user_id = ? AND status = ?", userID, StatusProcessing).Count(&processingCount)
	d.DB.Model(&File{}).Where("user_id = ? AND status = ?", userID, StatusFailed).Count(&failedCount)
	
	stats.ProcessedFiles = int(processedCount)
	stats.ProcessingFiles = int(processingCount)
	stats.FailedFiles = int(failedCount)
	
	// Файлы загруженные сегодня
	today := time.Now().Truncate(24 * time.Hour)
	var todayCount int64
	err = d.DB.Model(&File{}).
		Where("user_id = ? AND uploaded_at >= ?", userID, today).
		Count(&todayCount).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get today upload stats: %w", err)
	}
	stats.UploadedToday = int(todayCount)
	
	// Файлы обработанные сегодня
	var processedTodayCount int64
	err = d.DB.Model(&File{}).
		Where("user_id = ? AND status = ? AND processed_at >= ?", userID, StatusCompleted, today).
		Count(&processedTodayCount).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get today processing stats: %w", err)
	}
	stats.ProcessedToday = int(processedTodayCount)
	
	// Файлы загруженные на этой неделе
	weekStart := today.AddDate(0, 0, -int(today.Weekday()))
	var weekCount int64
	err = d.DB.Model(&File{}).
		Where("user_id = ? AND uploaded_at >= ?", userID, weekStart).
		Count(&weekCount).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get week upload stats: %w", err)
	}
	stats.UploadedThisWeek = int(weekCount)
	
	// Файлы обработанные на этой неделе
	var processedWeekCount int64
	err = d.DB.Model(&File{}).
		Where("user_id = ? AND status = ? AND processed_at >= ?", userID, StatusCompleted, weekStart).
		Count(&processedWeekCount).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get week processing stats: %w", err)
	}
	stats.ProcessedThisWeek = int(processedWeekCount)
	
	// Файлы загруженные в этом месяце
	monthStart := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location())
	var monthCount int64
	err = d.DB.Model(&File{}).
		Where("user_id = ? AND uploaded_at >= ?", userID, monthStart).
		Count(&monthCount).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get month upload stats: %w", err)
	}
	stats.UploadedThisMonth = int(monthCount)
	
	// Файлы обработанные в этом месяце
	var processedMonthCount int64
	err = d.DB.Model(&File{}).
		Where("user_id = ? AND status = ? AND processed_at >= ?", userID, StatusCompleted, monthStart).
		Count(&processedMonthCount).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get month processing stats: %w", err)
	}
	stats.ProcessedThisMonth = int(processedMonthCount)
	
	// Статистика по статусам файлов
	var statusStats []struct {
		Status string
		Count  int
	}
	err = d.DB.Model(&File{}).
		Select("status, COUNT(*) as count").
		Where("user_id = ?", userID).
		Group("status").
		Find(&statusStats).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get status stats: %w", err)
	}
	
	stats.FilesByStatus = make(map[string]int)
	for _, stat := range statusStats {
		stats.FilesByStatus[stat.Status] = stat.Count
	}
	
	// Статистика по типам файлов
	var typeStats []struct {
		MimeType string
		Count    int
	}
	err = d.DB.Model(&File{}).
		Select("mime_type, COUNT(*) as count").
		Where("user_id = ?", userID).
		Group("mime_type").
		Find(&typeStats).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get type stats: %w", err)
	}
	
	stats.FilesByType = make(map[string]int)
	for _, stat := range typeStats {
		// Упрощаем тип файла для отображения
		simpleType := getSimpleFileType(stat.MimeType)
		stats.FilesByType[simpleType] += stat.Count
	}
	
	// Последние 5 файлов
	var recentFiles []File
	err = d.DB.Where("user_id = ?", userID).
		Order("uploaded_at DESC").
		Limit(5).
		Find(&recentFiles).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get recent files: %w", err)
	}
	stats.RecentFiles = recentFiles
	
	return &stats, nil
}

// Вспомогательная функция для упрощения типов файлов
func getSimpleFileType(mimeType string) string {
	if strings.HasPrefix(mimeType, "image/") {
		return "image"
	}
	if strings.HasPrefix(mimeType, "video/") {
		return "video"
	}
	if strings.HasPrefix(mimeType, "audio/") {
		return "audio"
	}
	if strings.HasPrefix(mimeType, "text/") {
		return "text"
	}
	return "other"
}
