package cslogger

import (
	"context"
	"log/slog"
	"os"
)

const (
	colorReset   = "\033[0m"
	colorDebug   = "\033[36m" // Циан
	colorInfo    = "\033[32m" // Зелёный
	colorWarn    = "\033[33m" // Жёлтый
	colorError   = "\033[31m" // Красный
	colorTime    = "\033[90m" // Серый
	colorMessage = "\033[97m" // Белый
)

type ColorHandler struct {
	handler slog.Handler
	w       *os.File
}

func NewColorHandler(w *os.File) *ColorHandler {
	return &ColorHandler{
		handler: slog.NewTextHandler(w, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}),
		w: w,
	}
}

func (h *ColorHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *ColorHandler) Handle(ctx context.Context, r slog.Record) error {
	var levelColor string
	switch r.Level {
	case slog.LevelDebug:
		levelColor = colorDebug
	case slog.LevelInfo:
		levelColor = colorInfo
	case slog.LevelWarn:
		levelColor = colorWarn
	case slog.LevelError:
		levelColor = colorError
	default:
		levelColor = colorReset
	}

	// Форматируем время
	timeStr := r.Time.Format("15:04:05.000")

	// Формируем сообщение с цветами
	msg := colorTime + timeStr + colorReset + " " +
		levelColor + r.Level.String() + colorReset + " " +
		colorMessage + r.Message + colorReset

	// Добавляем атрибуты
	r.Attrs(func(attr slog.Attr) bool {
		msg += " " + levelColor + attr.Key + colorReset + "=" +
			colorMessage + attr.Value.String() + colorReset
		return true
	})

	// Выводим в поток
	_, err := h.w.WriteString(msg + "\n")
	return err
}

func (h *ColorHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ColorHandler{
		handler: h.handler.WithAttrs(attrs),
		w:       h.w,
	}
}

func (h *ColorHandler) WithGroup(name string) slog.Handler {
	return &ColorHandler{
		handler: h.handler.WithGroup(name),
		w:       h.w,
	}
}

// NewColorLogger создает новый логгер с цветным выводом
func NewColorLogger() *slog.Logger {
	return slog.New(NewColorHandler(os.Stdout))
}

// NewColorLoggerWithLevel создает цветной логгер с указанным уровнем
func NewColorLoggerWithLevel(level slog.Level) *slog.Logger {
	handler := NewColorHandler(os.Stdout)
	return slog.New(handler)
}
