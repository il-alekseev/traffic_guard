package v1

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"tg-an/internal/controllers/http/v1/dto"
	"tg-an/internal/models"
	"tg-an/internal/pkg/status"
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
// @Success 200 {object} dto.ListResponse "Успешный ответ"
// @Failure 400 {object} dto.ErrorResponse "Неверный формат параметров"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/sessions [get]
func (s *Server) GetSessions(c *gin.Context) {
	// Парсим временные метки
	from := c.DefaultQuery("from", "now-10m") // по умолчанию выдает последние 10 минут
	to := c.DefaultQuery("to", "now")

	parser := &trparser.TimeRangeParser{}
	now := time.Now()
	timeRange, err := parser.Parse(fmt.Sprintf("from=%s&to=%s", from, to), now)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный формат временного диапазона",
		})
		return
	}

	// Парсим по фильтрам
	filter := models.SessionFilter{
		HostName: c.DefaultQuery("hostname", ""),
		Category: c.DefaultQuery("category", ""),
		Type:     c.DefaultQuery("type", ""),
	}

	// Парсим параметры для поиска по части названия
	search := c.DefaultQuery("search", "")

	// Парсим параметры для пагинации
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный формат page",
		})
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный формат limit",
		})
		return
	}
	pagination := models.Pagination{
		Page:  page,
		Limit: limit,
	}

	// Парсим параметры для сортировки
	// TODO: добавить валидацию по названию колонки
	orderBy := c.DefaultQuery("order_by", "id")
	orderDir := c.DefaultQuery("order_dir", "desc")
	sorting := models.Sorting{
		OrderBy:  orderBy,
		OrderDir: orderDir,
	}

	sessions, total, err := s.u.GetSessions(c, timeRange, filter, search, pagination, sorting)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ошибка при получении сессий",
		})
		return
	}

	response := dto.ListResponse{
		Data: sessions,
		Meta: dto.PaginationMeta{
			Page:  page,
			Limit: limit,
			Total: total,
			Pages: int(math.Ceil(float64(total) / float64(limit))),
		},
	}
	c.JSON(http.StatusOK, response)
}

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
// @Success 200 {array} models.CategoryCount
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/dashboards/top-categories [get]
func (s *Server) GetTopCategories(c *gin.Context) {
	// Парсим временные метки
	from := c.DefaultQuery("from", "now-10m")
	to := c.DefaultQuery("to", "now")

	parser := &trparser.TimeRangeParser{}
	now := time.Now()
	timeRange, err := parser.Parse(fmt.Sprintf("from=%s&to=%s", from, to), now)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный формат временного диапазона",
		})
		return
	}

	// Парсим фильтры
	filter := models.CategoryFilter{
		HostName: c.DefaultQuery("hostname", ""),
		Type:     c.DefaultQuery("type", ""),
	}

	// Парсим количество категорий
	count, err := strconv.Atoi(c.DefaultQuery("count", "5"))
	if err != nil || count <= 0 {
		count = 5
	}
	if count > 50 {
		count = 50 // ограничение для защиты от больших запросов
	}

	categories, err := s.u.GetTopCategories(c.Request.Context(), timeRange, filter, count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка при получении топ категорий",
		})
		return
	}

	c.JSON(http.StatusOK, categories)
}

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
// @Success 200 {object} dto.ListResponse "Успешный ответ"
// @Failure 400 {object} object "Неверный формат параметров"
// @Failure 500 {object} object "Внутренняя ошибка сервера"
// @Router /api/v1/detections [get]
func (s *Server) GetTopDetections(c *gin.Context) {
	// Парсим временные метки
	from := c.DefaultQuery("from", "now-10m") // по умолчанию выдает последние 10 минут
	to := c.DefaultQuery("to", "now")

	parser := &trparser.TimeRangeParser{}
	now := time.Now()
	timeRange, err := parser.Parse(fmt.Sprintf("from=%s&to=%s", from, to), now)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный формат временного диапазона",
		})
		return
	}
	// Парсим фильтры
	filter := models.DetectionFilter{
		HostName:    c.DefaultQuery("hostname", ""),
		TopCategory: c.DefaultQuery("category", ""),
	}
	// Парсим параметры для пагинации
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный формат page",
		})
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный формат limit",
		})
		return
	}
	pagination := models.Pagination{
		Page:  page,
		Limit: limit,
	}
	detections, total, err := s.u.GetTopDetections(c, timeRange, filter, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ошибка при получении сессий",
		})
		return
	}

	response := dto.ListResponse{
		Data: detections,
		Meta: dto.PaginationMeta{
			Page:  page,
			Limit: limit,
			Total: total,
			Pages: int(math.Ceil(float64(total) / float64(limit))),
		},
	}
	c.JSON(http.StatusOK, response)
}

// @Summary Получение статистики по детекциям
// @Description Возвращает статистику детекций по категориям (обнаружено, принято, отклонено, неразрешено) за указанный период
// @Tags detections
// @Accept json
// @Produce json
// @Param from query string false "Начало временного диапазона (формат: now-10m, now-1h, 2024-01-01T00:00:00Z)" default(now-10m)
// @Param to query string false "Конец временного диапазона (формат: now, 2024-01-01T00:00:00Z)" default(now)
// @Param hostname query string false "Фильтр по имени хоста"
// @Param category query string false "Фильтр по категории" Enums(Агрессия, расизм, терроризм, Ботнеты, Веб-почта, Досуг и развлечения, Интернет-магазины, Компьютерные игры, Криптомайнинг, Наркотики, Порнография и секс, Прокси и анонимайзеры, Реестр запрещенных сайтов, Сайты для взрослых, Сайты, распространяющие вирусы, Социальные сети, Торренты и Р2Р-сети, Файловые архивы, Фильмы и видео онлайн, Фишинг, Чаты и мессенджеры, Дополнительно, Криптоджекинг, Реклама, Онлайн-игры, Игровые платформы, Вредоносное ПО, Азартные игры, Депресивный контент и суицид, Алкоголь, табак)
// @Success 200 {object} dto.DetectionStat "Статистика детекций"
// @Failure 400 {object} map[string]string "Неверный формат временного диапазона"
// @Failure 500 {object} map[string]string "Ошибка при получении статистики выявлений"
// @Router /api/v1/detections/stat [get]
func (s *Server) GetDetectionStat(c *gin.Context) {
	// Парсим временные метки
	from := c.DefaultQuery("from", "now-10m") // по умолчанию выдает последние 10 минут
	to := c.DefaultQuery("to", "now")

	parser := &trparser.TimeRangeParser{}
	now := time.Now()
	timeRange, err := parser.Parse(fmt.Sprintf("from=%s&to=%s", from, to), now)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный формат временного диапазона",
		})
		return
	}
	// Парсим фильтры
	filter := models.DetectionFilter{
		HostName:    c.DefaultQuery("hostname", ""),
		TopCategory: c.DefaultQuery("category", ""),
	}
	stat, err := s.u.GetDetectionStat(c, timeRange, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ошибка при получении статистики выявлений",
		})
		return
	}
	c.JSON(http.StatusOK, stat)
}

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
// @Failure 400 {object} dto.ErrorResponse "Неверный формат параметров"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/dashboards/requests [get]
func (s *Server) GetRequestsStat(c *gin.Context) {
	//Парсим временные метки
	from := c.DefaultQuery("from", "now-10m") // по умолчанию выдает последние 10 минут
	to := c.DefaultQuery("to", "now")

	parser := &trparser.TimeRangeParser{}
	now := time.Now()
	timeRange, err := parser.Parse(fmt.Sprintf("from=%s&to=%s", from, to), now)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный формат временного диапазона",
		})
		return
	}
	// Парсим тип запросов (по умолчанию - запрещено/prohibited)
	status, err := status.ParseStatus(c.DefaultQuery("status", "prohibited"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверныф тип запроса",
		})
		return
	}
	// Парсим фильтры
	filter := models.DashboardFilter{
		HostName:    c.DefaultQuery("hostname", ""),
		TopCategory: c.DefaultQuery("category", ""),
	}
	// Парсим колтчество точек
	count, err := strconv.ParseUint(c.DefaultQuery("count", "20"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный формат count",
		})
		return
	}
	// Получем данные
	stat, err := s.u.GetRequestsStat(c, timeRange, filter, status, uint(count))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ошибка при получении статистики запросов",
		})
		return
	}
	resp := dto.DataPointsResponse{
		Type:  fmt.Sprintf("requests_%s", status.ToJSONString()),
		Data:  stat,
		Count: uint(count),
	}
	c.JSON(http.StatusOK, resp)
}

// @Summary Получить список имен устройств
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

// @Summary Получить список категорий контента
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

// GetTrafficStat возвращает статистику трафика за указанный временной диапазон
// @Summary Получить статистику трафика
// @Description Возвращает статистику трафика за указанный временной диапазон с заданным количеством точек данных в Кб
// @Tags dashboards
// @Accept json
// @Produce json
// @Param from query string false "Начало временного диапазона в формате парсера времени (по умолчанию now-10m)" default(now-10m)
// @Param to query string false "Конец временного диапазона в формате парсера времени (по умолчанию now)" default(now)
// @Param count query integer false "Количество точек данных для возврата (по умолчанию 20)" minimum(1) default(20)
// @Success 200 {object} dto.TrafficStatResponse "Успешный ответ со статистикой трафика"
// @Failure 400 {object} dto.ErrorResponse "Неверный формат параметров запроса"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера при получении статистики"
// @Router /api/v1/dashboards/traffic [get]
func (s *Server) GetTrafficStat(c *gin.Context) {
	// Парсим временные метки
	from := c.DefaultQuery("from", "now-10m") // по умолчанию выдает последние 10 минут
	to := c.DefaultQuery("to", "now")

	parser := &trparser.TimeRangeParser{}
	now := time.Now()
	timeRange, err := parser.Parse(fmt.Sprintf("from=%s&to=%s", from, to), now)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный формат временного диапазона",
		})
		return
	}
	// Парсим колтчество точек
	count, err := strconv.ParseUint(c.DefaultQuery("count", "20"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный формат count",
		})
		return
	}
	stat, err := s.u.GetTrafficStat(c, timeRange, uint(count))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ошибка при получении статистики трафика",
		})
		return
	}
	resp := dto.TrafficStatResponse{
		Data:  stat,
		Count: uint(count),
	}
	c.JSON(http.StatusOK, resp)
}

//Устарело!

// @Summary Получение списка популярных ресурсов
// @Description Получение топ ресурсов с указанием количества обращений за указанный период
// @Tags not implemented
// @Produce application/json
// @Param start query int64 true "Начало периода в timestamp"
// @Param end query int64 true "Конец периода в timestamp"
// @Param count query int false "Количество возвращаемых ресурсов (по умолчанию 5)"
// @Success 200 {object} map[string]int "JSON объект, где ключ - URL ресурса, значение - количество обращений"
// @Failure 400 {object} dto.ErrorResponse "Неверный формат параметров"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/dashboards/resources [get]
func (s *Server) GetResources(c *gin.Context) {
	start, err := strconv.ParseInt(c.Query("start"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный формат start timestamp",
		})
		return
	}

	end, err := strconv.ParseInt(c.Query("end"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный формат end timestamp",
		})
		return
	}

	count, err := strconv.Atoi(c.DefaultQuery("count", "5"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный формат count",
		})
		return
	}

	startTime := time.Unix(start, 0)
	endTime := time.Unix(end, 0)

	response, err := s.u.GetResources(
		c,
		startTime,
		endTime,
		count,
		"",
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ошибка при получении ресурсов",
		})
		return
	}
	c.JSON(http.StatusOK, response)
}

// @Summary Получение статистики по устройствам
// @Description Получение агрегированной статистики по сетевым узлам за указанный период
// @Tags not implemented
// @Produce application/json
// @Param start query int64 true "Начало периода в timestamp"
// @Param end query int64 true "Конец периода в timestamp"
// @Success 200 {object} models.DeviceStat "Статистика по устройствам"
// @Failure 400 {object} dto.ErrorResponse "Неверный формат параметров"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/dashboards/devices [get]
func (s *Server) GetDevicesStat(c *gin.Context) {
	start, err := strconv.ParseInt(c.Query("start"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный формат start timestamp",
		})
		return
	}

	end, err := strconv.ParseInt(c.Query("end"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный формат end timestamp",
		})
		return
	}

	startTime := time.Unix(start, 0)
	endTime := time.Unix(end, 0)

	response, err := s.u.GetDevicesStat(
		c,
		startTime,
		endTime,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ошибка при получении статистики по узлам",
		})
		return
	}
	c.JSON(http.StatusOK, response)
}

// @Summary Получение графика запрещенной активности
// @Description Получение расписания запрещенной активности начиная с указанной даты
// @Tags not implemented
// @Produce application/json
// @Param start query int64 true "Дата начала в timestamp"
// @Success 200 {object} map[int64]int "График запрещенной активности"
// @Failure 400 {object} dto.ErrorResponse "Неверный формат параметров"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/dashboards/proh_activity [get]
func (s *Server) GetProhActivity(c *gin.Context) {
	start, err := strconv.ParseInt(c.Query("start"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный формат start timestamp",
		})
		return
	}

	startTime := time.Unix(start, 0)

	response, err := s.u.GetProhActSchedule(
		c,
		startTime,
		"",
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ошибка при получении графика запрещенной активности",
		})
		return
	}
	c.JSON(http.StatusOK, response)
}

// @Summary Получение информации об аномалиях
// @Description Получение списка обнаруженных аномалий в сетевом трафике за указанный период
// @Tags not implemented
// @Produce application/json
// @Param start query int64 true "Начало периода в timestamp"
// @Param end query int64 true "Конец периода в timestamp"
// @Success 200 {array} models.Anomaly "Список обнаруженных аномалий"
// @Failure 400 {object} dto.ErrorResponse "Неверный формат параметров"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/dashboards/anomalies [get]
func (s *Server) GetAnomalies(c *gin.Context) {
	start, err := strconv.ParseInt(c.Query("start"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный формат start timestamp",
		})
		return
	}

	end, err := strconv.ParseInt(c.Query("end"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный формат end timestamp",
		})
		return
	}

	startTime := time.Unix(start, 0)
	endTime := time.Unix(end, 0)

	response, err := s.u.GetAnomalies(
		c,
		startTime,
		endTime,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ошибка при получении статистики по аномалиям",
		})
		return
	}
	c.JSON(http.StatusOK, response)
}
