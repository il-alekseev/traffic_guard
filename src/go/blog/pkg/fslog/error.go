package fslog

import (
	"context"
	"log/slog"
	"runtime"
	"time"
)

type errorWithLogCtx struct {
	next error
	ctx  logCtx
}

func (e *errorWithLogCtx) Error() string {
	return e.next.Error()
}

func WrapError(ctx context.Context, err error) error {
	c := logCtx{}
	if x, ok := ctx.Value(key).(logCtx); ok {
		c = x
	}
	return &errorWithLogCtx{
		next: err,
		ctx:  c,
	}
}

func isErrorCtx(ctx context.Context, err error) context.Context {
	if e, ok := err.(*errorWithLogCtx); ok { // в реальной жизни используйте error.As
		return context.WithValue(ctx, key, e.ctx)
	}
	return ctx
}

func ErrorCtx(ctx context.Context, logger *slog.Logger, err error) {
	if !logger.Enabled(context.Background(), slog.LevelError) {
		return
	}

	var pcs [1]uintptr
	runtime.Callers(2, pcs[:]) // skip [Callers, ErrorCtx]

	r := slog.NewRecord(time.Now(), slog.LevelError, "Error: "+err.Error(), pcs[0])

	newCtx := isErrorCtx(ctx, err)

	_ = logger.Handler().Handle(newCtx, r)
}
