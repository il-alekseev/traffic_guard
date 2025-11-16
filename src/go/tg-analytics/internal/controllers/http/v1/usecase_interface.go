package v1

import (
	"context"
	"tg-an/internal/controllers/http/v1/dto"
	"tg-an/internal/models"
	"tg-an/pkg/trparser"

	"github.com/go-openapi/runtime"
)

type UseCaseInterface interface {
	// Common
	GetDevices(ctx context.Context, userMeta *models.UserMeta) ([]string, error)
	GetContentCategories(ctx context.Context) ([]string, error)

	// Sessions
	GetSessions(ctx context.Context, userMeta *models.UserMeta, tr *trparser.TimeRange, f models.SessionFilter, search string, count uint, s models.Sorting) ([]dto.Session, int64, error)

	// Dashboards
	GetTopCategories(ctx context.Context, userMeta *models.UserMeta, tr *trparser.TimeRange, f models.CategoryFilter, count int) ([]dto.Category, error)
	GetRequestStat(ctx context.Context, userMeta *models.UserMeta, tr *trparser.TimeRange, hostname, requestType string, count uint) (models.RequestStat, error)
	GetTrafficStat(ctx context.Context, userMeta *models.UserMeta, tr *trparser.TimeRange, hostName string, count uint) (models.TrafficStat, error)
	GetTopUnresolvedDetections(ctx context.Context, userMeta *models.UserMeta, tr *trparser.TimeRange, hostName string, count int) ([]dto.UnresolvedDetection, error)
	GetDeviceStat(ctx context.Context, userMeta *models.UserMeta, tr *trparser.TimeRange, count uint) (dto.DeviceStatResponse, error)
	GetAnomalies(ctx context.Context, userMeta *models.UserMeta, tr *trparser.TimeRange, hostname string) (dto.GetAnomaliesResponse, error)

	// Detections
	GetTopDetections(ctx context.Context, userMeta *models.UserMeta, tr *trparser.TimeRange, f models.DetectionFilter, action string, p models.Pagination) ([]dto.Detection, int64, error)
	GetDetectionStat(ctx context.Context, userMeta *models.UserMeta, tr *trparser.TimeRange, f models.DetectionFilter) (dto.DetectionStat, error)

	// Actions
	Act(ctx context.Context, userMeta *models.UserMeta, authInfo runtime.ClientAuthInfoWriter, action, path string) error

	// Reports
	CreateReport(ctx context.Context, userMeta *models.UserMeta, authInfo runtime.ClientAuthInfoWriter, tr *trparser.TimeRange) (models.Report, error)
	CreateReportForDevice(ctx context.Context, userMeta *models.UserMeta, authInfo runtime.ClientAuthInfoWriter, tr *trparser.TimeRange, hostname string) (models.ReportForDevice, error)
}
