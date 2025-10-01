package v1

import (
	"net/http"

	"dashboard_serv/internal/controllers/http/v1/dto"

	"github.com/gin-gonic/gin"
)

// DashboardHandler обработчик дашбордов
type DashboardHandler struct {
	// здесь будут зависимости (сервисы, репозитории)
}

// NewDashboardHandler создает новый экземпляр обработчика
func NewDashboardHandler() *DashboardHandler {
	return &DashboardHandler{}
}

// GetRequestStats возвращает статистику запросов
// @Summary Получить статистику запросов
// @Description Возвращает количество запросов (заблокированные, разрешенные) в целом и по отдельному NGFW
// @Tags dashboard
// @Accept json
// @Produce json
// @Param filter body dto.TimeRangeFilter true "Фильтр по времени и NGFW"
// @Success 200 {object} dto.RequestStatsResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/dashboard/requests [post]
func (h *DashboardHandler) GetRequestStats(c *gin.Context) {
	var filter dto.TimeRangeFilter
	if err := c.ShouldBindJSON(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: логика получения статистики запросов
	// response := h.service.GetRequestStats(filter)
	// c.JSON(http.StatusOK, response)
}

// GetTrafficStats возвращает статистику трафика
// @Summary Получить статистику трафика
// @Description Возвращает объем трафика (исходящий, входящий) в целом и по отдельному NGFW
// @Tags dashboard
// @Accept json
// @Produce json
// @Param filter body dto.TimeRangeFilter true "Фильтр по времени и NGFW"
// @Success 200 {object} dto.TrafficStatsResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/dashboard/traffic [post]
func (h *DashboardHandler) GetTrafficStats(c *gin.Context) {
	var filter dto.TimeRangeFilter
	if err := c.ShouldBindJSON(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: логика получения статистики трафика
}

// GetTopResources возвращает топ ресурсов и категорий
// @Summary Получить топ ресурсов и категорий
// @Description Возвращает топ 25 посещаемых ресурсов по доменным именам и топ 25 посещаемых категорий
// @Tags dashboard
// @Accept json
// @Produce json
// @Param filter body dto.TimeRangeFilter true "Фильтр по времени и NGFW"
// @Success 200 {object} dto.TopResourcesResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/dashboard/top-resources [post]
func (h *DashboardHandler) GetTopResources(c *gin.Context) {
	var filter dto.TimeRangeFilter
	if err := c.ShouldBindJSON(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: логика получения топа ресурсов
}

// GetBlockedCategories возвращает сработавшие запрещающие категории
// @Summary Получить запрещающие категории
// @Description Возвращает сработавшие запрещающие категории в целом и по отдельному NGFW
// @Tags dashboard
// @Accept json
// @Produce json
// @Param filter body dto.TimeRangeFilter true "Фильтр по времени и NGFW"
// @Success 200 {object} dto.BlockedCategoriesResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/dashboard/blocked-categories [post]
func (h *DashboardHandler) GetBlockedCategories(c *gin.Context) {
	var filter dto.TimeRangeFilter
	if err := c.ShouldBindJSON(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: логика получения запрещенных категорий
}

// GetBlockedResources возвращает ресурсы запрещенных категорий
// @Summary Получить ресурсы запрещенных категорий
// @Description Возвращает ресурсы которые относятся к запрещенным категориям с детальной информацией
// @Tags dashboard
// @Accept json
// @Produce json
// @Param filter body dto.TimeRangeFilter true "Фильтр по времени, NGFW и категориям"
// @Success 200 {object} dto.BlockedResourcesResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/dashboard/blocked-resources [post]
func (h *DashboardHandler) GetBlockedResources(c *gin.Context) {
	var filter dto.TimeRangeFilter
	if err := c.ShouldBindJSON(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: логика получения ресурсов запрещенных категорий
}

// GetNGFWList возвращает список доступных NGFW
// @Summary Получить список NGFW
// @Description Возвращает список всех доступных межсетевых экранов
// @Tags dashboard
// @Produce json
// @Success 200 {array} entity.NGFW
// @Failure 500 {object} map[string]string
// @Router /api/v1/dashboard/ngfw [get]
func (h *DashboardHandler) GetNGFWList(c *gin.Context) {
	// TODO: логика получения списка NGFW
}

// GetCategoriesList возвращает список категорий
// @Summary Получить список категорий
// @Description Возвращает список всех категорий ресурсов
// @Tags dashboard
// @Produce json
// @Success 200 {array} entity.Category
// @Failure 500 {object} map[string]string
// @Router /api/v1/dashboard/categories [get]
func (h *DashboardHandler) GetCategoriesList(c *gin.Context) {
	// TODO: логика получения списка категорий
}
