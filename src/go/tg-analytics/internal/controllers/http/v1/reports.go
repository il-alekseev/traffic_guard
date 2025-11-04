package v1

import (
	"fmt"
	"net/http"
	"tg-an/internal/controllers/http/v1/validation"
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
// @Success 200 {object} models.Report "Полный отчет по активности"
// @Failure 400 {object} dto.ErrorResponse "Неверный формат параметров"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера при генерации отчета"
// @Router /api/v1/reports [get]
func (s *Server) CreateReport(c *gin.Context) {
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
	report, err := s.u.CreateReport(c, timeRange)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error creating report",
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
// @Success 200 {object} models.ReportForDevice "Детализированный отчет по устройству"
// @Failure 400 {object} dto.ErrorResponse "Неверный формат параметров или устройство не найдено"
// @Failure 404 {object} dto.ErrorResponse "Устройство не найдено в базе данных"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера при генерации отчета"
// @Router /api/v1/reports/{hostname} [get]
func (s *Server) CreateReportForDevice(c *gin.Context) {
	// Проверяем, что устройство есть в БД
	hostname := c.Param("hostname")
	devices, err := s.u.GetDevices(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("error with getting devices: %v", err),
		})
		return
	}
	found := false
	for _, device := range devices {
		if device == hostname {
			found = true
			break
		}
	}

	if !found {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("device '%s' not found in DB", hostname),
		})
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
	report, err := s.u.CreateReportForDevice(c, timeRange, hostname)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error creating report for device",
		})
		return
	}
	c.JSON(http.StatusOK, report)
}
