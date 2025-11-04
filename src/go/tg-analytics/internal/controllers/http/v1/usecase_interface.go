package v1

import (
	"context"
	"tg-an/internal/controllers/http/v1/dto"
	"tg-an/internal/models"
	"tg-an/pkg/trparser"
	"time"
)

type UseCaseInterface interface {
	// Common
	GetDevices(ctx context.Context) ([]string, error)
	GetContentCategories(ctx context.Context) ([]string, error)
	// Sessions
	GetSessions(ctx context.Context, tr *trparser.TimeRange, f models.SessionFilter, search string, p models.Pagination, s models.Sorting) ([]dto.Session, int64, error)
	// Dashboards
	GetTopCategories(ctx context.Context, tr *trparser.TimeRange, f models.CategoryFilter, count int) ([]dto.Category, error)
	GetRequestStat(ctx context.Context, tr *trparser.TimeRange, hostname, requestType string, count uint) (models.RequestStat, error)
	GetTrafficStat(ctx context.Context, tr *trparser.TimeRange, hostName string, count uint) (models.TrafficStat, error)
	GetTopUnresolvedDetections(ctx context.Context, tr *trparser.TimeRange, hostName string, count int) ([]dto.UnresolvedDetection, error)
	// Detections
	GetTopDetections(ctx context.Context, tr *trparser.TimeRange, f models.DetectionFilter, action string, p models.Pagination) ([]dto.Detection, int64, error)
	GetDetectionStat(ctx context.Context, tr *trparser.TimeRange, f models.DetectionFilter) (dto.DetectionStat, error)
	// TODO:  Переделать

	GetDevicesStat(ctx context.Context,
		start time.Time,
		end time.Time,
	) ([]models.DeviceStat, error)

	GetProhActSchedule(ctx context.Context,
		start time.Time,
		filter string,
	) (map[int64]int, error)

	GetAnomalies(ctx context.Context,
		start time.Time,
		end time.Time,
	) (models.Anomaly, error)
}
