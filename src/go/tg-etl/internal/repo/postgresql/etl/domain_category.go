package postgresql

import (
	"context"
	"fmt"
	"strings"
	"tg-etl/internal/models"
	"tg-etl/pkg/slogger"

	"gorm.io/gorm"
)

// GetNegativeCategoriesByDomainID Получает список негативных категорий в виде строки
func (r *ELTRepoPG) GetNegativeCategoriesByDomainID(ctx context.Context, id uint) (string, error) {
	var categories []string
	err := r.db.GetDB().WithContext(ctx).
		Table("categories").
		Select("categories.name").
		Joins("JOIN domain_categories ON domain_categories.category_id = categories.id").
		Where("domain_categories.domain_id = ? AND categories.type = ?", id, models.CategoryTypeNegative.String()).
		Find(&categories).Error
	if err != nil {
		err = fmt.Errorf("failed to get categories: %w", err)
		return "", slogger.WrapError(ctx, err)
	}
	res := strings.Join(categories, " ")
	return res, nil
}

// GetNegativeCategoriesStatByDomainID получает статистику отрицательных категорий домена, отсортированные в порядке убывания %
func (r *ELTRepoPG) GetNegativeCategoriesStatByDomainID(ctx context.Context, id uint) (map[string]float32, error) {
	type CategoryCount struct {
		Name         string
		CategoryType string
		Count        float64
	}
	var categoryCounts []CategoryCount
	var totalCount float64
	// Получаем все категории домена с числом их встречаемости, отсортированные в порядке убывания встречаемости
	err := r.db.GetDB().WithContext(ctx).
		Table("domain_categories").
		Select("categories.name, categories.type, domain_categories.count").
		Joins("JOIN categories ON categories.id = domain_categories.category_id").
		Where("domain_categories.domain_id = ?", id).
		Order("domain_categories.count DESC").
		Scan(&categoryCounts).Error
	if err != nil {
		err = fmt.Errorf("failed to get categories stats for domain %d: %w", id, err)
		return nil, slogger.WrapError(ctx, err)
	}
	// Если нет категорий, возвращаем пустую карту
	if len(categoryCounts) == 0 {
		return make(map[string]float32), nil
	}
	// Считаем общую сумму count всех категорий
	for _, cc := range categoryCounts {
		totalCount += cc.Count
	}
	// Если totalCount == 0, избегаем деления на ноль
	if totalCount == 0 {
		return make(map[string]float32), nil
	}
	// Формируем результат только для отрицательных категорий
	result := make(map[string]float32)
	for _, cc := range categoryCounts {
		if cc.CategoryType == models.CategoryTypeNegative.String() {
			percent := float32(cc.Count / totalCount)
			result[cc.Name] = percent
		}
	}
	return result, nil
}

// GetMostNegativeCategoryByDomainID - получает самую часто встречающуюся негативную категорию домена
func (r *ELTRepoPG) GetMostNegativeCategoryByDomainID(ctx context.Context, id uint) (*models.Category, float32, error) {
	type Result struct {
		Name  string
		Count float32
	}

	var result Result

	err := r.db.GetDB().WithContext(ctx).
		Table("domain_categories").
		Select("categories.name, domain_categories.count").
		Joins("JOIN categories ON categories.id = domain_categories.category_id").
		Where("domain_categories.domain_id = ? AND categories.type = ?", id, models.CategoryTypeNegative.String()).
		Order("domain_categories.count DESC").
		Limit(1).
		Scan(&result).Error

	if err != nil {
		err = fmt.Errorf("failed to get most negative category for domain %d: %w", id, err)
		return nil, 0, slogger.WrapError(ctx, err)
	}
	cat, err := r.GetCategoryByName(ctx, result.Name)
	if err != nil {
		err = fmt.Errorf("failed to get most negative category for domain %d: %w", id, err)
		return nil, 0, slogger.WrapError(ctx, err)
	}
	return cat, result.Count, nil
}

// AddDomainCategory прибавляет 1, если категория уже встречалась, или создает новую категорию
func (r *ELTRepoPG) AddDomainCategory(ctx context.Context, c models.Category, id uint) error {
	err := r.db.WithTx(ctx, func(tx *gorm.DB) error {
		// Используем OnConflict для PostgreSQL
		sql := `
			INSERT INTO domain_categories (domain_id, category_id, count)
			VALUES (?, ?, 1)
			ON CONFLICT (domain_id, category_id) 
			DO UPDATE SET count = domain_categories.count + 1
		`
		err := tx.WithContext(ctx).
			Exec(sql, id, c.ID).Error

		if err != nil {
			return fmt.Errorf("failed to upsert domain_category: %w", err)
		}

		return nil
	})

	if err != nil {
		return slogger.WrapError(ctx, err)
	}
	return nil
}

// GetNegativeCategoriesTotalByDomainID - получает общее число поученных негативных категорий домена с учетом их встречаемости
func (r *ELTRepoPG) GetNegativeCategoriesTotalByDomainID(ctx context.Context, id uint) (int, error) {
	var total int

	err := r.db.GetDB().WithContext(ctx).
		Table("domain_categories").
		Select("COALESCE(SUM(domain_categories.count), 0)").
		Joins("JOIN categories ON categories.id = domain_categories.category_id").
		Where("domain_categories.domain_id = ? AND categories.type = ?", id, models.CategoryTypeNegative.String()).
		Scan(&total).Error

	if err != nil {
		err = fmt.Errorf("failed to get total negative categories for domain %d: %w", id, err)
		return 0, slogger.WrapError(ctx, err)
	}

	return total, nil
}
