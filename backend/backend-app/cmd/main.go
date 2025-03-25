package main

import (
	"flag"
	"log"

    "github.com/BurntSushi/toml"
	"obscura.app/backend/pkg/logger"

	apiserver "obscura.app/backend/internal/api-server"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/api-server.toml", "config file path")
}

func main() {
    logger, err := logger.NewLogger("cmd/app.log")
    if err != nil {
        log.Fatalf("Could not initialize logger: %v", err)
    }
    defer logger.Close()

    logger.Info("This is an info message")
    logger.Debug("This is a debug message")
    logger.Warning("This is a warning message")
    logger.Error("This is an error message")
    logger.Fatal("This is a fatal message")

    flag.Parse()

	config, _ := apiserver.LoadConfig(configPath, logger)
	_, err = toml.DecodeFile(configPath, config)
	if err != nil {
		logger.Fatal("", err)
	}

	if err := apiserver.Start(config); err != nil {
		logger.Fatal("", err)
	}
}