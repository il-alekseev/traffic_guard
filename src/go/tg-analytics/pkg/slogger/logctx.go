package slogger

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
