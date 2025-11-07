package v1

import (
	"api-gateway/internal/controllers/http/v1/utils"
	"api-gateway/pkg/ctxcontrol/contexts"
	"api-gateway/pkg/models"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary Получение списка всех контекстов в системе
// @Description Получение списка всех контекстов системы с возможностью пагинации и фильтрации
// @Tags contexts
// @Accept json
// @Produce json
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Количество элементов на странице" default(10)
// @Param search query string false "Фильтр по названию/id контекста"
// @Param order_by query string false "Поле для сортировки default(context_id)" Enums(context_id, created_at) default(context_id)
// @Param order_dir query string false "Направление сортировки default(desc)" Enums(asc,desc) default(desc)
// @Security BearerAuth
// @Success 200 {object} models.DtoContextListResponse
// @Failure 400 {object} models.DtoErrorResponse
// @Failure 401 {object} models.DtoErrorResponse
// @Failure 403 {object} models.DtoErrorResponse
// @Failure 500 {object} models.DtoErrorResponse
// @Router /api/v1/contexts [get]
func (s *Server) getContexts(c *gin.Context) {
	page, err := strconv.ParseInt(c.Query("page"), 10, 64)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "'page' value must be int64", err)
		return
	}
	limit, err := strconv.ParseInt(c.Query("limit"), 10, 64)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "'limit' value must be int64", err)
		return
	}
	// TODO: обработка ошибки
	search := c.Query("search")
	orderBy := c.Query("order_by")
	orderDir := c.Query("order_dir")

	// Создаем authInfoWriter для передачи токена
	authInfo, err := utils.GetAuthInfo(c)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "utils.GetAuthInfo(c)", err)
		return
	}

	resp, err := s.ctxCl.Contexts.GetV1Contexts(
		&contexts.GetV1ContextsParams{
			Page:     &page,
			Limit:    &limit,
			Search:   &search,
			OrderBy:  &orderBy,
			OrderDir: &orderDir,
			Context:  c,
		},
		authInfo,
	)

	if err != nil {
		if conflictErr, ok := err.(ResponseErrorInterface); ok {
			c.JSON(conflictErr.Code(), conflictErr.GetPayload())
			return
		}

		s.ErrorResponse(c, http.StatusBadRequest, "GetV1Contexts", err)
		return
	}
	c.JSON(resp.Code(), resp.GetPayload())
}

// @Summary Добавление нового контекста в систему
// @Description Создание нового контекста с указанным названием
// @Tags contexts
// @Accept plain
// @Produce json
// @Param input body models.DtoCreateContextRequest true "Название контекста"
// @Security BearerAuth
// @Success 200 {object} models.ModelsContext
// @Failure 400 {object} models.DtoErrorResponse
// @Failure 401 {object} models.DtoErrorResponse
// @Failure 403 {object} models.DtoErrorResponse
// @Failure 500 {object} models.DtoErrorResponse
// @Router /api/v1/contexts [post]
func (s *Server) createContext(c *gin.Context) {
	var req models.DtoCreateContextRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.DtoErrorResponse{Error: err.Error()})
		return
	}
	// TODO: валидация запроса на корректность полей

	// Создаем authInfoWriter для передачи токена
	authInfo, err := utils.GetAuthInfo(c)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "utils.GetAuthInfo(c)", err)
		return
	}

	resp, err := s.ctxCl.Contexts.PostV1Contexts(
		&contexts.PostV1ContextsParams{
			Input:   &req,
			Context: c,
		},
		authInfo,
	)

	if err != nil {
		if conflictErr, ok := err.(ResponseErrorInterface); ok {
			c.JSON(conflictErr.Code(), conflictErr.GetPayload())
			return
		}

		s.ErrorResponse(c, http.StatusBadRequest, "PostV1Contexts", err)
		return
	}
	c.JSON(resp.Code(), resp.GetPayload())
}

// @Summary Получение контекста по его идентификатору
// @Description Получение информации о контексте по его ID
// @Tags contexts
// @Accept json
// @Produce json
// @Param context_id path string true "ID контекста"
// @Security BearerAuth
// @Success 200 {object} models.ModelsContext
// @Failure 400 {object} models.DtoErrorResponse
// @Failure 401 {object} models.DtoErrorResponse
// @Failure 403 {object} models.DtoErrorResponse
// @Failure 404 {object} models.DtoErrorResponse
// @Failure 500 {object} models.DtoErrorResponse
// @Router /api/v1/contexts/{context_id} [get]
func (s *Server) getContextByID(c *gin.Context) {
	id := c.Param("context_id")
	if id == "" {
		s.ErrorResponse(c, http.StatusBadRequest, "Param", errors.New("context_id is empty"))
		return
	}

	authInfo, err := utils.GetAuthInfo(c)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "GetAuthInfo", err)
		return
	}

	resp, err := s.ctxCl.Contexts.GetV1ContextsContextID(
		&contexts.GetV1ContextsContextIDParams{
			ContextID: id,
			Context:   c,
		},
		authInfo,
	)

	if err != nil {
		if conflictErr, ok := err.(ResponseErrorInterface); ok {
			c.JSON(conflictErr.Code(), conflictErr.GetPayload())
			return
		}

		s.ErrorResponse(c, http.StatusBadRequest, "GetV1ContextsContextID", err)
		return
	}
	c.JSON(resp.Code(), resp.GetPayload())
}

// @Summary Удаление контекста
// @Description Удаление контекста по его ID
// @Tags contexts
// @Accept json
// @Produce json
// @Param context_id path string true "ID контекста"
// @Security BearerAuth
// @Success 200 {object} models.DtoSuccessResponse
// @Failure 400 {object} models.DtoErrorResponse
// @Failure 401 {object} models.DtoErrorResponse
// @Failure 403 {object} models.DtoErrorResponse
// @Failure 404 {object} models.DtoErrorResponse
// @Failure 500 {object} models.DtoErrorResponse
// @Router /api/v1/contexts/{context_id} [delete]
func (s *Server) deleteContext(c *gin.Context) {
	id := c.Param("context_id")
	if id == "" {
		s.ErrorResponse(c, http.StatusBadRequest, "Param", errors.New("context_id is empty"))
		return
	}

	authInfo, err := utils.GetAuthInfo(c)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "utils.GetAuthInfo(c)", err)
		return
	}

	resp, err := s.ctxCl.Contexts.DeleteV1ContextsContextID(
		&contexts.DeleteV1ContextsContextIDParams{
			ContextID: id,
			Context:   c,
		},
		authInfo,
	)

	if err != nil {
		if conflictErr, ok := err.(ResponseErrorInterface); ok {
			c.JSON(conflictErr.Code(), conflictErr.GetPayload())
			return
		}

		s.ErrorResponse(c, http.StatusBadRequest, "DeleteV1ContextsContextID", err)
		return
	}
	c.JSON(resp.Code(), resp.GetPayload())
}

// @Summary Получение числа всех контекстов системы
// @Description Получение общего количества контекстов в системе
// @Tags contexts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.DtoCountResponse
// @Failure 400 {object} models.DtoErrorResponse
// @Failure 401 {object} models.DtoErrorResponse
// @Failure 403 {object} models.DtoErrorResponse
// @Failure 500 {object} models.DtoErrorResponse
// @Router /api/v1/contexts/count [get]
func (s *Server) getContextsCount(c *gin.Context) {
	authInfo, err := utils.GetAuthInfo(c)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "utils.GetAuthInfo(c)", err)
		return
	}

	resp, err := s.ctxCl.Contexts.GetV1ContextsCount(
		&contexts.GetV1ContextsCountParams{
			Context: c,
		},
		authInfo,
	)

	if err != nil {
		if conflictErr, ok := err.(ResponseErrorInterface); ok {
			c.JSON(conflictErr.Code(), conflictErr.GetPayload())
			return
		}

		s.ErrorResponse(c, http.StatusBadRequest, "GetV1ContextsCount", err)
		return
	}
	c.JSON(resp.Code(), resp.GetPayload())
}

// @Summary Обновление контекста по его идентификатору
// @Description Обновление названия контекста и/или описания
// @Tags contexts
// @Accept json
// @Produce json
// @Param context_id path string true "ID контекста"
// @Param input body models.DtoUpdateContextRequest true "Информация о контексте"
// @Security BearerAuth
// @Success 200 {object} models.ModelsContext
// @Success 200 {object} models.DtoCountResponse
// @Failure 400 {object} models.DtoErrorResponse
// @Failure 401 {object} models.DtoErrorResponse
// @Failure 403 {object} models.DtoErrorResponse
// @Failure 500 {object} models.DtoErrorResponse
// @Router /api/v1/contexts/{context_id} [put]
func (s *Server) updateContextByID(c *gin.Context) {
	id := c.Param("context_id")
	if id == "" {
		s.ErrorResponse(c, http.StatusBadRequest, "Param", errors.New("context_id is empty"))
		return
	}

	var input models.DtoUpdateContextRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "ShouldBindJSON", err)
		return
	}

	authInfo, err := utils.GetAuthInfo(c)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "utils.GetAuthInfo(c)", err)
		return
	}

	resp, err := s.ctxCl.Contexts.PutV1ContextsContextID(
		&contexts.PutV1ContextsContextIDParams{
			ContextID: id,
			Input:     &input,
			Context:   c,
		},
		authInfo,
	)
	if err != nil {
		if conflictErr, ok := err.(ResponseErrorInterface); ok {
			c.JSON(conflictErr.Code(), conflictErr.GetPayload())
			return
		}

		s.ErrorResponse(c, http.StatusBadRequest, "PutV1ContextsContextID", err)
		return
	}
	c.JSON(resp.Code(), resp.GetPayload())
}

// @Summary Получение списка свободных портов
// @Description Выводит массив портов доступных для получения логов от МЭ
// @Tags contexts
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.DtoVectroFreePortsResponse
// @Failure 400 {object} models.DtoErrorResponse
// @Failure 401 {object} models.DtoErrorResponse
// @Failure 403 {object} models.DtoErrorResponse
// @Failure 500 {object} models.DtoErrorResponse
// @Router /api/v1/contexts/free-ports [get]
func (s *Server) getFreePorts(c *gin.Context) {
	authInfo, err := utils.GetAuthInfo(c)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "utils.GetAuthInfo(c)", err)
		return
	}

	resp, err := s.ctxCl.Contexts.GetV1ContextsFreePorts(
		&contexts.GetV1ContextsFreePortsParams{
			Context: c,
		},
		authInfo,
	)

	if err != nil {
		if conflictErr, ok := err.(ResponseErrorInterface); ok {
			c.JSON(conflictErr.Code(), conflictErr.GetPayload())
			return
		}

		s.ErrorResponse(c, http.StatusBadRequest, "GetV1ContextsFreePorts", err)
		return
	}
	c.JSON(resp.Code(), resp.GetPayload())
}
