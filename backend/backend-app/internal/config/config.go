package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	MaxUploadSize int64
	UploadPath    string
	LogFilePath   string
	StaticDir     string
	Port          string
	DBConfig      *DBConfig
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func NewConfig() *Config {
	// Получаем абсолютный путь к директории загрузок
	uploadPath, err := filepath.Abs("./uploads")
	if err != nil {
		uploadPath = "./uploads" // Fallback к относительному пути
	}

	return &Config{
		MaxUploadSize: 50 * 1024 * 1024, // 50MB
		UploadPath:    uploadPath,
		LogFilePath:   "./app.log",
		StaticDir:     "../../../frontend/static",
		Port:          getEnv("PORT", "8080"),
		DBConfig: &DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "obscura"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
