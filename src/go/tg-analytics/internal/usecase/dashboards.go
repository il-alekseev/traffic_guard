package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"tg-an/internal/controllers/http/v1/dto"
	"tg-an/internal/models"
	"tg-an/pkg/slogger/wsl"
	"tg-an/pkg/trparser"
)

func (u *Usecase) GetTopCategories(ctx context.Context, tr *trparser.TimeRange, f models.CategoryFilter, count int) ([]dto.Category, error) {
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

// GetRequestsStat возвращает статистику запросов за указанный период
func (u *Usecase) GetRequestStat(ctx context.Context,
	tr *trparser.TimeRange,
	hostname, requestType string, count uint) (models.RequestStat, error) {
	method := "GetRequestsStat"
	u.l.InfoContext(ctx,
		method,
		slog.Any("time_range", tr),
		slog.Any("hostname", hostname),
		slog.String("requestType", requestType),
		slog.Uint64("count", uint64(count)),
	)

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
func (u *Usecase) GetTrafficStat(ctx context.Context, tr *trparser.TimeRange, hostName string, count uint) (models.TrafficStat, error) {
	method := "GetTrafficStat"
	u.l.InfoContext(ctx,
		method,
		slog.Any("time_range", tr),
		slog.Uint64("count", uint64(count)),
	)
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

func (u *Usecase) GetTopUnresolvedDetections(ctx context.Context, tr *trparser.TimeRange, hostName string, count int) ([]dto.UnresolvedDetection, error) {
	method := "GetTopUnresolvedDetections"
	u.l.InfoContext(ctx,
		method,
		slog.Any("time_range", tr),
		slog.Any("hostname", hostName),
		slog.Int("count", count),
	)

	ud, err := u.db.GetTopUnresolvedDetections(ctx, tr, hostName, count)
	if err != nil {
		err = fmt.Errorf("%s: failed to get top unresolved detecrtions: %w", method, err)
		u.l.ErrorContext(ctx, "Database operation failed",
			wsl.String("method", method),
			wsl.String("error", err.Error()),
		)
		return nil, err
	}

	u.l.InfoContext(ctx, "Top unresolved detections retrieved",
		slog.String("method", method),
		slog.Int("categories_count", len(ud)),
		slog.Int("requested_count", count),
	)
	return ud, nil
}

func (u *Usecase) GetDeviceStat(ctx context.Context, tr *trparser.TimeRange, count uint) (dto.DeviceStatResponse, error) {
	method := "GetDeviceStat"
	u.l.InfoContext(ctx,
		method,
		slog.Any("time_range", tr),
		slog.Uint64("count", uint64(count)),
	)

	deviceStat, err := u.db.GetDeviceStat(ctx, tr, count)
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

func (u *Usecase) GetAnomalies(ctx context.Context, tr *trparser.TimeRange) (dto.GetAnomaliesResponse, error) {
	method := "GetAnomalies"
	u.l.InfoContext(ctx,
		method,
		slog.Any("time_range", tr),
	)

	anomalies, err := u.db.GetAnomalies(ctx, tr)
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
	)
	return anomalies, nil
}

func (u *Usecase) Act(ctx context.Context, action, path string) error {
	method := "Act"
	u.l.InfoContext(ctx,
		method,
		wsl.String("action", action),
		slog.Any("path", path),
	)

	err := u.db.Act(ctx, action, path)
	if err != nil {
		err = fmt.Errorf("%s: failed act with domain: %w", method, err)
		u.l.ErrorContext(ctx, "Database operation failed",
			wsl.String("method", method),
			wsl.String("error", err.Error()),
		)
		return err
	}

	u.l.InfoContext(ctx, "Act success",
		slog.String("method", method),
	)
	return nil
}
