package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	httpSwagger "github.com/swaggo/http-swagger"
	"obscura.app/pkg/logger"
)

type Server struct {
	config      *Config
	db          *Database
	logger      *logger.Logger
	router      *http.ServeMux
	rateLimiter *RateLimiter
	validator   *Validator
	fileCleaner *FileCleaner
}

func NewServer(config *Config, db *Database, logger *logger.Logger) *Server {
	rateLimiter := NewRateLimiter(config.MaxAttemptsHandled, time.Duration(config.HandlerTimeout)*time.Hour)
	validator := NewValidator(config.MaxFileSize)
	fileCleaner := NewFileCleaner(config.UploadPath, logger)

	server := &Server{
		config:      config,
		db:          db,
		logger:      logger,
		router:      http.NewServeMux(),
		rateLimiter: rateLimiter,
		validator:   validator,
		fileCleaner: fileCleaner,
	}

	fileCleaner.Start()
	return server
}

func (s *Server) SetupRoutes() {
	// Swagger UI
	s.router.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	// Middleware для CORS
	s.router.HandleFunc("/", s.corsMiddleware(s.handleRoot))

	// Health check
	s.router.HandleFunc("/health", s.corsMiddleware(s.handleHealth))

	// API маршруты
	s.router.HandleFunc("/api/register", s.corsMiddleware(s.handleRegister))
	s.router.HandleFunc("/api/login", s.corsMiddleware(s.handleLogin))

	// Профиль пользователя - только для авторизованных
	s.router.HandleFunc("/api/user/profile", s.corsMiddleware(s.authMiddleware(s.handleUserProfile)))
	s.router.HandleFunc("/api/user/profile/update", s.corsMiddleware(s.authMiddleware(s.handleUpdateProfile)))

	// Загрузка файлов
	s.router.HandleFunc("/api/upload", s.corsMiddleware(s.optionalAuthMiddleware(s.handleUpload)))

	// Получение списка файлов
	s.router.HandleFunc("/api/files", s.corsMiddleware(s.authMiddleware(s.handleGetFiles)))

	// Действия с файлами
	s.router.HandleFunc("/api/files/", s.corsMiddleware(s.optionalAuthMiddleware(s.handleFileActions)))

	// Статистика для профиля
	s.router.HandleFunc("/api/user/stats", s.corsMiddleware(s.authMiddleware(s.handleUserStats)))

	// Административная информация
	s.router.HandleFunc("/api/admin/stats", s.corsMiddleware(s.handleAdminStats))
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

// Опциональный middleware для аутентификации
func (s *Server) optionalAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.logger.Debug("Optional auth middleware for %s %s", r.Method, r.URL.Path)

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			s.logger.Debug("No authorization header - proceeding as anonymous user")
			r.Header.Set("X-User-ID", "0")
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

// @Summary API information
// @Description Get basic API information and version
// @Tags system
// @Produce json
// @Success 200 {object} map[string]string
// @Router / [get]
func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	s.sendJSON(w, map[string]string{
		"message": "Obscura API",
		"author": "bambutcha (Yagolnik Daniil)",
		"version": "1.0.0",
	})
}

// @Summary Health check
// @Description Check system health including database connection, ML service and file system stats
// @Tags system
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 503 {object} ErrorResponse
// @Router /health [get]
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sqlDB, err := s.db.DB.DB()
	if err != nil {
		s.sendError(w, "Database connection error", http.StatusServiceUnavailable)
		return
	}

	if err := sqlDB.Ping(); err != nil {
		s.sendError(w, "Database ping failed", http.StatusServiceUnavailable)
		return
	}

	fileStats, err := s.fileCleaner.GetStats()
	if err != nil {
		s.logger.Warning("Failed to get file stats: %v", err)
		fileStats = map[string]interface{}{"error": "failed to get stats"}
	}

	rateLimiterStats := s.rateLimiter.GetStats()

	// Проверка ML сервиса
	mlStatus := "disabled"
	if s.config.MLServiceEnabled {
		if s.checkMLService() {
			mlStatus = "healthy"
		} else {
			mlStatus = "unhealthy"
		}
	}

	health := map[string]interface{}{
		"status":       "healthy",
		"timestamp":    time.Now(),
		"version":      "1.0.0",
		"database":     "connected",
		"ml_service":   mlStatus,
		"file_system":  fileStats,
		"rate_limiter": rateLimiterStats,
	}

	s.sendJSON(w, health)
}

// @Summary Register new user
// @Description Create a new user account with email, password and name
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "User registration data"
// @Success 200 {object} AuthResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /api/register [post]
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

	if validationErrors := s.validator.ValidateRegistration(req); len(validationErrors) > 0 {
		s.logger.Warning("Validation failed for registration: %s", req.Email)
		s.sendValidationErrors(w, validationErrors)
		return
	}

	if _, err := s.db.GetUserByEmail(req.Email); err == nil {
		s.logger.Warning("User already exists: %s", req.Email)
		s.sendError(w, "User already exists", http.StatusConflict)
		return
	}

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

// @Summary User login
// @Description Authenticate user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "User login credentials"
// @Success 200 {object} AuthResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/login [post]
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

	if validationErrors := s.validator.ValidateLogin(req); len(validationErrors) > 0 {
		s.logger.Warning("Validation failed for login")
		s.sendValidationErrors(w, validationErrors)
		return
	}

	s.logger.Info("Login attempt for email: %s", req.Email)

	user, err := s.db.GetUserByEmail(req.Email)
	if err != nil {
		s.logger.Warning("Login failed - user not found: %s", req.Email)
		s.sendError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if !user.CheckPassword(req.Password) {
		s.logger.Warning("Login failed - invalid password for user: %s", req.Email)
		s.sendError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

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

// @Summary Get user profile
// @Description Get current user profile information
// @Tags user
// @Produce json
// @Security BearerAuth
// @Success 200 {object} SuccessResponse{data=User}
// @Failure 401 {object} ErrorResponse
// @Router /api/user/profile [get]
func (s *Server) handleUserProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.logger.Warning("Invalid method %s for profile endpoint", r.Method)
		s.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := strconv.Atoi(r.Header.Get("X-User-ID"))
	if err != nil {
		s.logger.Error("Invalid user ID in profile request: %v", err)
		s.sendError(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	s.logger.Debug("Getting profile for user %d", userID)

	user, err := s.db.GetUserByID(uint(userID))
	if err != nil {
		s.logger.Error("Failed to get user profile for user %d: %v", userID, err)
		s.sendError(w, "User not found", http.StatusNotFound)
		return
	}

	s.logger.Info("Profile retrieved for user %d", userID)

	s.sendJSON(w, SuccessResponse{
		Message: "Profile retrieved successfully",
		Data:    user,
	})
}

// @Summary Update user profile
// @Description Update user profile information (name, email, password)
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UpdateProfileRequest true "Profile update data"
// @Success 200 {object} SuccessResponse{data=User}
// @Failure 401 {object} ErrorResponse
// @Router /api/user/profile/update [put]
func (s *Server) handleUpdateProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		s.logger.Warning("Invalid method %s for update profile endpoint", r.Method)
		s.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := strconv.Atoi(r.Header.Get("X-User-ID"))
	if err != nil {
		s.logger.Error("Invalid user ID in update profile request: %v", err)
		s.sendError(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.logger.Warning("Invalid JSON in update profile request: %v", err)
		s.sendError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if validationErrors := s.validator.ValidateProfileUpdate(req); len(validationErrors) > 0 {
		s.logger.Warning("Validation failed for profile update for user %d", userID)
		s.sendValidationErrors(w, validationErrors)
		return
	}

	s.logger.Info("Profile update attempt for user %d", userID)

	user, err := s.db.GetUserByID(uint(userID))
	if err != nil {
		s.logger.Error("Failed to get user for profile update %d: %v", userID, err)
		s.sendError(w, "User not found", http.StatusNotFound)
		return
	}

	updated := false
	if req.Name != "" && req.Name != user.Name {
		user.Name = req.Name
		updated = true
	}

	if req.Email != "" && req.Email != user.Email {
		if _, err := s.db.GetUserByEmail(req.Email); err == nil {
			s.logger.Warning("Email already taken during profile update: %s", req.Email)
			s.sendError(w, "Email already taken", http.StatusConflict)
			return
		}
		user.Email = req.Email
		updated = true
	}

	if req.Password != "" {
		user.Password = req.Password
		if err := user.HashPassword(); err != nil {
			s.logger.Error("Failed to hash new password for user %d: %v", userID, err)
			s.sendError(w, "Failed to update password", http.StatusInternalServerError)
			return
		}
		updated = true
	}

	if !updated {
		s.logger.Debug("No changes in profile update for user %d", userID)
		s.sendJSON(w, SuccessResponse{
			Message: "No changes made",
			Data:    user,
		})
		return
	}

	if err := s.db.UpdateUser(user); err != nil {
		s.logger.Error("Failed to update user profile %d: %v", userID, err)
		s.sendError(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	s.logger.Info("Profile updated successfully for user %d", userID)

	s.sendJSON(w, SuccessResponse{
		Message: "Profile updated successfully",
		Data:    user,
	})
}

// @Summary Get user statistics
// @Description Get detailed statistics about user's files, uploads and processing status
// @Tags user
// @Produce json
// @Security BearerAuth
// @Success 200 {object} SuccessResponse{data=UserStats}
// @Failure 401 {object} ErrorResponse
// @Router /api/user/stats [get]
func (s *Server) handleUserStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.logger.Warning("Invalid method %s for user stats endpoint", r.Method)
		s.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := strconv.Atoi(r.Header.Get("X-User-ID"))
	if err != nil {
		s.logger.Error("Invalid user ID in user stats request: %v", err)
		s.sendError(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	s.logger.Debug("Getting stats for user %d", userID)

	stats, err := s.db.GetUserStats(uint(userID))
	if err != nil {
		s.logger.Error("Failed to get user stats for user %d: %v", userID, err)
		s.sendError(w, "Failed to get stats", http.StatusInternalServerError)
		return
	}

	s.logger.Info("Stats retrieved for user %d", userID)

	s.sendJSON(w, SuccessResponse{
		Message: "Stats retrieved successfully",
		Data:    stats,
	})
}

// @Summary Upload and process file
// @Description Upload a file (images and videos) and automatically send it for ML processing. Available for both authenticated and anonymous users with rate limiting
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "File to upload"
// @Param blur_type formData string false "Type of blur to apply" Enums(gaussian, motion, pixelate) default(gaussian)
// @Param intensity formData integer false "Effect intensity (1-10)" minimum(1) maximum(10) default(5)
// @Param object_types formData string false "Comma-separated list of objects to blur" example("face,person,car")
// @Success 200 {object} SuccessResponse{data=File} "File uploaded and processing started"
// @Failure 400 {object} ErrorResponse
// @Failure 429 {object} ErrorResponse
// @Router /api/upload [post]
func (s *Server) handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.logger.Warning("Invalid method %s for upload endpoint", r.Method)
		s.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userIDStr := r.Header.Get("X-User-ID")
	userID, _ := strconv.Atoi(userIDStr)
	isAnonymous := userID == 0

	// Rate limiting для анонимных пользователей
	if isAnonymous {
		allowed, count, waitTime := s.rateLimiter.IsAllowed(r)
		if !allowed {
			s.logger.Warning("Rate limit exceeded for anonymous user (count: %d, wait: %v)", count, waitTime)
			w.Header().Set("X-RateLimit-Limit", "3")
			w.Header().Set("X-RateLimit-Remaining", "0")
			w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(waitTime).Unix()))
			s.sendError(w, fmt.Sprintf("Rate limit exceeded. Try again in %v", waitTime.Round(time.Minute)), http.StatusTooManyRequests)
			return
		}

		w.Header().Set("X-RateLimit-Limit", "3")
		w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", 3-count))
		s.logger.Info("Anonymous file upload started (usage: %d/3)", count)
	} else {
		s.logger.Info("File upload started for user %d", userID)
	}

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

	if err := s.validator.ValidateFile(header); err != nil {
		if ve, ok := err.(ValidationError); ok {
			s.logger.Warning("File validation failed: %v", ve)
			s.sendValidationErrors(w, []ValidationError{ve})
		} else {
			s.logger.Warning("File validation failed: %v", err)
			s.sendError(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	// Парсим опции обработки
	options := s.parseProcessingOptions(r)

	s.logger.Info("Processing file upload: %s (%d bytes) %s", header.Filename, header.Size,
		func() string {
			if isAnonymous {
				return "for anonymous user"
			}
			return fmt.Sprintf("for user %d", userID)
		}())

	fileID := uuid.New().String()
	ext := filepath.Ext(header.Filename)
	fileName := fileID + ext
	filePath := filepath.Join(s.config.UploadPath, fileName)

	s.logger.Debug("Saving file to: %s", filePath)

	dst, err := os.Create(filePath)
	if err != nil {
		s.logger.Error("Failed to create file %s: %v", filePath, err)
		s.sendError(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	size, err := io.Copy(dst, file)
	if err != nil {
		s.logger.Error("Failed to copy file content: %v", err)
		os.Remove(filePath)
		s.sendError(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	s.logger.Info("File saved to disk: %s (%d bytes)", filePath, size)

	mimeType := s.determineMimeType(header, ext)

	fileRecord := &File{
		ID:           fileID,
		UserID:       uint(userID),
		OriginalName: header.Filename,
		FileName:     fileName,
		FileSize:     size,
		MimeType:     mimeType,
		Status:       StatusUploaded,
		UploadedAt:   time.Now(),
	}

	if !isAnonymous {
		if err := s.db.CreateFile(fileRecord); err != nil {
			s.logger.Error("Failed to save file record for user %d: %v", userID, err)
			os.Remove(filePath)
			s.sendError(w, "Failed to save file record", http.StatusInternalServerError)
			return
		}
		s.logger.Info("File upload completed and saved to history: %s (ID: %s) for user %d", header.Filename, fileID, userID)
	} else {
		s.logger.Info("Anonymous file upload completed: %s (ID: %s) - not saved to history", header.Filename, fileID)
	}

	// Запускаем обработку файла в фоне
	go s.processFileAsync(fileID, filePath, mimeType, options, isAnonymous)

	// Обновляем статус на "processing" если это не анонимный пользователь
	if !isAnonymous {
		s.db.UpdateFileStatus(fileID, StatusProcessing)
		fileRecord.Status = StatusProcessing
	}

	s.sendJSON(w, SuccessResponse{
		Message: func() string {
			if isAnonymous {
				return "File uploaded successfully and processing started (not saved to history)"
			}
			return "File uploaded successfully and processing started"
		}(),
		Data: fileRecord,
	})
}

// @Summary Get user files
// @Description Get list of all files uploaded by the authenticated user with their processing status
// @Tags files
// @Produce json
// @Security BearerAuth
// @Success 200 {object} SuccessResponse{data=[]File}
// @Failure 401 {object} ErrorResponse
// @Router /api/files [get]
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

// @Summary File operations
// @Description Handle file operations: GET for file info/download, DELETE for removal. Use ?type=original or ?type=processed query parameter for downloads
// @Tags files
// @Param id path string true "File ID"
// @Param type query string false "Download type" Enums(original, processed)
// @Security BearerAuth
// @Success 200 {object} SuccessResponse{data=File} "File information"
// @Success 200 {file} binary "File download (when type parameter is used)"
// @Success 200 {object} SuccessResponse "File deleted"
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/files/{id} [get]
// @Router /api/files/{id} [delete]
func (s *Server) handleFileActions(w http.ResponseWriter, r *http.Request) {
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

	// Получаем query parameter type для определения типа операции
	downloadType := r.URL.Query().Get("type")

	s.logger.Debug("File action %s for file %s by %s (type: %s)", r.Method, fileID,
		func() string {
			if isAnonymous {
				return "anonymous user"
			}
			return fmt.Sprintf("user %d", userID)
		}(), downloadType)

	switch r.Method {
	case http.MethodGet:
		if downloadType == "original" || downloadType == "processed" {
			// Скачивание файла
			s.handleDownloadFileByID(w, r, fileID, userID, isAnonymous, downloadType == "processed")
		} else {
			// Получение информации о файле
			s.handleGetFileInfo(w, r, fileID, userID, isAnonymous)
		}
	case http.MethodDelete:
		s.handleDeleteFileByID(w, r, fileID, userID, isAnonymous)
	default:
		s.logger.Warning("Method %s not allowed for file actions", r.Method)
		s.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Скачивание файла по ID
func (s *Server) handleDownloadFileByID(w http.ResponseWriter, r *http.Request, fileID string, userID int, isAnonymous bool, isProcessed bool) {
	if isAnonymous {
		fileName := fileID + "*"
		if isProcessed {
			fileName = fileID + "_processed*"
		}
		matches, err := filepath.Glob(filepath.Join(s.config.UploadPath, fileName))
		if err != nil || len(matches) == 0 {
			s.logger.Warning("Anonymous file not found on disk: %s", fileID)
			s.sendError(w, "File not found", http.StatusNotFound)
			return
		}

		filePath := matches[0]
		s.logger.Info("Serving anonymous file: %s (processed: %v)", fileID, isProcessed)

		mimeType := s.determineMimeTypeFromPath(filePath)
		w.Header().Set("Content-Type", mimeType)
		http.ServeFile(w, r, filePath)
		return
	}

	file, err := s.db.GetFileByID(fileID)
	if err != nil {
		s.logger.Warning("File not found in database: %s for user %d", fileID, userID)
		s.sendError(w, "File not found", http.StatusNotFound)
		return
	}

	if file.UserID != uint(userID) {
		s.logger.Warning("Access denied: user %d tried to access file %s owned by user %d", userID, fileID, file.UserID)
		s.sendError(w, "Access denied", http.StatusForbidden)
		return
	}

	s.handleDownloadFile(w, r, file, isProcessed)
}

// Получение информации о файле
func (s *Server) handleGetFileInfo(w http.ResponseWriter, r *http.Request, fileID string, userID int, isAnonymous bool) {
	if isAnonymous {
		s.sendError(w, "File info not available for anonymous users", http.StatusForbidden)
		return
	}

	file, err := s.db.GetFileByID(fileID)
	if err != nil {
		s.sendError(w, "File not found", http.StatusNotFound)
		return
	}

	if file.UserID != uint(userID) {
		s.sendError(w, "Access denied", http.StatusForbidden)
		return
	}

	s.sendJSON(w, SuccessResponse{
		Message: "File info retrieved",
		Data:    file,
	})
}

// Удаление файла по ID
func (s *Server) handleDeleteFileByID(w http.ResponseWriter, r *http.Request, fileID string, userID int, isAnonymous bool) {
	if isAnonymous {
		s.logger.Warning("Anonymous user attempted to delete file: %s", fileID)
		s.sendError(w, "Anonymous users cannot delete files", http.StatusForbidden)
		return
	}

	file, err := s.db.GetFileByID(fileID)
	if err != nil {
		s.logger.Warning("File not found for deletion: %s for user %d", fileID, userID)
		s.sendError(w, "File not found", http.StatusNotFound)
		return
	}

	if file.UserID != uint(userID) {
		s.logger.Warning("Access denied: user %d tried to delete file %s owned by user %d", userID, fileID, file.UserID)
		s.sendError(w, "Access denied", http.StatusForbidden)
		return
	}

	s.handleDeleteFile(w, r, file)
}

// Скачивание файла
func (s *Server) handleDownloadFile(w http.ResponseWriter, r *http.Request, file *File, isProcessed bool) {
	var filePath string
	var fileName string

	if isProcessed && file.IsProcessed() {
		filePath = filepath.Join(s.config.UploadPath, file.ProcessedName)
		fileName = "processed_" + file.OriginalName
	} else {
		filePath = filepath.Join(s.config.UploadPath, file.FileName)
		fileName = file.OriginalName
	}

	s.logger.Debug("Downloading file: %s (path: %s, processed: %v)", file.ID, filePath, isProcessed)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		s.logger.Error("File not found on disk: %s", filePath)
		s.sendError(w, "File not found on disk", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	w.Header().Set("Content-Type", file.MimeType)

	s.logger.Info("Serving file: %s (%s) to user %d (processed: %v)", fileName, file.ID, file.UserID, isProcessed)
	http.ServeFile(w, r, filePath)
}

// Удаление файла
func (s *Server) handleDeleteFile(w http.ResponseWriter, r *http.Request, file *File) {
	s.logger.Info("Deleting file: %s (%s) for user %d", file.OriginalName, file.ID, file.UserID)

	// Удаляем оригинальный файл
	filePath := filepath.Join(s.config.UploadPath, file.FileName)
	if err := os.Remove(filePath); err != nil {
		s.logger.Warning("Failed to remove original file from disk: %s - %v", filePath, err)
	} else {
		s.logger.Debug("Original file removed from disk: %s", filePath)
	}

	// Удаляем обработанный файл, если он есть
	if file.ProcessedName != "" {
		processedPath := filepath.Join(s.config.UploadPath, file.ProcessedName)
		if err := os.Remove(processedPath); err != nil {
			s.logger.Warning("Failed to remove processed file from disk: %s - %v", processedPath, err)
		} else {
			s.logger.Debug("Processed file removed from disk: %s", processedPath)
		}
	}

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

// @Summary Admin statistics
// @Description Get administrative statistics about server, file system, ML service and rate limiter
// @Tags admin
// @Produce json
// @Success 200 {object} SuccessResponse
// @Router /api/admin/stats [get]
func (s *Server) handleAdminStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fileStats, _ := s.fileCleaner.GetStats()
	rateLimiterStats := s.rateLimiter.GetStats()

	// Статистика обработки файлов
	processingStats := make(map[string]interface{})
	if uploadedFiles, err := s.db.GetFilesByStatus(StatusUploaded); err == nil {
		processingStats["pending_files"] = len(uploadedFiles)
	}
	if processingFiles, err := s.db.GetFilesByStatus(StatusProcessing); err == nil {
		processingStats["processing_files"] = len(processingFiles)
	}
	if completedFiles, err := s.db.GetFilesByStatus(StatusCompleted); err == nil {
		processingStats["completed_files"] = len(completedFiles)
	}
	if failedFiles, err := s.db.GetFilesByStatus(StatusFailed); err == nil {
		processingStats["failed_files"] = len(failedFiles)
	}

	stats := map[string]interface{}{
		"server_uptime":    time.Since(time.Now().Add(-time.Hour)),
		"file_system":      fileStats,
		"rate_limiter":     rateLimiterStats,
		"processing_stats": processingStats,
		"ml_service": map[string]interface{}{
			"enabled": s.config.MLServiceEnabled,
			"url":     s.config.MLServiceURL,
			"healthy": s.checkMLService(),
		},
	}

	s.sendJSON(w, SuccessResponse{
		Message: "Admin stats retrieved",
		Data:    stats,
	})
}

// Парсинг опций обработки из формы
func (s *Server) parseProcessingOptions(r *http.Request) ProcessingOptions {
	options := ProcessingOptions{
		BlurType:  "gaussian",
		Intensity: 5,
	}

	if blurType := r.FormValue("blur_type"); blurType != "" {
		options.BlurType = blurType
	}

	if intensity := r.FormValue("intensity"); intensity != "" {
		if intVal, err := strconv.Atoi(intensity); err == nil && intVal >= 1 && intVal <= 10 {
			options.Intensity = intVal
		}
	}

	if objectTypes := r.FormValue("object_types"); objectTypes != "" {
		options.ObjectTypes = strings.Split(objectTypes, ",")
		for i, obj := range options.ObjectTypes {
			options.ObjectTypes[i] = strings.TrimSpace(obj)
		}
	}

	return options
}

// Определение MIME типа
func (s *Server) determineMimeType(header *multipart.FileHeader, ext string) string {
	mimeType := header.Header.Get("Content-Type")
	if mimeType == "" || mimeType == "application/octet-stream" {
		mimeType = s.determineMimeTypeFromExtension(ext)
	}
	return mimeType
}

// Определение MIME типа по расширению
func (s *Server) determineMimeTypeFromExtension(ext string) string {
	switch strings.ToLower(ext) {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".bmp":
		return "image/bmp"
	case ".tiff", ".tif":
		return "image/tiff"
	case ".mp4":
		return "video/mp4"
	case ".avi":
		return "video/avi"
	case ".mov":
		return "video/quicktime"
	case ".wmv":
		return "video/x-ms-wmv"
	case ".flv":
		return "video/x-flv"
	case ".webm":
		return "video/webm"
	case ".mkv":
		return "video/x-matroska"
	default:
		return "application/octet-stream"
	}
}

// Определение MIME типа по пути
func (s *Server) determineMimeTypeFromPath(filePath string) string {
	ext := filepath.Ext(filePath)
	return s.determineMimeTypeFromExtension(ext)
}

// Асинхронная обработка файла (эмуляция ML)
func (s *Server) processFileAsync(fileID, filePath, mimeType string, options ProcessingOptions, isAnonymous bool) {
	s.logger.Info("Starting processing for file %s (anonymous: %v)", fileID, isAnonymous)

	if s.config.MLServiceEnabled {
		// Когда ML будет готов, здесь будет реальный вызов сервиса
		s.processWithMLService(fileID, filePath, mimeType, options, isAnonymous)
	} else {
		// Эмуляция обработки
		s.processWithEmulation(fileID, filePath, mimeType, options, isAnonymous)
	}
}

// Эмуляция обработки файла
func (s *Server) processWithEmulation(fileID, filePath, mimeType string, options ProcessingOptions, isAnonymous bool) {
	s.logger.Debug("Emulating ML processing for file %s", fileID)

	// Эмитируем время обработки (2-5 секунд)
	processingTime := time.Duration(2+time.Now().UnixNano()%3) * time.Second
	time.Sleep(processingTime)

	// Создаем "обработанный" файл (копируем оригинал)
	ext := filepath.Ext(filePath)
	processedName := fileID + "_processed" + ext
	processedPath := filepath.Join(s.config.UploadPath, processedName)

	// Копируем файл
	if err := s.copyFile(filePath, processedPath); err != nil {
		s.logger.Error("Failed to create processed file %s: %v", processedPath, err)
		if !isAnonymous {
			s.db.UpdateFileProcessing(fileID, "", 0, StatusFailed, "Failed to create processed file")
		}
		return
	}

	// Получаем размер обработанного файла
	processedSize := int64(0)
	if stat, err := os.Stat(processedPath); err == nil {
		processedSize = stat.Size()
	}

	if !isAnonymous {
		err := s.db.UpdateFileProcessing(fileID, processedName, processedSize, StatusCompleted, "")
		if err != nil {
			s.logger.Error("Failed to update file processing status for %s: %v", fileID, err)
		} else {
			s.logger.Info("File processing completed successfully: %s", fileID)
		}
	} else {
		s.logger.Info("Anonymous file processing completed: %s", fileID)
	}
}

// Обработка с помощью ML сервиса (заглушка для будущего)
func (s *Server) processWithMLService(fileID, filePath, mimeType string, options ProcessingOptions, isAnonymous bool) {
	s.logger.Debug("Processing file %s with ML service", fileID)

	// Подготавливаем запрос к ML сервису
	request := ProcessingRequest{
		FileID:   fileID,
		FilePath: filePath,
		MimeType: mimeType,
		Options:  options,
	}

	// TODO: Когда ML будет готов, здесь будет реальный HTTP запрос
	_ = request

	// Пока эмулируем
	s.processWithEmulation(fileID, filePath, mimeType, options, isAnonymous)
}

// Проверка доступности ML сервиса
func (s *Server) checkMLService() bool {
	if !s.config.MLServiceEnabled {
		return false
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(s.config.MLServiceURL + "/health")
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// Копирование файла
func (s *Server) copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

// Генерация JWT токена
func (s *Server) generateJWT(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Minute * 30).Unix(),
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

// Отправка ошибок валидации
func (s *Server) sendValidationErrors(w http.ResponseWriter, errors []ValidationError) {
	s.logger.Warning("Validation errors: %v", errors)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":  "Validation failed",
		"errors": errors,
	})
}

// Остановка сервера и очистка ресурсов
func (s *Server) Stop() {
	s.logger.Info("Stopping server...")
	if s.fileCleaner != nil {
		s.fileCleaner.Stop()
	}
	s.logger.Info("Server stopped")
}
