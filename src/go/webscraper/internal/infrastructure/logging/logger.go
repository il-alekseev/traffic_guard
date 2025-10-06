package logging

import (
	"log"
	"os"
	"path/filepath"
	"sync"
)

// Logger wraps a standard logger with file output and synchronization for concurrent writes.
type Logger struct {
	mu sync.Mutex
	l  *log.Logger
	f  *os.File
}

// NewFileLogger creates (and ensures directories for) a file-backed logger.
func NewFileLogger(path string) (*Logger, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, err
	}

	logger := log.New(file, "", log.LstdFlags|log.LUTC)

	return &Logger{
		l: logger,
		f: file,
	}, nil
}

// Close closes the underlying file handle.
func (l *Logger) Close() error {
	if l == nil || l.f == nil {
		return nil
	}
	return l.f.Close()
}

// Logf writes a formatted message in a thread-safe manner.
func (l *Logger) Logf(format string, args ...interface{}) {
	if l == nil || l.l == nil {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.l.Printf(format, args...)
}
