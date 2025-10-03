package v1

import (
	"context"
	"tg-dbd/internal/models"
)

type UseCaseInterface interface {
	GetRequestStat(dateStart string,
		dateStop string,
		// TODO: спросить про возможные фильтры
		filter string) (models.RequestStat, error)
	GetNotifications(сtx context.Context,
		dateStart string,
		dateStop string,
		// TODO: спросить про возможные фильтры
		filter string,
		topCount int64) ([]models.Notification, error)
	GetTopCategories(сtx context.Context,
		dateStart string,
		dateStop string,
		// TODO: спросить про возможные фильтры
		filter string,
		topCount int64) (map[string]int, error)
	GetTopResources(сtx context.Context,
		dateStart string,
		dateStop string,
		// TODO: спросить про возможные фильтры
		filter string,
		topCount int64) (map[string]int, error)
	GetDetectionCount(сtx context.Context,
		dateStart string,
		dateStop string,
		// TODO: спросить про возможные фильтры
		filter string) (int64, error)
	GetAcceptCount(сtx context.Context,
		dateStart string,
		dateStop string,
		// TODO: спросить про возможные фильтры
		filter string) (int64, error)
	GetDenyCount(сtx context.Context,
		dateStart string,
		dateStop string,
		// TODO: спросить про возможные фильтры
		filter string) (int64, error)
	GetUnresolvedDetectionCount(сtx context.Context,
		dateStart string,
		dateStop string,
		// TODO: спросить про возможные фильтры
		filter string) (int64, error)
	GetTopDetectionsInfo(сtx context.Context,
		dateStart string,
		dateStop string,
		// TODO: спросить про возможные фильтры
		filter string,
		topCount int64) ([]models.Detection, error)
	GetSessions(ctx context.Context,
		dateStart string,
		dateStop string,
		page int,
		limit int,
		filter string,
		orderBy string,
		orderDir string) ([]models.Session, int64, error)
}
