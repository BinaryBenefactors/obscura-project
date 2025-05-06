package config

import (
	"os"
	"path/filepath"
	"strconv"
)

type Config struct {
	// Настройки сервера
	Port          string
	MaxUploadSize int64
	UploadPath    string
	LogFilePath   string
	StaticDir     string

	// Настройки базы данных
	DB struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
	}
}

func NewConfig() *Config {
	// Получаем абсолютный путь к директории загрузок
	uploadPath, err := filepath.Abs("./uploads")
	if err != nil {
		uploadPath = "./uploads" // Fallback к относительному пути
	}

	cfg := &Config{
		Port:          getEnv("PORT", "8080"),
		MaxUploadSize: getEnvAsInt64("MAX_UPLOAD_SIZE", 10*1024*1024), // 10MB по умолчанию
		UploadPath:    uploadPath,
		LogFilePath:   getEnv("LOG_FILE", "app.log"),
		StaticDir:     "../../../frontend/static",
	}

	// Настройки базы данных
	cfg.DB.Host = getEnv("DB_HOST", "localhost")
	cfg.DB.Port = getEnv("DB_PORT", "5432")
	cfg.DB.User = getEnv("DB_USER", "postgres")
	cfg.DB.Password = getEnv("DB_PASSWORD", "postgres")
	cfg.DB.Name = getEnv("DB_NAME", "obscura")

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt64(key string, defaultValue int64) int64 {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}
