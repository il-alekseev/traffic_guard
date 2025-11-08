package utils

import (
	"tg-an/internal/controllers/http/v1/values"

	"github.com/gin-gonic/gin"
)

// GetRequestID извлекает ID запроса из контекста
func GetRequestID(c *gin.Context) string {
	if id, exists := c.Get(values.RequestIDContextKey); exists {
		if str, ok := id.(string); ok {
			return str
		}
	}
	return ""
}
