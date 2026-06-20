package v1

import (
	"fmt"
	"math"
	"net/http"
	"tg-an/internal/controllers/http/v1/dto"
	"tg-an/internal/controllers/http/v1/utils"
	"tg-an/internal/controllers/http/v1/validation"
	"tg-an/internal/models"
	"tg-an/pkg/slogger"
	"tg-an/pkg/trparser"
	"time"

	"github.com/gin-gonic/gin"
)

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
// @Success 200 {object} dto.GetDetectionsResponseV2 "Успешный ответ"
// @Failure 400 {object} dto.ErrorResponse "Неверный формат параметров"
// @Failure 403 {object} dto.ErrorResponse "Недостаточно прав для доступа к выявлениям"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v2/detections [get]
func (s *Server) GetDetections_v2(c *gin.Context) {
	userMeta, err := utils.GetUserMeta(c)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "utils.GetUserMeta", slogger.WrapError(c.Request.Context(), err))
		return
	}

	// Валидация запроса
	var req validation.GetTopDetectionsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "c.ShouldBindQuery", slogger.WrapError(c.Request.Context(), err))
		return
	}

	// Нормализация и валидация
	if err := req.ValidateAndNormalize(); err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "req.ValidateAndNormalize", slogger.WrapError(c.Request.Context(), err))
		return
	}

	// Парсим временной диапазон
	parser := &trparser.TimeRangeParser{}
	now := time.Now()
	timeRange, err := parser.Parse(fmt.Sprintf("from=%s&to=%s", req.From, req.To), now)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "parser.Parse", slogger.WrapError(c.Request.Context(), err))
		return
	}

	// Преобразуем в доменные модели
	filter := models.DetectionFilter{
		HostName:    req.HostName,
		TopCategory: req.Category,
		Status:      req.DetectionStatus,
	}

	pagination := models.Pagination{
		Page:  req.Page,
		Limit: req.Limit,
	}

	sorting := models.Sorting{
		OrderBy:  req.OrderBy,
		OrderDir: req.OrderDir,
	}

	// Получаем данные из usecase (версия 2)
	detections, total, err := s.u.GetTopDetections_v2(c.Request.Context(), userMeta, timeRange, filter, req.Action, pagination, req.Search, sorting)
	if err != nil {
		s.ErrorResponse(c, http.StatusInternalServerError, "s.u.GetTopDetections_v2", slogger.WrapError(c.Request.Context(), err))
		return
	}

	// Формируем ответ
	response := dto.GetDetectionsResponseV2{
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
