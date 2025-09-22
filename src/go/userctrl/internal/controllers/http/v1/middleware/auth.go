package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// AuthMiddleware - middleware для аутентификации запросов по JWT-токену
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Извлекает токен из заголовка "Token", проверяет его подпись и валидность.
		tokenString := c.GetHeader("Token")
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// Здесь должна быть проверка токена через Keycloak
			// Временная заглушка - возвращаем ключ для примера
			return []byte("secret"), nil
		})

		// При отсутствии токена или его невалидности возвращается ошибка 401.
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		c.Next()
	}
}

// VerifyToken - обработчик для проверки JWT-токена, подписанного RSA-ключом
func VerifyToken(c *gin.Context) {
	SecretKey := "-----BEGIN CERTIFICATE-----\n" +
		"MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA3WYfbPJX6" +
		"o1EH0mB5Uj7csO1ZfdHK3yCPZVKMaUi034xRoDNrqwgOIKEqyv7Sy" +
		"iTY+3+0CGwBdjAdtb+zvZr0H4PD9pRtKEutla/NnvuqKF1Au5k/YK" +
		"Iyrlvl8J4guBNAH43LBLVEUzNXFubvYUIxlWeTcT6e78qUaDeVaC1" +
		"pePDsGeTtECJrwbzvBBNAOez+mfeo1H+pegmtTUYjuWLuX+3eLX9k" +
		"StjfjWimSXtpZ+O2cmHP9By1oqgtmfEHfZwSmnpFF//EwDxDAx7hb" +
		"lyy0i4nRkRkud6RkzMqGBEAoj+r0jFi8NyARSjjsrOBXqeBBz69i4" +
		"EJDI6afx01ym5jwIDAQAB" + "\n-----END CERTIFICATE-----"

	reqToken := c.GetHeader("Authorization")

	key, er := jwt.ParseRSAPublicKeyFromPEM([]byte(SecretKey))
	if er != nil {
		fmt.Println(er)
		c.Abort()
		c.Writer.WriteHeader(http.StatusUnauthorized)
		c.Writer.Write([]byte("Unauthorized"))
		return
	}

	token, err := jwt.Parse(reqToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return key, nil
	})

	if err != nil {
		fmt.Println(err)
		c.Abort()
		c.Writer.WriteHeader(http.StatusUnauthorized)
		c.Writer.Write([]byte("Unauthorized"))
		return
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println("token is valid")
	}
}
