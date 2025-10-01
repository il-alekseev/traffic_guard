package fslog

import (
	"context"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

const (
	localEnv      = "local"
	debugEnv      = "debug"
	productionEnv = "prod"
)

type FSlogger struct {
	next slog.Handler
}

func newHandlerMiddleware(next slog.Handler) *FSlogger {
	return &FSlogger{next: next}
}

func (h *FSlogger) Enabled(ctx context.Context, rec slog.Level) bool {
	return h.next.Enabled(ctx, rec)
}

func (h *FSlogger) Handle(ctx context.Context, rec slog.Record) error {
	if c, ok := ctx.Value(key).(logCtx); ok {
		for k, v := range c.Values {
			rec.Add(k, v)
		}
	}
	return h.next.Handle(ctx, rec)
}

func (h *FSlogger) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &FSlogger{next: h.next.WithAttrs(attrs)} // не забыть обернуть, но осторожно
}

func (h *FSlogger) WithGroup(name string) slog.Handler {
	return &FSlogger{next: h.next.WithGroup(name)} // не забыть обернуть, но осторожно
}

func NewFSlogger(env string) *slog.Logger {

	logOption := &slog.HandlerOptions{
		ReplaceAttr: customAttr,
	}

	var handler slog.Handler
	switch env {
	case localEnv:
		logOption.AddSource = true
		logOption.Level = slog.LevelDebug
		handler = slog.NewTextHandler(os.Stdout, logOption)
	case debugEnv:
		logOption.AddSource = true
		logOption.Level = slog.LevelDebug
		handler = slog.NewJSONHandler(os.Stdout, logOption)
	case productionEnv:
		logOption.Level = slog.LevelInfo
		handler = slog.NewJSONHandler(os.Stdout, logOption)
	default:
		return nil
	}

	handler = newHandlerMiddleware(handler)
	return slog.New(handler)
}

func customAttr(groups []string, a slog.Attr) slog.Attr {
	wd, _ := os.Getwd()

	switch a.Key {
	case slog.TimeKey:
		tm := a.Value.Time().Format("2006-01-02 15:04:05")
		a.Value = slog.StringValue(tm)
	case slog.SourceKey:
		source := a.Value.Any().(*slog.Source)
		rel := strings.TrimPrefix(source.File, wd+string(os.PathSeparator))
		a.Value = slog.StringValue(rel + ":" + strconv.Itoa(source.Line))
	}
	return a
}
