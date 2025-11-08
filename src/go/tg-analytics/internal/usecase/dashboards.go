package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"tg-an/internal/controllers/http/v1/dto"
	"tg-an/internal/controllers/http/v1/values"
	"tg-an/internal/models"
	"tg-an/pkg/slogger/wsl"
	"tg-an/pkg/trparser"
)

func (u *Usecase) GetTopCategories(ctx context.Context, userMeta *models.UserMeta, tr *trparser.TimeRange, f models.CategoryFilter, count int) ([]dto.Category, error) {
	method := "GetTopCategories"
	u.l.InfoContext(ctx,
		method,
		slog.Any("time_range", tr),
		slog.Any("filter", f),
		slog.Int("count", count),
	)

	// Устанавливаем значение по умолчанию для count
	if count <= 0 {
		count = 5
	}

	// Проверяем роль пользователя
	// И если она CA, то фильтруем по хосту
	if userMeta.ShortRole == values.ContextAdmin {
		f.HostName = userMeta.ContextID
	}

	categories, err := u.db.GetTopCategories(ctx, tr, f, count)
	if err != nil {
		err = fmt.Errorf("%s: failed to get top categories: %w", method, err)
		u.l.ErrorContext(ctx, "Database operation failed",
			wsl.String("method", method),
			wsl.String("error", err.Error()),
		)
		return nil, err
	}

	u.l.InfoContext(ctx, "Top categories retrieved",
		slog.String("method", method),
		slog.Int("categories_count", len(categories)),
		slog.Int("requested_count", count),
	)
	return categories, nil
}

// GetRequestStat возвращает статистику запросов за указанный период
func (u *Usecase) GetRequestStat(ctx context.Context, userMeta *models.UserMeta, tr *trparser.TimeRange, hostname, requestType string, count uint) (models.RequestStat, error) {
	method := "GetRequestStat"
	u.l.InfoContext(ctx,
		method,
		slog.Any("time_range", tr),
		slog.Any("hostname", hostname),
		slog.String("requestType", requestType),
		slog.Uint64("count", uint64(count)),
	)

	// Проверяем роль пользователя
	// И если она CA, то фильтруем по хосту
	if userMeta.ShortRole == values.ContextAdmin {
		hostname = userMeta.ContextID
	}

	stat, err := u.db.GetRequestStat(ctx, tr, hostname, requestType, count)
	if err != nil {
		err = fmt.Errorf("%s: failed to get requests statistics: %w", method, err)
		u.l.ErrorContext(ctx, "Database operation failed",
			wsl.String("method", method),
			wsl.String("error", err.Error()),
		)
		return models.RequestStat{}, err
	}

	u.l.InfoContext(ctx, "Requests statistics retrieved",
		slog.String("method", method),
	)
	return stat, nil
}

// GetTrafficStat возвращает статистику трафика за указанный период
func (u *Usecase) GetTrafficStat(ctx context.Context, userMeta *models.UserMeta, tr *trparser.TimeRange, hostName string, count uint) (models.TrafficStat, error) {
	method := "GetTrafficStat"
	u.l.InfoContext(ctx,
		method,
		slog.Any("time_range", tr),
		slog.Uint64("count", uint64(count)),
	)

	// Проверяем роль пользователя
	// И если она CA, то фильтруем по хосту
	if userMeta.ShortRole == values.ContextAdmin {
		hostName = userMeta.ContextID
	}

	stat, err := u.mdb.GetTrafficStat(ctx, tr, hostName, count)
	if err != nil {
		err = fmt.Errorf("%s: failed to get traffic statistics: %w", method, err)
		u.l.ErrorContext(ctx, "Database operation failed",
			wsl.String("method", method),
			wsl.String("error", err.Error()),
		)
		return models.TrafficStat{}, err
	}

	u.l.InfoContext(ctx, "Traffic statistics retrieved",
		slog.String("method", method),
	)
	return stat, nil
}

// GetTopUnresolvedDetections возвращает топ нерешенных выявлений
func (u *Usecase) GetTopUnresolvedDetections(ctx context.Context, userMeta *models.UserMeta, tr *trparser.TimeRange, hostName string, count int) ([]dto.UnresolvedDetection, error) {
	method := "GetTopUnresolvedDetections"
	u.l.InfoContext(ctx,
		method,
		slog.Any("time_range", tr),
		slog.Any("hostname", hostName),
		slog.Int("count", count),
	)

	// Проверяем роль пользователя
	// И если она CA, то фильтруем по хосту
	if userMeta.ShortRole == values.ContextAdmin {
		hostName = userMeta.ContextID
	}

	ud, err := u.db.GetTopUnresolvedDetections(ctx, tr, hostName, count)
	if err != nil {
		err = fmt.Errorf("%s: failed to get top unresolved detections: %w", method, err)
		u.l.ErrorContext(ctx, "Database operation failed",
			wsl.String("method", method),
			wsl.String("error", err.Error()),
		)
		return nil, err
	}

	u.l.InfoContext(ctx, "Top unresolved detections retrieved",
		slog.String("method", method),
		slog.Int("detections_count", len(ud)),
		slog.Int("requested_count", count),
	)
	return ud, nil
}

// GetDeviceStat возвращает статистику по устройствам
func (u *Usecase) GetDeviceStat(ctx context.Context, userMeta *models.UserMeta, tr *trparser.TimeRange, count uint) (dto.DeviceStatResponse, error) {
	method := "GetDeviceStat"
	u.l.InfoContext(ctx,
		method,
		slog.Any("time_range", tr),
		slog.Uint64("count", uint64(count)),
	)

	// Для контекстного администратора получаем статистику только по его хосту
	var hostname string
	if userMeta.ShortRole == values.ContextAdmin {
		hostname = userMeta.ContextID
	}

	deviceStat, err := u.db.GetDeviceStat(ctx, tr, hostname, count)
	if err != nil {
		err = fmt.Errorf("%s: failed to get device statistics: %w", method, err)
		u.l.ErrorContext(ctx, "Database operation failed",
			wsl.String("method", method),
			wsl.String("error", err.Error()),
		)
		return deviceStat, err
	}

	u.l.InfoContext(ctx, "Device statistics retrieved",
		slog.String("method", method),
		slog.Int("devices_count", len(deviceStat.Data)),
		slog.Uint64("points_count", uint64(count)),
	)
	return deviceStat, nil
}

// GetAnomalies возвращает аномалии
func (u *Usecase) GetAnomalies(ctx context.Context, userMeta *models.UserMeta, tr *trparser.TimeRange, hostname string) (dto.GetAnomaliesResponse, error) {
	method := "GetAnomalies"
	u.l.InfoContext(ctx,
		method,
		slog.Any("time_range", tr),
	)

	// Проверяем роль пользователя
	// И если она CA, то фильтруем по хосту
	if userMeta.ShortRole == values.ContextAdmin {
		hostname = userMeta.ContextID
	}

	anomalies, err := u.db.GetAnomalies(ctx, tr, hostname)
	if err != nil {
		err = fmt.Errorf("%s: failed to get anomalies: %w", method, err)
		u.l.ErrorContext(ctx, "Database operation failed",
			wsl.String("method", method),
			wsl.String("error", err.Error()),
		)
		return anomalies, err
	}

	u.l.InfoContext(ctx, "Anomalies retrieved",
		slog.String("method", method),
		slog.Int("hosts_count", int(anomalies.HostCount)),
		slog.Int("anomalies_count", len(anomalies.HostAnomalies)),
	)
	return anomalies, nil
}
