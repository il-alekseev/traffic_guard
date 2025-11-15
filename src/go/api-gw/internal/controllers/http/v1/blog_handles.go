package v1

import (
	"api-gateway/internal/controllers/http/v1/utils"
	"api-gateway/pkg/blogserv/operations"
	"api-gateway/pkg/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// users godoc
// @Summary      Получение бизнес логов
// @Description  Получение бизнес логов с возможностью фильтрации и пагинации
// @Tags logs
// @Param        page    		query     int  		false  "Номер страницы с 1" default(1)
// @Param        limit   		query     int  		false  "Колличесвто отображаемых элементов на странице" default(10)
// @Param        role    		query     string  	false  "Фильтр по роли"
// @Param        contextID    	query     string  	false  "Фильтр по contextID"
// @Param        search    		query     string  	false  "фильтр по username/entity/description"
// @Security 	 BearerAuth
// @Success      200  {object}  models.ModelsLogs "Отфильтрованные логи"
// @Failure		 400 {object} 	models.ModelsAPIError "query params is not valid"
// @Failure 	 500 {object} 	models.ModelsAPIError "internal error"
// @Router       /api/v1/blog/logs [get]
func (s *Server) logs(c *gin.Context) {
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

	role := c.Query("role")
	contextID := c.Query("contextID")
	search := c.Query("search")

	// Создаем authInfoWriter для передачи токена
	authInfo, err := utils.GetAuthInfo(c)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "utils.GetAuthInfo(c)", err)
		return
	}

	resp, err := s.blogCL.Operations.GetAPIV1Logs(
		&operations.GetAPIV1LogsParams{
			Page:      &page,
			Limit:     &limit,
			Role:      &role,
			ContextID: &contextID,
			Search:    &search,
			Context:   c,
		},
		authInfo,
	)

	if err != nil {
		s.ErrorResponse(c, http.StatusInternalServerError, "s.blogCL.Operations.GetAPIV1Logs", err)
		return
	}

	c.JSON(resp.Code(), resp.GetPayload())
}

// postAddRecord - ручка для вставки записей в БД
// @Summary      Добавление логов в БД
// @Tags logs
// @Security BearerAuth
// @Param record body models.DtoDtoBusinessLog true "Структура записи бизнес лога"
// @Success      200  {object}  models.DtoSuccessResponse "Успех"
// @Failure 500 {object} models.ModelsAPIError "internal error"
// @Router       /api/v1/blog/add [post]
func (s *Server) postAddRecord(c *gin.Context) {
	// Создаем authInfoWriter для передачи токена
	authInfo, err := utils.GetAuthInfo(c)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "utils.GetAuthInfo(c)", err)
		return
	}

	var body models.DtoBusinessLog
	if err := c.ShouldBindJSON(&body); err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "invalid json body", err)
		return
	}

	status, err := s.blogCL.Operations.PostAPIV1Add(&operations.PostAPIV1AddParams{
		Context: c,
		Record:  &body,
	}, authInfo)
	if err != nil {
		s.ErrorResponse(c, http.StatusInternalServerError, "s.blogCL.Operations.PostAPIV1Add", err)
		return
	}
	c.JSON(status.Code(), status.GetPayload())
}
