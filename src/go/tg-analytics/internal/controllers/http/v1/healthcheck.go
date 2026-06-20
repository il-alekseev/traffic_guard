package v1

import (
	"net/http"
	"tg-an/internal/controllers/http/v1/dto"

	"github.com/gin-gonic/gin"
)

// @Summary Проверка работоспособности сервера
// @Description Проверка, что сервер работает
// @Tags utils
// @Produce json
// @Success 200 {object} dto.SuccessResponse
// @Router /api/v1/healthcheck [get]
func (s *Server) Healthcheck(c *gin.Context) {
	response := dto.SuccessResponse{
		Message: "ready",
	}
	c.JSON(http.StatusOK, response)
	s.l.Debug("Healthcheck")
}
