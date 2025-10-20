package postresql

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"tg-dbd/internal/controllers/http/v1/dto"
	"tg-dbd/internal/models"
	"tg-dbd/pkg/pgorm/pgorm"
	"tg-dbd/pkg/trparser"
)

type RepoPG struct {
	db pgorm.Interface
	l  slog.Logger
}

func New(db pgorm.Interface, l slog.Logger) *RepoPG {
	return &RepoPG{
		db: db,
		l:  l,
	}
}

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

// GetTopCategories возвращает топ категорий по количеству доступов
func (r *RepoPG) GetTopCategories(ctx context.Context, tr *trparser.TimeRange, f models.CategoryFilter, count int) ([]dto.Category, error) {
	var categories []dto.Category

	query := r.db.GetDB().WithContext(ctx).Table("sessions").
		Select("decisions.decision as name, COUNT(*) as access_count").
		Joins("LEFT JOIN domains ON sessions.domain_id = domains.id").
		Joins("LEFT JOIN decisions ON domains.decision_id = decisions.id").
		Where("decisions.decision IS NOT NULL AND decisions.decision != ''")

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
		Group("decisions.decision").
		Order("access_count DESC").
		Limit(count).
		Find(&categories).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get top categories: %w", err)
	}

	return categories, nil
}

// GetTopDetections возвращает список детекций с пагинацией
func (r *RepoPG) GetTopDetections(ctx context.Context, tr *trparser.TimeRange, f models.DetectionFilter, p models.Pagination) ([]dto.Detection, int64, error) {
	var detections []dto.Detection
	var total int64

	// Базовый запрос с джойнами
	query := r.db.GetDB().WithContext(ctx).Table("sessions").
		Select(`
			domains.ip,
			domains.port,
			domains.country,
			domains.path as url,
			COUNT(*) as access_count,
			devices.host_name,
			decisions.decision as category,
			domains.path as description,
			decisions.decision,
			MAX(sessions.datetime_utc) as last_access_datetime
		`).
		Joins("LEFT JOIN devices ON sessions.device_id = devices.id").
		Joins("LEFT JOIN domains ON sessions.domain_id = domains.id").
		Joins("LEFT JOIN decisions ON domains.decision_id = decisions.id").
		Where("domains.ip IS NOT NULL AND domains.ip != ''")

	// Применяем временной диапазон
	if tr != nil && !tr.From.IsZero() && !tr.To.IsZero() {
		query = query.Where("sessions.datetime_utc BETWEEN ? AND ?", tr.From, tr.To)
	}

	// Применяем фильтры
	if f.HostName != "" {
		query = query.Where("devices.host_name ILIKE ?", "%"+f.HostName+"%")
	}
	if f.TopCategory != "" {
		query = query.Where("decisions.decision ILIKE ?", "%"+f.TopCategory+"%")
	}

	// Группируем по уникальным детекциям
	query = query.Group(`
		domains.ip, 
		domains.port, 
		domains.country, 
		domains.path, 
		devices.host_name, 
		decisions.decision
	`)

	// Получаем общее количество записей (до пагинации)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count detections: %w", err)
	}

	// Сортируем по количеству доступов (по убыванию) и последнему доступу
	query = query.Order("access_count DESC, last_access_datetime DESC")

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
	var stat dto.DetectionStat

	// Базовый запрос
	query := r.db.GetDB().WithContext(ctx).Table("sessions").
		Joins("LEFT JOIN devices ON sessions.device_id = devices.id").
		Joins("LEFT JOIN domains ON sessions.domain_id = domains.id").
		Joins("LEFT JOIN decisions ON domains.decision_id = decisions.id").
		Where("domains.ip IS NOT NULL AND domains.ip != ''")

	// Применяем временной диапазон
	if tr != nil && !tr.From.IsZero() && !tr.To.IsZero() {
		query = query.Where("sessions.datetime_utc BETWEEN ? AND ?", tr.From, tr.To)
	}

	// Применяем фильтры
	if f.HostName != "" {
		query = query.Where("devices.host_name ILIKE ?", "%"+f.HostName+"%")
	}
	if f.TopCategory != "" {
		query = query.Where("decisions.decision ILIKE ?", "%"+f.TopCategory+"%")
	}

	// Подсчитываем статистику за один запрос
	var result struct {
		Total      int64
		Accepted   int64
		Denied     int64
		Unresolved int64
	}

	err := query.
		Select(`
			COUNT(*) as total,
			SUM(CASE 
				WHEN LOWER(decisions.decision) LIKE '%accept%' OR 
					 LOWER(decisions.decision) LIKE '%allow%' THEN 1 
				ELSE 0 
			END) as accepted,
			SUM(CASE 
				WHEN LOWER(decisions.decision) LIKE '%deny%' OR 
					 LOWER(decisions.decision) LIKE '%block%' OR 
					 LOWER(decisions.decision) LIKE '%malware%' THEN 1 
				ELSE 0 
			END) as denied,
			SUM(CASE 
				WHEN LOWER(decisions.decision) LIKE '%unresolved%' OR 
					 LOWER(decisions.decision) LIKE '%unknown%' OR
					 decisions.decision IS NULL OR 
					 decisions.decision = '' THEN 1 
				ELSE 0 
			END) as unresolved
		`).
		Scan(&result).Error

	if err != nil {
		return stat, fmt.Errorf("failed to get detection statistics: %w", err)
	}

	stat = dto.DetectionStat{
		Detected:   int(result.Total),
		Accepted:   int(result.Accepted),
		Denied:     int(result.Denied),
		Unresolved: int(result.Unresolved),
	}

	return stat, nil
}
