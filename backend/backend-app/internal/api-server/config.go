package apiserver

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/pelletier/go-toml/v2"
	"obscura.app/backend/pkg/logger"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`

	*ServerConfig
	*DatabaseConfig
}

type ServerConfig struct {
	Port string `yaml:"port"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

func LoadConfig(configPath string, log *logger.Logger) (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Warning(fmt.Sprintf("Error loading .env file: %v", err))
	}

	file, err := os.Open(configPath)
	if err != nil {
		log.Error("Failed to open config file", err)
		return nil, fmt.Errorf("open config file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Error("Failed to close config file", err)
		}
	}()

	decoder := toml.NewDecoder(file)
	config := &Config{}
	if err := decoder.Decode(config); err != nil {
		log.Error("Failed to decode config file", err)
		return nil, fmt.Errorf("decode config: %w", err)
	}

	if port := os.Getenv("SERVER_PORT"); port != "" {
		config.Server.Port = port
		log.Debug(fmt.Sprintf("Using SERVER_PORT from env: %s", port))
	}
	if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		config.Database.Host = dbHost
		log.Debug(fmt.Sprintf("Using DB_HOST from env: %s", dbHost))
	}

	log.Info("Configuration loaded successfully")
	return config, nil
}