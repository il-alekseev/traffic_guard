package v1

import (
	"api-gateway/pkg/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary      Проверка состояния сервера
// @Tags healthcheck
// @Success      200  {object}  models.DtoSuccessResponse "Отфильтрованные логи"
// @Router       /api/v1/healthcheck [get]
func (s *Server) healthcheck(c *gin.Context) {
	response := models.DtoSuccessResponse{
		Message: "ready",
	}
	c.JSON(http.StatusOK, response)
}

// @Summary Получение версии сервиса
// @Description Возвращает информацию о версии
// @Tags utils
// @Produce json
// @Success 200 {object} models.DtoSuccessResponse
// @Router /api/v1/version [get]
func (s *Server) version(c *gin.Context) {
	response := models.DtoSuccessResponse{
		Message: s.devVersion,
	}
	c.JSON(http.StatusOK, response)
}
