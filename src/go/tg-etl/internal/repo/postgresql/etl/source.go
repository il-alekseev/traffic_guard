package postgresql

import (
	"context"
	"fmt"
	"tg-etl/internal/models"
	"tg-etl/pkg/slogger"

	"gorm.io/gorm"
)

func (r *ELTRepoPG) GetSourceByAddr(ctx context.Context, ip string) (*models.Source, error) {
	var source models.Source
	err := r.db.GetDB().WithContext(ctx).Where("ip = ?", ip).First(&source).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		err = fmt.Errorf("failed to get source by addr %s: %w", ip, err)
		return nil, slogger.WrapError(ctx, err)
	}
	return &source, nil
}

func (r *ELTRepoPG) GetSourceByID(ctx context.Context, id uint) (*models.Source, error) {
	var source models.Source
	err := r.db.GetDB().WithContext(ctx).First(&source, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		err = fmt.Errorf("failed to get source by id %d: %w", id, err)
		return nil, slogger.WrapError(ctx, err)
	}
	return &source, nil
}

func (r *ELTRepoPG) CreateSource(ctx context.Context, source models.Source) error {
	err := r.db.WithTx(ctx, func(tx *gorm.DB) error {
		if err := tx.Create(&source).Error; err != nil {
			return fmt.Errorf("failed to create source %s: %w", source.IP, err)
		}
		return nil
	})
	if err != nil {
		return slogger.WrapError(ctx, err)
	}
	return nil
}

func (r *ELTRepoPG) UpdateSource(ctx context.Context, source models.Source) error {
	err := r.db.WithTx(ctx, func(tx *gorm.DB) error {
		result := tx.Save(&source)
		if result.Error != nil {
			return fmt.Errorf("failed to update source %d: %w", source.ID, result.Error)
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("source with id %d not found", source.ID)
		}
		return nil
	})
	if err != nil {
		return slogger.WrapError(ctx, err)
	}
	return nil
}

func (r *ELTRepoPG) GetSources(ctx context.Context) ([]models.Source, error) {
	var sources []models.Source
	err := r.db.GetDB().WithContext(ctx).Find(&sources).Error
	if err != nil {
		err = fmt.Errorf("failed to get sources: %w", err)
		return nil, slogger.WrapError(ctx, err)
	}
	return sources, nil
}

func (r *ELTRepoPG) DeleteSource(ctx context.Context, id uint) error {
	err := r.db.WithTx(ctx, func(tx *gorm.DB) error {
		result := tx.Delete(&models.Source{}, id)
		if result.Error != nil {
			return fmt.Errorf("failed to delete source %d: %w", id, result.Error)
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("source with id %d not found", id)
		}
		return nil
	})
	if err != nil {
		return slogger.WrapError(ctx, err)
	}
	return nil
}
