package utils

import (
	"api-gateway/internal/controllers/http/v1/metadata"
	"api-gateway/internal/controllers/http/v1/values"
	"errors"

	"github.com/gin-gonic/gin"
)

// Устанавливаем userMeta
// Данные необходимые для аторизации пользователя к эндпоинтам
var ErrMetaUserNotFound = errors.New("user_meta not found")
var ErrMetaUserinvadil = errors.New("user_meta invalid")

func SetUserMeta(c *gin.Context, userMeta metadata.UserMeta) {
	c.Set(values.UserMetaKey, userMeta)
}

func GetUserMeta(c *gin.Context) (*metadata.UserMeta, error) {
	data, ok := c.Get(values.UserMetaKey)
	if !ok {
		return nil, ErrMetaUserNotFound
	}
	userMeta, ok := data.(metadata.UserMeta)
	if !ok {
		return nil, ErrMetaUserinvadil
	}
	return &userMeta, nil
}
