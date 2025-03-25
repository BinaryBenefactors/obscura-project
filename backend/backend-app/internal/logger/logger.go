package logger

import (
	"log"
	"os"
	"sync"

	"github.com/fatih/color"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARNING
	ERROR
	FATAL
)

type Logger struct {
	infoLogger    *log.Logger
	errorLogger   *log.Logger
	debugLogger   *log.Logger
	warningLogger *log.Logger
	fatalLogger   *log.Logger
	fileLogger    *log.Logger // Логгер для файла
	logFile       *os.File
	mu            sync.Mutex
}

func NewLogger(logFilePath string) (*Logger, error) {
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}

	return &Logger{
		infoLogger:    log.New(os.Stdout, color.CyanString("[INFO] "), log.Ltime),
		errorLogger:   log.New(os.Stdout, color.RedString("[ERROR] "), log.Ltime),
		debugLogger:   log.New(os.Stdout, color.GreenString("[DEBUG] "), log.Ltime),
		warningLogger: log.New(os.Stdout, color.YellowString("[WARNING] "), log.Ltime),
		fatalLogger:   log.New(os.Stdout, color.MagentaString("[FATAL] "), log.Ltime),
		fileLogger:    log.New(logFile, "", log.Ldate|log.Ltime), // Даты + время в файл
		logFile:       logFile,
	}, nil
}

func (l *Logger) log(level LogLevel, msg string) {
	// Вывод в консоль с цветами
	switch level {
	case DEBUG:
		l.debugLogger.Println(msg)
	case INFO:
		l.infoLogger.Println(msg)
	case WARNING:
		l.warningLogger.Println(msg)
	case ERROR:
		l.errorLogger.Println(msg)
	case FATAL:
		l.fatalLogger.Println(msg)
	}

	// Запись в файл БЕЗ цветовых кодов
	var levelStr string
	switch level {
	case DEBUG:
		levelStr = "[DEBUG] "
	case INFO:
		levelStr = "[INFO] "
	case WARNING:
		levelStr = "[WARNING] "
	case ERROR:
		levelStr = "[ERROR] "
	case FATAL:
		levelStr = "[FATAL] "
	}

	l.mu.Lock()
	defer l.mu.Unlock()
	l.fileLogger.Println(levelStr + msg) // Дата и время добавятся автоматически
}

func (l *Logger) Debug(msg string) {
	l.log(DEBUG, msg)
}

func (l *Logger) Info(msg string) {
	l.log(INFO, msg)
}

func (l *Logger) Warning(msg string) {
	l.log(WARNING, msg)
}

func (l *Logger) Error(msg string) {
	l.log(ERROR, msg)
}

func (l *Logger) Fatal(msg string) {
	l.log(FATAL, msg)
	os.Exit(1)
}

func (l *Logger) Close() error {
	return l.logFile.Close()
}
