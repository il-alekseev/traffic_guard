package postgresql

import (
	"context"
	"tg-an/internal/models"
	"tg-an/pkg/trparser"
)

// GetCategories возвращает категории для главной страницы отчета
func (r *RepoPG) GetCategories(ctx context.Context, tr *trparser.TimeRange) ([]models.CategoryStat, error) {
	// TODO: реализовать логику получения категорий
	return []models.CategoryStat{}, nil
}

// GetResourses возвращает ресурсы для главной страницы отчета
func (r *RepoPG) GetResourses(ctx context.Context, tr *trparser.TimeRange) ([]models.ResourceStat, error) {
	// TODO: реализовать логику получения ресурсов
	return []models.ResourceStat{}, nil
}

// GetDevicesAnalytics возвращает аналитику по устройствам
func (r *RepoPG) GetDevicesAnalytics(ctx context.Context, tr *trparser.TimeRange, hostname string) ([]models.DeviceReport, error) {
	// TODO: реализовать логику получения аналитики устройств
	return []models.DeviceReport{}, nil
}

// GetAnomaliesList возвращает список аномалий
func (r *RepoPG) GetAnomaliesList(ctx context.Context, tr *trparser.TimeRange, hostname string) ([]models.DeviceAnomaly, error) {
	// TODO: реализовать логику получения списка аномалий
	return []models.DeviceAnomaly{}, nil
}

// GetTopAnomalies возвращает топ аномалий
func (r *RepoPG) GetTopAnomalies(ctx context.Context, tr *trparser.TimeRange) ([]models.DeviceAnomalyAnalytics, error) {
	// TODO: реализовать логику получения топ аномалий
	return []models.DeviceAnomalyAnalytics{}, nil
}

// GetTopCategoriesForReport возвращает топ категорий для отчета
func (r *RepoPG) GetTopCategoriesForReport(ctx context.Context, tr *trparser.TimeRange, hostname string) ([]models.TopCategory, error) {
	// TODO: реализовать логику получения топ категорий для отчета
	return []models.TopCategory{}, nil
}
