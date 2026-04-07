package postgresRepo

import (
	"context"
	"fiermon-blog/internal/models"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// GetRecord - Получение записей из таблицы BusinessLog c фильтрацией и пагинацией
func (p *PostgresDB) GetRecord(ctx context.Context, userMeta *models.UserMeta, page, limit int, role, contextID, search string) ([]models.BusinessLog, models.Meta, error) {
	businessLog := []models.BusinessLog{}
	var total int64

	offset := (page - 1) * limit
	err := p.db.Transaction(func(tx *gorm.DB) error {
		query := p.db.Model(&models.BusinessLog{})

		// Добавляем условия доступа
		if userMeta != nil && userMeta.ContextID != "" {
			ca := fmt.Sprintf("\"role\": \"CA-%s\"", userMeta.ContextID)
			co := fmt.Sprintf("\"role\": \"CO-%s\"", userMeta.ContextID)

			query = query.Where(
				p.db.Where("context = ?", userMeta.ContextID).
					Or("(entity = ? AND entity_id = ?)", "context", userMeta.ContextID).
					Or("(entity = ? AND (old_value::text LIKE ? OR old_value::text LIKE ? OR new_value::text LIKE ? OR new_value::text LIKE ?))",
						"user", "%"+co+"%", "%"+ca+"%", "%"+co+"%", "%"+ca+"%"),
			)
		}

		// Добавляем дополнительные фильтры
		if role != "" {
			query = query.Where("user_role = ?", role)
		}

		if contextID != "" {
			query = query.Where("context = ?", contextID)
		}

		if search != "" {
			searchTerm := "%" + strings.ToLower(search) + "%"
			query = query.Where(
				p.db.Where("LOWER(user_name) LIKE ?", searchTerm).
					Or("LOWER(entity) LIKE ?", searchTerm).
					//Or("LOWER(event_type) LIKE ?", searchTerm).
					//Or("LOWER(entity_id) LIKE ?", searchTerm).
					//Or("LOWER(context_str) LIKE ?", searchTerm).
					//Or("LOWER((old_value #>> '{}')) LIKE ?", searchTerm).
					//Or("LOWER((new_value #>> '{}')) LIKE ?", searchTerm).
					Or("LOWER(description) LIKE ?", searchTerm),
			)
		}

		if err := query.Count(&total).Error; err != nil {
			return err
		}

		err := query.Order("timestamp DESC").Limit(limit).Offset(offset).Find(&businessLog).Error
		if err != nil {
			return err
		}

		return nil
	})

	pages := int((total + int64(limit) - 1) / int64(limit))

	meta := models.Meta{
		Limit: limit,
		Page:  page,
		Pages: pages,
		Total: int(total),
	}

	return businessLog, meta, err
}
