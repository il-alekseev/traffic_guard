package postgresql

import (
	"context"
	"fmt"
	"tg-etl/internal/models"
	"tg-etl/pkg/slogger"

	"gorm.io/gorm"
)

func (r *ELTRepoPG) CreateCategories(ctx context.Context, categories []models.Category) error {
	if len(categories) == 0 {
		return nil
	}
	err := r.db.GetDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Пакетная вставка категорий
		result := tx.Create(&categories)
		if result.Error != nil {
			return fmt.Errorf("failed to create categories: %w", result.Error)
		}
		return nil
	})
	if err != nil {
		return slogger.WrapError(ctx, err)
	}
	return nil
}

func (r *ELTRepoPG) GetCategoryByID(ctx context.Context, id uint) (*models.Category, error) {
	var category models.Category
	err := r.db.GetDB().WithContext(ctx).First(&category, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		err = fmt.Errorf("failed to get category by id %d: %w", id, err)
		return nil, slogger.WrapError(ctx, err)
	}
	return &category, nil
}

func (r *ELTRepoPG) GetCategoryByName(ctx context.Context, name string) (*models.Category, error) {
	var category models.Category
	err := r.db.GetDB().WithContext(ctx).Where("name = ?", name).First(&category).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		err = fmt.Errorf("failed to get category by name %s: %w", name, err)
		return nil, slogger.WrapError(ctx, err)
	}
	return &category, nil
}

func (r *ELTRepoPG) GetCategories(ctx context.Context) ([]models.Category, error) {
	var categories []models.Category
	err := r.db.GetDB().WithContext(ctx).Find(&categories).Error
	if err != nil {
		err = fmt.Errorf("failed to get categories: %w", err)
		return nil, slogger.WrapError(ctx, err)
	}
	return categories, nil
}
