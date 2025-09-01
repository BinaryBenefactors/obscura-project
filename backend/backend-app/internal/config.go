package internal

import (
	"os"
	"strconv"
)

type Config struct {
	Port        string
	UploadPath  string
	MaxFileSize int64
	JWTSecret   string
	MaxAttemptsHandled int
	HandlerTimeout int

	// База данных
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// ML сервис
	MLServiceURL       string
	MLServiceTimeout   int    // в секундах
	MLServiceEnabled   bool
}

func NewConfig() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		UploadPath:  getEnv("UPLOAD_PATH", "./uploads"),
		MaxFileSize: getEnvAsInt64("MAX_FILE_SIZE", 52428800), // 50MB
		JWTSecret:   getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		MaxAttemptsHandled:  getEnvAsInt("MAX_ATTEMPTS_HANDLED", 3),
		HandlerTimeout: getEnvAsInt("HANDLER_TIMEOUT", 24),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "obscura"),

		MLServiceURL:       getEnv("ML_SERVICE_URL", "http://ml:5000"),
		MLServiceTimeout:   getEnvAsInt("ML_SERVICE_TIMEOUT", 300), // 5 минут
		MLServiceEnabled:   getEnvAsBool("ML_SERVICE_ENABLED", false), // пока отключен
	}
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

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
