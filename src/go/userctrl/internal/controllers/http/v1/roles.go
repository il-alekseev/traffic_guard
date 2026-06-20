package v1

import (
	"errors"
	"userctrl/internal/controllers/http/v1/dto"
	"userctrl/internal/controllers/http/v1/utils"
	"userctrl/internal/controllers/http/v1/values"
	"userctrl/pkg/slogger/wsl"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// getUsers godoc
// @Summary Создание ролей и групп ролей для контекста
// @Description Создает роли CA и CO для клиента Grafana. Создает группу с именем контекста и в ней подругппы ca и co.
// @Description В группы добавляет соответствующие роли.
// @Tags roles
// @Param contextID query string false "ID контекста"
// @Security BearerAuth
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/roles/context [put]
func (s *Server) createRolesAndGroupsForContext(c *gin.Context) {
	userMeta, err := utils.GetUserMeta(c)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "utils.GetUserMeta", err)
		return
	}

	if userMeta.ClientRole != values.SystemAdmin {
		s.ErrorResponse(c, http.StatusForbidden, "", ErrAccessIsDenied)
		return
	}

	contextID := c.Query("contextID")
	if contextID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "contextID is empty"})
		return
	}

	// Создание ролей и групп ролей для контекста
	err = s.u.СreateRolesAndGroupsForContext(c.Request.Context(), userMeta, contextID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "success"})
}

// getUsers godoc
// @Summary Удаление ролей и групп ролей для контекста
// @Description Обратное действие к <i>/v1/roles/context [put]</i>
// @Tags roles
// @Param contextID query string false "ID контекста"
// @Security BearerAuth
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/roles/context [delete]
func (s *Server) deleteRolesAndGroupsForContext(c *gin.Context) {
	userMeta, err := utils.GetUserMeta(c)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "utils.GetUserMeta", err)
		return
	}

	if userMeta.ClientRole != values.SystemAdmin {
		s.ErrorResponse(c, http.StatusForbidden, "", ErrAccessIsDenied)
		return
	}

	contextID := c.Query("contextID")
	if contextID == "" {
		s.ErrorResponse(c, http.StatusBadRequest, "Query", errors.New("contextID is empty"))
		return
	}

	// Удаление ролей и групп ролей для контекста
	err = s.u.DeleteRolesAndGroupsForContext(c.Request.Context(), userMeta, contextID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "success"})
}

// getUsers godoc
// @Summary Получение списка ролей с пагинацией и фильтрацией
// @Description Возвращает список ролей с возможностью пагинации и фильтрации по названию
// @Tags roles
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Количество элементов на странице" default(10)
// @Param search query string false "Поиск по названию роли"
// @Security BearerAuth
// @Success 200 {object} dto.RolesListResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/roles [get]
func (s *Server) getRoles(c *gin.Context) {
	log := s.logger.With(wsl.Label("method", "getRoles"))
	userMeta, err := utils.GetUserMeta(c)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "utils.GetUserMeta", err)
		return
	}

	if userMeta.ClientRole != values.SystemAdmin {
		s.ErrorResponse(c, http.StatusForbidden, "", ErrAccessIsDenied)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	// Валидация параметров
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Возвращает список ролей с возможностью пагинации и фильтрации по названию
	roles, total, err := s.u.GetRoles(c.Request.Context(), userMeta, page, limit, search)
	if err != nil {
		log.ErrorContext(c.Request.Context(), "Failed to get roles", wsl.Err(err))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Internal server error"})
		return
	}

	// Формируем ответ
	response := dto.RolesListResponse{
		Data: *roles,
		Meta: dto.PaginationMeta{
			Page:  page,
			Limit: limit,
			Total: total,
			Pages: int(math.Ceil(float64(total) / float64(limit))),
		},
	}

	c.JSON(http.StatusOK, response)
}

// getRolesCount godoc
// @Summary Получение числа всех ролей системы
// @Description Получение числа всех ролей системы.
// @Tags roles
// @Security BearerAuth
// @Success 200 {object} dto.CountResponse "Количество ролей"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/roles/count [get]
func (s *Server) getRolesCount(c *gin.Context) {
	userMeta, err := utils.GetUserMeta(c)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "utils.GetUserMeta", err)
		return
	}

	if userMeta.ClientRole != values.SystemAdmin {
		s.ErrorResponse(c, http.StatusForbidden, "", ErrAccessIsDenied)
		return
	}

	// Получение числа всех ролей системы
	count, err := s.u.GetRolesCount(c.Request.Context(), userMeta)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}
	res := dto.CountResponse{
		Count: count,
	}
	c.JSON(http.StatusOK, res)
}
