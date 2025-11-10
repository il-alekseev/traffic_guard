package fslog

import "context"

const key = keyType("logCtxKey") // ключ, по которому лежит контекст логирования

type logCtx struct {
	Values map[string]string
}

func newLogCtx() logCtx {
	return logCtx{
		Values: make(map[string]string),
	}
}

type keyType string

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
