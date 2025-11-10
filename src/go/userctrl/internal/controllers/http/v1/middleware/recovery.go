package middleware

import "github.com/gin-gonic/gin"

// RecoveryMiddleware - предотвращение паники
func RecoveryMiddleware() gin.HandlerFunc {
	return gin.Recovery()
}
