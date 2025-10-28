package postresql

import (
	"context"
	"fmt"
	"tg-an/internal/controllers/http/v1/dto"
	models "tg-an/internal/models"
	pkg "tg-an/pkg/models"
	"tg-an/pkg/trparser"
)

// GetTopDetections возвращает список детекций с пагинацией
func (r *RepoPG) GetTopDetections(ctx context.Context, tr *trparser.TimeRange, f models.DetectionFilter, p models.Pagination) ([]dto.Detection, int64, error) {
	var detections []dto.Detection
	var total int64

	// TODO: пока поле "Описание" дублирует категорию
	query := r.db.GetDB().WithContext(ctx).Table("sessions").
		Select(`
			domains.ip,
			domains.port,
			domains.country as location,
			domains.path as domain,
			domains.categorized_at as categorized_at,
			COUNT(*) as request_count,
			devices.host_name,
			categories.name as category,
			categories.name as description,
			actions.action as action
		`).
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
		query = query.Where("devices.host_name = ?", f.HostName)
	}
	if f.TopCategory != "" {
		query = query.Where("categories.name = ?", f.TopCategory)
	}

	// Группируем по уникальным детекциям
	query = query.Group(`
		domains.ip, domains.port, domains.country, domains.path, domains.categorized_at,
		devices.host_name, categories.name, actions.action
	`)
	// Получаем общее количество записей (до пагинации)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count detections: %w", err)
	}

	// Сортируем по количеству доступов (по убыванию) и по времени определения категории
	query = query.Order("request_count DESC, categorized_at DESC")

	// Применяем пагинацию
	if p.Limit > 0 {
		offset := (p.Page - 1) * p.Limit
		if offset < 0 {
			offset = 0
		}
		query = query.Offset(offset).Limit(p.Limit)
	}

	// Выполняем запрос
	if err := query.Find(&detections).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get detections: %w", err)
	}

	return detections, total, nil
}

// GetDetectionStat возвращает статистику по детекциям
func (r *RepoPG) GetDetectionStat(ctx context.Context, tr *trparser.TimeRange, f models.DetectionFilter) (dto.DetectionStat, error) {
	var result dto.DetectionStat
	// Базовый запрос
	query := r.db.GetDB().WithContext(ctx).Table("sessions").
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
		query = query.Where("devices.host_name = ?", f.HostName)
	}
	if f.TopCategory != "" {
		query = query.Where("categories.name = ?", f.TopCategory)
	}

	// Подсчитываем все выявления (уникальные по domains.id)
	detectedQuery := query.Select("COUNT(DISTINCT domains.id)")
	if err := detectedQuery.Count(&result.Detected).Error; err != nil {
		return result, fmt.Errorf("failed to get detected events: %w", err)
	}

	// Подсчитываем разрешенные выявления
	allowedQuery := detectedQuery.Where("actions.action = ?", pkg.ActionTypeAllowed.String())
	if err := allowedQuery.Count(&result.Allowed).Error; err != nil {
		return result, fmt.Errorf("failed to get accepted events: %w", err)
	}

	// Подсчитываем запрещенные выявления
	deniedQuery := detectedQuery.Where("actions.action = ?", pkg.ActionTypeDenied.String())
	if err := deniedQuery.Count(&result.Denied).Error; err != nil {
		return result, fmt.Errorf("failed to get denied events: %w", err)
	}
	// Вычисляем неразрешенны
	result.Unresolved = result.Detected - result.Allowed - result.Denied
	return result, nil
}
