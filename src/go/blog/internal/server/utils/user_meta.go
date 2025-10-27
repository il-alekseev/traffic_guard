package utils

import (
	"errors"
	"fiermon-blog/internal/models"
	"fiermon-blog/internal/server/values"

	"github.com/gin-gonic/gin"
)

// Устанавливаем userMeta
// Данные необходимые для аторизации пользователя к эндпоинтам
var ErrMetaUserNotFound = errors.New("user_meta not found")
var ErrMetaUserinvadil = errors.New("user_meta invalid")

func SetUserMeta(c *gin.Context, userMeta models.UserMeta) {
	c.Set(values.UserMetaKey, userMeta)
}

func GetUserMeta(c *gin.Context) (*models.UserMeta, error) {
	data, ok := c.Get(values.UserMetaKey)
	if !ok {
		return nil, ErrMetaUserNotFound
	}
	userMeta, ok := data.(models.UserMeta)
	if !ok {
		return nil, ErrMetaUserinvadil
	}
	return &userMeta, nil
}
