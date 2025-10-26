package postgresql

import (
	"context"
	"fmt"
	"tg-etl/internal/models"
	"tg-etl/pkg/slogger"

	"gorm.io/gorm"
)

// AddDomainToList добавляет домен в указанный список (whitelist/blacklist)
func (r *ELTRepoPG) AddDomainToList(ctx context.Context, domain models.Domain, list models.ListType) error {
	// Проверяем валидность типа списка
	if !list.IsValid() {
		return fmt.Errorf("invalid list type: %s", list)
	}

	domainList := models.DomainControlLists{
		DomainID: domain.ID,
		Type:     list.String(),
	}

	err := r.db.GetDB().WithContext(ctx).
		Create(&domainList).Error

	if err != nil {
		err = fmt.Errorf("failed to add domain %d to list %s: %w", domain.ID, list, err)
		return slogger.WrapError(ctx, err)
	}

	return nil
}

// GetListByDomainID получает тип списка (whitelist/blacklist) по ID домена
func (r *ELTRepoPG) GetListByDomainID(ctx context.Context, id uint) (*string, error) {
	var domainList models.DomainControlLists

	err := r.db.GetDB().WithContext(ctx).
		Model(&models.DomainControlLists{}).
		Select("type").
		Where("domain_id = ?", id).
		First(&domainList).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		err = fmt.Errorf("failed to get list by domain id %d: %w", id, err)
		return nil, slogger.WrapError(ctx, err)
	}

	return &domainList.Type, nil
}
