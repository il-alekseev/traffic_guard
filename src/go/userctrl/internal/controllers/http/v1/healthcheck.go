package v1

import (
	"userctrl/internal/controllers/http/v1/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

// healthcheck - обработчик для проверки работоспособности микросервиса
func (s *Server) healthcheck(c *gin.Context) {
	response := dto.SuccessResponse{
		Message: "ready",
	}
	c.JSON(http.StatusOK, response)
}
