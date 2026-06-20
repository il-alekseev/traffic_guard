package middleware

import (
	"api-gateway/internal/controllers/http/v1/values"
	"api-gateway/pkg/slogger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestIDMiddleware добавляет UUID для каждого запроса
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем или генерируем Request ID
		requestID := c.GetHeader(values.RequestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Сохраняем в контексте Gin
		c.Set(values.RequestIDContextKey, requestID)

		// Добавляем в заголовки ответа
		c.Writer.Header().Set(values.RequestIDHeader, requestID)

		// Передаем управление следующему обработчику
		c.Request = c.Request.WithContext(
			slogger.WithLog(
				c.Request.Context(), "request_id", requestID,
			),
		)
		c.Next()
	}
}
