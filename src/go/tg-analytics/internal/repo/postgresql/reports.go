package postresql

import (
	"context"
	"tg-an/internal/models"
	"tg-an/pkg/trparser"
)

// GetCategories возвращает категории для главной страницы отчета
func (r *RepoPG) GetCategories(ctx context.Context, tr *trparser.TimeRange) (map[string]models.RequestReport, error) {
	// TODO: реализовать логику получения категорий
	return map[string]models.RequestReport{}, nil
}

// GetResourses возвращает ресурсы для главной страницы отчета
func (r *RepoPG) GetResourses(ctx context.Context, tr *trparser.TimeRange) (map[string]models.ResourceStat, error) {
	// TODO: реализовать логику получения ресурсов
	return map[string]models.ResourceStat{}, nil
}

// GetDevicesAnalytics возвращает аналитику по устройствам
func (r *RepoPG) GetDevicesAnalytics(ctx context.Context, tr *trparser.TimeRange, hostname string) (models.DevicesAnalyticsPage, error) {
	// TODO: реализовать логику получения аналитики устройств
	return models.DevicesAnalyticsPage{
		Analytics: make(map[string]models.DeviceReport),
	}, nil
}

// GetAnomaliesList возвращает список аномалий
func (r *RepoPG) GetAnomaliesList(ctx context.Context, tr *trparser.TimeRange, hostname string) (map[string]models.AnomaliesListPage, error) {
	// TODO: реализовать логику получения списка аномалий
	return make(map[string]models.AnomaliesListPage), nil
}

// GetTopAnomalies возвращает топ аномалий
func (r *RepoPG) GetTopAnomalies(ctx context.Context, tr *trparser.TimeRange) (models.TopAnomaliesPage, error) {
	// TODO: реализовать логику получения топ аномалий
	return models.TopAnomaliesPage{
		Anomalies: make(map[string]models.TopAnomaly),
	}, nil
}

// GetTopCategoriesForReport возвращает топ категорий для отчета
func (r *RepoPG) GetTopCategoriesForReport(ctx context.Context, tr *trparser.TimeRange, hostname string) (map[string]models.TopCategoriesPage, error) {
	// TODO: реализовать логику получения топ категорий для отчета
	return make(map[string]models.TopCategoriesPage), nil
}
