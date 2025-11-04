package postresql

import (
	"context"
	"fmt"
	"tg-an/internal/controllers/http/v1/dto"
	"tg-an/internal/models"
	"tg-an/internal/pkg/status"
	pkg "tg-an/pkg/models"
	"tg-an/pkg/trparser"
	"time"
)

// GetTopCategories возвращает топ категорий по количеству доступов
func (r *RepoPG) GetTopCategories(ctx context.Context, tr *trparser.TimeRange, f models.CategoryFilter, count int) ([]dto.Category, error) {
	var categories []dto.Category

	query := r.db.GetDB().WithContext(ctx).Table("sessions").
		Select("categories.name as name, COUNT(*) as access_count").
		Joins("LEFT JOIN domains ON sessions.domain_id = domains.id").
		Joins("LEFT JOIN categories ON domains.category_id = categories.id").
		Where("categories.name IS NOT NULL AND categories.name != 'Неизвестный класс'")

	// Применяем временной диапазон
	if tr != nil && !tr.From.IsZero() && !tr.To.IsZero() {
		query = query.Where("sessions.datetime_utc BETWEEN ? AND ?", tr.From, tr.To)
	}

	// Применяем фильтр по hostname
	if f.HostName != "" {
		query = query.Joins("JOIN devices ON sessions.device_id = devices.id").
			Where("devices.host_name = ?", f.HostName)
	}

	// Применяем фильтр по типу сессии
	if f.Type != "" {
		query = query.Where("sessions.type = ?", f.Type)
	}

	err := query.
		Group("categories.name").
		Order("access_count DESC").
		Limit(count).
		Find(&categories).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get top categories: %w", err)
	}

	return categories, nil
}

// GetRequestsStat возвращает статистику запросов по временным интервалам
func (r *RepoPG) GetRequestStat(ctx context.Context, tr *trparser.TimeRange, hostname, requestType string, count uint) (models.RequestStat, error) {
	result := models.RequestStat{
		Time: make([]time.Time, count),
		Data: make([]uint, count),
	}

	if !tr.IsValid() {
		return result, fmt.Errorf("invalid time range: %s", tr.String())
	}

	if count == 0 {
		return result, fmt.Errorf("count must be greater than 0")
	}

	totalDuration := tr.Duration()
	intervalDuration := totalDuration / time.Duration(count)

	// Создаем временные интервалы заранее
	for i := uint(0); i < count; i++ {
		start := tr.From.Add(time.Duration(i) * intervalDuration)
		result.Time[i] = start.Add(intervalDuration / 2)
	}

	// Получаем все записи используя GORM
	var statRecords []models.Session
	query := r.db.GetDB().WithContext(ctx).Table("sessions").
		Where("sessions.datetime_utc BETWEEN ? AND ?", tr.From, tr.To).
		Joins("LEFT JOIN domains ON sessions.domain_id = domains.id").
		Joins("LEFT JOIN categories ON domains.category_id = categories.id").
		Joins("LEFT JOIN actions ON domains.action_id = actions.id")

	// Применяем фильтр по hostname
	if hostname != "" {
		query = query.Joins("JOIN devices ON sessions.device_id = devices.id").
			Where("devices.host_name = ?", hostname)
	}

	rt, err := models.ParseRequestStatus(requestType)
	if err != nil {
		return result, fmt.Errorf("failed to parse request type: %w", err)
	}

	if rt == models.RequestStatusBeforeBlock {
		query = query.Where("sessions.status = ? or sessions.status = ?", models.RequestStatusPending.String(), status.StatusAnomaly.String())
	} else {
		query = query.Where("sessions.status = ?", rt.String())
	}
	if err := query.Find(&statRecords).Error; err != nil {
		return result, fmt.Errorf("failed to get request stat: %w", err)
	}

	// Обрабатываем каждую запись и распределяем по интервалам
	for _, record := range statRecords {
		// Определяем индекс интервала
		timeDiff := record.DatetimeUTC.Sub(tr.From)
		intervalIndex := int(timeDiff / intervalDuration)

		if intervalIndex >= 0 && intervalIndex < int(count) {
			result.Data[intervalIndex]++
		}
	}

	return result, nil
}

// GetTopUnresolvedDetections возвращает топ нерешенных выявлений по числу запросов
func (r *RepoPG) GetTopUnresolvedDetections(ctx context.Context, tr *trparser.TimeRange, hostName string, count int) ([]dto.UnresolvedDetection, error) {
	var detections []dto.UnresolvedDetection

	query := r.db.GetDB().WithContext(ctx).Table("sessions").
		Select(`
			domains.path as domain,
			COUNT(*) as requests_all,
			COUNT(CASE WHEN sessions.datetime_utc < domains.categorized_at THEN 1 END) as requests_before,
			COUNT(CASE WHEN sessions.datetime_utc >= domains.categorized_at THEN 1 END) as requests_after
		`).
		Joins("LEFT JOIN domains ON sessions.domain_id = domains.id").
		Joins("LEFT JOIN categories ON domains.category_id = categories.id").
		Where("categories.type = ?", pkg.CategoryTypeNegative).       // признак выявления
		Where("(domains.action_id IS NULL OR domains.action_id = 0)") // признак того, что выявление нерешенное

	// Применяем временной диапазон
	if tr != nil && !tr.To.IsZero() {
		query = query.Where("sessions.datetime_utc BETWEEN ? AND ?", tr.From, tr.To)
	}

	// Применяем фильтр по hostname
	if hostName != "" {
		query = query.Joins("JOIN devices ON sessions.device_id = devices.id").
			Where("devices.host_name = ?", hostName)
	}

	err := query.
		Group("domains.path").
		Order("requests_after DESC, requests_before DESC").
		Limit(count).
		Find(&detections).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get top unresolved detections: %w", err)
	}

	return detections, nil
}
