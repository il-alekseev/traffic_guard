package postresql

import (
	"context"
	"fmt"
	"tg-an/internal/controllers/http/v1/dto"
	"tg-an/internal/models"
	"tg-an/internal/pkg/status"
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
func (r *RepoPG) GetRequestsStat(ctx context.Context, tr *trparser.TimeRange, f models.DashboardFilter, s status.Status, count uint) ([]uint, error) {
	var stats []uint

	// Вычисляем интервал
	duration := tr.To.Sub(tr.From)
	interval := duration / time.Duration(count)

	// Подготавливаем параметры для фильтров
	hostNameFilter := f.HostName != ""
	categoryFilter := f.TopCategory != ""
	statusFilter := s != status.Status(0) // Проверяем, что статус не нулевой

	// Базовый SQL запрос
	sqlQuery := `
		WITH time_buckets AS (
			SELECT generate_series($1::timestamp, $2::timestamp, $3::interval) as bucket_start
		)
		SELECT 
			COALESCE(COUNT(sessions.id), 0) as count
		FROM time_buckets
		LEFT JOIN sessions ON sessions.datetime_utc >= bucket_start AND sessions.datetime_utc < bucket_start + $3::interval
		LEFT JOIN devices ON sessions.device_id = devices.id
		LEFT JOIN domains ON sessions.domain_id = domains.id
		LEFT JOIN decisions ON domains.decision_id = decisions.id
		WHERE 1=1
	`

	// Добавляем условия фильтрации
	args := []interface{}{tr.From, tr.To, interval}
	argCount := 3

	if hostNameFilter {
		argCount++
		sqlQuery += fmt.Sprintf(" AND ($%d OR devices.host_name = $%d)", argCount, argCount+1)
		args = append(args, !hostNameFilter, f.HostName)
		argCount++
	}

	if categoryFilter {
		argCount++
		sqlQuery += fmt.Sprintf(" AND ($%d OR decisions.decision ILIKE $%d)", argCount, argCount+1)
		args = append(args, !categoryFilter, "%"+f.TopCategory+"%")
		argCount++
	}

	if statusFilter {
		argCount++
		sqlQuery += fmt.Sprintf(" AND ($%d OR ", argCount)

		// Используем строковое представление статуса для сравнения
		switch s {
		case status.Allowed:
			sqlQuery += "decisions.decision IN ('allow', 'accept', 'Разрешено'))"
		case status.Blocked:
			sqlQuery += "decisions.decision IN ('block', 'deny', 'blocked', 'Заблокировано'))"
		case status.Prohibited:
			sqlQuery += "decisions.decision IN ('prohibited', 'deny', 'block', 'malware', 'phishing', 'Запрещено'))"
		case status.Waiting:
			sqlQuery += "decisions.decision IS NULL OR decisions.decision = '' OR decisions.decision NOT IN ('allow', 'accept', 'block', 'deny', 'prohibited', 'malware', 'phishing'))"
		default:
			sqlQuery += "true)"
		}
		args = append(args, !statusFilter)
	}

	sqlQuery += " GROUP BY bucket_start ORDER BY bucket_start"

	// Выполняем запрос
	query := r.db.GetDB().WithContext(ctx).Raw(sqlQuery, args...)
	if err := query.Pluck("count", &stats).Error; err != nil {
		return nil, fmt.Errorf("failed to get requests statistics: %w", err)
	}

	// Если количество результатов меньше запрошенного, дополняем нулями
	if len(stats) < int(count) {
		for i := len(stats); i < int(count); i++ {
			stats = append(stats, 0)
		}
	}

	return stats, nil
}
