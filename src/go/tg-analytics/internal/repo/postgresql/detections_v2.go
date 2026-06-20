package postgresql

import (
	"context"
	"fmt"
	"tg-an/internal/controllers/http/v1/dto"
	models "tg-an/internal/models"
	pkg "tg-an/pkg/models"
	"tg-an/pkg/slogger/wsl"
	"tg-an/pkg/trparser"
	"time"
)

func (r *RepoPG) GetTopDetectionsv2(ctx context.Context, tr *trparser.TimeRange, f models.DetectionFilter, thresh float32, action string, p models.Pagination, search string, sorting models.Sorting) ([]dto.Detection_v2, int64, error) {
	var detections []dto.Detection_v2
	var total int64

	// Основной запрос для получения детекций
	type DetectionResult struct {
		IP            string    `gorm:"column:ip"`
		Port          int       `gorm:"column:port"`
		Location      string    `gorm:"column:location"`
		Domain        string    `gorm:"column:domain"`
		RequestCount  int       `gorm:"column:request_count"`
		HostName      string    `gorm:"column:host_name"`
		Category      string    `gorm:"column:category"`
		Description   string    `gorm:"column:description"`
		Action        string    `gorm:"column:action"`
		CategorizedAt time.Time `gorm:"column:categorized_at"`
		NegRate       float32   `gorm:"column:neg_rate"`
		DomainID      uint      `gorm:"column:domain_id"`
	}

	var results []DetectionResult

	query := r.db.GetDB().WithContext(ctx).Table("sessions").
		Select(fmt.Sprintf(`
			domains.id as domain_id,
			domains.ip,
			domains.port,
			domains.country as location,
			domains.path as domain,
			domains.categorized_at as categorized_at,
			domains.neg_rate,
			COUNT(*) as request_count,
			devices.hostname as host_name,
			categories.name as category,
			COALESCE((
				SELECT string_agg(c2.name, ', ' ORDER BY dc2.count DESC)
				FROM domain_categories dc2
				JOIN categories c2 ON c2.id = dc2.category_id
				WHERE dc2.domain_id = domains.id AND c2.type = '%s'
			), '') as description,
			COALESCE(actions.action, '%s') as action
		`, pkg.CategoryTypeNegative.String(), pkg.ActionTypeUnresolved.String())).
		Joins("LEFT JOIN devices ON sessions.device_id = devices.id").
		Joins("LEFT JOIN domains ON sessions.domain_id = domains.id").
		Joins("LEFT JOIN actions ON domains.action_id = actions.id").
		Joins("LEFT JOIN categories ON domains.category_id = categories.id").
		Where("categories.type = ?", pkg.CategoryTypeNegative.String())

	// Применяем временной диапазон
	if tr != nil && !tr.To.IsZero() {
		query = query.Where("sessions.datetime_utc BETWEEN ? AND ?", tr.From, tr.To)
	}

	// Применяем фильтры
	if f.HostName != "" {
		query = query.Where("devices.hostname = ?", f.HostName)
	}
	if f.TopCategory != "" {
		query = query.Where("categories.name = ?", f.TopCategory)
	}

	// Фильтрация по статусу выявления
	switch f.Status {
	case "Все":
		break
	case "Рекомендуется блокировка":
		query = query.Where("domains.neg_rate >= ?", thresh)
		query = query.Where("actions.action IS NULL")
	case "Требуется проверка":
		query = query.Where("domains.neg_rate < ?", thresh)
		query = query.Where("actions.action IS NULL")
	case "Заблокирован":
		query = query.Where("actions.action = ?", "Заблокировано")
	default:
		break
	}

	// Применяем фильтр по действию
	if action != "" {
		if action == pkg.ActionTypeUnresolved.String() {
			query = query.Where("actions.action IS NULL")
		} else {
			query = query.Where("actions.action = ?", action)
		}
	}

	// Применяем поиск по полям path и ip
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("domains.path ILIKE ? OR domains.ip ILIKE ?", searchPattern, searchPattern)
	}

	// Группируем по уникальным детекциям
	query = query.Group(fmt.Sprintf(`
		domains.id, domains.ip, domains.port, domains.country, domains.path, domains.categorized_at,
		devices.hostname, categories.name, COALESCE(actions.action, '%s')
	`, pkg.ActionTypeUnresolved.String()))

	// Получаем общее количество записей (до пагинации)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count detections: %w", err)
	}

	// Применяем сортировку
	if sorting.OrderBy != "" {
		orderField := getDetectionOrderField(sorting.OrderBy)
		orderDirection := "ASC"
		if sorting.OrderDir == "desc" {
			orderDirection = "DESC"
		}
		query = query.Order(orderField + " " + orderDirection)
	} else {
		query = query.Order("request_count DESC, categorized_at DESC")
	}

	// Применяем пагинацию
	if p.Limit > 0 {
		offset := (p.Page - 1) * p.Limit
		if offset < 0 {
			offset = 0
		}
		query = query.Offset(offset).Limit(p.Limit)
	}

	// Выполняем основной запрос
	if err := query.Find(&results).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get detections: %w", err)
	}

	// Если нет результатов, возвращаем пустой слайс
	if len(results) == 0 {
		return []dto.Detection_v2{}, total, nil
	}

	// Собираем все domain_id для которых нужно получить статистику
	domainIDs := make([]uint, 0, len(results))
	for _, result := range results {
		domainIDs = append(domainIDs, result.DomainID)
	}

	// Получаем статистику для всех доменов одним запросом
	categoriesStatsMap, err := r.getCategoriesStatsForDomains(ctx, domainIDs)
	if err != nil {
		r.l.WarnContext(ctx, "failed to get categories stats for domains", wsl.Err(err))
		categoriesStatsMap = make(map[uint][]dto.CategoryStat)
	}

	// Формируем результат
	for _, result := range results {
		detections = append(detections, dto.Detection_v2{
			IP:            result.IP,
			Port:          result.Port,
			Location:      result.Location,
			Domain:        result.Domain,
			RequestCount:  result.RequestCount,
			HostName:      result.HostName,
			Category:      result.Category,
			Description:   result.Description,
			Action:        result.Action,
			CategorizedAt: result.CategorizedAt,
			NegRate:       result.NegRate,
			Categories:    categoriesStatsMap[result.DomainID],
		})
	}

	return detections, total, nil
}

// getCategoriesStatsForDomains получает статистику категорий для нескольких доменов одним запросом
func (r *RepoPG) getCategoriesStatsForDomains(ctx context.Context, domainIDs []uint) (map[uint][]dto.CategoryStat, error) {
	if len(domainIDs) == 0 {
		return make(map[uint][]dto.CategoryStat), nil
	}

	type CategoryStatResult struct {
		DomainID   uint
		Name       string
		Count      float64
		TotalCount float64
	}

	var results []CategoryStatResult

	err := r.db.GetDB().WithContext(ctx).
		Table("domain_categories").
		Select(`
			domain_categories.domain_id,
			categories.name,
			domain_categories.count,
			SUM(domain_categories.count) OVER (PARTITION BY domain_categories.domain_id) as total_count
		`).
		Joins("JOIN categories ON categories.id = domain_categories.category_id").
		Where("domain_categories.domain_id IN ? AND categories.type = ?", domainIDs, pkg.CategoryTypeNegative.String()).
		Order("domain_categories.domain_id, domain_categories.count DESC").
		Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get categories stats: %w", err)
	}

	// Группируем результаты по domain_id
	statsMap := make(map[uint][]dto.CategoryStat)
	for _, res := range results {
		var rate float32 = 0
		if res.TotalCount > 0 {
			rate = float32(res.Count / res.TotalCount)
		}
		statsMap[res.DomainID] = append(statsMap[res.DomainID], dto.CategoryStat{
			Name: res.Name,
			Rate: rate,
		})
	}

	// Заполняем пустым слайсом для доменов без категорий
	for _, domainID := range domainIDs {
		if _, exists := statsMap[domainID]; !exists {
			statsMap[domainID] = []dto.CategoryStat{}
		}
	}

	return statsMap, nil
}
