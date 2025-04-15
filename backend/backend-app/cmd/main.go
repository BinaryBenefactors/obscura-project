package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
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

type UploadedFile struct {
    ID         string    `json:"id"`
    Original   string    `json:"original_name"`
    Processed  string    `json:"processed_path"`  // Путь к обработанному файлу от ML
    UploadedAt time.Time `json:"uploaded_at"`
    Status     string    `json:"status"`          // "processing", "completed", "failed"
}

type Application struct {
    config  *Config
    logger  *logger.Logger
    uploads map[string][]UploadedFile  // Пока храним в памяти. Позже заменим на БД
    mu      sync.RWMutex               // Для безопасного доступа к map из горутин
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

type contextKey string

const (
    userIDKey contextKey = "userID"
)

// Инициализация (теперь все зависимости создаются внутри)
func NewApplication(cfg *Config) (*Application, error) {
    l, err := logger.NewLogger(cfg.LogFilePath)
    if err != nil {
        return nil, fmt.Errorf("logger init failed: %w", err)
    }

    return &Application{
        config:  cfg,
        logger:  l,
        uploads: make(map[string][]UploadedFile),
    }, nil
}

func (app *Application) Run() error {
    defer app.logger.Close()

    if err := app.setupUploadDir(); err != nil {
        return fmt.Errorf("setup failed: %w", err)
    }

    // Обновленные роуты:
    http.HandleFunc("/upload", app.corsMiddleware(app.optionalAuthMiddleware(app.errorHandler(app.uploadFileHandler))))
    http.HandleFunc("/api/history", app.corsMiddleware(app.authMiddleware(app.errorHandler(app.historyHandler))))
    http.HandleFunc("/api/upload/", app.corsMiddleware(app.optionalAuthMiddleware(app.errorHandler(app.uploadStatusHandler))))
    
    // Статические файлы
    http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir(app.config.UploadPath))))
    
    // В прошлом для хоста index.html через сервер Golang
    // http.Handle("/", http.FileServer(http.Dir(app.config.StaticDir)))

    app.logger.Info("Server starting on http://localhost:%s", app.config.Port)
    return http.ListenAndServe(":"+app.config.Port, nil)
}

// Необязательная аутентификация
func (app *Application) optionalAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var userID string
        
        // Пробуем получить токен
        authHeader := r.Header.Get("Authorization")
        if authHeader != "" {
            token := strings.TrimPrefix(authHeader, "Bearer ")
            if token == "demo_token" { // Валидный токен
                userID = "user123"
            }
        }

        // Если не авторизован - создаем анонимную сессию
        if userID == "" {
            sessionID := r.Header.Get("X-Session-ID")
            if sessionID == "" {
                sessionID = app.generateFileID() // Генерируем временный ID
            }
            userID = "anon_" + sessionID
        }

        ctx := context.WithValue(r.Context(), userIDKey, userID)
        next.ServeHTTP(w, r.WithContext(ctx))
    }
}

func (app *Application) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Временная заглушка - в реальности проверяем JWT или сессию
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            w.WriteHeader(http.StatusUnauthorized)
            json.NewEncoder(w).Encode(ErrorResponse{Error: "Authorization required"})
            return
        }

        // Эмуляция проверки токена (в реальности будем использовать jwt.Parse и т.д.)
        token := strings.TrimPrefix(authHeader, "Bearer ")
        if token != "demo_token" { // Заменим на реальную проверку
            w.WriteHeader(http.StatusForbidden)
            json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid token"})
            return
        }

        // Добавляем userID в контекст
        ctx := context.WithValue(r.Context(), userIDKey, "user123") // Временное значение
        next.ServeHTTP(w, r.WithContext(ctx))
    }
}

func (app *Application) corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if r.Method == "OPTIONS" {
            return
        }
        
        next.ServeHTTP(w, r)
    }
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

// GET /api/upload/{id} - статус конкретной загрузки
func (app *Application) uploadStatusHandler(w http.ResponseWriter, r *http.Request) error {
    if r.Method != http.MethodGet {
        return app.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
    }

    fileID := strings.TrimPrefix(r.URL.Path, "/api/upload/")
    if fileID == "" {
        return app.sendError(w, "File ID required", http.StatusBadRequest)
    }

    app.mu.RLock()
    defer app.mu.RUnlock()
    
    for _, file := range app.uploads["user123"] {
        if file.ID == fileID {
            w.Header().Set("Content-Type", "application/json")
            return json.NewEncoder(w).Encode(file)
        }
    }

    return app.sendError(w, "File not found", http.StatusNotFound)
}

// GET /api/history - список всех загрузок
func (app *Application) historyHandler(w http.ResponseWriter, r *http.Request) error {
    if r.Method != http.MethodGet {
        return app.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
    }

    userID, ok := r.Context().Value("userID").(string)
    if !ok || strings.HasPrefix(userID, "anon_") {
        return app.sendError(w, "Authentication required", http.StatusUnauthorized)
    }

    app.mu.RLock()
    defer app.mu.RUnlock()
    
    userUploads := app.uploads[userID] // Теперь используем реальный userID

    w.Header().Set("Content-Type", "application/json")
    return json.NewEncoder(w).Encode(userUploads)
}

func (app *Application) uploadFileHandler(w http.ResponseWriter, r *http.Request) error {
    userID, _ := r.Context().Value("userID").(string)
	
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

	fileID := app.generateFileID()
    processedPath := fmt.Sprintf("/processed/%s_result.mp4", fileID)

    app.mu.Lock()
    defer app.mu.Unlock()
    
    app.uploads[userID] = append(app.uploads[userID], UploadedFile{
        ID:         fileID,
        Original:   fileHeader.Filename,
        Processed:  processedPath,
        UploadedAt: time.Now(),
        Status:     "processing", // Статус изменится при получении ответа от ML
    })

    // Эмуляция async-обработки ML (заглушка)
    go app.simulateMLProcessing(fileID, "user123")

    // Ответ фронту
    response := struct {
        SuccessResponse
        FileID string `json:"file_id"`
    }{
        SuccessResponse: SuccessResponse{
            FileURL:  "/files/" + newFileName,
            FileType: fileType,
            Message:  "File is being processed",
        },
        FileID: fileID,
    }

	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated) // Явно указываем статус 201
    return json.NewEncoder(w).Encode(response)
}



func (app *Application) simulateMLProcessing(fileID, userID string) {
    time.Sleep(5 * time.Second) // Эмуляция работы ML

    app.mu.Lock()
    defer app.mu.Unlock()
    
    for i, file := range app.uploads[userID] {
        if file.ID == fileID {
            app.uploads[userID][i].Status = "completed"
            app.uploads[userID][i].Processed = fmt.Sprintf("/processed/%s_processed.mp4", fileID)
            break
        }
    }
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

func (app *Application) generateFileID() string {
    return fmt.Sprintf("%d", time.Now().UnixNano())
}