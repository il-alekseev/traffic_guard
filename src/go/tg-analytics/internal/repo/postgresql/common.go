package postgresql

import (
	"context"
	"fmt"
)

func (r *RepoPG) GetDevices(ctx context.Context, hostname string) ([]string, error) {
	var devices []string
	query := r.db.GetDB().WithContext(ctx).Table("devices").
		Select(`hostname`)

	// Если указан конкретный hostname, проверяем его существование
	if hostname != "" {
		query = query.Where("hostname = ?", hostname)
	}

	if err := query.Find(&devices).Error; err != nil {
		return nil, fmt.Errorf("failed to get devices: %w", err)
	}

	return devices, nil
}

func (r *RepoPG) GetContentCategories(ctx context.Context) ([]string, error) {
	var categories []string
	query := r.db.GetDB().WithContext(ctx).Table("categories").
		Select(`name`)

	if err := query.Find(&categories).Error; err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}

	return categories, nil
}
