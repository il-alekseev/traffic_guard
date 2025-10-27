package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CorsMiddleware — middleware для настройки CORS (Cross-Origin Resource Sharing)
func CorsMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		// Разрешает запросы со всех источников (*), поддерживаемые методы: GET, POST, PUT, DELETE, OPTIONS, HEAD
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
		// Разрешены заголовки, включая авторизационные и кастомные; разрешены учётные данные (куки, авторизация)
		AllowHeaders:     []string{"Origin", "Content-type", "Accept", "Authorization", "X-Refresh-Token", "X-Request-ID"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		// Максимальное время кэширования предварительного запроса (preflight) — 12 часов
		MaxAge: 12 * time.Hour,
	})
}
