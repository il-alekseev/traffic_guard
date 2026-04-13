package postgresql

import (
	"context"
	"fmt"
	"tg-etl/internal/models"
	"tg-etl/pkg/slogger"

	"gorm.io/gorm"
)

func (r *ELTRepoPG) GetActionByDomainID(ctx context.Context, id uint) (*models.Action, error) {
	var action models.Action
	err := r.db.GetDB().WithContext(ctx).Where("domain_id = ?", id).First(&action).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		err = fmt.Errorf("failed to get action by domain id %d: %w", id, err)
		return nil, slogger.WrapError(ctx, err)
	}
	return &action, nil
}
