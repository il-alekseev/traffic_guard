package v1

import (
	"api-gateway/internal/controllers/http/v1/utils"
	"api-gateway/pkg/analytics/actions"
	"api-gateway/pkg/analytics/common"
	"api-gateway/pkg/analytics/dashboards"
	"api-gateway/pkg/analytics/detections"
	"api-gateway/pkg/analytics/reports"
	"api-gateway/pkg/analytics/sessions"
	"api-gateway/pkg/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

//---------------------common---------------------

// getCategories -
// @Summary Получение списка категорий контента
// @Description Возвращает список всех уникальных категорий контента из системы
// @Tags common
// @Produce json
// @Security BearerAuth
// @Success 200 {array} string
// @Failure 403 {object} models.DtoErrorResponse "Недостаточно прав для доступа к статистике выявлений"
// @Failure 500 {object} models.DtoErrorResponse "Ошибка при получении списка категорий"
// @Router /api/v1/analytics/categories [get]
func (s *Server) getCategories(c *gin.Context) {
	//Создаем authInfoWriter для передачи токена
	//authInfo, err := utils.GetAuthInfo(c)
	//if err != nil {
	//	s.ErrorResponse(c, http.StatusBadRequest, "utils.GetAuthInfo(c)", err)
	//	return
	//}

	resp, err := s.analyticsCL.Common.GetAPIV1Categories(&common.GetAPIV1CategoriesParams{})
	if err != nil {
		if conflictErr, ok := err.(ResponseErrorInterface); ok {
			c.JSON(conflictErr.Code(), conflictErr.GetPayload())
			return
		}

		s.ErrorResponse(c, http.StatusBadRequest, "Common.GetCategories", err)
		return
	}

	c.JSON(resp.Code(), resp.GetPayload())
}

// getDevices -
// @Summary Получить список имен устройств
// @Description Возвращает список всех уникальных имен устройств (хостов) из системы
// @Tags common
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} string "Список имен устройств"
// @Failure 403 {object} models.DtoErrorResponse "Недостаточно прав для доступа к статистике выявлений"
// @Failure 500 {object} models.DtoErrorResponse "Ошибка при получении имен устройств"
// @Router /api/v1/analytics/devices [get]
func (s *Server) getDevices(c *gin.Context) {
	// Создаем authInfoWriter для передачи токена
	authInfo, err := utils.GetAuthInfo(c)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "utils.GetAuthInfo(c)", err)
		return
	}

	resp, err := s.analyticsCL.Common.GetAPIV1Devices(&common.GetAPIV1DevicesParams{}, authInfo)
	if err != nil {
		if conflictErr, ok := err.(ResponseErrorInterface); ok {
			c.JSON(conflictErr.Code(), conflictErr.GetPayload())
			return
		}

		s.ErrorResponse(c, http.StatusBadRequest, "Common.GetDevices", err)
		return
	}

	c.JSON(resp.Code(), resp.GetPayload())
}

//---------------------dashboards---------------------

// getV1DashboardsAnomalies -
// @Summary Получение информации об аномалиях
// @Description Получение статистики об аномалиях за указанный период
// @Tags dashboards
// @Produce application/json
// @Security BearerAuth
// @Param from query string false "Начало временного диапазона (формат: now-10m, 2023-12-01T10:00:00Z)" default(now-10m)
// @Param to query string false "Конец временного диапазона (формат: now, 2023-12-01T12:00:00Z)" default(now)
// @Param hostname query string false "Фильтр по имени хоста"
// @Success 200 {object} models.DtoGetAnomaliesResponse "Список обнаруженных аномалий"
// @Failure 400 {object} models.DtoErrorResponse "Неверный формат параметров"
// @Failure 403 {object} models.DtoErrorResponse "Недостаточно прав для доступа к статистике выявлений"
// @Failure 500 {object} models.DtoErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/analytics/dashboards/anomalies [get]
func (s *Server) getV1DashboardsAnomalies(c *gin.Context) {
	// Создаем authInfoWriter для передачи токена
	authInfo, err := utils.GetAuthInfo(c)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "utils.GetAuthInfo(c)", err)
		return
	}

	//Парсим входные данные
	from := c.DefaultQuery("from", "now-10m")
	to := c.DefaultQuery("to", "now")
	hostname := c.Query("hostname")

	resp, err := s.analyticsCL.Dashboards.GetAPIV1DashboardsAnomalies(&dashboards.GetAPIV1DashboardsAnomaliesParams{
		From:     &from,
		To:       &to,
		Hostname: &hostname,
	}, authInfo)
	if err != nil {
		if conflictErr, ok := err.(ResponseErrorInterface); ok {
			c.JSON(conflictErr.Code(), conflictErr.GetPayload())
			return
		}

		s.ErrorResponse(c, http.StatusBadRequest, "Dashboards.GetAPIV1DashboardsAnomalies", err)
		return
	}

	c.JSON(resp.Code(), resp.GetPayload())
}

// getV1DashboardsDevices -
// @Summary Получение статистики по устройствам
// @Description Получение агрегированной статистики по сетевым узлам за указанный период
// @Tags dashboards
// @Produce application/json
// @Security BearerAuth
// @Param from query string false "Начало временного диапазона (формат: now-10m, 2023-12-01T10:00:00Z)" default(now-10m)
// @Param to query string false "Конец временного диапазона (формат: now, 2023-12-01T12:00:00Z)" default(now)
// @Param count query int false "Количество точек измерений" default(20) minimum(1)
// @Success 200 {object} models.DtoDeviceStatResponse "Статистика по устройствам"
// @Failure 400 {object} models.DtoErrorResponse "Неверный формат параметров"
// @Failure 403 {object} models.DtoErrorResponse "Недостаточно прав для доступа к статистике выявлений"
// @Failure 500 {object} models.DtoErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/analytics/dashboards/devices [get]
func (s *Server) getV1DashboardsDevices(c *gin.Context) {
	// Создаем authInfoWriter для передачи токена
	authInfo, err := utils.GetAuthInfo(c)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "utils.GetAuthInfo(c)", err)
		return
	}

	//Парсим входные данные
	from := c.DefaultQuery("from", "now-10m")
	to := c.DefaultQuery("to", "now")
	count, err := strconv.ParseInt(c.Query("count"), 10, 64)
	if err != nil {
		if conflictErr, ok := err.(ResponseErrorInterface); ok {
			c.JSON(conflictErr.Code(), conflictErr.GetPayload())
			return
		}

		s.ErrorResponse(c, http.StatusBadRequest, "Parse int count", err)
		return
	}

	resp, err := s.analyticsCL.Dashboards.GetAPIV1DashboardsDevices(&dashboards.GetAPIV1DashboardsDevicesParams{
		From:  &from,
		To:    &to,
		Count: &count,
	}, authInfo)
	if err != nil {
		if conflictErr, ok := err.(ResponseErrorInterface); ok {
			c.JSON(conflictErr.Code(), conflictErr.GetPayload())
			return
		}

		s.ErrorResponse(c, http.StatusBadRequest, "Dashboards.GetAPIV1DashboardsDevices", err)
		return
	}

	c.JSON(resp.Code(), resp.GetPayload())
}

// getDashboardsRequests -
// @Summary Получение статистики запросов
// @Description Получение статистики запросов за указанный период с фильтрацией по хосту и типу запросов
// @Tags dashboards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param from query string false "Начало временного диапазона (формат: now-10m, now-1h, 2024-01-01T00:00:00Z)" default(now-10m)
// @Param to query string false "Конец временного диапазона (формат: now, 2024-01-01T00:00:00Z)" default(now)
// @Param request_type query string false "Тип запроса" Enums(allowed, blocked, before_block, pending) default(pending)
// @Param hostname query string false "Фильтр по имени хоста"
// @Param count query integer false "Количество интервалов" default(10)
// @Success 200 {object} models.DtoRequestStatResponse "Статистика запросов (массив чисел)"
// @Failure 400 {object} models.DtoErrorResponse "Неверный формат параметров"
// @Failure 403 {object} models.DtoErrorResponse "Недостаточно прав для доступа к статистике выявлений"
// @Failure 500 {object} models.DtoErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/analytics/dashboards/requests [get]
func (s *Server) getDashboardsRequests(c *gin.Context) {
	// Создаем authInfoWriter для передачи токена
	authInfo, err := utils.GetAuthInfo(c)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "utils.GetAuthInfo(c)", err)
		return
	}

	//Парсим входные данные
	from := c.DefaultQuery("from", "now-10m") // по умолчанию выдает последние 10 минут
	to := c.DefaultQuery("to", "now")
	requestType := c.Query("request_type")
	hostname := c.Query("hostname")
	count, err := strconv.ParseInt(c.Query("count"), 10, 64)
	if err != nil {
		if conflictErr, ok := err.(ResponseErrorInterface); ok {
			c.JSON(conflictErr.Code(), conflictErr.GetPayload())
			return
		}

		s.ErrorResponse(c, http.StatusBadRequest, "Parse int count", err)
		return
	}

	resp, err := s.analyticsCL.Dashboards.GetAPIV1DashboardsRequests(&dashboards.GetAPIV1DashboardsRequestsParams{
		From:        &from,
		To:          &to,
		RequestType: &requestType,
		Hostname:    &hostname,
		Count:       &count,
	}, authInfo)
	if err != nil {
		if conflictErr, ok := err.(ResponseErrorInterface); ok {
			c.JSON(conflictErr.Code(), conflictErr.GetPayload())
			return
		}

		s.ErrorResponse(c, http.StatusBadRequest, "Dashboards.GetAPIV1DashboardsRequests", err)
		return
	}

	c.JSON(resp.Code(), resp.GetPayload())
}

// getDashboardsTopCategories -
// @Summary Получение списка самых запрашиваемых категорий
// @Description Возвращает наиболее часто встречаемые категории в сессиях с возможностью фильтрации
// @Tags dashboards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param from query string false "Начало временного диапазона" default(now-24h)
// @Param to query string false "Конец временного диапазона" default(now)
// @Param hostname query string false "Фильтр по имени хоста"
// @Param type query string false "Фильтр по типу сессии" Enums(Разрешен, Заблокирован, VPN)
// @Param count query int false "Количество возвращаемых категорий" default(5) minimum(1) maximum(50)
// @Success 200 {array} models.DtoCategory
// @Failure 400 {object} models.DtoErrorResponse
// @Failure 403 {object} models.DtoErrorResponse "Недостаточно прав для доступа к статистике выявлений"
// @Failure 500 {object} models.DtoErrorResponse
// @Router /api/v1/analytics/dashboards/top-categories [get]
func (s *Server) getDashboardsTopCategories(c *gin.Context) {
	// Создаем authInfoWriter для передачи токена
	authInfo, err := utils.GetAuthInfo(c)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "utils.GetAuthInfo(c)", err)
		return
	}

	//Парсим временные метки
	from := c.DefaultQuery("from", "now-10m") // по умолчанию выдает последние 10 минут
	to := c.DefaultQuery("to", "now")
	hostname := c.Query("hostname")
	typeStr := c.Query("type")
	count, err := strconv.ParseInt(c.Query("count"), 10, 64)
	if err != nil {
		if conflictErr, ok := err.(ResponseErrorInterface); ok {
			c.JSON(conflictErr.Code(), conflictErr.GetPayload())
			return
		}

		s.ErrorResponse(c, http.StatusBadRequest, "Parse int count", err)
		return
	}

	resp, err := s.analyticsCL.Dashboards.GetAPIV1DashboardsTopCategories(&dashboards.GetAPIV1DashboardsTopCategoriesParams{
		From:     &from,
		To:       &to,
		Hostname: &hostname,
		Type:     &typeStr,
		Count:    &count,
	}, authInfo)
	if err != nil {
		if conflictErr, ok := err.(ResponseErrorInterface); ok {
			c.JSON(conflictErr.Code(), conflictErr.GetPayload())
			return
		}

		s.ErrorResponse(c, http.StatusBadRequest, "Dashboards.GetAPIV1DashboardsTopCategories", err)
		return
	}

	c.JSON(resp.Code(), resp.GetPayload())
}

// getDashboardsTopUnresolvedDetections -
// @Summary Получение списка топ нерешенных выявлений
// @Description Возвращает список наиболее частых нерешенных выявлений за указанный временной период с возможностью фильтрации
// @Tags dashboards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param from query string false "Начало временного диапазона (формат: now-10m, 2023-12-01T10:00:00Z)" default(now-10m)
// @Param to query string false "Конец временного диапазона (формат: now, 2023-12-01T12:00:00Z)" default(now)
// @Param hostname query string false "Фильтр по имени хоста"
// @Param count query int false "Количество возвращаемых записей" default(5) minimum(1)
// @Success 200 {array} models.DtoUnresolvedDetection
// @Failure 400 {object} models.DtoErrorResponse
// @Failure 403 {object} models.DtoErrorResponse "Недостаточно прав для доступа к статистике выявлений"
// @Failure 500 {object} models.DtoErrorResponse
// @Router /api/v1/analytics/dashboards/top-unresolved_detections [get]
func (s *Server) getDashboardsTopUnresolvedDetections(c *gin.Context) {
	// Создаем authInfoWriter для передачи токена
	authInfo, err := utils.GetAuthInfo(c)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "utils.GetAuthInfo(c)", err)
		return
	}

	//Парсим временные метки
	from := c.DefaultQuery("from", "now-10m") // по умолчанию выдает последние 10 минут
	to := c.DefaultQuery("to", "now")
	hostname := c.Query("hostname")
	count, err := strconv.ParseInt(c.Query("count"), 10, 64)
	if err != nil {
		if conflictErr, ok := err.(ResponseErrorInterface); ok {
			c.JSON(conflictErr.Code(), conflictErr.GetPayload())
			return
		}

		s.ErrorResponse(c, http.StatusBadRequest, "error Parse int count", err)
		return
	}

	resp, err := s.analyticsCL.Dashboards.GetAPIV1DashboardsTopUnresolvedDetections(&dashboards.GetAPIV1DashboardsTopUnresolvedDetectionsParams{
		From:     &from,
		To:       &to,
		Hostname: &hostname,
		Count:    &count,
	}, authInfo)
	if err != nil {
		if conflictErr, ok := err.(ResponseErrorInterface); ok {
			c.JSON(conflictErr.Code(), conflictErr.GetPayload())
			return
		}

		s.ErrorResponse(c, http.StatusBadRequest, "Dashboards.GetAPIV1DashboardsTopUnresolvedDetections", err)
		return
	}

	c.JSON(resp.Code(), resp.GetPayload())
}

// getDashboardsTraffic - возвращает статистику трафика за указанный временной диапазон
// @Summary Получение статистики трафика
// @Description Возвращает статистику трафика за указанный временной диапазон с заданным количеством точек данных в Кб
// @Tags dashboards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param from query string false "Начало временного диапазона в формате парсера времени (по умолчанию now-10m)" default(now-10m)
// @Param to query string false "Конец временного диапазона в формате парсера времени (по умолчанию now)" default(now)
// @Param hostname query string false "Фильтр по имени хоста"
// @Param count query integer false "Количество точек данных для возврата (по умолчанию 20)" minimum(1) default(20)
// @Success 200 {object} models.DtoTrafficStatResponse "Успешный ответ со статистикой трафика"
// @Failure 400 {object} models.DtoErrorResponse "Неверный формат параметров запроса"
// @Failure 403 {object} models.DtoErrorResponse "Недостаточно прав для доступа к статистике выявлений"
// @Failure 500 {object} models.DtoErrorResponse "Внутренняя ошибка сервера при получении статистики"
// @Router /api/v1/analytics/dashboards/traffic [get]
func (s *Server) getDashboardsTraffic(c *gin.Context) {
	// Создаем authInfoWriter для передачи токена
	authInfo, err := utils.GetAuthInfo(c)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "utils.GetAuthInfo(c)", err)
		return
	}

	//Парсим временные метки
	from := c.DefaultQuery("from", "now-10m") // по умолчанию выдает последние 10 минут
	to := c.DefaultQuery("to", "now")
	hostname := c.Query("hostname")
	count, err := strconv.ParseInt(c.Query("count"), 10, 64)
	if err != nil {
		if conflictErr, ok := err.(ResponseErrorInterface); ok {
			c.JSON(conflictErr.Code(), conflictErr.GetPayload())
			return
		}

		s.ErrorResponse(c, http.StatusBadRequest, "Parse int count", err)
		return
	}

	resp, err := s.analyticsCL.Dashboards.GetAPIV1DashboardsTraffic(&dashboards.GetAPIV1DashboardsTrafficParams{
		From:     &from,
		To:       &to,
		Hostname: &hostname,
		Count:    &count,
	}, authInfo)
	if err != nil {
		if conflictErr, ok := err.(ResponseErrorInterface); ok {
			c.JSON(conflictErr.Code(), conflictErr.GetPayload())
			return
		}

		s.ErrorResponse(c, http.StatusBadRequest, "Dashboards.GetAPIV1DashboardsTraffic", err)
		return
	}

	c.JSON(resp.Code(), resp.GetPayload())
}

// @Summary Получение графика запрещенной активности за год
// @Description Возвращает статистику запрещенной активности за указанный год с возможностью фильтрации по имени устройства
// @Description Временной диапазон автоматически формируется от 1 января 00:00:00 UTC до 31 декабря 23:59:59.999 UTC указанного года
// @Tags dashboards
// @Accept json
// @Produce json
// @Param year query int false "год, за который нужно получить данные для графика запрещенной активности" default(2025)
// @Param hostname query string false "Фильтр по имени хоста"
// @Security BearerAuth
// @Success 200 {object} models.DtoGetProhActivityResponse "Данные графика запрещенной активности за год"
// @Failure 400 {object} models.DtoErrorResponse "Неверный формат параметров, некорректный год или временной диапазон"
// @Failure 403 {object} models.DtoErrorResponse "Недостаточно прав для доступа к графику запрещенной активности"
// @Failure 500 {object} models.DtoErrorResponse "Внутренняя ошибка сервера при получении данных графика"
// @Router /api/v1/analytics/dashboards/proh-activity [get]
func (s *Server) getDashboardsProhActivity(c *gin.Context) {
	// Создаем authInfoWriter для передачи токена
	authInfo, err := utils.GetAuthInfo(c)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "utils.GetAuthInfo(c)", err)
		return
	}

	//Парсим временные метки
	year, err := strconv.ParseInt(c.Query("year"), 10, 64)
	if err != nil {
		if conflictErr, ok := err.(ResponseErrorInterface); ok {
			c.JSON(conflictErr.Code(), conflictErr.GetPayload())
			return
		}

		s.ErrorResponse(c, http.StatusBadRequest, "Parse int year", err)
		return
	}
	hostname := c.Query("hostname")

	resp, err := s.analyticsCL.Dashboards.GetAPIV1DashboardsProhActivity(&dashboards.GetAPIV1DashboardsProhActivityParams{
		Year:     &year,
		Hostname: &hostname,
	}, authInfo)
	if err != nil {
		if conflictErr, ok := err.(ResponseErrorInterface); ok {
			c.JSON(conflictErr.Code(), conflictErr.GetPayload())
			return
		}

		s.ErrorResponse(c, http.StatusBadRequest, "Dashboards.GetAPIV1DashboardsProhActivity", err)
		return
	}

	c.JSON(resp.Code(), resp.GetPayload())
}

//---------------------actions---------------------

// @Summary Выполнение действия над выявлением
// @Description Устанавливает действие (разрешить/заблокировать) для указанного домена.
// @Tags actions
// @Accept json
// @Produce json
// @Param request body models.DtoDetectionActRequest true "Данные для выполнения действия над доменом"
// @Security BearerAuth
// @Success 200 {object} models.DtoSuccessResponse "Действие успешно применено к домену"
// @Failure 400 {object} models.DtoErrorResponse "Неверные параметры запроса"
// @Failure 403 {object} models.DtoErrorResponse "Недостаточно прав для выполнения действия"
// @Failure 404 {object} models.DtoErrorResponse "Домен не найден"
// @Failure 500 {object} models.DtoErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/analytics/detections/act [patch]
func (s *Server) patchV1DashboardsAct(c *gin.Context) {
	// Создаем authInfoWriter для передачи токена
	authInfo, err := utils.GetAuthInfo(c)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "utils.GetAuthInfo(c)", err)
		return
	}

	request := &models.DtoDetectionActRequest{}
	if err = c.ShouldBindJSON(request); err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "request body is not valid", err)
		return
	}

	resp, err := s.analyticsCL.Actions.PatchAPIV1DetectionsAct(&actions.PatchAPIV1DetectionsActParams{
		Request: request,
	}, authInfo)
	if err != nil {
		if conflictErr, ok := err.(ResponseErrorInterface); ok {
			c.JSON(conflictErr.Code(), conflictErr.GetPayload())
			return
		}

		s.ErrorResponse(c, http.StatusBadRequest, "Actions.GetAPIV1DetectionsAct", err)
		return
	}

	c.JSON(resp.Code(), resp.GetPayload())
}

//---------------------detections---------------------

// getDetections -
// @Summary Получение списка выявлений
// @Description Возвращает список выявлений за указанный временной период с пагинацией и фильтрацией
// @Tags detections
// @Accept json
// @Produce json
// @Param from query string false "Начало временного диапазона (формат: now-10m, 2023-12-01T10:00:00Z)" default(now-10m)
// @Param to query string false "Конец временного диапазона (формат: now, 2023-12-01T11:00:00Z)" default(now)
// @Param hostname query string false "Фильтр по имени хоста"
// @Param status query string false "Фильтр по статусу выявления" Enums(Рекомендуется_блокировка, Требуется_проверка, Заблокирован)
// @Param category query string false "Фильтр по категории" Enums(Агрессия, расизм, терроризм, Ботнеты, Веб-почта, Досуг и развлечения, Интернет-магазины, Компьютерные игры, Криптомайнинг, Наркотики, Порнография и секс, Прокси и анонимайзеры, Реестр запрещенных сайтов, Сайты для взрослых, Сайты распространяющие вирусы, Социальные сети, Торренты и Р2Р-сети, Файловые архивы, Фильмы и видео онлайн, Фишинг, Чаты и мессенджеры, Дополнительно, Криптоджекинг, Реклама, Онлайн-игры, Игровые платформы, Вредоносное ПО, Азартные игры, Депресивный контент и суицид, Алкоголь, табак)
// @Param action query string false "Действие пользователя" Enums(Разрешено, Заблокировано, Не решено)
// @Param page query int false "Номер страницы" default(1) minimum(1)
// @Param limit query int false "Количество записей на странице" default(10) minimum(1) maximum(100)
// @Param search query string false "Поиск по URL или IP адресу домена"
// @Param order_by query string false "Поле для сортировки" default(categorized_at) Enums(domain, request_count, categorized_at)
// @Param order_dir query string false "Направление сортировки (asc/desc)" default(desc) Enums(asc, desc)
// @Security BearerAuth
// @Success 200 {object} models.DtoGetDetectionsResponse "Успешный ответ"
// @Failure 400 {object} models.DtoErrorResponse "Неверный формат параметров"
// @Failure 403 {object}  models.DtoErrorResponse "Недостаточно прав для доступа к выявлениям"
// @Failure 500 {object} models.DtoErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/analytics/detections [get]
func (s *Server) getDetections(c *gin.Context) {
	// Создаем authInfoWriter для передачи токена
	authInfo, err := utils.GetAuthInfo(c)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "utils.GetAuthInfo(c)", err)
		return
	}

	//Парсим входные данные
	from := c.DefaultQuery("from", "now-10m") // по умолчанию выдает последние 10 минут
	to := c.DefaultQuery("to", "now")
	hostname := c.Query("hostname")
	status := c.Query("_status")
	category := c.Query("category")
	action := c.Query("action")
	page, err := strconv.ParseInt(c.Query("page"), 10, 64)
	if err != nil {
		if conflictErr, ok := err.(ResponseErrorInterface); ok {
			c.JSON(conflictErr.Code(), conflictErr.GetPayload())
			return
		}

		s.ErrorResponse(c, http.StatusBadRequest, "Parse int page", err)
		return
	}
	limit, err := strconv.ParseInt(c.Query("limit"), 10, 64)
	if err != nil {
		if conflictErr, ok := err.(ResponseErrorInterface); ok {
			c.JSON(conflictErr.Code(), conflictErr.GetPayload())
			return
		}

		s.ErrorResponse(c, http.StatusBadRequest, "Parse int limit", err)
		return
	}
	search := c.Query("search")
	orderBy := c.Query("order_by")
	orderDir := c.Query("order_dir")

	resp, err := s.analyticsCL.Detections.GetAPIV1Detections(&detections.GetAPIV1DetectionsParams{
		From:     &from,
		To:       &to,
		Hostname: &hostname,
		Category: &category,
		Status:   &status,
		Action:   &action,
		Page:     &page,
		Limit:    &limit,
		Search:   &search,
		OrderBy:  &orderBy,
		OrderDir: &orderDir,
	}, authInfo)
	if err != nil {
		if conflictErr, ok := err.(ResponseErrorInterface); ok {
			c.JSON(conflictErr.Code(), conflictErr.GetPayload())
			return
		}

		s.ErrorResponse(c, http.StatusBadRequest, "Detections.GetAPIV1Detections", err)
		return
	}

	c.JSON(resp.Code(), resp.GetPayload())
}

// getDetectionsV2 -
// @Summary Получение списка выявлений (версия 2)
// @Description Возвращает расширенный список выявлений с детальной статистикой по категориям за указанный временной период с пагинацией и фильтрацией
// @Tags detections
// @Accept json
// @Produce json
// @Param from query string false "Начало временного диапазона (формат: now-10m, 2023-12-01T10:00:00Z)" default(now-10m)
// @Param to query string false "Конец временного диапазона (формат: now, 2023-12-01T11:00:00Z)" default(now)
// @Param hostname query string false "Фильтр по имени хоста"
// @Param status query string false "Фильтр по статусу выявления" Enums(Рекомендуется_блокировка, Требуется_проверка, Заблокирован)
// @Param category query string false "Фильтр по категории" Enums(Агрессия, расизм, терроризм, Ботнеты, Веб-почта, Досуг и развлечения, Интернет-магазины, Компьютерные игры, Криптомайнинг, Наркотики, Порнография и секс, Прокси и анонимайзеры, Реестр запрещенных сайтов, Сайты для взрослых, Сайты распространяющие вирусы, Социальные сети, Торренты и Р2Р-сети, Файловые архивы, Фильмы и видео онлайн, Фишинг, Чаты и мессенджеры, Дополнительно, Криптоджекинг, Реклама, Онлайн-игры, Игровые платформы, Вредоносное ПО, Азартные игры, Депресивный контент и суицид, Алкоголь, табак)
// @Param action query string false "Действие пользователя" Enums(Разрешено, Заблокировано, Не решено)
// @Param page query int false "Номер страницы" default(1) minimum(1)
// @Param limit query int false "Количество записей на странице" default(10) minimum(1) maximum(100)
// @Param search query string false "Поиск по URL или IP адресу домена"
// @Param order_by query string false "Поле для сортировки" default(categorized_at) Enums(domain, request_count, categorized_at)
// @Param order_dir query string false "Направление сортировки (asc/desc)" default(desc) Enums(asc, desc)
// @Security BearerAuth
// @Success 200 {object} models.DtoGetDetectionsResponseV2 "Успешный ответ"
// @Failure 400 {object} models.DtoErrorResponse "Неверный формат параметров"
// @Failure 403 {object} models.DtoErrorResponse "Недостаточно прав для доступа к выявлениям"
// @Failure 500 {object} models.DtoErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/analytics/detections_v2 [get]
func (s *Server) getDetectionsV2(c *gin.Context) {
	// Создаем authInfoWriter для передачи токена
	authInfo, err := utils.GetAuthInfo(c)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "utils.GetAuthInfo(c)", err)
		return
	}

	// Парсим входные данные
	from := c.DefaultQuery("from", "now-10m") // по умолчанию выдает последние 10 минут
	to := c.DefaultQuery("to", "now")
	hostname := c.Query("hostname")
	status := c.Query("status")
	category := c.Query("category")
	action := c.Query("action")

	page, err := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "Parse int page", err)
		return
	}

	limit, err := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "Parse int limit", err)
		return
	}

	// Валидация limit
	if limit < 1 || limit > 100 {
		s.ErrorResponse(c, http.StatusBadRequest, "limit must be between 1 and 100", fmt.Errorf("invalid limit: %d", limit))
		return
	}

	// Валидация page
	if page < 1 {
		s.ErrorResponse(c, http.StatusBadRequest, "page must be greater than 0", fmt.Errorf("invalid page: %d", page))
		return
	}

	search := c.Query("search")
	orderBy := c.DefaultQuery("order_by", "categorized_at")
	orderDir := c.DefaultQuery("order_dir", "desc")

	// Валидация orderDir
	if orderDir != "asc" && orderDir != "desc" {
		s.ErrorResponse(c, http.StatusBadRequest, "order_dir must be asc or desc", fmt.Errorf("invalid order_dir: %s", orderDir))
		return
	}

	// Вызов API клиента (версия 2)
	resp, err := s.analyticsCL.Detections.GetAPIV2Detections(&detections.GetAPIV2DetectionsParams{
		From:     &from,
		To:       &to,
		Hostname: &hostname,
		Status:   &status,
		Category: &category,
		Action:   &action,
		Page:     &page,
		Limit:    &limit,
		Search:   &search,
		OrderBy:  &orderBy,
		OrderDir: &orderDir,
	}, authInfo)
	if err != nil {
		if conflictErr, ok := err.(ResponseErrorInterface); ok {
			c.JSON(conflictErr.Code(), conflictErr.GetPayload())
			return
		}

		s.ErrorResponse(c, http.StatusBadRequest, "Detections.GetAPIV2Detections", err)
		return
	}

	c.JSON(resp.Code(), resp.GetPayload())
}

// getDetectionsStat -
// @Summary Получение статистики по выявлениям
// @Description Возвращает статистику выявлений за указанный период с фильтрацией
// @Tags detections
// @Accept json
// @Produce json
// @Param from query string false "Начало временного диапазона (формат: now-10m, now-1h, 2024-01-01T00:00:00Z)" default(now-10m)
// @Param to query string false "Конец временного диапазона (формат: now, 2024-01-01T00:00:00Z)" default(now)
// @Param hostname query string false "Фильтр по имени хоста"
// @Param category query string false "Фильтр по категории" Enums(Агрессия, расизм, терроризм, Ботнеты, Веб-почта, Досуг и развлечения, Интернет-магазины, Компьютерные игры, Криптомайнинг, Наркотики, Порнография и секс, Прокси и анонимайзеры, Реестр запрещенных сайтов, Сайты для взрослых, Сайты распространяющие вирусы, Социальные сети, Торренты и Р2Р-сети, Файловые архивы, Фильмы и видео онлайн, Фишинг, Чаты и мессенджеры, Дополнительно, Криптоджекинг, Реклама, Онлайн-игры, Игровые платформы, Вредоносное ПО, Азартные игры, Депресивный контент и суицид, Алкоголь, табак)
// @Security BearerAuth
// @Success 200 {object} models.DtoDetectionStat "Статистика детекций"
// @Failure 400 {object} models.DtoErrorResponse "Неверный формат временного диапазона"
// @Failure 403 {object} models.DtoErrorResponse "Недостаточно прав для доступа к статистике выявлений"
// @Failure 500 {object} models.DtoErrorResponse "Ошибка при получении статистики выявлений"
// @Router /api/v1/analytics/detections/stat [get]
func (s *Server) getDetectionsStat(c *gin.Context) {
	// Создаем authInfoWriter для передачи токена
	authInfo, err := utils.GetAuthInfo(c)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "utils.GetAuthInfo(c)", err)
		return
	}

	//Парсим входные данные
	from := c.DefaultQuery("from", "now-10m") // по умолчанию выдает последние 10 минут
	to := c.DefaultQuery("to", "now")
	hostname := c.Query("hostname")
	category := c.Query("category")

	resp, err := s.analyticsCL.Detections.GetAPIV1DetectionsStat(&detections.GetAPIV1DetectionsStatParams{
		From:     &from,
		To:       &to,
		Hostname: &hostname,
		Category: &category,
	}, authInfo)
	if err != nil {
		if conflictErr, ok := err.(ResponseErrorInterface); ok {
			c.JSON(conflictErr.Code(), conflictErr.GetPayload())
			return
		}

		s.ErrorResponse(c, http.StatusBadRequest, "Detections.GetAPIV1DetectionsStat", err)
		return
	}

	c.JSON(resp.Code(), resp.GetPayload())
}

//---------------------reports---------------------

// getV1Reports -
// @Summary Создание отчета
// @Description Генерирует полный отчет по активности за указанный временной период, включая аналитику по устройствам, категориям и аномалиям
// @Tags reports
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param from query string false "Начало временного диапазона (формат: now-24h, 2023-12-01T10:00:00Z)" default(now-24h)
// @Param to query string false "Конец временного диапазона (формат: now, 2023-12-01T12:00:00Z)" default(now)
// @Success 200 {object} models.ModelsReport "Полный отчет по активности"
// @Failure 400 {object} models.DtoErrorResponse "Неверный формат параметров"
// @Failure 403 {object} models.DtoErrorResponse "Недостаточно прав для доступа к статистике выявлений"
// @Failure 500 {object} models.DtoErrorResponse "Внутренняя ошибка сервера при генерации отчета"
// @Router /api/v1/analytics/reports [get]
func (s *Server) getV1Reports(c *gin.Context) {
	// Создаем authInfoWriter для передачи токена
	authInfo, err := utils.GetAuthInfo(c)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "utils.GetAuthInfo(c)", err)
		return
	}

	//Парсим входные данные
	from := c.DefaultQuery("from", "now-24h")
	to := c.DefaultQuery("to", "now")

	resp, err := s.analyticsCL.Reports.GetAPIV1Reports(&reports.GetAPIV1ReportsParams{
		From: &from,
		To:   &to,
	}, authInfo)
	if err != nil {
		if conflictErr, ok := err.(ResponseErrorInterface); ok {
			c.JSON(conflictErr.Code(), conflictErr.GetPayload())
			return
		}

		s.ErrorResponse(c, http.StatusBadRequest, "Reports.GetAPIV1Reports", err)
		return
	}

	c.JSON(resp.Code(), resp.GetPayload())
}

// @Summary Создание отчета по конкретному устройству
// @Description Генерирует детализированный отчет по конкретному сетевому устройству за указанный временной период, включая статистику трафика, аномалии и категории запросов
// @Tags reports
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param hostname path string true "Имя сетевого устройства (хоста)"
// @Param from query string false "Начало временного диапазона (формат: now-24h, 2023-12-01T10:00:00Z)" default(now-24h)
// @Param to query string false "Конец временного диапазона (формат: now, 2023-12-01T12:00:00Z)" default(now)
// @Success 200 {object} models.ModelsReportForDevice "Детализированный отчет по устройству"
// @Failure 400 {object} models.DtoErrorResponse "Неверный формат параметров или устройство не найдено"
// @Failure 404 {object} models.DtoErrorResponse "Устройство не найдено в базе данных"
// @Failure 500 {object} models.DtoErrorResponse "Внутренняя ошибка сервера при генерации отчета"
// @Router /api/v1/analytics/reports/{hostname} [get]
func (s *Server) getV1ReportsHostname(c *gin.Context) {
	// Создаем authInfoWriter для передачи токена
	authInfo, err := utils.GetAuthInfo(c)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "utils.GetAuthInfo(c)", err)
		return
	}

	//Парсим входные данные
	hostname := c.Param("hostname")
	from := c.DefaultQuery("from", "now-24h")
	to := c.DefaultQuery("to", "now")

	resp, err := s.analyticsCL.Reports.GetAPIV1ReportsHostname(&reports.GetAPIV1ReportsHostnameParams{
		Hostname: hostname,
		From:     &from,
		To:       &to,
	}, authInfo)
	if err != nil {
		if conflictErr, ok := err.(ResponseErrorInterface); ok {
			c.JSON(conflictErr.Code(), conflictErr.GetPayload())
			return
		}

		s.ErrorResponse(c, http.StatusBadRequest, "Reports.GetAPIV1ReportsHostname", err)
		return
	}

	c.JSON(resp.Code(), resp.GetPayload())
}

//---------------------sessions---------------------

// getSessions -
// @Summary Получение списка сессий
// @Description Возвращает список сессий с возможностью фильтрации, поиска, сортировки и пагинации
// @Tags sessions
// @Accept json
// @Produce json
// @Param from query string false "Начало временного диапазона (формат: now-10m, 2023-12-01T10:00:00Z). По умолчанию: now-10m" default(now-10m)
// @Param to query string false "Конец временного диапазона (формат: now, 2023-12-01T12:00:00Z). По умолчанию: now" default(now)
// @Param hostname query string false "Фильтр по имени хоста"
// @Param category query string false "Фильтр по категории" Enums(Неизвестный класс, Агрессия, расизм, терроризм, Ботнеты, Веб-почта, Досуг и развлечения, Интернет магазины, Компьютерные игры, Криптомайнинг, Наркотики, Порнография и секс, Прокси и анонимайзеры, Реестр запрещенных сайтов, Сайты для взрослых, Сайты распространяющие вирусы, Социальные сети, Торренты и Р2Р-сети, Файловые архивы, Фильмы и видео онлайн, Фишинг, Чаты и мессенджеры, Криптоджекинг, Реклама, Онлайн-игры, Игровые платформы, Вредоносное ПО, Азартные игры, Депрессивный контент, Алкоголь и табак, Положительная категория)
// @Param type query string false "Фильтр по типу сессии" Enums(Разрешен, Запрещен, VPN)
// @Param status query string false "Фильтр по статусу" Enums(Разрешен, Запрещен, Ожидает, Аномалия)
// @Param search query string false "Поиск по URL, IP адресу пользователя или IP адресу домена"
// @Param page query int false "Номер страницы" default(1) minimum(1)
// @Param count query int false "Количество записей на странице" default(10) minimum(1) maximum(100)
// @Param order_by query string false "Поле для сортировки" default(datetime_utc) Enums(id, datetime_utc, type, status, url, proto, hostname, src_ip, src_country, username, dst_ip, dst_port, dst_country, category)
// @Param order_dir query string false "Направление сортировки (asc/desc)" default(desc) Enums(asc, desc)
// @Security BearerAuth
// @Success 200 {object} models.DtoGetSessionsResponse "Успешный ответ"
// @Failure 400 {object} models.DtoErrorResponse "Неверный формат параметров"
// @Failure 403 {object} models.DtoErrorResponse "Недостаточно прав для доступа к статистике выявлений"
// @Failure 500 {object} models.DtoErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/analytics/sessions [get]
func (s *Server) getSessions(c *gin.Context) {
	// Создаем authInfoWriter для передачи токена
	authInfo, err := utils.GetAuthInfo(c)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "utils.GetAuthInfo(c)", err)
		return
	}

	//Парсим входные данные
	from := c.DefaultQuery("from", "now-10m") // по умолчанию выдает последние 10 минут
	to := c.DefaultQuery("to", "now")
	hostname := c.Query("hostname")
	category := c.Query("category")
	status := c.Query("status")
	typeStr := c.Query("type")
	search := c.Query("search")
	page, err := strconv.ParseInt(c.Query("page"), 10, 64)
	if err != nil {
		if conflictErr, ok := err.(ResponseErrorInterface); ok {
			c.JSON(conflictErr.Code(), conflictErr.GetPayload())
			return
		}

		s.ErrorResponse(c, http.StatusBadRequest, "Parse int page", err)
		return
	}
	// TODO: count -> limit
	limit, err := strconv.ParseInt(c.Query("count"), 10, 64)
	if err != nil {
		if conflictErr, ok := err.(ResponseErrorInterface); ok {
			c.JSON(conflictErr.Code(), conflictErr.GetPayload())
			return
		}

		s.ErrorResponse(c, http.StatusBadRequest, "Parse int count", err)
		return
	}
	orderBy := c.Query("order_by")
	orderDir := c.Query("order_dir")

	resp, err := s.analyticsCL.Sessions.GetAPIV1Sessions(&sessions.GetAPIV1SessionsParams{
		From:     &from,
		To:       &to,
		Hostname: &hostname,
		Category: &category,
		Status:   &status,
		Type:     &typeStr,
		Search:   &search,
		Page:     &page,
		Count:    &limit,
		OrderBy:  &orderBy,
		OrderDir: &orderDir,
	}, authInfo)
	if err != nil {
		if conflictErr, ok := err.(ResponseErrorInterface); ok {
			c.JSON(conflictErr.Code(), conflictErr.GetPayload())
			return
		}

		s.ErrorResponse(c, http.StatusBadRequest, "Sessions.GetAPIV1Sessions", err)
		return
	}

	c.JSON(resp.Code(), resp.GetPayload())
}
