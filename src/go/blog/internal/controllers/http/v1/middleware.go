package v1

import (
	"fiermon-blog/internal/controllers/http/v1/dto"
	"fiermon-blog/internal/models"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	allowedRoleSA = "SA"
	allowedRoleCA = "CA"
	allowedRoleCO = "CO"
)

// validateParams - Middleware валидация параметров
func (s *Server) validateParams() gin.HandlerFunc {
	log := s.logger.With("middleware", "validateParams")

	return func(c *gin.Context) {

		// считывание параметров из запроса
		pageStr := c.Request.URL.Query().Get("page")
		limitStr := c.Request.URL.Query().Get("limit")
		role := c.Request.URL.Query().Get("role")
		contextID := c.Request.URL.Query().Get("context")
		search := c.Request.URL.Query().Get("search")

		var (
			page, limit int
			err         error
		)

		// валидация полей, установка default значений
		if pageStr == "" {
			page = 1
		} else {
			page, err = strconv.Atoi(pageStr)
			if err != nil || page < 1 {
				log.WarnContext(c.Request.Context(), "error parsing page query parameter or value a negative number")
				page = 1
			}
		}

		if limitStr == "" {
			limit = 10
		} else {
			limit, err = strconv.Atoi(limitStr)
			if err != nil || limit < 1 {
				log.WarnContext(c.Request.Context(), "error parsing page query parameter or value a negative number")
				limit = 10
			}
		}

		// объединяем переменные в структуру и записываем для следующих обработчиков
		if role == allowedRoleSA || role == allowedRoleCA || role == "" || role == allowedRoleCO { //|| role == "CO"
			dtoObj := dto.FilterLogs{
				Page:      page,
				Limit:     limit,
				Role:      role,
				ContextID: contextID,
				Search:    search,
			}

			c.Set("dto", dtoObj)

			c.Next()
		} else {
			log.WarnContext(c.Request.Context(), "query params is not valid", slog.String("role", role))
			c.AbortWithStatusJSON(http.StatusBadRequest, &models.APIError{
				Code:    http.StatusBadRequest,
				Message: "user role is not valid",
			})
			return
		}
	}
}
