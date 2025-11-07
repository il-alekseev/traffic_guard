package logging

import (
	"log/slog"
	"os"
	"strings"
)

// Logger wraps slog.Logger providing convenience helpers for structured output.
type Logger struct {
	base *slog.Logger
}

// New creates a stdout logger with the provided level (debug/info/warn/error).
func New(level string) *Logger {
	var lvl slog.Level
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		lvl = slog.LevelDebug
	case "warn", "warning":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: lvl,
	})

	return &Logger{base: slog.New(handler)}
}

// With returns a child logger with the given attributes bound.
func (l *Logger) With(args ...any) *Logger {
	if l == nil || l.base == nil {
		return l
	}
	return &Logger{base: l.base.With(args...)}
}

func (l *Logger) Debug(msg string, args ...any) {
	if l == nil || l.base == nil {
		return
	}
	l.base.Debug(msg, args...)
}

func (l *Logger) Info(msg string, args ...any) {
	if l == nil || l.base == nil {
		return
	}
	l.base.Info(msg, args...)
}

func (l *Logger) Warn(msg string, args ...any) {
	if l == nil || l.base == nil {
		return
	}
	l.base.Warn(msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	if l == nil || l.base == nil {
		return
	}
	l.base.Error(msg, args...)
}
