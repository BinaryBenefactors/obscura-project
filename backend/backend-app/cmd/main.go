package main

import (
	"fmt"
	"net/http"
	"os"

	"obscura.app/backend/internal/config"
	"obscura.app/backend/internal/domain/repository"
	"obscura.app/backend/internal/handlers"
	"obscura.app/backend/internal/handlers/middleware"
	"obscura.app/backend/internal/services"
	"obscura.app/backend/pkg/logger"
)

type Application struct {
	config      *config.Config
	logger      *logger.Logger
	fileHandler *handlers.FileHandler
}

func NewApplication(cfg *config.Config) (*Application, error) {
	l, err := logger.NewLogger(cfg.LogFilePath)
	if err != nil {
		return nil, fmt.Errorf("logger init failed: %w", err)
	}

	// Создаем репозиторий в памяти
	fileRepo := repository.NewMemoryFileRepository()
	fileService := services.NewFileService(fileRepo, cfg.UploadPath, cfg.MaxUploadSize)
	fileHandler := handlers.NewFileHandler(fileService)

	return &Application{
		config:      cfg,
		logger:      l,
		fileHandler: fileHandler,
	}, nil
}

func (app *Application) Run() error {
	defer app.logger.Close()

	// Роуты
	http.HandleFunc("/upload",
		middleware.CORSMiddleware(
			middleware.OptionalAuthMiddleware(
				middleware.ErrorHandler(app.logger, app.fileHandler.UploadFile))))

	http.HandleFunc("/api/history",
		middleware.CORSMiddleware(
			middleware.AuthMiddleware(
				middleware.ErrorHandler(app.logger, app.fileHandler.GetUserFiles))))

	http.HandleFunc("/api/upload/",
		middleware.CORSMiddleware(
			middleware.OptionalAuthMiddleware(
				middleware.ErrorHandler(app.logger, app.fileHandler.GetFileStatus))))

	// Статические файлы
	http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir(app.config.UploadPath))))

	app.logger.Info("Server starting on http://localhost:%s", app.config.Port)
	return http.ListenAndServe(":"+app.config.Port, nil)
}

func main() {
	cfg := config.NewConfig()

	app, err := NewApplication(cfg)
	if err != nil {
		fmt.Printf("FATAL: %v\n", err)
		os.Exit(1)
	}

	if err := app.Run(); err != nil {
		app.logger.Fatal("Server failed: %v", err)
	}
}
