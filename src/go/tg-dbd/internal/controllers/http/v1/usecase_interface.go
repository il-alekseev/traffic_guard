package v1

import (
	"context"
	"tg-dbd/internal/controllers/http/v1/dto"
	"tg-dbd/internal/models"
	"tg-dbd/pkg/trparser"
	"time"
)

type UseCaseInterface interface {
	// Сессии
	GetSessions(ctx context.Context, tr *trparser.TimeRange, f models.SessionFilter, search string, p models.Pagination, s models.Sorting) ([]dto.Session, int64, error)
	// Dashboard
	GetTopCategories(ctx context.Context, tr *trparser.TimeRange, f models.CategoryFilter, count int) ([]dto.Category, error)
	GetTopDetections(ctx context.Context, tr *trparser.TimeRange, f models.DetectionFilter, p models.Pagination) ([]dto.Detection, int64, error)
	GetDetectionStat(ctx context.Context, tr *trparser.TimeRange, f models.DetectionFilter) (dto.DetectionStat, error)
	// TODO:  Переделать

	GetResources(ctx context.Context,
		start time.Time,
		end time.Time,
		// TODO: спросить про возможные фильтры
		count int,
		filter string,
	) ([]models.Resource, error)

	// TODO: Логика добавления в уведомления бизнес-логов
	GetEvents(ctx context.Context,
		start time.Time,
		end time.Time,
		// TODO: спросить про возможные фильтры
		count int,
		filter string,
	) ([]models.Notification, error)

	GetDevicesStat(ctx context.Context,
		start time.Time,
		end time.Time,
	) ([]models.DeviceStat, error)

	GetTrafficStat(ctx context.Context,
		start time.Time,
		end time.Time,
		filter string,
	) ([]models.TrafficPoint, error)

	GetRequestsStat(ctx context.Context,
		start time.Time,
		end time.Time,
		filter string,
	) (models.RequestsStat, error)

	GetProhActSchedule(ctx context.Context,
		start time.Time,
		filter string,
	) (map[int64]int, error)

	GetAnomalies(ctx context.Context,
		start time.Time,
		end time.Time,
	) (models.Anomaly, error)
}
