package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"obscura.app/backend/pkg/logger"
)

type Config struct {
	MaxUploadSize int64
	UploadPath    string
	LogFilePath   string
	StaticDir     string
	Port          string
}

type Application struct {
	config *Config
	logger *logger.Logger
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

type SuccessResponse struct {
	FileURL  string `json:"fileUrl"`
	FileType string `json:"fileType"`
	Message  string `json:"message"`
}

// Инициализация (теперь все зависимости создаются внутри)
func NewApplication(cfg *Config) (*Application, error) {
	l, err := logger.NewLogger(cfg.LogFilePath)
	if err != nil {
		return nil, fmt.Errorf("logger init failed: %w", err)
	}

	return &Application{
		config: cfg,
		logger: l,
	}, nil
}

func (app *Application) Run() error {
	defer app.logger.Close()

	if err := app.setupUploadDir(); err != nil {
		return fmt.Errorf("setup failed: %w", err)
	}

	http.HandleFunc("/upload", app.errorHandler(app.uploadFileHandler))
	http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir(app.config.UploadPath))))
	http.Handle("/", http.FileServer(http.Dir(app.config.StaticDir)))

	app.logger.Info("Server starting on http://localhost:%s", app.config.Port)
	return http.ListenAndServe(":"+app.config.Port, nil)
}

func (app *Application) setupUploadDir() error {
	if err := os.MkdirAll(app.config.UploadPath, os.ModePerm); err != nil {
		app.logger.Error("Upload dir creation failed: %v", err)
		return fmt.Errorf("mkdir failed: %w", err)
	}
	return nil
}

func main() {
	cfg := &Config{
		MaxUploadSize: 50 * 1024 * 1024,
		UploadPath:    "./uploads",
		LogFilePath:   "./app.log",
		StaticDir:     "../../../frontend/static",
		Port:          "8080",
	}

	app, err := NewApplication(cfg)
	if err != nil {
		fmt.Printf("FATAL: %v\n", err)
		os.Exit(1)
	}

    // app.logger.Info("This is an info message")
    // app.logger.Debug("This is a debug message")
    // app.logger.Warning("This is a warning message")
    // app.logger.Error("This is an error message")
    // app.logger.Fatal("This is a fatal message")

	if err := app.Run(); err != nil {
		app.logger.Fatal("Server failed: %v", err)
	}
}


func (app *Application) errorHandler(h func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        err := h(w, r)
        if err != nil {
            app.logger.Error("Request handling error: %v", err)
            
            response := ErrorResponse{
                Error:   "Internal Server Error",
                Details: err.Error(),
            }
            
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusInternalServerError)
            if err := json.NewEncoder(w).Encode(response); err != nil {
                app.logger.Error("Failed to encode error response: %v", err)
            }
        }
    }
}

func (app *Application) uploadFileHandler(w http.ResponseWriter, r *http.Request) error {
    // Логирование начала обработки запроса
    app.logger.Debug("Starting file upload processing")

    if r.Method != http.MethodPost {
        return app.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
    }

    // Проверка размера файла с использованием конфига
    if r.ContentLength > app.config.MaxUploadSize {
        return app.sendError(w, 
            fmt.Sprintf("File too large (max %dMB)", 
                app.config.MaxUploadSize/(1024*1024)), 
            http.StatusBadRequest)
    }

    if err := r.ParseMultipartForm(10 << 20); err != nil {
        app.logger.Warning("Failed to parse multipart form: %v", err)
        return app.sendError(w, "Failed to parse form data", http.StatusBadRequest)
    }

    file, fileHeader, err := r.FormFile("mediaFile")
    if err != nil {
        app.logger.Warning("No file uploaded: %v", err)
        return app.sendError(w, "No file uploaded", http.StatusBadRequest)
    }
    defer file.Close()

    fileType := fileHeader.Header.Get("Content-Type")
    if !app.isAllowedFileType(fileType) {
        app.logger.Warning("Attempt to upload forbidden file type: %s", fileType)
        return app.sendError(w, "Only JPEG, PNG images and MP4 videos are allowed", http.StatusBadRequest)
    }

    ext := filepath.Ext(fileHeader.Filename)
    if !app.isAllowedExtension(ext) {
        app.logger.Warning("Attempt to upload forbidden file extension: %s", ext)
        return app.sendError(w, "Invalid file extension", http.StatusBadRequest)
    }

    newFileName := app.generateFilename(fileHeader.Filename)
    filePath := filepath.Join(app.config.UploadPath, newFileName)

    dst, err := os.Create(filePath)
    if err != nil {
        app.logger.Error("Failed to create file %s: %v", filePath, err)
        return fmt.Errorf("failed to create file: %w", err)
    }
    defer dst.Close()

    if _, err := io.Copy(dst, file); err != nil {
        app.logger.Error("Failed to save file %s: %v", filePath, err)
        return fmt.Errorf("failed to save file: %w", err)
    }

    app.logger.Info("File uploaded successfully: %s (%s)", newFileName, fileType)

    response := SuccessResponse{
        FileURL:  "/files/" + newFileName,
        FileType: fileType,
        Message:  "File uploaded successfully",
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    return json.NewEncoder(w).Encode(response)
}

func (app *Application) sendError(w http.ResponseWriter, message string, statusCode int) error {
    response := ErrorResponse{
        Error: message,
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    return json.NewEncoder(w).Encode(response)
}

func (app *Application) isAllowedFileType(fileType string) bool {
    allowedTypes := []string{
        "image/jpeg",
        "image/png",
        "video/mp4",
        "video/quicktime",
    }

    for _, t := range allowedTypes {
        if fileType == t {
            return true
        }
    }
    return false
}

func (app *Application) isAllowedExtension(ext string) bool {
    allowedExts := []string{".jpg", ".jpeg", ".png", ".mp4", ".mov"}
    for _, e := range allowedExts {
        if ext == e {
            return true
        }
    }
    return false
}

func (app *Application) generateFilename(original string) string {
    ext := filepath.Ext(original)
    return fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
}