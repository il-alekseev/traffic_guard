package postgresql

import (
	"context"
	"fmt"
	"tg-etl/internal/models"
	"tg-etl/pkg/slogger"

	"gorm.io/gorm"
)

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
