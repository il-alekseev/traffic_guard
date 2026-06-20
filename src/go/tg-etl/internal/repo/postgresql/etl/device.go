package postgresql

import (
	"context"
	"fmt"
	"tg-etl/internal/models"
	"tg-etl/pkg/slogger"

	"gorm.io/gorm"
)

func (r *ELTRepoPG) CreateDevice(ctx context.Context, d models.Device) error {
	err := r.db.WithTx(ctx, func(tx *gorm.DB) error {
		if err := tx.Create(&d).Error; err != nil {
			return fmt.Errorf("failed to create device: %w", err)
		}
		return nil
	})
	if err != nil {
		return slogger.WrapError(ctx, err)
	}
	return nil
}

func (r *ELTRepoPG) GetDeviceByID(ctx context.Context, id int) (*models.Device, error) {
	var d models.Device
	err := r.db.GetDB().WithContext(ctx).Where("id = ?", id).First(&d).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		err = fmt.Errorf("failed to get device: %w", err)
		return nil, slogger.WrapError(ctx, err)
	}
	return &d, nil
}

func (r *ELTRepoPG) GetDevices(ctx context.Context) ([]models.Device, error) {
	var devices []models.Device
	err := r.db.GetDB().WithContext(ctx).Find(&devices).Error
	if err != nil {
		err = fmt.Errorf("failed to get devices: %w", err)
		return nil, slogger.WrapError(ctx, err)
	}
	return devices, nil
}
