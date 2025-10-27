package postresql

import (
	"context"
	"fmt"
	"strings"
	"tg-dbd/internal/controllers/http/v1/dto"
	"tg-dbd/internal/models"
	"tg-dbd/pkg/trparser"
)

// GetSessions возвращает список сессий с пагинацией и фильтрацией
func (r *RepoPG) GetSessions(ctx context.Context, tr *trparser.TimeRange, f models.SessionFilter, search string, p models.Pagination, s models.Sorting) ([]dto.Session, int64, error) {
	var sessions []dto.Session
	var total int64

	// Базовый запрос с джойнами и явным SELECT
	query := r.db.GetDB().WithContext(ctx).Table("sessions").
		Select(`
			sessions.id,
			sessions.datetime_utc,
			sessions.type,
			sessions.status,
			sessions.url,
			sessions.proto,
			devices.host_name,
			sources.ip as src_ip,
			sources.port as src_port,
			sources.country as src_country,
			sources.username,
			domains.ip as dst_ip,
			domains.port as dst_port,
			domains.country as dst_country,
			decisions.decision as category
		`).
		Joins("LEFT JOIN devices ON sessions.device_id = devices.id").
		Joins("LEFT JOIN sources ON sessions.src_id = sources.id").
		Joins("LEFT JOIN domains ON sessions.domain_id = domains.id").
		Joins("LEFT JOIN decisions ON domains.decision_id = decisions.id")

	// Применяем временной диапазон
	if tr != nil && !tr.From.IsZero() && !tr.To.IsZero() {
		query = query.Where("sessions.datetime_utc BETWEEN ? AND ?", tr.From, tr.To)
	}

	// Применяем фильтры
	if f.HostName != "" {
		query = query.Where("devices.host_name ILIKE ?", "%"+f.HostName+"%")
	}
	if f.Category != "" {
		query = query.Where("decisions.decision ILIKE ?", "%"+f.Category+"%")
	}
	if f.Type != "" {
		query = query.Where("sessions.type = ?", f.Type)
	}

	// Применяем пользовательский поиск
	if search != "" {
		searchPattern := "%" + strings.ToLower(search) + "%"
		query = query.Where(
			"LOWER(sessions.url) LIKE ? OR LOWER(sources.username) LIKE ?",
			searchPattern, searchPattern,
		)
	}

	// Получаем общее количество записей (до пагинации)
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

	// Применяем пагинацию
	if p.Limit > 0 {
		offset := (p.Page - 1) * p.Limit
		if offset < 0 {
			offset = 0
		}
		query = query.Offset(offset).Limit(p.Limit)
	}

	// Выполняем запрос
	if err := query.Find(&sessions).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get sessions: %w", err)
	}

	return sessions, total, nil
}
