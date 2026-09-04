package logger

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type SlogLogger struct{
	l *slog.Logger
	File *os.File
	mu   sync.Mutex
}

// todo: cambiar a .env
const logDir = "./logs"

func NewSlogLogger() (*SlogLogger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	fileName := filepath.Join(logDir, fmt.Sprintf("app-%s.log", time.Now().Format("2006-01-02")))

	logFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err!=nil{
		return nil, fmt.Errorf("failed to open app.log file: %w", err)
	}
	
	jsonHandler := slog.NewJSONHandler(logFile, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	logger := slog.New(jsonHandler)
	return &SlogLogger{l: logger, File: logFile}, nil
}

func (s *SlogLogger) RotateLogFile() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	newLogger, err := NewSlogLogger()
	if err != nil {
		return err
	}

	oldFile := s.File

	s.l = newLogger.l
	s.File = newLogger.File

	if oldFile != nil {
		err := oldFile.Close()
		if err != nil {
			return fmt.Errorf("failed to close old log file: %w", err)
		}
	}

	return nil
}

func (s *SlogLogger) CleanupOldLogFiles() error {
	const retentionDays = 60

	cutoff := time.Now().AddDate(0, 0, -retentionDays)

	entries, err := os.ReadDir(logDir)
	if err != nil {
		return fmt.Errorf("failed to read log directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()

		// app-2026-06-16.log
		if !strings.HasPrefix(name, "app-") || !strings.HasSuffix(name, ".log") {
			continue
		}

		dateStr := strings.TrimSuffix(
			strings.TrimPrefix(name, "app-"),
			".log",
		)

		logDate, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			// Ignoramos archivos con nombres inesperados.
			continue
		}

		if logDate.Before(cutoff) {
			path := filepath.Join(logDir, name)

			if err := os.Remove(path); err != nil {
				return fmt.Errorf("failed to remove old log file %q: %w", name, err)
			}
		}
	}

	return nil
}


func (s *SlogLogger) Info(msg string, args ...any){
	s.l.Info(msg, args...)
}

func (s *SlogLogger) Error(msg string, args ...any){
	s.l.Error(msg, args...)
}

func (s *SlogLogger) Debug(msg string, args ...any){
	s.l.Debug(msg, args...)
}

