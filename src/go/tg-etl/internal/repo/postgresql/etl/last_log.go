package postgresql

import (
	"context"
	"fmt"
	"tg-etl/internal/models"
	"tg-etl/pkg/slogger"

	"gorm.io/gorm"
)

func (r *ELTRepoPG) GetLastLog(ctx context.Context) (*models.LastLog, error) {
	var log models.LastLog
	query := r.db.GetDB().WithContext(ctx).Model(&models.LastLog{})
	if err := query.First(&log).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		err = fmt.Errorf("failed to get last log: %w", err)
		return nil, slogger.WrapError(ctx, err)
	}
	return &log, nil
}

func (r ELTRepoPG) CreateOrUpdateLastLog(ctx context.Context, log models.LastLog) error {
	oldLog, err := r.GetLastLog(ctx)
	if err != nil {
		err = fmt.Errorf("failed to get last log: %w", err)
		return slogger.WrapError(ctx, err)
	}
	// Лога еще не было
	if oldLog == nil {
		return r.createLastLog(ctx, log)
	} else {
		return r.updateLastLog(ctx, oldLog.ID, log)
	}
}

func (r ELTRepoPG) createLastLog(ctx context.Context, log models.LastLog) error {

	err := r.db.WithTx(ctx, func(tx *gorm.DB) error {
		if err := tx.Create(&log).Error; err != nil {
			return fmt.Errorf("failed to create last log %d: %w", log.ID, err)
		}
		return nil
	})
	if err != nil {
		return slogger.WrapError(ctx, err)
	}
	return nil
}

func (r ELTRepoPG) updateLastLog(ctx context.Context, id uint, log models.LastLog) error {
	err := r.db.WithTx(ctx, func(tx *gorm.DB) error {
		result := tx.Model(&models.LastLog{}).Where("id = ?", id).Updates(log)
		if result.Error != nil {
			return fmt.Errorf("failed to update last log %d: %w", id, result.Error)
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("domain with id %d not found", id)
		}
		return nil
	})
	if err != nil {
		return slogger.WrapError(ctx, err)
	}
	return nil
}
