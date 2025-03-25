package logger

import (
	"log"
	"os"

	"github.com/fatih/color"
)

type Logger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
}

func NewLogger() *Logger {
	return &Logger{
		infoLogger:  log.New(os.Stdout, color.CyanString("[INFO] "), log.Ltime),
		errorLogger: log.New(os.Stdout, color.RedString("[ERROR] "), log.Ltime),
	}
}

func (l *Logger) Info(msg string) {
	l.infoLogger.Println(msg)
}

func (l *Logger) Error(msg string) {
	l.errorLogger.Println(msg)
}
