package utils

import (
	"errors"
	"fiermon-blog/internal/controllers/http/v1/values"
	"fiermon-blog/internal/models"
	"github.com/gin-gonic/gin"
)

// Устанавливаем userMeta
// Данные необходимые для аторизации пользователя к эндпоинтам
var ErrMetaUserNotFound = errors.New("user_meta not found")
var ErrMetaUserinvadil = errors.New("user_meta invalid")

var ErrMetaServiceNotFound = errors.New("service_meta not found")
var ErrMetaServiceInvalid = errors.New("service_meta invalid")

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

func SetServiceMeta(c *gin.Context, serviceMeta models.ServiceMeta) {
	c.Set(values.ServiceMetaKey, serviceMeta)
}

func GetServiceMeta(c *gin.Context) (*models.ServiceMeta, error) {
	data, ok := c.Get(values.ServiceMetaKey)
	if !ok {
		return nil, ErrMetaServiceNotFound
	}
	serviceMeta, ok := data.(models.ServiceMeta)
	if !ok {
		return nil, ErrMetaServiceInvalid
	}
	return &serviceMeta, nil
}
