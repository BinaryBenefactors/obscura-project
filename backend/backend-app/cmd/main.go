package main

import (
	"log"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"obscura.app/backend/internal/logger"
)

const (
	maxUploadSize = 50 * 1024 * 1024 // 50 MB
	uploadPath    = "./uploads"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

type SuccessResponse struct {
	FileURL  string `json:"fileUrl"`
	FileType string `json:"fileType"`
	Message  string `json:"message"`
}

func main() {
	// Создаем папку для загрузок
	if err := setupUploadDir(); err != nil {
		log.Fatalf("Failed to setup upload directory: %v", err)
	}

	// Настраиваем маршруты
	http.HandleFunc("/upload", errorHandler(uploadFileHandler))
	http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir(uploadPath))))

	// Статический сервер
	http.Handle("/", http.FileServer(http.Dir("../../../frontend/static")))

	port := "8080"
	log.Printf("Server started on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

	logger, err := logger.NewLogger("cmd/app.log")
    if err != nil {
        log.Fatalf("Could not initialize logger: %v", err)
    }
    defer logger.Close()

    // logger.Info("This is an info message")
    // logger.Debug("This is a debug message")
    // logger.Warning("This is a warning message")
    // logger.Error("This is an error message")
    // logger.Fatal("This is a fatal message")
}

func setupUploadDir() error {
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		return fmt.Errorf("could not create upload directory: %w", err)
	}
	return nil
}

// Обработчик ошибок для всех HTTP хендлеров
func errorHandler(h func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
		if err != nil {
			log.Printf("Error handling request: %v", err)
			
			response := ErrorResponse{
				Error:   "Internal Server Error",
				Details: err.Error(),
			}
			
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
		}
	}
}

func uploadFileHandler(w http.ResponseWriter, r *http.Request) error {
	// Проверка метода
	if r.Method != "POST" {
		return sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	// Проверка размера файла
	if r.ContentLength > maxUploadSize {
		return sendError(w, fmt.Sprintf("File too large (max %dMB)", maxUploadSize/(1024*1024)), http.StatusBadRequest)
	}

	// Парсинг multipart формы
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		return sendError(w, "Failed to parse form data", http.StatusBadRequest)
	}

	// Получение файла из формы
	file, fileHeader, err := r.FormFile("mediaFile")
	if err != nil {
		return sendError(w, "No file uploaded", http.StatusBadRequest)
	}
	defer file.Close()

	// Проверка типа файла
	fileType := fileHeader.Header.Get("Content-Type")
	if !isAllowedFileType(fileType) {
		return sendError(w, "Only JPEG, PNG images and MP4 videos are allowed", http.StatusBadRequest)
	}

	// Проверка расширения файла
	ext := filepath.Ext(fileHeader.Filename)
	if !isAllowedExtension(ext) {
		return sendError(w, "Invalid file extension", http.StatusBadRequest)
	}

	// Создание файла на сервере
	newFileName := generateFilename(fileHeader.Filename)
	filePath := filepath.Join(uploadPath, newFileName)

	dst, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	// Копирование содержимого файла
	if _, err := io.Copy(dst, file); err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	// Формирование успешного ответа
	response := SuccessResponse{
		FileURL:  "/files/" + newFileName,
		FileType: fileType,
		Message:  "File uploaded successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(response)
}

func sendError(w http.ResponseWriter, message string, statusCode int) error {
	response := ErrorResponse{
		Error: message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(response)
}

func isAllowedFileType(fileType string) bool {
	allowedTypes := []string{
		"image/jpeg",
		"image/png",
		"video/mp4",
		"video/quicktime", // для MOV файлов
	}

	for _, t := range allowedTypes {
		if fileType == t {
			return true
		}
	}
	return false
}

func isAllowedExtension(ext string) bool {
	allowedExts := []string{".jpg", ".jpeg", ".png", ".mp4", ".mov"}
	for _, e := range allowedExts {
		if ext == e {
			return true
		}
	}
	return false
}

func generateFilename(original string) string {
	ext := filepath.Ext(original)
	return fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
}