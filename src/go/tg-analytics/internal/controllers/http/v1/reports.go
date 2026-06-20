package v1

import (
	"fmt"
	"net/http"
	"tg-an/internal/controllers/http/v1/utils"
	"tg-an/internal/controllers/http/v1/validation"
	"tg-an/internal/controllers/http/v1/values"
	"tg-an/pkg/slogger"
	"tg-an/pkg/trparser"
	"time"

	"github.com/gin-gonic/gin"
)

// @Summary Создание отчета
// @Description Генерирует полный отчет по активности за указанный временной период, включая аналитику по устройствам, категориям и аномалиям
// @Tags reports
// @Accept json
// @Produce json
// @Param from query string false "Начало временного диапазона (формат: now-24h, 2023-12-01T10:00:00Z)" default(now-24h)
// @Param to query string false "Конец временного диапазона (формат: now, 2023-12-01T12:00:00Z)" default(now)
// @Security BearerAuth
// @Success 200 {object} models.Report "Полный отчет по активности"
// @Failure 400 {object} dto.ErrorResponse "Неверный формат параметров"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера при генерации отчета"
// @Router /api/v1/reports [get]
func (s *Server) CreateReport(c *gin.Context) {
	userMeta, err := utils.GetUserMeta(c)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "utils.GetUserMeta", slogger.WrapError(c.Request.Context(), err))
		return
	}

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

	report, err := s.u.CreateReport(c.Request.Context(), userMeta, timeRange)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка при создании отчета",
		})
		return
	}
	c.JSON(http.StatusOK, report)
}

// @Summary Создание отчета по конкретному устройству
// @Description Генерирует детализированный отчет по конкретному сетевому устройству за указанный временной период, включая статистику трафика, аномалии и категории запросов
// @Tags reports
// @Accept json
// @Produce json
// @Param hostname path string true "Имя сетевого устройства (хоста)"
// @Param from query string false "Начало временного диапазона (формат: now-24h, 2023-12-01T10:00:00Z)" default(now-24h)
// @Param to query string false "Конец временного диапазона (формат: now, 2023-12-01T12:00:00Z)" default(now)
// @Security BearerAuth
// @Success 200 {object} models.ReportForDevice "Детализированный отчет по устройству"
// @Failure 400 {object} dto.ErrorResponse "Неверный формат параметров или устройство не найдено"
// @Failure 403 {object} dto.ErrorResponse "Недостаточно прав для доступа к устройству"
// @Failure 404 {object} dto.ErrorResponse "Устройство не найдено в базе данных"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера при генерации отчета"
// @Router /api/v1/reports/{hostname} [get]
func (s *Server) CreateReportForDevice(c *gin.Context) {
	userMeta, err := utils.GetUserMeta(c)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "utils.GetUserMeta", slogger.WrapError(c.Request.Context(), err))
		return
	}

	// Проверяем, что устройство есть в БД
	hostname := c.Param("hostname")
	devices, err := s.u.GetDevices(c.Request.Context(), userMeta)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка при получении списка устройств",
		})
		return
	}

	// Проверяем существование устройства
	found := false
	for _, device := range devices {
		if device == hostname {
			found = true
			break
		}
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("Устройство '%s' не найдено", hostname),
		})
		return
	}

	// Проверяем права доступа для контекстных администраторов
	if userMeta.ShortRole == values.ContextAdmin && hostname != userMeta.ClientRole {
		c.JSON(http.StatusForbidden, gin.H{
			"error": fmt.Sprintf("Недостаточно прав для доступа к устройству '%s'", hostname),
		})
		return
	}

	// Валидация запроса
	var req validation.GetAnomaliesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Неверные параметры запроса: %v", err),
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

	report, err := s.u.CreateReportForDevice(c.Request.Context(), userMeta, timeRange, hostname)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка при создании отчета для устройства",
		})
		return
	}
	c.JSON(http.StatusOK, report)
}
