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
	
	// Загрузка файлов - доступна для всех (с опциональной авторизацией)
	s.router.HandleFunc("/api/upload", s.corsMiddleware(s.optionalAuthMiddleware(s.handleUpload)))
	
	// Получение списка файлов - только для авторизованных
	s.router.HandleFunc("/api/files", s.corsMiddleware(s.authMiddleware(s.handleGetFiles)))
	
	// Действия с файлами - доступны для всех (анонимные файлы по ID, пользовательские с проверкой)
	s.router.HandleFunc("/api/files/", s.corsMiddleware(s.optionalAuthMiddleware(s.handleFileActions)))
}

// GetRouter возвращает HTTP роутер сервера
func (s *Server) GetRouter() *http.ServeMux {
	return s.router
}

// CORS middleware
func (s *Server) corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.logger.Debug("Request: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			s.logger.Debug("Handling OPTIONS request for %s", r.URL.Path)
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

// Опциональный middleware для аутентификации (не требует авторизации, но проверяет если токен есть)
func (s *Server) optionalAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.logger.Debug("Optional auth middleware for %s %s", r.Method, r.URL.Path)
		
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			// Нет авторизации - продолжаем как анонимный пользователь
			s.logger.Debug("No authorization header - proceeding as anonymous user")
			r.Header.Set("X-User-ID", "0") // 0 означает анонимный пользователь
			next(w, r)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			s.logger.Debug("Invalid Bearer token format - proceeding as anonymous user")
			r.Header.Set("X-User-ID", "0")
			next(w, r)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(s.config.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			s.logger.Debug("Invalid token - proceeding as anonymous user: %v", err)
			r.Header.Set("X-User-ID", "0")
			next(w, r)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			s.logger.Debug("Invalid token claims - proceeding as anonymous user")
			r.Header.Set("X-User-ID", "0")
			next(w, r)
			return
		}

		userID, err := strconv.Atoi(fmt.Sprintf("%.0f", claims["user_id"]))
		if err != nil {
			s.logger.Debug("Invalid user ID in token - proceeding as anonymous user: %v", err)
			r.Header.Set("X-User-ID", "0")
			next(w, r)
			return
		}

		s.logger.Debug("Authenticated user %d for %s", userID, r.URL.Path)
		r.Header.Set("X-User-ID", strconv.Itoa(userID))
		next(w, r)
	}
}

// Обязательный middleware для аутентификации
func (s *Server) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.logger.Debug("Auth middleware for %s %s", r.Method, r.URL.Path)
		
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			s.logger.Warning("Missing Authorization header for %s", r.URL.Path)
			s.sendError(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			s.logger.Warning("Invalid Bearer token format for %s", r.URL.Path)
			s.sendError(w, "Bearer token required", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(s.config.JWTSecret), nil
		})

		if err != nil {
			s.logger.Warning("JWT parse error for %s: %v", r.URL.Path, err)
			s.sendError(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			s.logger.Warning("Invalid JWT token for %s", r.URL.Path)
			s.sendError(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			s.logger.Warning("Invalid token claims for %s", r.URL.Path)
			s.sendError(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		userID, err := strconv.Atoi(fmt.Sprintf("%.0f", claims["user_id"]))
		if err != nil {
			s.logger.Warning("Invalid user ID in token for %s: %v", r.URL.Path, err)
			s.sendError(w, "Invalid user ID in token", http.StatusUnauthorized)
			return
		}

		s.logger.Debug("Authenticated user %d for %s", userID, r.URL.Path)
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
		s.logger.Warning("Invalid method %s for register endpoint", r.Method)
		s.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.logger.Warning("Invalid JSON in register request: %v", err)
		s.sendError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	s.logger.Info("Registration attempt for email: %s", req.Email)

	// Валидация
	if req.Email == "" || req.Password == "" || req.Name == "" {
		s.logger.Warning("Missing required fields in register request")
		s.sendError(w, "All fields are required", http.StatusBadRequest)
		return
	}

	// Проверяем, существует ли пользователь
	if _, err := s.db.GetUserByEmail(req.Email); err == nil {
		s.logger.Warning("User already exists: %s", req.Email)
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
		s.logger.Error("Failed to create user %s: %v", req.Email, err)
		s.sendError(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Генерируем токен
	token, err := s.generateJWT(user.ID)
	if err != nil {
		s.logger.Error("Failed to generate token for user %d: %v", user.ID, err)
		s.sendError(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	s.logger.Info("User registered successfully: %s (ID: %d)", user.Email, user.ID)
	s.sendJSON(w, AuthResponse{
		Token: token,
		User:  *user,
	})
}

// Авторизация пользователя
func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.logger.Warning("Invalid method %s for login endpoint", r.Method)
		s.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.logger.Warning("Invalid JSON in login request: %v", err)
		s.sendError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	s.logger.Info("Login attempt for email: %s", req.Email)

	// Находим пользователя
	user, err := s.db.GetUserByEmail(req.Email)
	if err != nil {
		s.logger.Warning("Login failed - user not found: %s", req.Email)
		s.sendError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Проверяем пароль
	if !user.CheckPassword(req.Password) {
		s.logger.Warning("Login failed - invalid password for user: %s", req.Email)
		s.sendError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Генерируем токен
	token, err := s.generateJWT(user.ID)
	if err != nil {
		s.logger.Error("Failed to generate token for user %d: %v", user.ID, err)
		s.sendError(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	s.logger.Info("User logged in successfully: %s (ID: %d)", user.Email, user.ID)
	s.sendJSON(w, AuthResponse{
		Token: token,
		User:  *user,
	})
}

// Загрузка файла (доступна для всех)
func (s *Server) handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.logger.Warning("Invalid method %s for upload endpoint", r.Method)
		s.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userIDStr := r.Header.Get("X-User-ID")
	userID, _ := strconv.Atoi(userIDStr)
	isAnonymous := userID == 0

	if isAnonymous {
		s.logger.Info("Anonymous file upload started")
	} else {
		s.logger.Info("File upload started for user %d", userID)
	}

	// Парсим multipart форму
	err := r.ParseMultipartForm(s.config.MaxFileSize)
	if err != nil {
		s.logger.Warning("Failed to parse multipart form: %v", err)
		s.sendError(w, "File too large or invalid form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		s.logger.Warning("No file provided in upload request: %v", err)
		s.sendError(w, "No file provided", http.StatusBadRequest)
		return
	}
	defer file.Close()

	s.logger.Info("Processing file upload: %s (%d bytes) %s", header.Filename, header.Size, 
		func() string {
			if isAnonymous {
				return "for anonymous user"
			}
			return fmt.Sprintf("for user %d", userID)
		}())

	// Генерируем уникальное имя файла
	fileID := uuid.New().String()
	ext := filepath.Ext(header.Filename)
	fileName := fileID + ext
	filePath := filepath.Join(s.config.UploadPath, fileName)

	s.logger.Debug("Saving file to: %s", filePath)

	// Создаем файл на диске
	dst, err := os.Create(filePath)
	if err != nil {
		s.logger.Error("Failed to create file %s: %v", filePath, err)
		s.sendError(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Копируем содержимое
	size, err := io.Copy(dst, file)
	if err != nil {
		s.logger.Error("Failed to copy file content: %v", err)
		os.Remove(filePath) // Удаляем частично созданный файл
		s.sendError(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	s.logger.Info("File saved to disk: %s (%d bytes)", filePath, size)

	// Определяем MIME-тип
	mimeType := header.Header.Get("Content-Type")
	if mimeType == "" {
		switch strings.ToLower(ext) {
		case ".jpg", ".jpeg":
			mimeType = "image/jpeg"
		case ".png":
			mimeType = "image/png"
		case ".gif":
			mimeType = "image/gif"
		case ".mp4":
			mimeType = "video/mp4"
		case ".avi":
			mimeType = "video/avi"
		case ".mov":
			mimeType = "video/quicktime"
		case ".webm":
			mimeType = "video/webm"
		default:
			mimeType = "application/octet-stream"
		}
	}

	// Создаем запись о файле
	fileRecord := &File{
		ID:           fileID,
		UserID:       uint(userID), // 0 для анонимных пользователей
		OriginalName: header.Filename,
		FileName:     fileName,
		FileSize:     size,
		MimeType:     mimeType,
		Status:       StatusUploaded,
		UploadedAt:   time.Now(),
	}

	// Сохраняем в БД только для авторизованных пользователей
	if !isAnonymous {
		if err := s.db.CreateFile(fileRecord); err != nil {
			s.logger.Error("Failed to save file record for user %d: %v", userID, err)
			os.Remove(filePath) // Удаляем файл при ошибке БД
			s.sendError(w, "Failed to save file record", http.StatusInternalServerError)
			return
		}
		s.logger.Info("File upload completed and saved to history: %s (ID: %s) for user %d", header.Filename, fileID, userID)
	} else {
		s.logger.Info("Anonymous file upload completed: %s (ID: %s) - not saved to history", header.Filename, fileID)
	}

	s.sendJSON(w, SuccessResponse{
		Message: func() string {
			if isAnonymous {
				return "File uploaded successfully (not saved to history)"
			}
			return "File uploaded successfully"
		}(),
		Data: fileRecord,
	})
}

// Получение списка файлов пользователя (только для авторизованных)
func (s *Server) handleGetFiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.logger.Warning("Invalid method %s for get files endpoint", r.Method)
		s.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := strconv.Atoi(r.Header.Get("X-User-ID"))
	if err != nil {
		s.logger.Error("Invalid user ID in get files request: %v", err)
		s.sendError(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	s.logger.Debug("Getting files list for user %d", userID)

	files, err := s.db.GetUserFiles(uint(userID))
	if err != nil {
		s.logger.Error("Failed to get user files for user %d: %v", userID, err)
		s.sendError(w, "Failed to get files", http.StatusInternalServerError)
		return
	}

	s.logger.Info("Retrieved %d files for user %d", len(files), userID)

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
		s.logger.Warning("Empty file ID in file actions request")
		s.sendError(w, "File ID required", http.StatusBadRequest)
		return
	}

	userIDStr := r.Header.Get("X-User-ID")
	userID, _ := strconv.Atoi(userIDStr)
	isAnonymous := userID == 0

	s.logger.Debug("File action %s for file %s by %s", r.Method, fileID, 
		func() string {
			if isAnonymous {
				return "anonymous user"
			}
			return fmt.Sprintf("user %d", userID)
		}())

	switch r.Method {
	case http.MethodGet:
		s.handleDownloadFileByID(w, r, fileID, userID, isAnonymous)
	case http.MethodDelete:
		s.handleDeleteFileByID(w, r, fileID, userID, isAnonymous)
	default:
		s.logger.Warning("Method %s not allowed for file actions", r.Method)
		s.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Скачивание файла по ID
func (s *Server) handleDownloadFileByID(w http.ResponseWriter, r *http.Request, fileID string, userID int, isAnonymous bool) {
	if isAnonymous {
		// Для анонимных пользователей просто проверяем существование файла на диске
		fileName := fileID + "*" // Ищем файл с любым расширением
		matches, err := filepath.Glob(filepath.Join(s.config.UploadPath, fileName))
		if err != nil || len(matches) == 0 {
			s.logger.Warning("Anonymous file not found on disk: %s", fileID)
			s.sendError(w, "File not found", http.StatusNotFound)
			return
		}
		
		filePath := matches[0]
		s.logger.Info("Serving anonymous file: %s", fileID)
		
		// Определяем MIME-тип по расширению
		ext := strings.ToLower(filepath.Ext(filePath))
		var mimeType string
		switch ext {
		case ".jpg", ".jpeg":
			mimeType = "image/jpeg"
		case ".png":
			mimeType = "image/png"
		case ".gif":
			mimeType = "image/gif"
		case ".mp4":
			mimeType = "video/mp4"
		case ".avi":
			mimeType = "video/avi"
		case ".mov":
			mimeType = "video/quicktime"
		case ".webm":
			mimeType = "video/webm"
		default:
			mimeType = "application/octet-stream"
		}
		
		w.Header().Set("Content-Type", mimeType)
		http.ServeFile(w, r, filePath)
		return
	}

	// Для авторизованных пользователей проверяем БД
	file, err := s.db.GetFileByID(fileID)
	if err != nil {
		s.logger.Warning("File not found in database: %s for user %d", fileID, userID)
		s.sendError(w, "File not found", http.StatusNotFound)
		return
	}

	// Проверяем, что файл принадлежит пользователю
	if file.UserID != uint(userID) {
		s.logger.Warning("Access denied: user %d tried to access file %s owned by user %d", userID, fileID, file.UserID)
		s.sendError(w, "Access denied", http.StatusForbidden)
		return
	}

	s.handleDownloadFile(w, r, file)
}

// Удаление файла по ID
func (s *Server) handleDeleteFileByID(w http.ResponseWriter, r *http.Request, fileID string, userID int, isAnonymous bool) {
	if isAnonymous {
		// Анонимные пользователи не могут удалять файлы
		s.logger.Warning("Anonymous user attempted to delete file: %s", fileID)
		s.sendError(w, "Anonymous users cannot delete files", http.StatusForbidden)
		return
	}

	// Для авторизованных пользователей
	file, err := s.db.GetFileByID(fileID)
	if err != nil {
		s.logger.Warning("File not found for deletion: %s for user %d", fileID, userID)
		s.sendError(w, "File not found", http.StatusNotFound)
		return
	}

	// Проверяем, что файл принадлежит пользователю
	if file.UserID != uint(userID) {
		s.logger.Warning("Access denied: user %d tried to delete file %s owned by user %d", userID, fileID, file.UserID)
		s.sendError(w, "Access denied", http.StatusForbidden)
		return
	}

	s.handleDeleteFile(w, r, file)
}

// Скачивание файла
func (s *Server) handleDownloadFile(w http.ResponseWriter, r *http.Request, file *File) {
	filePath := filepath.Join(s.config.UploadPath, file.FileName)

	s.logger.Debug("Downloading file: %s (path: %s)", file.ID, filePath)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		s.logger.Error("File not found on disk: %s", filePath)
		s.sendError(w, "File not found on disk", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", file.OriginalName))
	w.Header().Set("Content-Type", file.MimeType)

	s.logger.Info("Serving file: %s (%s) to user %d", file.OriginalName, file.ID, file.UserID)
	http.ServeFile(w, r, filePath)
}

// Удаление файла
func (s *Server) handleDeleteFile(w http.ResponseWriter, r *http.Request, file *File) {
	s.logger.Info("Deleting file: %s (%s) for user %d", file.OriginalName, file.ID, file.UserID)

	// Удаляем файл с диска
	filePath := filepath.Join(s.config.UploadPath, file.FileName)
	if err := os.Remove(filePath); err != nil {
		s.logger.Warning("Failed to remove file from disk: %s - %v", filePath, err)
	} else {
		s.logger.Debug("File removed from disk: %s", filePath)
	}

	// Удаляем запись из БД
	if err := s.db.DeleteFile(file.ID); err != nil {
		s.logger.Error("Failed to delete file record %s: %v", file.ID, err)
		s.sendError(w, "Failed to delete file", http.StatusInternalServerError)
		return
	}

	s.logger.Info("File deleted successfully: %s", file.ID)

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
	if err := json.NewEncoder(w).Encode(data); err != nil {
		s.logger.Error("Failed to encode JSON response: %v", err)
	}
}

// Отправка ошибки
func (s *Server) sendError(w http.ResponseWriter, message string, status int) {
	s.logger.Warning("Sending error response: %d - %s", status, message)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}
