package logger

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARNING
	ERROR
)

type LogEntry struct {
	Timestamp time.Time
	Level     LogLevel
	Message   string
}

var (
	logFile    *os.File
	logEntries []LogEntry
	mu         sync.Mutex
	logLevel   LogLevel
	logger     *log.Logger
)

func Init(level string, filePath string) error {
	switch level {
	case "debug":
		logLevel = DEBUG
	case "info":
		logLevel = INFO
	case "warning":
		logLevel = WARNING
	case "error":
		logLevel = ERROR
	default:
		return fmt.Errorf("invalid log level: %s", level)
	}

	var err error
	logFile, err = os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	logger = log.New(logFile, "", log.LstdFlags)
	return nil
}

func logMessage(level LogLevel, message string) {
	if level < logLevel {
		return
	}

	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
	}

	mu.Lock()
	logEntries = append(logEntries, entry)
	mu.Unlock()

	logger.Printf("[%s] %s", getLevelString(level), message)
}

func Debug(message string) {
	logMessage(DEBUG, message)
}

func Info(message string) {
	logMessage(INFO, message)
}

func Warning(message string) {
	logMessage(WARNING, message)
}

func Error(message string) {
	logMessage(ERROR, message)
}

func GetLogs(level LogLevel) []LogEntry {
	mu.Lock()
	defer mu.Unlock()

	filteredLogs := make([]LogEntry, 0)
	for _, entry := range logEntries {
		if entry.Level >= level {
			filteredLogs = append(filteredLogs, entry)
		}
	}

	return filteredLogs
}

func getLevelString(level LogLevel) string {
	switch level {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARNING:
		return "WARNING"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

func GetLevelString(level LogLevel) string {
	return getLevelString(level)
}
