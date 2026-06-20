package slogger

import (
	"context"
	"log/slog"
	"os"
)

type HandlerMiddlware struct {
	next slog.Handler
}

func NewHandlerMiddlware(next slog.Handler) *HandlerMiddlware {
	return &HandlerMiddlware{next: next}
}

func (h *HandlerMiddlware) Enabled(ctx context.Context, rec slog.Level) bool {
	return h.next.Enabled(ctx, rec)
}

func (h *HandlerMiddlware) Handle(ctx context.Context, rec slog.Record) error {
	if c, ok := ctx.Value(key).(logCtx); ok {
		for k, v := range c.Values {
			rec.Add(k, v)
		}
	}
	return h.next.Handle(ctx, rec)
}

func (h *HandlerMiddlware) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &HandlerMiddlware{next: h.next.WithAttrs(attrs)} // осторожно обернуть
}

func (h *HandlerMiddlware) WithGroup(name string) slog.Handler {
	return &HandlerMiddlware{next: h.next.WithGroup(name)} // осторожно обернуть
}

// WithLog adds key-value pair to logging context
func WithLog(ctx context.Context, k string, v string) context.Context {
	if c, ok := ctx.Value(key).(logCtx); ok {
		c.Values[k] = v
		return context.WithValue(ctx, key, c)
	}
	lc := newLogCtx()
	lc.Values[k] = v
	return context.WithValue(ctx, key, lc)
}

func InitLogging(logLevel slog.Level) {
	handler := slog.Handler(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: true,
	}))
	handler = NewHandlerMiddlware(handler)
	slog.SetDefault(slog.New(handler))
}
