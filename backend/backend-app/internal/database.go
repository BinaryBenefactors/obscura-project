package internal

import (
	"fmt"

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

// Методы для работы с файлами
func (d *Database) CreateFile(file *File) error {
	return d.DB.Create(file).Error
}

func (d *Database) GetFileByID(id string) (*File, error) {
	var file File
	err := d.DB.Preload("User").First(&file, "id = ?", id).Error
	return &file, err
}

func (d *Database) GetUserFiles(userID uint) ([]File, error) {
	var files []File
	err := d.DB.Where("user_id = ?", userID).Order("uploaded_at DESC").Find(&files).Error
	return files, err
}

func (d *Database) UpdateFileStatus(id string, status string) error {
	return d.DB.Model(&File{}).Where("id = ?", id).Update("status", status).Error
}

func (d *Database) DeleteFile(id string) error {
	return d.DB.Delete(&File{}, "id = ?", id).Error
}
