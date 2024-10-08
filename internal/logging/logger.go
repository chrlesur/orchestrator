package logging

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	DebugLogger *log.Logger
)

func InitLoggers(logFile string, toConsole bool, debugMode bool) error {
	// Ensure the logs directory exists
	logDir := filepath.Dir(logFile)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	var writers []io.Writer
	writers = append(writers, file)
	if toConsole {
		writers = append(writers, os.Stdout)
	}
	multiWriter := io.MultiWriter(writers...)

	InfoLogger = log.New(multiWriter, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(multiWriter, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	if debugMode {
		DebugLogger = log.New(multiWriter, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		DebugLogger = log.New(io.Discard, "", 0)
	}

	return nil
}
