package v1

import (
	"net/http"
	"tg-etl/internal/controllers/http/v1/dto"

	"github.com/gin-gonic/gin"
)

// @Summary Проверка работоспособности сервера
// @Description Проверка, что сервер работает
// @Tags utils
// @Produce json
// @Success 200 {object} dto.SuccessResponse
// @Router /healthcheck [get]
func (s *Server) Healthcheck(c *gin.Context) {
	response := dto.SuccessResponse{
		Message: "ready",
	}
	c.JSON(http.StatusOK, response)
	s.l.Debug("Healthcheck")
}

// @Summary Получение версии ETL
// @Description Возвращает информацию о версии, времени сборки и коммите
// @Tags utils
// @Produce json
// @Success 200 {object} dto.SuccessResponse
// @Router /version [get]
func (s *Server) Version(c *gin.Context) {
	response := dto.SuccessResponse{
		Message: s.devVersion,
	}
	c.JSON(http.StatusOK, response)
	s.l.Debug("version")
}
