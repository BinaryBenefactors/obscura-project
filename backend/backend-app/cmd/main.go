package main

import (
	"log"

	"obscura.app/backend/internal/logger"
)

func main() {
    logger, err := logger.NewLogger("app.log")
    if err != nil {
        log.Fatalf("Could not initialize logger: %v", err)
    }
    defer logger.Close()

    logger.Info("This is an info message")
    logger.Debug("This is a debug message")
    logger.Warning("This is a warning message")
    logger.Error("This is an error message")
    logger.Fatal("This is a fatal message")
}