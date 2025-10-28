package usecase

import (
	"context"
	"tg-an/internal/controllers/http/v1/dto"
	"tg-an/internal/models"
	"tg-an/internal/pkg/status"
	"tg-an/pkg/trparser"
)

type RepoPGInterface interface {
	// Common
	GetDevices(ctx context.Context) ([]string, error)
	GetContentCategories(ctx context.Context) ([]string, error)
	// Sessions
	GetSessions(ctx context.Context, tr *trparser.TimeRange, f models.SessionFilter, search string, p models.Pagination, s models.Sorting) ([]dto.Session, int64, error)
	// Dashboards
	GetTopCategories(ctx context.Context, tr *trparser.TimeRange, f models.CategoryFilter, count int) ([]dto.Category, error)
	GetRequestsStat(ctx context.Context, tr *trparser.TimeRange, f models.DashboardFilter, s status.Status, count uint) ([]uint, error)
	// Detections
	GetTopDetections(ctx context.Context, tr *trparser.TimeRange, f models.DetectionFilter, action string, p models.Pagination) ([]dto.Detection, int64, error)
	GetDetectionStat(ctx context.Context, tr *trparser.TimeRange, f models.DetectionFilter) (dto.DetectionStat, error)
}

type RepoMetricsPGInterface interface {
	GetTrafficStat(ctx context.Context, tr *trparser.TimeRange, hostname string, count uint) (models.TrafficStat, error)
}
