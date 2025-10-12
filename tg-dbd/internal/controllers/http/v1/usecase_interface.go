package v1

import (
	"context"
	"tg-dbd/internal/models"
	"time"
)

type UseCaseInterface interface {
	GetSessions(ctx context.Context,
		start time.Time,
		end time.Time,
		page int,
		limit int,
		filter string,
		orderBy string,
		orderDir string) ([]models.Session, int64, error)

	GetDetections(ctx context.Context,
		start time.Time,
		end time.Time,
		// TODO: спросить про возможные фильтры
		count int,
		filter string,
	) ([]models.Detection, error)

	GetDetectionsStat(ctx context.Context,
		start time.Time,
		end time.Time,
		// TODO: спросить про возможные фильтры
		filter string,
	) (models.DetectionStat, error)

	GetCategories(ctx context.Context,
		start time.Time,
		end time.Time,
		// TODO: спросить про возможные фильтры
		count int,
		filter string,
	) (map[string]int, error)

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
