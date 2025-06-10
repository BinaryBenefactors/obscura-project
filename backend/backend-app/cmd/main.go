package main

import (
	"log"
	"net/http"
	"os"

	"obscura.app/internal"
	"obscura.app/pkg/logger"
)

func main() {
	// Создаем логгер
	appLogger, err := logger.NewLogger("app.log")
	if err != nil {
		log.Fatal("Failed to create logger:", err)
	}
	defer appLogger.Close()

	// Загружаем конфигурацию
	cfg := internal.NewConfig()

	// Создаем папку для загрузок
	if err := os.MkdirAll(cfg.UploadPath, 0755); err != nil {
		appLogger.Fatal("Failed to create upload directory: %v", err)
	}

	// Подключаемся к базе данных
	db, err := internal.NewDatabase(cfg, appLogger)
	if err != nil {
		appLogger.Fatal("Failed to connect to database: %v", err)
	}

	// Создаем сервер
	server := internal.NewServer(cfg, db, appLogger)

	// Настраиваем роуты
	server.SetupRoutes()

	appLogger.Info("Server starting on port %s", cfg.Port)
	appLogger.Fatal("Server failed: %v", http.ListenAndServe(":"+cfg.Port, server.GetRouter()))
}