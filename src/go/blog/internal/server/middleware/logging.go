package middleware

import (
	"fiermon-blog/pkg/fslog/wsl"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func LoggingMiddleware(l *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем или генерируем X-Request-ID
		requestID := c.GetHeader(RequestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
			c.Request.Header.Set(RequestIDHeader, requestID)
		}
		c.Set(RequestIDContextKey, requestID)

		// Логируем входящий запрос
		startTime := time.Now()
		l.InfoContext(c.Request.Context(), "Incoming request",
			wsl.Label("method", c.Request.Method),
			wsl.Label("path", c.Request.URL.Path),
			wsl.Label("ip", c.ClientIP()),
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
		l.InfoContext(c.Request.Context(), "Request completed",
			wsl.Label("method", c.Request.Method),
			wsl.Label("path", c.Request.URL.Path),
			wsl.Int("status", rw.statusCode),
			wsl.Label("duration", duration.String()),
			wsl.Int("size", rw.size),
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
