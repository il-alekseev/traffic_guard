package utils

import (
	"api-gateway/internal/controllers/http/v1/values"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

func GetAuthInfo(c *gin.Context) (runtime.ClientAuthInfoWriterFunc, error) {
	// Проверка токена осуществляется в middleware.CheckAuthHeader()
	// В apiRouter.Use(middleware.VerifyToken()) мы токен валидируем
	// Тут мы его спокойно забираем
	authToken := c.GetHeader("Authorization")
	authInfo := runtime.ClientAuthInfoWriterFunc(func(r runtime.ClientRequest, _ strfmt.Registry) error {
		if err := r.SetHeaderParam("Authorization", authToken); err != nil {
			return err
		}
		// Ещё и свеху устанавливаем заголовок X-Request-ID
		if err := r.SetHeaderParam(values.RequestIDHeader, GetRequestID(c)); err != nil {
			return err
		}
		return nil
	})

	return authInfo, nil
}
