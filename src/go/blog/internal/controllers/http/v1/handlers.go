package v1

import (
	"encoding/json"
	"errors"
	"fiermon-blog/internal/controllers/http/v1/dto"
	"fiermon-blog/internal/controllers/http/v1/utils"
	"io"
	"net/http"

	_ "fiermon-blog/docs"
	"fiermon-blog/internal/models"
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
func (s *Server) getLogs(c *gin.Context) {
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

// postAddRecord - ручка для вставки записей в БД
// @Summary      Добавление логов в БД
// @Param record body dto.DtoBusinessLog true "Структура записи бизнес лога"
// @Success      200  {object}  dto.SuccessResponse "Успех"
// @Security BearerAuth
// @Failure 400 {object} models.APIError "query params is not valid"
// @Failure 500 {object} models.APIError "internal error"
// @Router       /api/v1/add [post]
func (s *Server) postAddRecord(c *gin.Context) {
	log := s.logger.With("method", "postAddRecord")

	userMeta, err := utils.GetUserMeta(c)
	if err != nil {
		log.ErrorContext(c.Request.Context(), "failed to get userMeta", wsl.Err(err))
		s.ErrorResponse(c, http.StatusBadRequest, "utils.GetUserMeta", err)
		return
	}

	if userMeta.UUID == "" {
		log.ErrorContext(c.Request.Context(), "failed to get userMeta.UUID", wsl.Err(err))
		s.ErrorResponse(c, http.StatusBadRequest, "token param", errors.New("userID in token is not exists"))
		return
	}

	dtoBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.ErrorContext(c.Request.Context(), "failed to read body", wsl.Err(err))
		s.ErrorResponse(c, http.StatusBadRequest, "utils.ReadAll", err)
		return
	}

	var dtoBuisnesLog dto.DtoBusinessLog
	err = json.Unmarshal(dtoBody, &dtoBuisnesLog)
	if err != nil {
		log.ErrorContext(c.Request.Context(), "failed to unmarshal body", wsl.Err(err))
		s.ErrorResponse(c, http.StatusBadRequest, "utils.Unmarshal", err)
		return
	}

	var oldValueMap, newValueMap models.JSONB
	if dtoBuisnesLog.OldValue != (dto.Value{}) {
		data, err := json.Marshal(dtoBuisnesLog.OldValue)
		if err != nil {
			s.ErrorResponse(c, http.StatusInternalServerError, "failed to marshal OldValue", err)
			return
		}
		err = json.Unmarshal(data, &oldValueMap)
		if err != nil {
			s.ErrorResponse(c, http.StatusInternalServerError, "failed to unmarshal OldValue", err)
			return
		}
	}

	if dtoBuisnesLog.NewValue != (dto.Value{}) {
		data, err := json.Marshal(dtoBuisnesLog.NewValue)
		if err != nil {
			s.ErrorResponse(c, http.StatusInternalServerError, "failed to marshal NewValue", err)
			return
		}
		err = json.Unmarshal(data, &newValueMap)
		if err != nil {
			s.ErrorResponse(c, http.StatusInternalServerError, "failed to unmarshal NewValue", err)
			return
		}
	}

	buisnessLog := &models.BusinessLog{
		EventType:   dtoBuisnesLog.EventType,
		Entity:      dtoBuisnesLog.Entity,
		Username:    dtoBuisnesLog.Username,
		UserRole:    dtoBuisnesLog.UserRole,
		ContextID:   dtoBuisnesLog.ContextID,
		EntityID:    dtoBuisnesLog.EntityID,
		OldValue:    oldValueMap,
		NewValue:    newValueMap,
		Description: dtoBuisnesLog.Description,
		ContextStr:  dtoBuisnesLog.ContextStr,
	}

	err = s.serv.AddRecordToRepo(c.Request.Context(), buisnessLog)
	if err != nil {
		log.ErrorContext(c.Request.Context(), "failed to add record", wsl.Err(err))
		s.ErrorResponse(c, http.StatusInternalServerError, "utils.AddRecordToRepo", err)
		return
	}

	response := models.DtoSuccessResponse{
		Message: "ok",
	}
	c.JSON(http.StatusOK, response)
}

// healthcheck godoc
// @Summary      Проверка здоровья сервиса
// @Success      200  {object}  dto.SuccessResponse "ready"
// @Router       /api/v1/healthcheck [get]
func (s *Server) healthcheck(c *gin.Context) {
	response := models.DtoSuccessResponse{
		Message: "ready",
	}
	c.JSON(http.StatusOK, response)
}
