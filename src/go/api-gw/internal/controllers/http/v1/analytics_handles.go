package v1

import (
	"api-gateway/pkg/analytics/common"
	"api-gateway/pkg/analytics/dashboards"
	"api-gateway/pkg/analytics/detections"
	"api-gateway/pkg/analytics/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/runtime"
	"net/http"
	"strconv"
)

// getCategories -
// @Summary Получение списка категорий контента
// @Description Получение списка всех категорий контента
// @Tags dashboard
// @Produce json
// @Security BearerAuth
// @Success 200 {object} []string
// @Failure 500 {object} models.DtoErrorResponse "Ошибка при получении списка категорий"
// @Router /v1/analytics/categories [get]
func (s *Server) getCategories(c *gin.Context) {
	//Создаем authInfoWriter для передачи токена
	//authInfo, err := utils.GetAuthInfo(c)
	//if err != nil {
	//	s.ErrorResponse(c, http.StatusBadRequest, "utils.GetAuthInfo(c)", err)
	//	return
	//}

	resp, err := s.analyticsCL.Common.GetCategories(&common.GetCategoriesParams{}, func(operation *runtime.ClientOperation) {})
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
// @Success 200 {array} string "Список имен устройств"
// @Failure 500 {object} models.DtoErrorResponse "Ошибка при получении имен устройств"
// @Router /v1/analytics/devices [get]
func (s *Server) getDevices(c *gin.Context) {
	// Создаем authInfoWriter для передачи токена
	//authInfo, err := utils.GetAuthInfo(c)
	//if err != nil {
	//	s.ErrorResponse(c, http.StatusBadRequest, "utils.GetAuthInfo(c)", err)
	//	return
	//}

	resp, err := s.analyticsCL.Common.GetDevices(&common.GetDevicesParams{}, nil)
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

// getDashboardsRequests -
// @Summary Получить статистику запросов
// @Description Возвращает статистику запросов за указанный период с фильтрацией по статусу, хосту и категории
// @Tags dashboards
// @Accept json
// @Produce json
// @Param from query string false "Начало временного диапазона (формат: now-10m, now-1h, 2024-01-01T00:00:00Z)" default(now-10m)
// @Param to query string false "Конец временного диапазона (формат: now, 2024-01-01T00:00:00Z)" default(now)
// @Param status query string false "Статус запросов (allowed, blocked, prohibited, waiting)" default(prohibited)
// @Param hostname query string false "Фильтр по имени хоста"
// @Param category query string false "Фильтр по категории" Enums(Агрессия, расизм, терроризм, Ботнеты, Веб-почта, Досуг и развлечения, Интернет-магазины, Компьютерные игры, Криптомайнинг, Наркотики, Порнография и секс, Прокси и анонимайзеры, Реестр запрещенных сайтов, Сайты для взрослых, Сайты, распространяющие вирусы, Социальные сети, Торренты и Р2Р-сети, Файловые архивы, Фильмы и видео онлайн, Фишинг, Чаты и мессенджеры, Дополнительно, Криптоджекинг, Реклама, Онлайн-игры, Игровые платформы, Вредоносное ПО, Азартные игры, Депресивный контент и суицид, Алкоголь, табак)
// @Param count query integer false "Количество интервалов" default(10)
// @Success 200 {array} integer "Статистика запросов (массив чисел)"
// @Failure 400 {object} models.DtoErrorResponse "Неверный формат параметров"
// @Failure 500 {object} models.DtoErrorResponse "Внутренняя ошибка сервера"
// @Router /v1/analytics/dashboards/requests [get]
func (s *Server) getDashboardsRequests(c *gin.Context) {
	// Создаем authInfoWriter для передачи токена
	//authInfo, err := utils.GetAuthInfo(c)
	//if err != nil {
	//	s.ErrorResponse(c, http.StatusBadRequest, "utils.GetAuthInfo(c)", err)
	//	return
	//}

	//Парсим входные данные
	from := c.DefaultQuery("from", "now-10m") // по умолчанию выдает последние 10 минут
	to := c.DefaultQuery("to", "now")
	status := c.Query("status")
	hostname := c.Query("hostname")
	category := c.Query("category")
	count, err := strconv.ParseInt(c.Query("count"), 10, 64)
	if err != nil {
		if conflictErr, ok := err.(ResponseErrorInterface); ok {
			c.JSON(conflictErr.Code(), conflictErr.GetPayload())
			return
		}

		s.ErrorResponse(c, http.StatusBadRequest, "Parse int", err)
		return
	}

	resp, err := s.analyticsCL.Dashboards.GetDashboardsRequests(&dashboards.GetDashboardsRequestsParams{
		From:     &from,
		To:       &to,
		Status:   &status,
		Hostname: &hostname,
		Category: &category,
		Count:    &count,
	}, nil)
	if err != nil {
		if conflictErr, ok := err.(ResponseErrorInterface); ok {
			c.JSON(conflictErr.Code(), conflictErr.GetPayload())
			return
		}

		s.ErrorResponse(c, http.StatusBadRequest, "Dashboards.GetDashboardsRequests", err)
		return
	}

	c.JSON(resp.Code(), resp.GetPayload())
}

// getDashboardsTopCategories -
// @Summary Получить топ категорий сессий
// @Description Возвращает наиболее часто встречаемые категории в сессиях с возможностью фильтрации
// @Tags dashboards
// @Accept json
// @Produce json
// @Param from query string false "Начало временного диапазона" default(now-24h)
// @Param to query string false "Конец временного диапазона" default(now)
// @Param hostname query string false "Фильтр по имени хоста"
// @Param type query string false "Фильтр по типу сессии" Enums(Заблокирован, Запрещен, Ожидает, Разрешен)
// @Param count query int false "Количество возвращаемых категорий" default(5) minimum(1) maximum(50)
// @Success 200 {array} []models.ModelsCategoryCount
// @Failure 400 {object} models.DtoErrorResponse
// @Failure 500 {object} models.DtoErrorResponse
// @Router /v1/analytics/dashboards/top-categories [get]
func (s *Server) getDashboardsTopCategories(c *gin.Context) {
	// Создаем authInfoWriter для передачи токена
	//authInfo, err := utils.GetAuthInfo(c)
	//if err != nil {
	//	s.ErrorResponse(c, http.StatusBadRequest, "utils.GetAuthInfo(c)", err)
	//	return
	//}

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

		s.ErrorResponse(c, http.StatusBadRequest, "Parse int", err)
		return
	}

	resp, err := s.analyticsCL.Dashboards.GetDashboardsTopCategories(&dashboards.GetDashboardsTopCategoriesParams{
		From:     &from,
		To:       &to,
		Hostname: &hostname,
		Type:     &typeStr,
		Count:    &count,
	}, nil)
	if err != nil {
		if conflictErr, ok := err.(ResponseErrorInterface); ok {
			c.JSON(conflictErr.Code(), conflictErr.GetPayload())
			return
		}

		s.ErrorResponse(c, http.StatusBadRequest, "Dashboards.GetDashboardsTopCategories", err)
		return
	}

	c.JSON(resp.Code(), resp.GetPayload())
}

// getDashboardsTraffic возвращает статистику трафика за указанный временной диапазон
// @Summary Получить статистику трафика
// @Description Возвращает статистику трафика за указанный временной диапазон с заданным количеством точек данных в Кб
// @Tags dashboards
// @Accept json
// @Produce json
// @Param from query string false "Начало временного диапазона в формате парсера времени (по умолчанию now-10m)" default(now-10m)
// @Param to query string false "Конец временного диапазона в формате парсера времени (по умолчанию now)" default(now)
// @Param count query integer false "Количество точек данных для возврата (по умолчанию 20)" minimum(1) default(20)
// @Success 200 {object} models.DtoTrafficStatResponse "Успешный ответ со статистикой трафика"
// @Failure 400 {object} models.DtoErrorResponse "Неверный формат параметров запроса"
// @Failure 500 {object} models.DtoErrorResponse "Внутренняя ошибка сервера при получении статистики"
// @Router /v1/analytics/dashboards/traffic [get]
func (s *Server) getDashboardsTraffic(c *gin.Context) {
	// Создаем authInfoWriter для передачи токена
	//authInfo, err := utils.GetAuthInfo(c)
	//if err != nil {
	//	s.ErrorResponse(c, http.StatusBadRequest, "utils.GetAuthInfo(c)", err)
	//	return
	//}

	//Парсим временные метки
	from := c.DefaultQuery("from", "now-10m") // по умолчанию выдает последние 10 минут
	to := c.DefaultQuery("to", "now")
	count, err := strconv.ParseInt(c.Query("count"), 10, 64)
	if err != nil {
		if conflictErr, ok := err.(ResponseErrorInterface); ok {
			c.JSON(conflictErr.Code(), conflictErr.GetPayload())
			return
		}

		s.ErrorResponse(c, http.StatusBadRequest, "Parse int", err)
		return
	}

	resp, err := s.analyticsCL.Dashboards.GetDashboardsTraffic(&dashboards.GetDashboardsTrafficParams{
		From:  &from,
		To:    &to,
		Count: &count,
	}, nil)
	if err != nil {
		if conflictErr, ok := err.(ResponseErrorInterface); ok {
			c.JSON(conflictErr.Code(), conflictErr.GetPayload())
			return
		}

		s.ErrorResponse(c, http.StatusBadRequest, "Dashboards.GetDashboardsTraffic", err)
		return
	}

	c.JSON(resp.Code(), resp.GetPayload())
}

// getDetections -
// @Summary Получить список топ выявлений
// @Description Возвращает список наиболее частых детекций за указанный временной период с пагинацией
// @Tags detections
// @Accept json
// @Produce json
// @Param from query string false "Начало временного диапазона (формат: now-10m, 2023-12-01T10:00:00Z)" default(now-10m)
// @Param to query string false "Конец временного диапазона (формат: now, 2023-12-01T11:00:00Z)" default(now)
// @Param hostname query string false "Фильтр по имени хоста"
// @Param category query string false "Фильтр по категории" Enums(Агрессия, расизм, терроризм, Ботнеты, Веб-почта, Досуг и развлечения, Интернет-магазины, Компьютерные игры, Криптомайнинг, Наркотики, Порнография и секс, Прокси и анонимайзеры, Реестр запрещенных сайтов, Сайты для взрослых, Сайты, распространяющие вирусы, Социальные сети, Торренты и Р2Р-сети, Файловые архивы, Фильмы и видео онлайн, Фишинг, Чаты и мессенджеры, Дополнительно, Криптоджекинг, Реклама, Онлайн-игры, Игровые платформы, Вредоносное ПО, Азартные игры, Депресивный контент и суицид, Алкоголь, табак)
// @Param page query int false "Номер страницы" default(1) minimum(1)
// @Param limit query int false "Количество записей на странице" default(10) minimum(1) maximum(100)
// @Success 200 {object} models.DtoListResponse "Успешный ответ"
// @Failure 400 {object} models.DtoErrorResponse "Неверный формат параметров"
// @Failure 500 {object} models.DtoErrorResponse "Внутренняя ошибка сервера"
// @Router /v1/analytics/detections [get]
func (s *Server) getDetections(c *gin.Context) {
	// Создаем authInfoWriter для передачи токена
	//authInfo, err := utils.GetAuthInfo(c)
	//if err != nil {
	//	s.ErrorResponse(c, http.StatusBadRequest, "utils.GetAuthInfo(c)", err)
	//	return
	//}

	//Парсим входные данные
	from := c.DefaultQuery("from", "now-10m") // по умолчанию выдает последние 10 минут
	to := c.DefaultQuery("to", "now")
	hostname := c.Query("hostname")
	category := c.Query("category")
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

	resp, err := s.analyticsCL.Detections.GetDetections(&detections.GetDetectionsParams{
		From:     &from,
		To:       &to,
		Hostname: &hostname,
		Category: &category,
		Page:     &page,
		Limit:    &limit,
	}, nil)
	if err != nil {
		if conflictErr, ok := err.(ResponseErrorInterface); ok {
			c.JSON(conflictErr.Code(), conflictErr.GetPayload())
			return
		}

		s.ErrorResponse(c, http.StatusBadRequest, "Detections.GetDetections", err)
		return
	}

	c.JSON(resp.Code(), resp.GetPayload())
}

// getDetectionsStat -
// @Summary Получение статистики по детекциям
// @Description Возвращает статистику детекций по категориям (обнаружено, принято, отклонено, неразрешено) за указанный период
// @Tags detections
// @Accept json
// @Produce json
// @Param from query string false "Начало временного диапазона (формат: now-10m, now-1h, 2024-01-01T00:00:00Z)" default(now-10m)
// @Param to query string false "Конец временного диапазона (формат: now, 2024-01-01T00:00:00Z)" default(now)
// @Param hostname query string false "Фильтр по имени хоста"
// @Param category query string false "Фильтр по категории" Enums(Агрессия, расизм, терроризм, Ботнеты, Веб-почта, Досуг и развлечения, Интернет-магазины, Компьютерные игры, Криптомайнинг, Наркотики, Порнография и секс, Прокси и анонимайзеры, Реестр запрещенных сайтов, Сайты для взрослых, Сайты, распространяющие вирусы, Социальные сети, Торренты и Р2Р-сети, Файловые архивы, Фильмы и видео онлайн, Фишинг, Чаты и мессенджеры, Дополнительно, Криптоджекинг, Реклама, Онлайн-игры, Игровые платформы, Вредоносное ПО, Азартные игры, Депресивный контент и суицид, Алкоголь, табак)
// @Success 200 {object} models.DtoDetectionStat "Статистика детекций"
// @Failure 400 {object} models.DtoErrorResponse "Неверный формат временного диапазона"
// @Failure 500 {object} models.DtoErrorResponse "Ошибка при получении статистики выявлений"
// @Router /v1/analytics/detections/stat [get]
func (s *Server) getDetectionsStat(c *gin.Context) {
	// Создаем authInfoWriter для передачи токена
	//authInfo, err := utils.GetAuthInfo(c)
	//if err != nil {
	//	s.ErrorResponse(c, http.StatusBadRequest, "utils.GetAuthInfo(c)", err)
	//	return
	//}

	//Парсим входные данные
	from := c.DefaultQuery("from", "now-10m") // по умолчанию выдает последние 10 минут
	to := c.DefaultQuery("to", "now")
	hostname := c.Query("hostname")
	category := c.Query("category")

	resp, err := s.analyticsCL.Detections.GetDetectionsStat(&detections.GetDetectionsStatParams{
		From:     &from,
		To:       &to,
		Hostname: &hostname,
		Category: &category,
	}, nil)
	if err != nil {
		if conflictErr, ok := err.(ResponseErrorInterface); ok {
			c.JSON(conflictErr.Code(), conflictErr.GetPayload())
			return
		}

		s.ErrorResponse(c, http.StatusBadRequest, "Detections.GetDetectionsStat", err)
		return
	}

	c.JSON(resp.Code(), resp.GetPayload())
}

// getSessions -
// @Summary Получить список сессий
// @Description Возвращает список сессий с возможностью фильтрации, поиска, сортировки и пагинации
// @Tags sessions
// @Accept json
// @Produce json
// @Param from query string false "Начало временного диапазона (формат: now-10m, 2023-12-01T10:00:00Z). По умолчанию: now-10m" default(now-10m)
// @Param to query string false "Конец временного диапазона (формат: now, 2023-12-01T12:00:00Z). По умолчанию: now" default(now)
// @Param hostname query string false "Фильтр по имени хоста"
// @Param category query string false "Фильтр по категории" Enums(Агрессия, расизм, терроризм, Ботнеты, Веб-почта, Досуг и развлечения, Интернет-магазины, Компьютерные игры, Криптомайнинг, Наркотики, Порнография и секс, Прокси и анонимайзеры, Реестр запрещенных сайтов, Сайты для взрослых, Сайты, распространяющие вирусы, Социальные сети, Торренты и Р2Р-сети, Файловые архивы, Фильмы и видео онлайн, Фишинг, Чаты и мессенджеры, Дополнительно, Криптоджекинг, Реклама, Онлайн-игры, Игровые платформы, Вредоносное ПО, Азартные игры, Депресивный контент и суицид, Алкоголь, табак)
// @Param type query string false "Фильтр по типу сессии" Enums(Заблокирован, Запрещен, Ожидает, Разрешен)
// @Param search query string false "Поиск по частичному совпадению"
// @Param page query int false "Номер страницы" default(1) minimum(1)
// @Param limit query int false "Количество записей на странице" default(10) minimum(1) maximum(100)
// @Param order_by query string false "Поле для сортировки" default(id) Enums(id, datetime_utc, type, status, url, proto, host_name, src_ip, src_port, src_country, username, dst_ip, dst_port, dst_country, category)
// @Param order_dir query string false "Направление сортировки (asc/desc)" default(desc) Enums(asc, desc)
// @Success 200 {object} models.DtoListResponse "Успешный ответ"
// @Failure 400 {object} models.DtoErrorResponse "Неверный формат параметров"
// @Failure 500 {object} models.DtoErrorResponse "Внутренняя ошибка сервера"
// @Router /v1/analytics/sessions [get]
func (s *Server) getSessions(c *gin.Context) {
	// Создаем authInfoWriter для передачи токена
	//authInfo, err := utils.GetAuthInfo(c)
	//if err != nil {
	//	s.ErrorResponse(c, http.StatusBadRequest, "utils.GetAuthInfo(c)", err)
	//	return
	//}

	//Парсим входные данные
	from := c.DefaultQuery("from", "now-10m") // по умолчанию выдает последние 10 минут
	to := c.DefaultQuery("to", "now")
	hostname := c.Query("hostname")
	category := c.Query("category")
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
	limit, err := strconv.ParseInt(c.Query("limit"), 10, 64)
	if err != nil {
		if conflictErr, ok := err.(ResponseErrorInterface); ok {
			c.JSON(conflictErr.Code(), conflictErr.GetPayload())
			return
		}

		s.ErrorResponse(c, http.StatusBadRequest, "Parse int limit", err)
		return
	}
	orderBy := c.Query("order_by")
	orderDir := c.Query("order_dir")

	resp, err := s.analyticsCL.Sessions.GetSessions(&sessions.GetSessionsParams{
		From:     &from,
		To:       &to,
		Hostname: &hostname,
		Category: &category,
		Type:     &typeStr,
		Search:   &search,
		Page:     &page,
		Limit:    &limit,
		OrderBy:  &orderBy,
		OrderDir: &orderDir,
	}, nil)
	if err != nil {
		if conflictErr, ok := err.(ResponseErrorInterface); ok {
			c.JSON(conflictErr.Code(), conflictErr.GetPayload())
			return
		}

		s.ErrorResponse(c, http.StatusBadRequest, "Detections.GetDetections", err)
		return
	}

	c.JSON(resp.Code(), resp.GetPayload())
}
