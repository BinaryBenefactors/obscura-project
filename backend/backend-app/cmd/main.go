package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"obscura.app/internal"
	"obscura.app/pkg/logger"
	
	// Импорт сгенерированной документации - добавить после генерации
	_ "obscura.app/docs"
)

// @title Obscura API
// @version 1.0
// @description Advanced file upload and ML processing API for image and video blurring

// @contact.name API Support  
// @contact.email yagadanaga@ya.ru

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token

func main() {
	// Создаем логгер
	appLogger, err := logger.NewLogger("app.log")
	if err != nil {
		log.Fatal("Failed to create logger:", err)
	}
	defer appLogger.Close()

	// Загружаем конфигурацию
	cfg := internal.NewConfig()
	
	appLogger.Info("Starting Obscura API server...")
	appLogger.Info("Configuration loaded: ML service enabled: %v, URL: %s", cfg.MLServiceEnabled, cfg.MLServiceURL)

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
	
	appLogger.Info("Routes configured successfully")

	// Создаем HTTP сервер
	httpServer := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      server.GetRouter(),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Запускаем сервер в горутине
	go func() {
		appLogger.Info("Server starting on port %s", cfg.Port)
		appLogger.Info("Swagger UI available at: http://localhost:%s/swagger/", cfg.Port)
		appLogger.Info("API endpoints available at: http://localhost:%s/api/", cfg.Port)
		
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Fatal("Server failed: %v", err)
		}
	}()

	// Ожидание сигнала для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	
	<-quit
	appLogger.Info("Shutting down server...")

	// Graceful shutdown с таймаутом 30 секунд
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Останавливаем HTTP сервер
	if err := httpServer.Shutdown(ctx); err != nil {
		appLogger.Error("Server forced to shutdown: %v", err)
	}

	// Останавливаем внутренние сервисы
	server.Stop()

	appLogger.Info("Server exited gracefully")
}
