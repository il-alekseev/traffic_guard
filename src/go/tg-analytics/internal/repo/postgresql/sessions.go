package postgresql

import (
	"context"
	"fmt"
	"strings"
	"tg-an/internal/controllers/http/v1/dto"
	"tg-an/internal/models"
	"tg-an/pkg/trparser"
)

// GetSessions возвращает список сессий с пагинацией и фильтрацией
func (r *RepoPG) GetSessions(ctx context.Context, tr *trparser.TimeRange, f models.SessionFilter, search string, count uint, s models.Sorting) ([]dto.Session, int64, error) {
	var sessions []dto.Session
	var total int64

	// Базовый запрос с джойнами и явным SELECT
	query := r.db.GetDB().WithContext(ctx).Table("sessions").
		Select(`
			sessions.id,
			sessions.datetime_utc,
			sessions.type,
			sessions.status,
			urls.path as url,
			urls.proto,
			devices.hostname as host_name,
			sources.ip as src_ip,
			sources.country as src_country,
			sources.username,
			domains.ip as dst_ip,
			domains.port as dst_port,
			domains.country as dst_country,
			categories.name as category
		`).
		Joins("LEFT JOIN devices ON sessions.device_id = devices.id").
		Joins("LEFT JOIN sources ON sessions.src_id = sources.id").
		Joins("LEFT JOIN domains ON sessions.domain_id = domains.id").
		Joins("LEFT JOIN urls ON domains.id = urls.domain_id").
		Joins("LEFT JOIN categories ON domains.category_id = categories.id")

	// Применяем временной диапазон
	if tr != nil && !tr.From.IsZero() && !tr.To.IsZero() {
		query = query.Where("sessions.datetime_utc BETWEEN ? AND ?", tr.From, tr.To)
	}

	// Применяем фильтры
	if f.HostName != "" {
		query = query.Where("devices.hostname = ?", f.HostName)
	}
	if f.Category != "" {
		query = query.Where("categories.name = ?", f.Category)
	}
	if f.Type != "" {
		query = query.Where("sessions.type = ?", f.Type)
	}

	// Применяем пользовательский поиск
	if search != "" {
		searchPattern := "%" + strings.ToLower(search) + "%"
		query = query.Where(
			"LOWER(urls.path) LIKE ? OR LOWER(sources.ip) LIKE ? OR LOWER(domains.ip) LIKE ?",
			searchPattern, searchPattern,
		)
	}

	// Получаем общее количество записей
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count sessions: %w", err)
	}

	// Применяем сортировку
	if s.OrderBy != "" {
		orderClause := s.OrderBy
		if s.OrderDir != "" {
			orderClause += " " + s.OrderDir
		} else {
			orderClause += " DESC"
		}
		query = query.Order(orderClause)
	} else {
		// Сортировка по умолчанию
		query = query.Order("sessions.datetime_utc DESC")
	}

	// Уменьшаем число сессий до значения count
	query = query.Limit(int(count))

	// Выполняем запрос
	if err := query.Find(&sessions).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get sessions: %w", err)
	}

	return sessions, total, nil
}
