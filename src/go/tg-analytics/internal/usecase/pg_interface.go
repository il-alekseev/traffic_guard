package usecase

import (
	"context"
	"tg-an/internal/controllers/http/v1/dto"
	"tg-an/internal/models"
	"tg-an/pkg/trparser"
)

type RepoPGInterface interface {
	// Common
	GetDevices(ctx context.Context, hostname string) ([]string, error)
	GetContentCategories(ctx context.Context) ([]string, error)
	// Sessions
	GetSessions(ctx context.Context, tr *trparser.TimeRange, f models.SessionFilter, search string, count uint, s models.Sorting) ([]dto.Session, int64, error)
	// Dashboards
	GetTopCategories(ctx context.Context, tr *trparser.TimeRange, f models.CategoryFilter, count int) ([]dto.Category, error)
	GetRequestStat(ctx context.Context, tr *trparser.TimeRange, hostname, requestType string, count uint) (models.RequestStat, error)
	GetTopUnresolvedDetections(ctx context.Context, tr *trparser.TimeRange, hostName string, count int) ([]dto.UnresolvedDetection, error)
	GetDeviceStat(ctx context.Context, tr *trparser.TimeRange, hostname string, count uint) (dto.DeviceStatResponse, error)
	GetAnomalies(ctx context.Context, tr *trparser.TimeRange, hostname string) (dto.GetAnomaliesResponse, error)
	// Detections
	GetTopDetections(ctx context.Context, tr *trparser.TimeRange, f models.DetectionFilter, action string, p models.Pagination) ([]dto.Detection, int64, error)
	GetDetectionStat(ctx context.Context, tr *trparser.TimeRange, f models.DetectionFilter) (dto.DetectionStat, error)
	// Actions
	Act(ctx context.Context, username, action, path string) error
	GetDomainAction(ctx context.Context, path string) (string, error)
	// Reports
	GetCategories(ctx context.Context, tr *trparser.TimeRange) ([]models.CategoryStat, error)
	GetResourses(ctx context.Context, tr *trparser.TimeRange) ([]models.ResourceStat, error)
	GetDevicesAnalytics(ctx context.Context, tr *trparser.TimeRange, hostname string) ([]models.DeviceReport, error)
	GetAnomaliesList(ctx context.Context, tr *trparser.TimeRange, hostname string) ([]models.DeviceAnomaly, error)
	GetTopAnomalies(ctx context.Context, tr *trparser.TimeRange) ([]models.DeviceAnomalyAnalytics, error)
	GetTopCategoriesForReport(ctx context.Context, tr *trparser.TimeRange, hostname string) ([]models.TopCategory, error)
}

type RepoMetricsPGInterface interface {
	GetTrafficStat(ctx context.Context, tr *trparser.TimeRange, hostname string, count uint) (models.TrafficStat, error)
}
