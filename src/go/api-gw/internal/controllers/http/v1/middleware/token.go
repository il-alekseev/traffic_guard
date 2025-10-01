// Package middleware provides HTTP middleware functions for Gin framework.
package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// VerifyToken is a Gin middleware function that validates JWT tokens in the Authorization header.
//
// It expects:
//   - A valid RSA public key in PEM format as publicKey parameter
//   - An Authorization header with "Bearer " prefix followed by the JWT token
//
// The middleware performs the following validations:
//   - Checks for proper Authorization header format
//   - Verifies the token signature using the provided RSA public key
//   - Ensures the token uses RSA signing method
//   - Validates the token claims
//
// If any validation fails, it aborts the request with HTTP 401 Unauthorized status.
//
// Example usage:
//
//	r := gin.Default()
//	r.Use(middleware.VerifyToken(publicKey))
//
// Note: This middleware only validates the token but doesn't extract or use the claims.
func VerifyToken(publicKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		reqToken := c.GetHeader("Authorization")
		reqToken = strings.TrimPrefix(reqToken, "Bearer ")
		key, er := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
		if er != nil {
			slog.ErrorContext(c.Request.Context(), "ParseRSAPublicKeyFromPEM", "err", er.Error())
			c.Abort()
			c.Writer.WriteHeader(http.StatusUnauthorized)
			c.Writer.Write([]byte("Unauthorized: " + er.Error()))
			return
		}

		token, err := jwt.Parse(reqToken, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return key, nil
		})

		if err != nil {
			slog.ErrorContext(c.Request.Context(), "jwt.Parse", "err", err.Error())
			c.Abort()
			c.Writer.WriteHeader(http.StatusUnauthorized)
			c.Writer.Write([]byte("Unauthorized: " + err.Error()))
			return
		}

		if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			slog.DebugContext(c.Request.Context(), "token is valid")
		}

		c.Next()
	}
}
