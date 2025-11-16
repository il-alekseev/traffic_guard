package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"tg-an/internal/controllers/http/v1/dto"
	models "tg-an/internal/models"
	pkg "tg-an/pkg/models"
	"tg-an/pkg/trparser"
	"time"

	"gorm.io/gorm"
)

// GetTopDetections возвращает список детекций с пагинацией
func (r *RepoPG) GetTopDetections(ctx context.Context, tr *trparser.TimeRange, f models.DetectionFilter, action string, p models.Pagination) ([]dto.Detection, int64, error) {
	var detections []dto.Detection
	var total int64

	// TODO: пока поле "Описание" дублирует категорию
	query := r.db.GetDB().WithContext(ctx).Table("sessions").
		Select(fmt.Sprintf(`
			domains.ip,
			domains.port,
			domains.country as location,
			domains.path as domain,
			domains.categorized_at as categorized_at,
			COUNT(*) as request_count,
			devices.hostname as host_name,
			categories.name as category,
			categories.name as description,
			COALESCE(actions.action, '%s') as action
		`, pkg.ActionTypeUnresolved.String())).
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

	// Применяем фильтр по действию
	if action != "" {
		// Для действия "Не решено" ищем записи где actions.action IS NULL
		if action == pkg.ActionTypeUnresolved.String() {
			query = query.Where("actions.action IS NULL")
		} else {
			query = query.Where("actions.action = ?", action)
		}
	}

	// Группируем по уникальным детекциям
	query = query.Group(fmt.Sprintf(`
    domains.ip, domains.port, domains.country, domains.path, domains.categorized_at,
    devices.hostname, categories.name, COALESCE(actions.action, '%s')
`, pkg.ActionTypeUnresolved.String()))
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

	// Структура для результатов подсчета
	type countResult struct {
		Detected int64
		Allowed  int64
		Denied   int64
	}

	var counts countResult

	// Один запрос для подсчета всех статистик
	query := r.db.GetDB().WithContext(ctx).Table("sessions").
		Select(`
			COUNT(DISTINCT domains.id) as detected,
			COUNT(DISTINCT CASE WHEN actions.action = ? THEN domains.id END) as allowed,
			COUNT(DISTINCT CASE WHEN actions.action = ? THEN domains.id END) as denied
		`, pkg.ActionTypeAllowed.String(), pkg.ActionTypeDenied.String()).
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

	// Выполняем запрос
	if err := query.Scan(&counts).Error; err != nil {
		return result, fmt.Errorf("failed to get detection stats: %w", err)
	}

	// Заполняем результат
	result.Detected = counts.Detected
	result.Allowed = counts.Allowed
	result.Denied = counts.Denied
	result.Unresolved = counts.Detected - counts.Allowed - counts.Denied

	return result, nil
}

func (r *RepoPG) Act(ctx context.Context, username, action, path string) error {
	return r.db.GetDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Находим домен по пути
		var domain models.Domain
		err := tx.Where("path = ?", path).First(&domain).Error
		if err != nil {
			return fmt.Errorf("failed to find domain with path %s: %w", path, err)
		}

		// Обрабатываем Action
		detectionType, err := models.ParseDetectionStatus(action)
		if err != nil {
			return fmt.Errorf("incorrect action: %s", action)
		}

		// Получаем общее количество записей action
		var count int64
		if err := tx.Model(&models.Action{}).Count(&count).Error; err != nil {
			return fmt.Errorf("failed to count actions: %w", err)
		}
		// Cоздаем действие
		var actionRecord = models.Action{
			ID:        uint(count + 1),
			Action:    detectionType.String(),
			CreatedAt: time.Now(),
			CreatedBy: username,
			DomainID:  domain.ID,
		}
		if err := tx.Create(&actionRecord).Error; err != nil {
			return fmt.Errorf("failed to create action: %w", err)
		}
		// Обновляем домен
		result := tx.Model(&models.Domain{}).
			Where("id = ?", domain.ID).
			Update("action_id", actionRecord.ID)

		if result.Error != nil {
			return fmt.Errorf("failed to update domain action: %w", result.Error)
		}

		if result.RowsAffected == 0 {
			return fmt.Errorf("no domain was updated")
		}
		return nil
	})
}

func (r *RepoPG) GetDomainAction(ctx context.Context, path string) (string, error) {
	var action sql.NullString

	err := r.db.GetDB().WithContext(ctx).Table("domains").
		Joins("LEFT JOIN actions ON domains.action_id = actions.id").
		Where("domains.path = ?", path).
		Select("actions.action").
		Pluck("actions.action", &action).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", fmt.Errorf("domain with path %s not found", path)
		}
		return "", fmt.Errorf("failed to get domain action: %w", err)
	}

	if !action.Valid {
		return "", nil
	}

	return action.String, nil
}
