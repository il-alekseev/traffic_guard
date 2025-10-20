# Slogger - Structured Logger for Go

## Быстрый старт
```go
// Инициализация
slogger.InitLogging(slog.LevelInfo)

// Логирование с контекстом
ctx := slogger.WithLog(context.Background(), "user_id", 123)
slog.InfoContext(ctx, "User action", "action", "login")

// Обработка ошибок
err := doSomething(ctx)
if err != nil {
    slog.ErrorContext(slogger.ErrorCtx(ctx, err), "Operation failed")
}
```

## Init
```go
func Run(cfg *config.Config) {
    slogger.InitLogging(slog.LevelDebug)
    // Доступные уровни:
    // slog.LevelDebug
    // slog.LevelInfo
    // slog.LevelWarn
    // slog.LevelError
}
```

## Использование
 Перед передачей контекста в функции ниже, в объект ctx добавляются дополнительные атрибуты. Эти атрибуты автоматически включаются в логи, связанные с данным контекстом. Рекомендуется всегда вести логирование с использованием контекста для более информативных и структурированных записей.
```go
slog.DebugContext(ctx, msg, attrs...)
slog.InfoContext(ctx, msg, attrs...)
slog.WarnContext(ctx, msg, attrs...)
slog.ErrorContext(ctx, msg, attrs...)
```

## Добавление атрибутов в контекст 
Добавление атрибутов (по одному)
```go
ctx = slogger.WithLog(ctx, key, value)
```
Для `gin.Context`:
```go
// c *gin.Context
c.Request = c.Request.WithContext(
    slogger.WithLog(
		c.Request.Context(), "user_id", userID,
	),
)
```
## Логирование ошибок
В месте возникновения/созания/появления ошибки - сделать вписываение ошибки в контекст.

Использовать `slogger.WrapError(ctx, err)`:
```go
func processUser(ctx context.Context, userID int) error {
    if err := dbQuery(ctx); err != nil {
        // Дополняем информацию об ошибке
        err = fmt.Errorf("database query failed: %w", err)
        
        // Оборачиваем ошибку с контекстом
        return slogger.WrapError(ctx, err)
    }
    return nil
}
```
Функция `slogger.WrapError` вернёт структуру, реализующую интерфейс `error`.

Логирование ошибки:
```go
err := processUser(ctx, 123)
if err != nil {
    slog.ErrorContext(slogger.ErrorCtx(ctx, err), "err: "+err.Error())
}
```
