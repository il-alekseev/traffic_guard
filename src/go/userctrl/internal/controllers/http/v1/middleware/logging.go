package middleware

import (
	"userctrl/internal/controllers/http/v1/values"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// LoggingMiddleware - middleware для добавления метаданных в контекст логирования
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем или генерируем X-Request-ID
		requestID := c.GetHeader(values.RequestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
			c.Request.Header.Set(values.RequestIDHeader, requestID)
		}
		c.Set(values.RequestIDContextKey, requestID)

		// Логируем входящий запрос
		startTime := time.Now()
		slog.InfoContext(c.Request.Context(),
			"Incoming request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"ip", c.ClientIP(),
			"request_id", requestID,
		)

		// Перехватываем ответ
		rw := &responseWriter{
			ResponseWriter: gin.ResponseWriter(c.Writer),
			statusCode:     http.StatusOK,
		}
		c.Writer = rw

		// Передаем управление следующему middleware/handler
		c.Next()

		// Логируем исходящий ответ
		duration := time.Since(startTime)
		slog.InfoContext(c.Request.Context(),
			"HTTP request completed",
			slog.Group("request",
				slog.String("method", c.Request.Method),
				slog.String("path", c.Request.URL.Path),
				slog.String("id", requestID),
			),
			slog.Group("response",
				slog.Int("status", rw.statusCode),
				slog.String("duration", duration.String()),
				slog.Int("size", rw.size),
			),
		)
	}
}

// responseWriter перехватывает статус и размер ответа
type responseWriter struct {
	gin.ResponseWriter
	statusCode int
	size       int
}

func (w *responseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *responseWriter) Write(b []byte) (int, error) {
	size, err := w.ResponseWriter.Write(b)
	w.size += size
	return size, err
}
