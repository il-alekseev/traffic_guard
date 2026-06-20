package middleware

import (
	"net/http"
	"strings"
	"tg-an/internal/controllers/http/v1/dto"
	"tg-an/internal/controllers/http/v1/utils"
	"tg-an/internal/models"
	"tg-an/pkg/slogger"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// CheckAuthHeader проверяет наличие и корректность заголовка Authorization.
// Возвращает middleware Gin, который:
//   - Проверяет наличие заголовка Authorization
//   - Проверяет формат "Bearer <token>"
//   - Прерывает запрос с кодом 401 при ошибках валидации
//
// Пример использования:
//
//	router.Use(middleware.CheckAuthHeader())
func CheckAuthHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		authToken := c.GetHeader("Authorization")

		if authToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				dto.ErrorResponse{
					Error: "authorization header is required",
				},
			)
			return
		}

		if !strings.HasPrefix(authToken, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				dto.ErrorResponse{
					Error: "Invalid token format. Expected 'Bearer <token>'",
				},
			)
			return
		}

		c.Next()
	}
}

// SetUserMetaData извлекает метаданные пользователя из JWT токена.
// Возвращает middleware Gin, который:
//   - Парсит JWT токен из заголовка Authorization
//   - Извлекает user UUID, username и роли из claims
//   - Добавляет UserMeta в контекст Gin
//   - Добавляет информацию о пользователе в логгер
//
// Ожидает токен в формате:
//
//	Authorization: Bearer <token>
//
// Требуемые claims в токене:
//   - sub (user UUID)
//   - preferred_username (username)
//   - roles или resource_access.grafana-sso.roles (роли пользователя)
//
// Прерывает запрос с кодом 401 при ошибках валидации токена.
func SetUserMetaData() gin.HandlerFunc {
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
		// TODO: валидация токена
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
		userMeta := models.UserMeta{
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

		// Проброс значений для логирования
		c.Request = c.Request.WithContext(
			slogger.WithLog(
				c.Request.Context(), "user_id", userMeta.UUID,
			),
		)
		c.Request = c.Request.WithContext(
			slogger.WithLog(
				c.Request.Context(), "username", userMeta.Username,
			),
		)
		c.Request = c.Request.WithContext(
			slogger.WithLog(
				c.Request.Context(), "client_role", userMeta.ClientRole,
			),
		)

		utils.SetUserMeta(c, userMeta)
		c.Next()
	}
}

// extractRoles извлекает роли пользователя из JWT claims.
//
// Проверяет несколько возможных мест хранения ролей в токене:
//   - Прямое поле "roles"
//   - resource_access.grafana-sso.roles
//   - realm_access.roles
//
// Параметры:
//   - claims: MapClaims из JWT токена
//
// Возвращает:
//   - []string: список ролей пользователя
//   - bool: true если роли найдены, false если нет
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

// ParseRole содержит результат парсинга роли пользователя.
type ParseRole struct {
	ClientRole string // Полная роль (например "CA-prod-env")
	ShortRole  string // Сокращенная роль (например "CA")
	ContextID  string // ID контекста (если есть в роли)
}

// parseUser анализирует строку роли пользователя и разбивает на компоненты.
// Поддерживает форматы ролей:
//   - "SA" (системный администратор)
//   - "CA-<context_id>" (контекстный администратор)
//   - "CO-<context_id>" (контекстный оператор)
//
// Параметры:
//   - clientRole: строка роли из токена
//
// Возвращает:
//   - ParseRole: структуру с разобранными компонентами роли
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
