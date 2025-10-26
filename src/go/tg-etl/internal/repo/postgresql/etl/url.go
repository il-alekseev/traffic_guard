package postgresql

import (
	"context"
	"fmt"
	"tg-etl/internal/models"
	"tg-etl/pkg/slogger"

	"gorm.io/gorm"
)

// GetURLByPath получает URL по пути
func (r *ELTRepoPG) GetURLByPath(ctx context.Context, path string) (*models.URL, error) {
	var url models.URL

	err := r.db.GetDB().WithContext(ctx).
		Where("path = ?", path).
		First(&url).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		err = fmt.Errorf("failed to get URL by path %s: %w", path, err)
		return nil, slogger.WrapError(ctx, err)
	}

	return &url, nil
}

func (r *ELTRepoPG) GetURLByPathDomain(ctx context.Context, path string, id uint) (*models.URL, error) {
	var url models.URL

	err := r.db.GetDB().WithContext(ctx).
		Where("path = ? and domain_id = ?", path, id).
		First(&url).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		err = fmt.Errorf("failed to get URL by path %s: %w", path, err)
		return nil, slogger.WrapError(ctx, err)
	}

	return &url, nil
}

// CreateURL создает новый URL, если он еще не существует
func (r *ELTRepoPG) CreateURL(ctx context.Context, url models.URL) error {
	err := r.db.GetDB().WithContext(ctx).Create(&url).Error
	if err != nil {
		err = fmt.Errorf("failed to create URL: %w", err)
		return slogger.WrapError(ctx, err)
	}

	return nil
}
