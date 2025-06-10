package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"obscura.app/pkg/logger"
)

type Server struct {
	config *Config
	db     *Database
	logger *logger.Logger
	router *http.ServeMux
}

func NewServer(config *Config, db *Database, logger *logger.Logger) *Server {
	return &Server{
		config: config,
		db:     db,
		logger: logger,
		router: http.NewServeMux(),
	}
}

func (s *Server) SetupRoutes() {
	// Middleware для CORS
	s.router.HandleFunc("/", s.corsMiddleware(s.handleRoot))

	// API маршруты
	s.router.HandleFunc("/api/register", s.corsMiddleware(s.handleRegister))
	s.router.HandleFunc("/api/login", s.corsMiddleware(s.handleLogin))
	s.router.HandleFunc("/api/upload", s.corsMiddleware(s.authMiddleware(s.handleUpload)))
	s.router.HandleFunc("/api/files", s.corsMiddleware(s.authMiddleware(s.handleGetFiles)))
	s.router.HandleFunc("/api/files/", s.corsMiddleware(s.authMiddleware(s.handleFileActions)))
}

// GetRouter возвращает HTTP роутер сервера
func (s *Server) GetRouter() *http.ServeMux {
	return s.router
}

// CORS middleware
func (s *Server) corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

// Middleware для аутентификации
func (s *Server) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			s.sendError(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			s.sendError(w, "Bearer token required", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(s.config.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			s.sendError(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			s.sendError(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		userID, err := strconv.Atoi(fmt.Sprintf("%.0f", claims["user_id"]))
		if err != nil {
			s.sendError(w, "Invalid user ID in token", http.StatusUnauthorized)
			return
		}

		// Добавляем user_id в контекст запроса
		r.Header.Set("X-User-ID", strconv.Itoa(userID))
		next(w, r)
	}
}

// Обработчик корневого маршрута
func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	s.sendJSON(w, map[string]string{
		"message": "Obscura API",
		"version": "1.0.0",
	})
}

// Регистрация пользователя
func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Валидация
	if req.Email == "" || req.Password == "" || req.Name == "" {
		s.sendError(w, "All fields are required", http.StatusBadRequest)
		return
	}

	// Проверяем, существует ли пользователь
	if _, err := s.db.GetUserByEmail(req.Email); err == nil {
		s.sendError(w, "User already exists", http.StatusConflict)
		return
	}

	// Создаем пользователя
	user := &User{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	}

	if err := s.db.CreateUser(user); err != nil {
		s.logger.Error("Failed to create user: %v", err)
		s.sendError(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Генерируем токен
	token, err := s.generateJWT(user.ID)
	if err != nil {
		s.logger.Error("Failed to generate token: %v", err)
		s.sendError(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	s.sendJSON(w, AuthResponse{
		Token: token,
		User:  *user,
	})
}

// Авторизация пользователя
func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Находим пользователя
	user, err := s.db.GetUserByEmail(req.Email)
	if err != nil {
		s.sendError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Проверяем пароль
	if !user.CheckPassword(req.Password) {
		s.sendError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Генерируем токен
	token, err := s.generateJWT(user.ID)
	if err != nil {
		s.logger.Error("Failed to generate token: %v", err)
		s.sendError(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	s.sendJSON(w, AuthResponse{
		Token: token,
		User:  *user,
	})
}

// Загрузка файла
func (s *Server) handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := strconv.Atoi(r.Header.Get("X-User-ID"))
	if err != nil {
		s.sendError(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Парсим multipart форму
	err = r.ParseMultipartForm(s.config.MaxFileSize)
	if err != nil {
		s.sendError(w, "File too large", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		s.sendError(w, "No file provided", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Генерируем уникальное имя файла
	fileID := uuid.New().String()
	ext := filepath.Ext(header.Filename)
	fileName := fileID + ext
	filePath := filepath.Join(s.config.UploadPath, fileName)

	// Создаем файл на диске
	dst, err := os.Create(filePath)
	if err != nil {
		s.logger.Error("Failed to create file: %v", err)
		s.sendError(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Копируем содержимое
	size, err := io.Copy(dst, file)
	if err != nil {
		s.logger.Error("Failed to copy file: %v", err)
		s.sendError(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	// Сохраняем информацию в БД
	fileRecord := &File{
		ID:           fileID,
		UserID:       uint(userID),
		OriginalName: header.Filename,
		FileName:     fileName,
		FileSize:     size,
		MimeType:     header.Header.Get("Content-Type"),
		Status:       StatusUploaded,
		UploadedAt:   time.Now(),
	}

	if err := s.db.CreateFile(fileRecord); err != nil {
		s.logger.Error("Failed to save file record: %v", err)
		os.Remove(filePath) // Удаляем файл при ошибке
		s.sendError(w, "Failed to save file record", http.StatusInternalServerError)
		return
	}

	s.sendJSON(w, SuccessResponse{
		Message: "File uploaded successfully",
		Data:    fileRecord,
	})
}

// Получение списка файлов пользователя
func (s *Server) handleGetFiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := strconv.Atoi(r.Header.Get("X-User-ID"))
	if err != nil {
		s.sendError(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	files, err := s.db.GetUserFiles(uint(userID))
	if err != nil {
		s.logger.Error("Failed to get user files: %v", err)
		s.sendError(w, "Failed to get files", http.StatusInternalServerError)
		return
	}

	s.sendJSON(w, SuccessResponse{
		Message: "Files retrieved successfully",
		Data:    files,
	})
}

// Действия с файлами (получение, удаление)
func (s *Server) handleFileActions(w http.ResponseWriter, r *http.Request) {
	// Извлекаем ID файла из URL
	path := strings.TrimPrefix(r.URL.Path, "/api/files/")
	fileID := strings.Split(path, "/")[0]

	if fileID == "" {
		s.sendError(w, "File ID required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(r.Header.Get("X-User-ID"))
	if err != nil {
		s.sendError(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	file, err := s.db.GetFileByID(fileID)
	if err != nil {
		s.sendError(w, "File not found", http.StatusNotFound)
		return
	}

	// Проверяем, что файл принадлежит пользователю
	if file.UserID != uint(userID) {
		s.sendError(w, "Access denied", http.StatusForbidden)
		return
	}

	switch r.Method {
	case http.MethodGet:
		s.handleDownloadFile(w, r, file)
	case http.MethodDelete:
		s.handleDeleteFile(w, r, file)
	default:
		s.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Скачивание файла
func (s *Server) handleDownloadFile(w http.ResponseWriter, r *http.Request, file *File) {
	filePath := filepath.Join(s.config.UploadPath, file.FileName)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		s.sendError(w, "File not found on disk", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", file.OriginalName))
	w.Header().Set("Content-Type", file.MimeType)

	http.ServeFile(w, r, filePath)
}

// Удаление файла
func (s *Server) handleDeleteFile(w http.ResponseWriter, r *http.Request, file *File) {
	// Удаляем файл с диска
	filePath := filepath.Join(s.config.UploadPath, file.FileName)
	if err := os.Remove(filePath); err != nil {
		s.logger.Warning("Failed to remove file from disk: %v", err)
	}

	// Удаляем запись из БД
	if err := s.db.DeleteFile(file.ID); err != nil {
		s.logger.Error("Failed to delete file record: %v", err)
		s.sendError(w, "Failed to delete file", http.StatusInternalServerError)
		return
	}

	s.sendJSON(w, SuccessResponse{
		Message: "File deleted successfully",
	})
}

// Генерация JWT токена
func (s *Server) generateJWT(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 дней
	})

	return token.SignedString([]byte(s.config.JWTSecret))
}

// Отправка JSON ответа
func (s *Server) sendJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// Отправка ошибки
func (s *Server) sendError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}
