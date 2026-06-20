package utils

import (
	"errors"
	"userctrl/internal/controllers/http/v1/values"
	entity "userctrl/internal/models"

	"github.com/gin-gonic/gin"
)

// Устанавливаем userMeta
// Данные необходимые для аторизации пользователя к эндпоинтам
var ErrMetaUserNotFound = errors.New("user_meta not found")
var ErrMetaUserinvadil = errors.New("user_meta invalid")

func SetUserMeta(c *gin.Context, userMeta entity.UserMeta) {
	c.Set(values.UserMetaKey, userMeta)
}

func GetUserMeta(c *gin.Context) (*entity.UserMeta, error) {
	data, ok := c.Get(values.UserMetaKey)
	if !ok {
		return nil, ErrMetaUserNotFound
	}
	userMeta, ok := data.(entity.UserMeta)
	if !ok {
		return nil, ErrMetaUserinvadil
	}
	return &userMeta, nil
}
