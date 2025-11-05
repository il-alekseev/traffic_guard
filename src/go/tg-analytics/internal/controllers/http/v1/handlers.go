package v1

import (
	"fmt"
	"math"
	"net/http"
	"tg-an/internal/controllers/http/v1/dto"
	"tg-an/internal/controllers/http/v1/validation"
	"tg-an/internal/models"
	"tg-an/pkg/trparser"
	"time"

	"github.com/gin-gonic/gin"
)

// @Summary Получение версии сервиса
// @Description Возвращает информацию о версии
// @Tags utils
// @Produce json
// @Success 200 {object} dto.SuccessResponse
// @Router /api/v1/version [get]
func (s *Server) Version(c *gin.Context) {
	response := dto.SuccessResponse{
		Message: s.devVersion,
	}
	c.JSON(http.StatusOK, response)
	s.l.Debug("version")
}

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
// @Param search query string false "Поиск по URL или имени пользователя"
// @Param page query int false "Номер страницы" default(1) minimum(1)
// @Param limit query int false "Количество записей на странице" default(10) minimum(1) maximum(100)
// @Param order_by query string false "Поле для сортировки" default(datetime_utc) Enums(id, datetime_utc, type, status, url, proto, host_name, src_ip, src_country, username, dst_ip, dst_port, dst_country, category)
// @Param order_dir query string false "Направление сортировки (asc/desc)" default(desc) Enums(asc, desc)
// @Success 200 {object} dto.GetSessionsResponse "Успешный ответ"
// @Failure 400 {object} dto.ErrorResponse "Неверный формат параметров"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/sessions [get]
func (s *Server) GetSessions(c *gin.Context) {
	// Получение данных запроса
	var req validation.GetSessionsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Invalid query parameters: %v", err),
		})
		return
	}
	// Нормализация и валидация
	if err := req.ValidateAndNormalize(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	// Парсим временной диапазон
	parser := &trparser.TimeRangeParser{}
	now := time.Now()
	timeRange, err := parser.Parse(fmt.Sprintf("from=%s&to=%s", req.From, req.To), now)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный формат временного диапазона",
		})
		return
	}
	// Преобразуем в доменные модели
	filter := models.SessionFilter{
		HostName: req.HostName,
		Category: req.Category,
		Type:     req.Type,
	}
	pagination := models.Pagination{
		Page:  req.Page,
		Limit: req.Limit,
	}
	sorting := models.Sorting{
		OrderBy:  req.OrderBy,
		OrderDir: req.OrderDir,
	}
	// Получаем данные из usecase
	sessions, total, err := s.u.GetSessions(c, timeRange, filter, req.Search, pagination, sorting)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ошибка при получении сессий",
		})
		return
	}
	// Формируем ответ
	response := dto.ListResponse{
		Data: sessions,
		Meta: dto.PaginationMeta{
			Page:  req.Page,
			Limit: req.Limit,
			Total: total,
			Pages: int(math.Ceil(float64(total) / float64(req.Limit))),
		},
	}
	c.JSON(http.StatusOK, response)
}

// @Summary Получение списка самых запрашиваемых категорий
// @Description Возвращает наиболее часто встречаемые категории в сессиях с возможностью фильтрации
// @Tags dashboards
// @Accept json
// @Produce json
// @Param from query string false "Начало временного диапазона" default(now-24h)
// @Param to query string false "Конец временного диапазона" default(now)
// @Param hostname query string false "Фильтр по имени хоста"
// @Param type query string false "Фильтр по типу сессии" Enums(Разрешен, Заблокирован, VPN)
// @Param count query int false "Количество возвращаемых категорий" default(5) minimum(1) maximum(50)
// @Success 200 {object} []dto.Category
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/dashboards/top-categories [get]
func (s *Server) GetTopCategories(c *gin.Context) {
	// Валидация запроса
	var req validation.GetTopCategoriesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Invalid query parameters: %v", err),
		})
		return
	}
	// Нормализация и валидация
	if err := req.ValidateAndNormalize(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	// Парсим временной диапазон
	parser := &trparser.TimeRangeParser{}
	now := time.Now()
	timeRange, err := parser.Parse(fmt.Sprintf("from=%s&to=%s", req.From, req.To), now)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный формат временного диапазона",
		})
		return
	}
	// Преобразуем в доменные модели
	filter := models.CategoryFilter{
		HostName: req.HostName,
		Type:     req.Type,
	}
	// Получаем данные из usecase
	categories, err := s.u.GetTopCategories(c.Request.Context(), timeRange, filter, req.Count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка при получении топ категорий",
		})
		return
	}
	c.JSON(http.StatusOK, categories)
}

// @Summary Получение списка выявлений
// @Description Возвращает список выявлений за указанный временной период с пагинацией и фильтрацией
// @Tags detections
// @Accept json
// @Produce json
// @Param from query string false "Начало временного диапазона (формат: now-10m, 2023-12-01T10:00:00Z)" default(now-10m)
// @Param to query string false "Конец временного диапазона (формат: now, 2023-12-01T11:00:00Z)" default(now)
// @Param hostname query string false "Фильтр по имени хоста"
// @Param category query string false "Фильтр по категории" Enums(Агрессия, расизм, терроризм, Ботнеты, Веб-почта, Досуг и развлечения, Интернет-магазины, Компьютерные игры, Криптомайнинг, Наркотики, Порнография и секс, Прокси и анонимайзеры, Реестр запрещенных сайтов, Сайты для взрослых, Сайты распространяющие вирусы, Социальные сети, Торренты и Р2Р-сети, Файловые архивы, Фильмы и видео онлайн, Фишинг, Чаты и мессенджеры, Дополнительно, Криптоджекинг, Реклама, Онлайн-игры, Игровые платформы, Вредоносное ПО, Азартные игры, Депресивный контент и суицид, Алкоголь, табак)
// @Param action query string false "Действие пользователя" Enums(Разрешено, Заблокировано, Не решено)
// @Param page query int false "Номер страницы" default(1) minimum(1)
// @Param limit query int false "Количество записей на странице" default(10) minimum(1) maximum(100)
// @Success 200 {object} dto.GetDetectionsResponse "Успешный ответ"
// @Failure 400 {object} dto.ErrorResponse "Неверный формат параметров"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/detections [get]
func (s *Server) GetDetections(c *gin.Context) {
	// Валидация запроса
	var req validation.GetTopDetectionsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Invalid query parameters: %v", err),
		})
		return
	}

	// Нормализация и валидация
	if err := req.ValidateAndNormalize(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Парсим временной диапазон
	parser := &trparser.TimeRangeParser{}
	now := time.Now()
	timeRange, err := parser.Parse(fmt.Sprintf("from=%s&to=%s", req.From, req.To), now)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный формат временного диапазона",
		})
		return
	}

	// Преобразуем в доменные модели
	filter := models.DetectionFilter{
		HostName:    req.HostName,
		TopCategory: req.Category,
	}

	pagination := models.Pagination{
		Page:  req.Page,
		Limit: req.Limit,
	}
	// Получаем данные из usecase
	detections, total, err := s.u.GetTopDetections(c, timeRange, filter, req.Action, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ошибка при получении списка выявлений",
		})
		return
	}

	// Формируем ответ
	response := dto.GetDetectionsResponse{
		Data: detections,
		Meta: dto.PaginationMeta{
			Page:  req.Page,
			Limit: req.Limit,
			Total: total,
			Pages: int(math.Ceil(float64(total) / float64(req.Limit))),
		},
	}
	c.JSON(http.StatusOK, response)
}

// @Summary Получение статистики по выявлениям
// @Description Возвращает статистику выявлений за указанный период с фильтрацией
// @Tags detections
// @Accept json
// @Produce json
// @Param from query string false "Начало временного диапазона (формат: now-10m, now-1h, 2024-01-01T00:00:00Z)" default(now-10m)
// @Param to query string false "Конец временного диапазона (формат: now, 2024-01-01T00:00:00Z)" default(now)
// @Param hostname query string false "Фильтр по имени хоста"
// @Param category query string false "Фильтр по категории" Enums(Агрессия, расизм, терроризм, Ботнеты, Веб-почта, Досуг и развлечения, Интернет-магазины, Компьютерные игры, Криптомайнинг, Наркотики, Порнография и секс, Прокси и анонимайзеры, Реестр запрещенных сайтов, Сайты для взрослых, Сайты распространяющие вирусы, Социальные сети, Торренты и Р2Р-сети, Файловые архивы, Фильмы и видео онлайн, Фишинг, Чаты и мессенджеры, Дополнительно, Криптоджекинг, Реклама, Онлайн-игры, Игровые платформы, Вредоносное ПО, Азартные игры, Депресивный контент и суицид, Алкоголь, табак)
// @Success 200 {object} dto.DetectionStat "Статистика детекций"
// @Failure 400 {object} dto.ErrorResponse "Неверный формат временного диапазона"
// @Failure 500 {object} dto.ErrorResponse "Ошибка при получении статистики выявлений"
// @Router /api/v1/detections/stat [get]
func (s *Server) GetDetectionStat(c *gin.Context) {
	var req validation.GetDetectionStatRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Invalid query parameters: %v", err),
		})
		return
	}

	// Нормализация и валидация
	if err := req.ValidateAndNormalize(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Парсим временной диапазон
	parser := &trparser.TimeRangeParser{}
	now := time.Now()
	timeRange, err := parser.Parse(fmt.Sprintf("from=%s&to=%s", req.From, req.To), now)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный формат временного диапазона",
		})
		return
	}

	// Преобразуем в доменные модели
	filter := models.DetectionFilter{
		HostName:    req.HostName,
		TopCategory: req.Category,
	}

	// Получаем данные из usecase
	stat, err := s.u.GetDetectionStat(c, timeRange, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ошибка при получении статистики выявлений",
		})
		return
	}

	c.JSON(http.StatusOK, stat)
}

// @Summary Получение статистики запросов
// @Description Получение статистики запросов за указанный период с фильтрацией по хосту и типу запросов
// @Tags dashboards
// @Accept json
// @Produce json
// @Param from query string false "Начало временного диапазона (формат: now-10m, now-1h, 2024-01-01T00:00:00Z)" default(now-10m)
// @Param to query string false "Конец временного диапазона (формат: now, 2024-01-01T00:00:00Z)" default(now)
// @Param request_type query string false "Тип запроса" Enums(allowed, blocked, before_block, pending) default(pending)
// @Param hostname query string false "Фильтр по имени хоста"
// @Param count query integer false "Количество интервалов" default(10)
// @Success 200 {object} dto.RequestStatResponse "Статистика запросов (массив чисел)"
// @Failure 400 {object} dto.ErrorResponse "Неверный формат параметров"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/dashboards/requests [get]
func (s *Server) GetRequestStat(c *gin.Context) {
	var req validation.GetRequestStatRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Invalid query parameters: %v", err),
		})
		return
	}
	// Нормализация и валидация
	if err := req.ValidateAndNormalize(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Парсим временной диапазон
	parser := &trparser.TimeRangeParser{}
	now := time.Now()
	timeRange, err := parser.Parse(fmt.Sprintf("from=%s&to=%s", req.From, req.To), now)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный формат временного диапазона",
		})
		return
	}

	// Получем данные
	stat, err := s.u.GetRequestStat(c, timeRange, req.HostName, req.RequestType, req.Count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ошибка при получении статистики запросов",
		})
		return
	}
	resp := dto.RequestStatResponse{
		Type:  fmt.Sprintf("requests_%s", req.RequestType),
		Data:  stat,
		Count: req.Count,
	}
	c.JSON(http.StatusOK, resp)
}

// @Summary Полуечение списка имен устройств
// @Description Возвращает список всех уникальных имен устройств (хостов) из системы
// @Tags common
// @Accept json
// @Produce json
// @Success 200 {array} string "Список имен устройств"
// @Failure 500 {object} map[string]string "Ошибка при получении имен устройств"
// @Router /api/v1/devices [get]
func (s *Server) GetDevices(c *gin.Context) {
	// Получем данные
	devices, err := s.u.GetDevices(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ошибка при получении имен устройств",
		})
		return
	}
	c.JSON(http.StatusOK, devices)
}

// @Summary Получение списка категорий контента
// @Description Возвращает список всех уникальных категорий контента из системы
// @Tags common
// @Accept json
// @Produce json
// @Success 200 {array} string "Список категорий контента"
// @Failure 500 {object} map[string]string "Ошибка при получении списка категорий"
// @Router /api/v1/categories [get]
func (s *Server) GetContentCategories(c *gin.Context) {
	// Получем данные
	devices, err := s.u.GetContentCategories(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ошибка при получении списка всех категорий",
		})
		return
	}
	c.JSON(http.StatusOK, devices)
}

// @Summary Получение статистики трафика
// @Description Возвращает статистику трафика за указанный временной диапазон с заданным количеством точек данных в Кб
// @Tags dashboards
// @Accept json
// @Produce json
// @Param from query string false "Начало временного диапазона в формате парсера времени (по умолчанию now-10m)" default(now-10m)
// @Param to query string false "Конец временного диапазона в формате парсера времени (по умолчанию now)" default(now)
// @Param hostname query string false "Фильтр по имени хоста"
// @Param count query integer false "Количество точек данных для возврата (по умолчанию 20)" minimum(1) default(20)
// @Success 200 {object} dto.TrafficStatResponse "Успешный ответ со статистикой трафика"
// @Failure 400 {object} dto.ErrorResponse "Неверный формат параметров запроса"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера при получении статистики"
// @Router /api/v1/dashboards/traffic [get]
func (s *Server) GetTrafficStat(c *gin.Context) {
	// Валидация запроса
	var req validation.GetTrafficStatRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Invalid query parameters: %v", err),
		})
		return
	}

	// Нормализация и валидация
	if err := req.ValidateAndNormalize(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Парсим временной диапазон
	parser := &trparser.TimeRangeParser{}
	now := time.Now()
	timeRange, err := parser.Parse(fmt.Sprintf("from=%s&to=%s", req.From, req.To), now)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный формат временного диапазона",
		})
		return
	}

	// Получаем данные из usecase
	stat, err := s.u.GetTrafficStat(c.Request.Context(), timeRange, req.HostName, req.Count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ошибка при получении статистики трафика",
		})
		return
	}

	// Формируем ответ
	resp := dto.TrafficStatResponse{
		Data:  stat,
		Count: req.Count,
	}
	c.JSON(http.StatusOK, resp)
}

// @Summary Получение списка топ нерешенных выявлений
// @Description Возвращает список наиболее частых нерешенных выявлений за указанный временной период с возможностью фильтрации
// @Tags dashboards
// @Accept json
// @Produce json
// @Param from query string false "Начало временного диапазона (формат: now-10m, 2023-12-01T10:00:00Z)" default(now-10m)
// @Param to query string false "Конец временного диапазона (формат: now, 2023-12-01T12:00:00Z)" default(now)
// @Param hostname query string false "Фильтр по имени хоста"
// @Param count query int false "Количество возвращаемых записей" default(5) minimum(1)
// @Success 200 {array} dto.UnresolvedDetection "Успешный ответ со списком нерешенных выявлений"
// @Failure 400 {object} dto.ErrorResponse "Неверный формат параметров"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/dashboards/top-unresolved_detections [get]
func (s *Server) GetTopUnresolvedDetections(c *gin.Context) {
	// Валидация запроса
	var req validation.GetTopCategoriesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Invalid query parameters: %v", err),
		})
		return
	}
	// Нормализация и валидация
	if err := req.ValidateAndNormalize(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	// Парсим временной диапазон
	parser := &trparser.TimeRangeParser{}
	now := time.Now()
	timeRange, err := parser.Parse(fmt.Sprintf("from=%s&to=%s", req.From, req.To), now)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный формат временного диапазона",
		})
		return
	}
	// Получаем данные из usecase
	UnresolvedDetections, err := s.u.GetTopUnresolvedDetections(c.Request.Context(), timeRange, req.HostName, req.Count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка при получении топ нерешенных выявлений",
		})
		return
	}
	c.JSON(http.StatusOK, UnresolvedDetections)
}

// @Summary Получение статистики по устройствам
// @Description Получение агрегированной статистики по сетевым узлам за указанный период
// @Tags dashboards
// @Produce application/json
// @Param from query string false "Начало временного диапазона (формат: now-10m, 2023-12-01T10:00:00Z)" default(now-10m)
// @Param to query string false "Конец временного диапазона (формат: now, 2023-12-01T12:00:00Z)" default(now)
// @Param count query int false "Количество точек измерений" default(20) minimum(1)
// @Success 200 {object} dto.DeviceStatResponse "Статистика по устройствам"
// @Failure 400 {object} dto.ErrorResponse "Неверный формат параметров"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/dashboards/devices [get]
func (s *Server) GetDeviceStat(c *gin.Context) {
	// Валидация запроса
	var req validation.GetDeviceStatRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Invalid query parameters: %v", err),
		})
		return
	}
	// Нормализация и валидация
	if err := req.ValidateAndNormalize(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	// Парсим временной диапазон
	parser := &trparser.TimeRangeParser{}
	now := time.Now()
	timeRange, err := parser.Parse(fmt.Sprintf("from=%s&to=%s", req.From, req.To), now)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный формат временного диапазона",
		})
		return
	}
	// Получаем данные из usecase
	response, err := s.u.GetDeviceStat(c, timeRange, req.Count)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Ошибка при получении статистики сетевых узлов",
		})
		return
	}
	c.JSON(http.StatusOK, response)
}

// @Summary Получение информации об аномалиях
// @Description Получение статистики об аномалиях за указанный период
// @Tags dashboards
// @Produce application/json
// @Param from query string false "Начало временного диапазона (формат: now-10m, 2023-12-01T10:00:00Z)" default(now-10m)
// @Param to query string false "Конец временного диапазона (формат: now, 2023-12-01T12:00:00Z)" default(now)
// @Success 200 {object} dto.GetAnomaliesResponse "Статистика обнаруженных аномалий"
// @Failure 400 {object} dto.ErrorResponse "Неверный формат параметров"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/dashboards/anomalies [get]
func (s *Server) GetAnomalies(c *gin.Context) {
	// Валидация запроса
	var req validation.GetAnomaliesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Invalid query parameters: %v", err),
		})
		return
	}
	// Нормализация и валидация
	if err := req.ValidateAndNormalize(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	// Парсим временной диапазон
	parser := &trparser.TimeRangeParser{}
	now := time.Now()
	timeRange, err := parser.Parse(fmt.Sprintf("from=%s&to=%s", req.From, req.To), now)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный формат временного диапазона",
		})
		return
	}
	// Получаем данные из usecase
	response, err := s.u.GetAnomalies(c, timeRange)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Ошибка при получении статистики сетевых узлов",
		})
		return
	}
	c.JSON(http.StatusOK, response)
}

// @Summary Выполнение действия над доменом
// @Description Устанавливает действие (разрешить/заблокировать) для указанного домена
// @Tags actions
// @Accept json
// @Produce json
// @Param action query string true "Тип действия" Enums(allow, deny) default(allow)
// @Param path query string true "Путь домена"
// @Success 200 {object} dto.SuccessResponse "Действие успешно применено к домену"
// @Failure 400 {object} dto.ErrorResponse "Неверные параметры запроса"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/dashbords/act [get]
func (s *Server) Act(c *gin.Context) {
	// Валидация запроса
	var req validation.ActRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Invalid query parameters: %v", err),
		})
		return
	}
	// Нормализация и валидация
	if err := req.ValidateAndNormalize(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	err := s.u.Act(c, req.Path, req.Action)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Ошибка при установке значения действия к домену",
		})
		return
	}
	resp := dto.SuccessResponse{
		Message: fmt.Sprintf("Success %s to domain %s", req.Action, req.Path),
	}
	c.JSON(http.StatusOK, resp)
}
