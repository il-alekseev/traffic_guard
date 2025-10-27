package server

import (
	"errors"
	"net/http"

	_ "fiermon-blog/docs"
	"fiermon-blog/internal/models"
	"fiermon-blog/internal/server/dto"
	"fiermon-blog/internal/server/utils"
	"fiermon-blog/pkg/fslog/wsl"

	"github.com/gin-gonic/gin"
)

// logs godoc
// @Summary      Получение логов из БД с возможностью фильтрации и пагинации
// @Description  Возвращает массив логов
// @Param        page    query     int  false  "Номер страницы с 1" default(1)
// @Param        limit    query     int  false  "Колличесвто отображаемых элементов на странице" default(10)
// @Param        role    query     string  false  "Фильтр по роли"
// @Param        context_id    query     string  false  "Фильтр по contextID"
// @Param        search    query     string  false  "фильтр по username/entity/description"
// @Success      200  {object}  models.Logs "Отфильтрованные логи"
// @Failure 400 {object} models.APIError "query params is not valid"
// @Security BearerAuth
// @Failure 500 {object} models.APIError "internal error"
// @Router       /api/v1/logs [get]
func (s *Server) logs(c *gin.Context) {
	log := s.logger.With("method", "users")

	userMeta, err := utils.GetUserMeta(c)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "utils.GetUserMeta", err)
		return
	}

	if userMeta.UUID == "" {
		s.ErrorResponse(c, http.StatusBadRequest, "token param", errors.New("userID in token is not exists"))
		return
	}

	reqAny, exists := c.Get("dto")
	if !exists {
		log.ErrorContext(c.Request.Context(), "failed to validate dto")
		c.AbortWithStatusJSON(http.StatusInternalServerError, &models.APIError{
			Code:    http.StatusInternalServerError,
			Message: "failed to get records from database",
		})
	}

	validReq := reqAny.(dto.ValidateQuery)

	var (
		records []models.BusinessLog
		meta    models.Meta
	)

	if userMeta.ShortRole == allowedRoleCA {
		records, meta, err = s.serv.GetRecords(
			c.Request.Context(),
			userMeta,
			validReq.Page,
			validReq.Limit,
			allowedRoleCA,
			validReq.ContextID,
			validReq.Search,
		)
	} else if userMeta.ShortRole == allowedRoleCO {
		records, meta, err = s.serv.GetRecords(
			c.Request.Context(),
			userMeta,
			validReq.Page,
			validReq.Limit,
			allowedRoleCO,
			validReq.ContextID,
			validReq.Search,
		)
	} else {
		records, meta, err = s.serv.GetRecords(
			c.Request.Context(),
			userMeta,
			validReq.Page,
			validReq.Limit,
			validReq.Role,
			validReq.ContextID,
			validReq.Search,
		)
	}

	if err != nil {
		log.ErrorContext(c.Request.Context(), "failed to get records from database", wsl.Err(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, &models.APIError{
			Code:    http.StatusInternalServerError,
			Message: "failed to get records from database",
		})
		return
	}

	resp := models.Logs{
		Data: records,
		Meta: &meta,
	}

	c.JSON(http.StatusOK, &resp)
}

// healthcheck godoc
// @Summary      Проверка здоровья сервиса
// @Success      200  {object}  models.DtoSuccessResponse "ready"
// @Router       /api/v1/healthcheck [get]
func (s *Server) healthcheck(c *gin.Context) {
	response := models.DtoSuccessResponse{
		Message: "ready",
	}
	c.JSON(http.StatusOK, response)
}
