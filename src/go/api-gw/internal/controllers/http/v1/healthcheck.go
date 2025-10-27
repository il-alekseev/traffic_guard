package v1

import (
	"api-gateway/pkg/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) healthcheck(c *gin.Context) {
	response := models.DtoSuccessResponse{
		Message: "ready",
	}
	c.JSON(http.StatusOK, response)
}
