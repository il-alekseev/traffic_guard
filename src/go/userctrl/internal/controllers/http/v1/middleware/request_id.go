package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const RequestIDHeader = "X-Request-ID"
const RequestIDContextKey = "requestID"

// RequestIDMiddleware добавляет UUID для каждого запроса
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем или генерируем Request ID
		requestID := c.GetHeader(RequestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Сохраняем в контексте Gin
		c.Set(RequestIDContextKey, requestID)

		// Добавляем в заголовки ответа
		c.Writer.Header().Set(RequestIDHeader, requestID)

		// Передаем управление следующему обработчику
		c.Next()
	}
}

// GetRequestID извлекает ID запроса из контекста
func GetRequestID(c *gin.Context) string {
	if id, exists := c.Get(RequestIDContextKey); exists {
		if str, ok := id.(string); ok {
			return str
		}
	}
	return ""
}
