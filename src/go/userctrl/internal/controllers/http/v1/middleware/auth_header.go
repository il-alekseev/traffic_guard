package middleware

import (
	//"fiermon-ctxcontrol/internal/controllers/http/v1/metadata"
	"userctrl/internal/controllers/http/v1/dto"
	"userctrl/internal/controllers/http/v1/utils"
	entity "userctrl/internal/models"
	"userctrl/pkg/slogger/wsl"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// CheckAuthHeader - проверка наличия Berar авторизации
func CheckAuthHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		authToken := c.GetHeader("Authorization")

		if !strings.HasPrefix(authToken, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				dto.ErrorResponse{
					Error: "Invalid token format. Expected 'Bearer <token>'",
				},
			)
			return
		}

		if authToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				dto.ErrorResponse{
					Error: "authorization header is required",
				},
			)
			return
		}
		c.Next()
	}
}

// SetUserMetaData - установка метаданных пользователя
func SetUserMetaData(l slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем токен из заголовка Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		// Проверяем формат заголовка (должен быть "Bearer <token>")
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			return
		}

		tokenString := tokenParts[1]

		// Парсим токен без проверки подписи (если нужно проверять подпись, добавьте ключ)
		// TODO: валидация ьокена
		token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		// Извлекаем sub (userID)
		userUUID, ok := claims["sub"].(string)
		if !ok || userUUID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "sub claim is required (user uuid)"})
			return
		}

		// Извлекаем preferred_username (username)
		username, ok := claims["preferred_username"].(string)
		if !ok || username == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "preferred_username claim is required (username)"})
			return
		}

		// Передаем userID в контекст
		userMeta := entity.UserMeta{
			UUID:     userUUID,
			Username: username,
		}

		// Извлекаем роли и проверяем нужные
		if roles, ok := extractRoles(claims); ok {
			for _, role := range roles {
				if role == "SA" || strings.HasPrefix(role, "CA-") || strings.HasPrefix(role, "CO-") {
					parseRole := parseUser(role)
					userMeta.ClientRole = parseRole.ClientRole
					userMeta.ShortRole = parseRole.ShortRole
					userMeta.ContextID = parseRole.ContextID
					break
				}
			}
		}

		l.InfoContext(c.Request.Context(), "SetUserMetaData", wsl.Label("usermeta", fmt.Sprintf("username: %s, userID: %s, clientRole: %s, shortRole: %s, contextID: %s", userMeta.Username,
			userMeta.UUID,
			userMeta.ClientRole,
			userMeta.ShortRole,
			userMeta.ContextID)),
		)

		utils.SetUserMeta(c, userMeta)
		c.Next()
	}
}

// extractRoles извлекает роли из claims (проверяет несколько возможных мест)
func extractRoles(claims jwt.MapClaims) ([]string, bool) {
	// Попробуем получить из прямого поля roles
	if roles, ok := claims["roles"].([]interface{}); ok {
		var result []string
		for _, r := range roles {
			if role, ok := r.(string); ok {
				result = append(result, role)
			}
		}
		return result, true
	}

	// Попробуем получить из resource_access.grafana-sso.roles
	if resourceAccess, ok := claims["resource_access"].(map[string]interface{}); ok {
		if grafanaSSO, ok := resourceAccess["grafana-sso"].(map[string]interface{}); ok {
			if roles, ok := grafanaSSO["roles"].([]interface{}); ok {
				var result []string
				for _, r := range roles {
					if role, ok := r.(string); ok {
						result = append(result, role)
					}
				}
				return result, true
			}
		}
	}

	// Попробуем получить из realm_access.roles
	if realmAccess, ok := claims["realm_access"].(map[string]interface{}); ok {
		if roles, ok := realmAccess["roles"].([]interface{}); ok {
			var result []string
			for _, r := range roles {
				if role, ok := r.(string); ok {
					result = append(result, role)
				}
			}
			return result, true
		}
	}

	return nil, false
}

type ParseRole struct {
	ClientRole string
	ShortRole  string
	ContextID  string
}

// parseUser - вспомогательная функция парсинга мета данных
func parseUser(clientRole string) ParseRole {
	res := ParseRole{
		ClientRole: clientRole,
	}

	switch {
	case clientRole == "SA":
		res.ShortRole = "SA"
	case strings.HasPrefix(clientRole, "CA-"):
		context := strings.TrimPrefix(clientRole, "CA-")
		res.ShortRole = "CA"
		res.ContextID = context
	case strings.HasPrefix(clientRole, "CO-"):
		context := strings.TrimPrefix(clientRole, "CO-")
		res.ShortRole = "CO"
		res.ContextID = context
	}

	return res
}
