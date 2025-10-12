package v1

import (
	"math"
	"net/http"
	"strconv"
	"tg-dbd/internal/controllers/http/v1/dto"
	"time"

	"github.com/gin-gonic/gin"
)

// @Summary Получение списка сессий в заданном временном диапазоне
// @Description Метод возвращает список сессий с поддержкой пагинации и сортировки
// @Tags session
// @Accept  json
// @Produce  json
// @Param start query int64 true "Начальный timestamp диапазона"
// @Param end query int64 true "Конечный timestamp диапазона"
// @Param page query int false "Номер страницы" Default: 1
// @Param limit query int false "Количество записей на странице" Default: 10
// @Param order_by query string false "Поле для сортировки" Default: url
// @Param order_dir query string false "Направление сортировки" Enum: asc,desc Default: desc
// @Success 200 {object} dto.SessionListResponse "Успешный ответ"
// @Failure 400 {object} dto.ErrorResponse "Неверный формат параметров"
// @Failure 500 {object} dto.ErrorResponse "Ошибка сервера"
// @Router /v1/sessions [get]
func (s *Server) getSessions(c *gin.Context) {
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

	orderBy := c.DefaultQuery("order_by", "url")
	orderDir := c.DefaultQuery("order_dir", "desc")

	// Конвертация в time.Time
	startTime := time.Unix(start, 0)
	endTime := time.Unix(end, 0)

	sessions, total, err := s.u.GetSessions(c,
		startTime,
		endTime,
		page,
		limit,
		"", // TODO: добавить обработку фильтров
		orderBy,
		orderDir,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ошибка при получении сессий",
		})
		return
	}

	response := dto.SessionListResponse{
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

// @Summary Получение информации об обнаружениях
// @Description Получение списка обнаружений с фильтрацией по дате
// @Tags detection
// @Produce application/json
// @Param start query int64 true "Начало периода в timestamp"
// @Param end query int64 true "Конец периода в timestamp"
// @Param count query int false "Количество записей (по умолчанию 5)"
// @Success 200 {object} models.Detection
// @Failure 400 {object} dto.ErrorResponse "Неверный формат параметров"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /v1/detections [get]
func (s *Server) getDetections(c *gin.Context) {
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

	response, err := s.u.GetDetections(
		c,
		startTime,
		endTime,
		count,
		"",
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ошибка при получении выявлений",
		})
		return
	}
	c.JSON(http.StatusOK, response)
}

// @Summary Получение статистики обнаружений
// @Description Получение агрегированной статистики обнаружений за указанный период
// @Tags detection
// @Produce application/json
// @Param start query int64 true "Начало периода в timestamp"
// @Param end query int64 true "Конец периода в timestamp"
// @Success 200 {object} models.DetectionStat
// @Failure 400 {object} dto.ErrorResponse "Неверный формат параметров"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /v1/detections/stat [get]
func (s *Server) getDetectionStat(c *gin.Context) {
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

	response, err := s.u.GetDetectionsStat(
		c,
		startTime,
		endTime,
		"",
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ошибка при получении статистики выявлений",
		})
		return
	}
	c.JSON(http.StatusOK, response)
}

// @Summary Получение списка категорий
// @Description Получение топ категорий с указанием количества обращений за указанный период
// @Tags dashboard
// @Produce application/json
// @Param start query int64 true "Начало периода в timestamp"
// @Param end query int64 true "Конец периода в timestamp"
// @Param count query int false "Количество возвращаемых категорий (по умолчанию 5)"
// @Success 200 {object} map[string]int "JSON объект, где ключ - название категории, значение - количество обращений"
// @Failure 400 {object} dto.ErrorResponse "Неверный формат параметров"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /v1/dashboards/categories [get]
func (s *Server) getCategories(c *gin.Context) {
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

	response, err := s.u.GetCategories(
		c,
		startTime,
		endTime,
		count,
		"",
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ошибка при получении категорий",
		})
		return
	}
	c.JSON(http.StatusOK, response)
}

// @Summary Получение списка популярных ресурсов
// @Description Получение топ ресурсов с указанием количества обращений за указанный период
// @Tags dashboard
// @Produce application/json
// @Param start query int64 true "Начало периода в timestamp"
// @Param end query int64 true "Конец периода в timestamp"
// @Param count query int false "Количество возвращаемых ресурсов (по умолчанию 5)"
// @Success 200 {object} map[string]int "JSON объект, где ключ - URL ресурса, значение - количество обращений"
// @Failure 400 {object} dto.ErrorResponse "Неверный формат параметров"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /v1/dashboards/resources [get]
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

// @Summary Получение списка уведомлений
// @Description Получение списка уведомлений системы безопасности за указанный период
// @Tags dashboard
// @Produce application/json
// @Param start query int64 true "Начало периода в timestamp"
// @Param end query int64 true "Конец периода в timestamp"
// @Param count query int false "Количество возвращаемых уведомлений (по умолчанию 5)"
// @Success 200 {array} models.Notification "Список уведомлений"
// @Failure 400 {object} dto.ErrorResponse "Неверный формат параметров"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /v1/dashboards/events [get]
func (s *Server) GetEvents(c *gin.Context) {
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

	response, err := s.u.GetEvents(
		c,
		startTime,
		endTime,
		count,
		"",
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ошибка при получении уведомлений",
		})
		return
	}
	c.JSON(http.StatusOK, response)
}

// @Summary Получение статистики по устройствам
// @Description Получение агрегированной статистики по сетевым узлам за указанный период
// @Tags dashboard
// @Produce application/json
// @Param start query int64 true "Начало периода в timestamp"
// @Param end query int64 true "Конец периода в timestamp"
// @Success 200 {object} models.DeviceStat "Статистика по устройствам"
// @Failure 400 {object} dto.ErrorResponse "Неверный формат параметров"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /v1/dashboards/devices [get]
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

// @Summary Получение статистики сетевого трафика
// @Description Получение агрегированной статистики по сетевому трафику за указанный период
// @Tags dashboard
// @Produce application/json
// @Param start query int64 true "Начало периода в timestamp"
// @Param end query int64 true "Конец периода в timestamp"
// @Success 200 {object} []models.TrafficPoint "Статистика сетевого трафика"
// @Failure 400 {object} dto.ErrorResponse "Неверный формат параметров"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /v1/dashboards/traffic [get]
func (s *Server) GetTrafficStat(c *gin.Context) {
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

	response, err := s.u.GetTrafficStat(
		c,
		startTime,
		endTime,
		"",
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ошибка при получении статистики по трафику",
		})
		return
	}
	c.JSON(http.StatusOK, response)
}

// @Summary Получение статистики HTTP-запросов
// @Description Получение агрегированной статистики по HTTP-запросам за указанный период
// @Tags dashboard
// @Produce application/json
// @Param start query int64 true "Начало периода в timestamp"
// @Param end query int64 true "Конец периода в timestamp"
// @Success 200 {object} models.RequestsStat "Статистика HTTP-запросов"
// @Failure 400 {object} dto.ErrorResponse "Неверный формат параметров"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /v1/dashboards/requests [get]
func (s *Server) GetRequestsStat(c *gin.Context) {
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

	response, err := s.u.GetRequestsStat(
		c,
		startTime,
		endTime,
		"",
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ошибка при получении статистики по запросам",
		})
		return
	}
	c.JSON(http.StatusOK, response)
}

// @Summary Получение графика запрещенной активности
// @Description Получение расписания запрещенной активности начиная с указанной даты
// @Tags dashboard
// @Produce application/json
// @Param start query int64 true "Дата начала в timestamp"
// @Success 200 {object} map[int64]int "График запрещенной активности"
// @Failure 400 {object} dto.ErrorResponse "Неверный формат параметров"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /v1/dashboards/proh_activity [get]
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
// @Tags dashboard
// @Produce application/json
// @Param start query int64 true "Начало периода в timestamp"
// @Param end query int64 true "Конец периода в timestamp"
// @Success 200 {array} models.Anomaly "Список обнаруженных аномалий"
// @Failure 400 {object} dto.ErrorResponse "Неверный формат параметров"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /v1/dashboards/anomalies [get]
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
