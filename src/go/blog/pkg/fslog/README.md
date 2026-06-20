# FSlog - Форматированное логирование

## Начало использования

```go
// создание объекта логера
// поддерживает 3 уровня логирования `local, debug, prod`
log := fslog.NewFSlogger("debug")

// Передача логера в слой логики
layer := NewLayer(log)
```
```go
// Пример слоя логики

type Layer struct {
	log *slog.Logger
}

func NewLayer(log *slog.Logger) *Layer {
	log = log.With(wsl.Label("layer", "use_case"))
	return &Layer{log: log}
}

```
## Использование 
Контекст позволяет удобно накапливать необходимые мета данные ошибки
```go
// Логирование с контекстом
ctx := slogger.WithLog(context.Background(), "user_id", 111)
log.InfoContext(ctx, "Handler started")

//вывод - {"time":"2025-07-24 14:53:40","level":"INFO","source":"main.go:68","msg":"Handler started","UserID":"111"}

// создание ошибки с контекстом
err := errors.New("some kind error")
if err != nil {
return fslog.WrapError(ctx, err)
}

// Вывод ошибок с контекстом
err := doSomething(ctx)
if err != nil {
    fslog.ErrorCtx(ctx, log, err)
}
```
