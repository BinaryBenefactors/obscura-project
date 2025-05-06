package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"obscura.app/backend/internal/domain/models"
)

// NewPostgresDB создает новое подключение к PostgreSQL
func NewPostgresDB(host, port, user, password, dbname string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Автоматическая миграция схемы
	err = db.AutoMigrate(&models.User{}, &models.Session{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database connection established and migrations completed")
	return db, nil
}
