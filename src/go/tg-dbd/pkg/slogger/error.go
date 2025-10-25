package slogger

import (
	"context"
)

type errorWithLogCtx struct {
	next error
	ctx  logCtx
}

func (e *errorWithLogCtx) Error() string {
	return e.next.Error()
}

func WrapError(ctx context.Context, err error) error {
	c := newLogCtx()
	if x, ok := ctx.Value(key).(logCtx); ok {
		c = x
	}
	return &errorWithLogCtx{
		next: err,
		ctx:  c,
	}
}

// Используется там, где мы логием ошибку
// И логировать ошибку так
//
//	if err != nil {
//			slog.ErrorContext(ErrorCtx(ctx, err), "Error: "+err.Error())
//	     ...
//	}
//
// При этом в месте возникновения ошибки вернуть не ошибку, а так
//
//	...
//	return WrapError(ctx, err)
//	...
func ErrorCtx(ctx context.Context, err error) context.Context {
	if e, ok := err.(*errorWithLogCtx); ok { // в реальной жизни используется errors.As
		return context.WithValue(ctx, key, e.ctx)
	}
	return ctx
}
