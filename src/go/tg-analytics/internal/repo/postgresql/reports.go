package postgresql

import (
	"context"
	"fmt"
	"log/slog"
	"tg-an/internal/models"
	"tg-an/internal/pkg/status"
	pkg "tg-an/pkg/models"
	"tg-an/pkg/trparser"
	"time"
)

// GetCategories возвращает категории для главной страницы отчета
func (r *RepoPG) GetCategories(ctx context.Context, tr *trparser.TimeRange) ([]models.CategoryStat, error) {
	type categoryTemp struct {
		Category    string `json:"category"`
		Total       uint   `json:"total"`
		BeforeBlock uint   `json:"before_block"`
		AfterBlock  uint   `json:"after_block"`
		Pending     uint   `json:"pending"`
	}

	var tempResults []categoryTemp

	query := r.db.GetDB().WithContext(ctx).Table("sessions").
		Select(`
			categories.name as category,
			COUNT(*) as total,
			COUNT(CASE WHEN sessions.datetime_utc < domains.categorized_at THEN 1 END) as before_block,
			COUNT(CASE WHEN sessions.datetime_utc >= domains.categorized_at AND sessions.status != ? THEN 1 END) as after_block,
			COUNT(CASE WHEN sessions.status = ? THEN 1 END) as pending
		`, status.StatusPending.String(), status.StatusPending.String()).
		Joins("LEFT JOIN domains ON sessions.domain_id = domains.id").
		Joins("LEFT JOIN categories ON domains.category_id = categories.id").
		Where("categories.type = ?", pkg.CategoryTypeNegative)

	// Применяем временной диапазон
	if tr != nil && !tr.From.IsZero() && !tr.To.IsZero() {
		query = query.Where("sessions.datetime_utc BETWEEN ? AND ?", tr.From, tr.To)
	}

	err := query.
		Group("categories.name").
		Order("total DESC").
		Find(&tempResults).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get categories for report: %w", err)
	}

	// Преобразуем в конечный формат
	var categories []models.CategoryStat
	for _, temp := range tempResults {
		categories = append(categories, models.CategoryStat{
			Category: temp.Category,
			Stat: models.RequestReport{
				All:         temp.Total,
				BeforeBlock: temp.BeforeBlock,
				Pending:     temp.Pending,
				AfterBlock:  temp.AfterBlock,
			},
		})
	}

	return categories, nil
}

// GetResourses возвращает ресурсы для главной страницы отчета
func (r *RepoPG) GetResourses(ctx context.Context, tr *trparser.TimeRange) ([]models.ResourceStat, error) {
	type resourceTemp struct {
		Resource    string `json:"resource"`
		Category    string `json:"category"`
		Total       uint   `json:"total"`
		BeforeBlock uint   `json:"before_block"`
		AfterBlock  uint   `json:"after_block"`
		Pending     uint   `json:"pending"`
	}

	var tempResults []resourceTemp

	query := r.db.GetDB().WithContext(ctx).Table("sessions").
		Select(`
			domains.path as resource,
			categories.name as category,
			COUNT(*) as total,
			COUNT(CASE WHEN sessions.datetime_utc < domains.categorized_at THEN 1 END) as before_block,
			COUNT(CASE WHEN sessions.datetime_utc >= domains.categorized_at AND sessions.status != ? THEN 1 END) as after_block,
			COUNT(CASE WHEN sessions.status = ? THEN 1 END) as pending
		`, status.StatusPending.String(), status.StatusPending.String()).
		Joins("LEFT JOIN domains ON sessions.domain_id = domains.id").
		Joins("LEFT JOIN categories ON domains.category_id = categories.id").
		Where("categories.type = ?", pkg.CategoryTypeNegative)

	// Применяем временной диапазон
	if tr != nil && !tr.From.IsZero() && !tr.To.IsZero() {
		query = query.Where("sessions.datetime_utc BETWEEN ? AND ?", tr.From, tr.To)
	}

	err := query.
		Group("domains.path, categories.name").
		Order("total DESC").
		Find(&tempResults).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get resources for report: %w", err)
	}

	// Преобразуем результат в нужный формат
	var result []models.ResourceStat
	for _, res := range tempResults {
		result = append(result, models.ResourceStat{
			Resource:   res.Resource,
			Categories: []string{res.Category},
			Stat: models.RequestReport{
				All:         res.Total,
				BeforeBlock: res.BeforeBlock,
				Pending:     res.Pending,
				AfterBlock:  res.AfterBlock,
			},
		})
	}

	return result, nil
}

// GetDevicesAnalytics возвращает аналитику по устройствам
func (r *RepoPG) GetDevicesAnalytics(ctx context.Context, tr *trparser.TimeRange, hostname string) ([]models.DeviceReport, error) {
	type deviceAnalyticsTemp struct {
		HostName             string `json:"hostname"`
		Requests             uint   `json:"requests"`
		Anomalies            uint   `json:"anomalies"`
		Blocks               uint   `json:"blocks"`
		AllDetections        uint   `json:"all_detections"`
		UnresolvedDetections uint   `json:"unresolved_detections"`
		BlockedDetections    uint   `json:"blocked_detections"`
		AllowedDetections    uint   `json:"allowed_detections"`
	}

	var tempResults []deviceAnalyticsTemp

	query := r.db.GetDB().WithContext(ctx).Table("sessions").
		Select(`
			devices.hostname as host_name,
			COUNT(*) as requests,
			COUNT(CASE WHEN sessions.status = ? THEN 1 END) as anomalies,
			COUNT(CASE WHEN sessions.status = ? THEN 1 END) as blocks,
			COUNT(CASE WHEN domains.category_id IS NOT NULL AND categories.type = ? THEN 1 END) as all_detections,
			COUNT(CASE WHEN domains.category_id IS NOT NULL AND categories.type = ? AND (domains.action_id IS NULL OR domains.action_id = 0) THEN 1 END) as unresolved_detections,
			COUNT(CASE WHEN domains.category_id IS NOT NULL AND categories.type = ? AND domains.action_id > 0 AND actions.action = 'block' THEN 1 END) as blocked_detections,
			COUNT(CASE WHEN domains.category_id IS NOT NULL AND categories.type = ? AND domains.action_id > 0 AND actions.action = 'allow' THEN 1 END) as allowed_detections
		`,
			status.StatusAnomaly.String(),
			status.StatusBlocked.String(),
			pkg.CategoryTypeNegative,
			pkg.CategoryTypeNegative,
			pkg.CategoryTypeNegative,
			pkg.CategoryTypeNegative).
		Joins("LEFT JOIN devices ON sessions.device_id = devices.id").
		Joins("LEFT JOIN domains ON sessions.domain_id = domains.id").
		Joins("LEFT JOIN categories ON domains.category_id = categories.id").
		Joins("LEFT JOIN actions ON domains.action_id = actions.id").
		Where("devices.hostname IS NOT NULL")

	// Применяем временной диапазон
	if tr != nil && !tr.From.IsZero() && !tr.To.IsZero() {
		query = query.Where("sessions.datetime_utc BETWEEN ? AND ?", tr.From, tr.To)
	}

	// Применяем фильтр по hostname если указан
	if hostname != "" {
		query = query.Where("devices.hostname = ?", hostname)
	}

	err := query.
		Group("devices.hostname").
		Order("requests DESC").
		Find(&tempResults).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get devices analytics: %w", err)
	}

	// Преобразуем результат в нужный формат
	var result []models.DeviceReport
	for _, device := range tempResults {
		result = append(result, models.DeviceReport{
			HostName: device.HostName,
			Traffic: models.Traffic{
				Input:  0, // Заглушка
				Output: 0, // Заглушка
			},
			Requests: device.Requests,
			AnomalyBlockStat: models.AnomalyBlockStat{
				Anomalies: device.Anomalies,
				Blocks:    device.Blocks,
				All:       device.Anomalies + device.Blocks,
			},
			Detections: models.DetectionReport{
				All:        device.AllDetections,
				Unresolved: device.UnresolvedDetections,
				Blocked:    device.BlockedDetections,
				Allowed:    device.AllowedDetections,
			},
		})
	}

	return result, nil
}

// GetAnomaliesList возвращает список аномалий
func (r *RepoPG) GetAnomaliesList(ctx context.Context, tr *trparser.TimeRange, hostname string) ([]models.DeviceAnomaly, error) {
	type anomalyTemp struct {
		HostName        string     `json:"hostname"`
		URL             string     `json:"url"`
		CategorizedAt   time.Time  `json:"categorized_at"`
		ActionCreatedAt *time.Time `json:"action_created_at"`
		Action          string     `json:"action"`
		Total           uint       `json:"total"`
		BeforeBlock     uint       `json:"before_block"`
		AfterBlock      uint       `json:"after_block"`
		Pending         uint       `json:"pending"`
	}

	var tempResults []anomalyTemp

	query := r.db.GetDB().WithContext(ctx).Table("sessions").
		Select(`
			devices.hostname as hostname,
			domains.path as url,
			domains.categorized_at as categorized_at,
			actions.created_at as action_created_at,
			actions.action as action,
			COUNT(*) as total,
			COUNT(CASE WHEN sessions.datetime_utc < domains.categorized_at THEN 1 END) as before_block,
			COUNT(CASE WHEN sessions.datetime_utc >= domains.categorized_at AND sessions.status != ? THEN 1 END) as after_block,
			COUNT(CASE WHEN sessions.status = ? THEN 1 END) as pending
		`, status.StatusPending.String(), status.StatusPending.String()).
		Joins("LEFT JOIN devices ON sessions.device_id = devices.id").
		Joins("LEFT JOIN domains ON sessions.domain_id = domains.id").
		Joins("LEFT JOIN categories ON domains.category_id = categories.id").
		Joins("LEFT JOIN actions ON domains.action_id = actions.id").
		Where("sessions.status = ?", status.StatusAnomaly.String()). // Только аномальные сессии
		Where("categories.type = ?", pkg.CategoryTypeNegative)       // Только негативные категории

	// Применяем временной диапазон
	if tr != nil && !tr.From.IsZero() && !tr.To.IsZero() {
		query = query.Where("sessions.datetime_utc BETWEEN ? AND ?", tr.From, tr.To)
	}

	// Применяем фильтр по hostname если указан
	if hostname != "" {
		query = query.Where("devices.hostname = ?", hostname)
	}

	err := query.
		Group("devices.hostname, domains.path, domains.categorized_at, actions.created_at, actions.action").
		Order("total DESC").
		Find(&tempResults).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get anomalies list: %w", err)
	}

	// Группируем по устройствам
	deviceAnomalies := make(map[string][]models.AnomalyReport)
	for _, anomaly := range tempResults {
		// Вычисляем LiveCount как разницу в днях
		var liveCount int64
		var status string = "Не решено"
		if anomaly.ActionCreatedAt != nil && !anomaly.ActionCreatedAt.IsZero() {
			// Если есть действие, считаем разницу между действием и категоризацией
			duration := anomaly.ActionCreatedAt.Sub(anomaly.CategorizedAt)
			liveCount = int64(duration.Hours() / 24) // Переводим в дни
			// TODO: переделать под Enum
			if anomaly.Action == "allow" {
				status = "Разрешено"
			} else {
				status = "Заблокировано"
			}
		} else {
			// Если действия нет, считаем разницу от текущего времени до конца временного диапазона
			duration := time.Now().UTC().Sub(tr.To)
			liveCount = int64(duration.Hours() / 24) // Переводим в дни
		}

		anomalyReport := models.AnomalyReport{
			URL:       anomaly.URL,
			LiveCount: liveCount,
			Status:    status,
			Traffic: models.Traffic{
				Input:  0, // Заглушка
				Output: 0, // Заглушка
			},
			Stat: models.RequestReport{
				All:         anomaly.Total,
				BeforeBlock: anomaly.BeforeBlock,
				Pending:     anomaly.Pending,
				AfterBlock:  anomaly.AfterBlock,
			},
		}

		deviceAnomalies[anomaly.HostName] = append(deviceAnomalies[anomaly.HostName], anomalyReport)
	}

	// Преобразуем в конечный формат
	var result []models.DeviceAnomaly
	for hostname, anomalies := range deviceAnomalies {
		result = append(result, models.DeviceAnomaly{
			HostName:    hostname,
			AnomalyStat: anomalies,
		})
	}

	return result, nil
}

// GetTopAnomalies возвращает топ аномалий
func (r *RepoPG) GetTopAnomalies(ctx context.Context, tr *trparser.TimeRange) ([]models.DeviceAnomalyAnalytics, error) {
	type deviceAnomalyTemp struct {
		HostName              string `json:"hostname"`
		Requests              uint   `json:"requests"`
		Anomalies             uint   `json:"anomalies"`
		Blocks                uint   `json:"blocks"`
		AllDetections         uint   `json:"all_detections"`
		UnresolvedDetections  uint   `json:"unresolved_detections"`
		BlockedDetections     uint   `json:"blocked_detections"`
		AllowedDetections     uint   `json:"allowed_detections"`
		AnomalyResourcesCount uint   `json:"anomaly_resources_count"`
	}

	var tempResults []deviceAnomalyTemp

	// Получаем основную статистику по устройствам с аномалиями
	query := r.db.GetDB().WithContext(ctx).Table("sessions").
		Select(`
			devices.hostname as hostname,
			COUNT(*) as requests,
			COUNT(CASE WHEN sessions.status = ? THEN 1 END) as anomalies,
			COUNT(CASE WHEN sessions.status = ? THEN 1 END) as blocks,
			COUNT(CASE WHEN domains.category_id IS NOT NULL AND categories.type = ? THEN 1 END) as all_detections,
			COUNT(CASE WHEN domains.category_id IS NOT NULL AND categories.type = ? AND (domains.action_id IS NULL OR domains.action_id = 0) THEN 1 END) as unresolved_detections,
			COUNT(CASE WHEN domains.category_id IS NOT NULL AND categories.type = ? AND domains.action_id > 0 AND actions.action = 'block' THEN 1 END) as blocked_detections,
			COUNT(CASE WHEN domains.category_id IS NOT NULL AND categories.type = ? AND domains.action_id > 0 AND actions.action = 'allow' THEN 1 END) as allowed_detections
		`,
			status.StatusAnomaly.String(),
			status.StatusBlocked.String(),
			pkg.CategoryTypeNegative,
			pkg.CategoryTypeNegative,
			pkg.CategoryTypeNegative,
			pkg.CategoryTypeNegative).
		Joins("LEFT JOIN devices ON sessions.device_id = devices.id").
		Joins("LEFT JOIN domains ON sessions.domain_id = domains.id").
		Joins("LEFT JOIN categories ON domains.category_id = categories.id").
		Joins("LEFT JOIN actions ON domains.action_id = actions.id").
		Where("devices.hostname IS NOT NULL").
		Where("sessions.status = ?", status.StatusAnomaly.String()) // Только устройства с аномалиями

	// Применяем временной диапазон
	if tr != nil && !tr.From.IsZero() && !tr.To.IsZero() {
		query = query.Where("sessions.datetime_utc BETWEEN ? AND ?", tr.From, tr.To)
	}

	err := query.
		Group("devices.hostname").
		Having("COUNT(CASE WHEN sessions.status = ? THEN 1 END) > 0", status.StatusAnomaly.String()). // Используем исходное выражение вместо псевдонима
		Order("anomalies DESC").                                                                      // Сортируем по количеству аномалий
		Find(&tempResults).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get top anomalies: %w", err)
	}

	// Получаем количество уникальных аномальных ресурсов для каждого устройства
	for i := range tempResults {
		var resourcesCount int64
		err := r.db.GetDB().WithContext(ctx).Table("sessions").
			Distinct("domains.path").
			Joins("LEFT JOIN devices ON sessions.device_id = devices.id").
			Joins("LEFT JOIN domains ON sessions.domain_id = domains.id").
			Joins("LEFT JOIN categories ON domains.category_id = categories.id").
			Where("devices.hostname = ?", tempResults[i].HostName).
			Where("sessions.status = ?", status.StatusAnomaly.String()).
			Where("categories.type = ?", pkg.CategoryTypeNegative).
			Count(&resourcesCount).Error

		if err != nil {
			r.l.WarnContext(ctx, "failed to count anomaly resources", slog.String("hostname", tempResults[i].HostName), slog.Any("error", err))
			tempResults[i].AnomalyResourcesCount = 0
		} else {
			tempResults[i].AnomalyResourcesCount = uint(resourcesCount)
		}
	}

	// Преобразуем результат в нужный формат
	var result []models.DeviceAnomalyAnalytics
	for _, device := range tempResults {
		result = append(result, models.DeviceAnomalyAnalytics{
			HostName: device.HostName,
			Traffic: models.Traffic{
				Input:  0, // Заглушка
				Output: 0, // Заглушка
			},
			Requests: device.Requests,
			AnomalyBlockStat: models.AnomalyBlockStat{
				Anomalies: device.Anomalies,
				Blocks:    device.Blocks,
				All:       device.Anomalies + device.Blocks,
			},
			Detections: models.DetectionReport{
				All:        device.AllDetections,
				Unresolved: device.UnresolvedDetections,
				Blocked:    device.BlockedDetections,
				Allowed:    device.AllowedDetections,
			},
		})
	}

	return result, nil
}

// GetTopCategoriesForReport возвращает топ категорий для отчета
func (r *RepoPG) GetTopCategoriesForReport(ctx context.Context, tr *trparser.TimeRange, hostname string) ([]models.TopCategory, error) {
	type categoryTemp struct {
		Category    string `json:"category"`
		Total       uint   `json:"total"`
		BeforeBlock uint   `json:"before_block"`
		AfterBlock  uint   `json:"after_block"`
		Pending     uint   `json:"pending"`
	}

	var tempResults []categoryTemp

	query := r.db.GetDB().WithContext(ctx).Table("sessions").
		Select(`
			categories.name as category,
			COUNT(*) as total,
			COUNT(CASE WHEN sessions.datetime_utc < domains.categorized_at THEN 1 END) as before_block,
			COUNT(CASE WHEN sessions.datetime_utc >= domains.categorized_at AND sessions.status != ? THEN 1 END) as after_block,
			COUNT(CASE WHEN sessions.status = ? THEN 1 END) as pending
		`, status.StatusPending.String(), status.StatusPending.String()).
		Joins("LEFT JOIN domains ON sessions.domain_id = domains.id").
		Joins("LEFT JOIN categories ON domains.category_id = categories.id").
		Where("categories.name IS NOT NULL AND categories.name != ''")

	// Применяем временной диапазон
	if tr != nil && !tr.From.IsZero() && !tr.To.IsZero() {
		query = query.Where("sessions.datetime_utc BETWEEN ? AND ?", tr.From, tr.To)
	}

	// Применяем фильтр по hostname если указан
	if hostname != "" {
		query = query.Joins("LEFT JOIN devices ON sessions.device_id = devices.id").
			Where("devices.hostname = ?", hostname)
	}

	err := query.
		Group("categories.name").
		Order("total DESC"). // Сортируем по общему количеству запросов
		Find(&tempResults).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get top categories for report: %w", err)
	}

	// Преобразуем результат в нужный формат
	var result []models.TopCategory
	for _, category := range tempResults {
		result = append(result, models.TopCategory{
			Category: category.Category,
			Traffic: models.Traffic{
				Input:  0, // Заглушка
				Output: 0, // Заглушка
			},
			Stat: models.RequestReport{
				All:         category.Total,
				BeforeBlock: category.BeforeBlock,
				Pending:     category.Pending,
				AfterBlock:  category.AfterBlock,
			},
		})
	}

	return result, nil
}
