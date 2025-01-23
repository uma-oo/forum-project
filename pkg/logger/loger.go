package logger

import (
	"log"
	"os"
)

// Logger is a reusable logger with file output
type Logger struct {
	LogFile *os.File
}

// NewLogger initializes a logger writing to a file
func Create_Logger() (*Logger, error) {
	file, err := os.OpenFile("server.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	log.SetOutput(file) // Set global logging output to the file
	return &Logger{LogFile: file}, nil
}

// Close closes the log file
func (l *Logger) Close() {
	if l.LogFile != nil {
		l.LogFile.Close()
	}
}
