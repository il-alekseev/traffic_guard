package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"tg-an/internal/controllers/http/v1/dto"
	"tg-an/internal/models"
	"tg-an/internal/pkg/status"
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
func (u *Usecase) GetRequestsStat(ctx context.Context, tr *trparser.TimeRange, f models.DashboardFilter, s status.Status, count uint) ([]uint, error) {
	method := "GetRequestsStat"
	u.l.InfoContext(ctx,
		method,
		slog.Any("time_range", tr),
		slog.Any("filter", f),
		slog.String("status", s.String()),
		slog.Uint64("count", uint64(count)),
	)

	stat, err := u.db.GetRequestsStat(ctx, tr, f, s, count)
	if err != nil {
		err = fmt.Errorf("%s: failed to get requests statistics: %w", method, err)
		u.l.ErrorContext(ctx, "Database operation failed",
			wsl.String("method", method),
			wsl.String("error", err.Error()),
		)
		return nil, err
	}

	u.l.InfoContext(ctx, "Requests statistics retrieved",
		slog.String("method", method),
	)
	return stat, nil
}

// GetTrafficStat возвращает статистику трафика за указанный период
func (u *Usecase) GetTrafficStat(ctx context.Context, tr *trparser.TimeRange, count uint) (models.TrafficStat, error) {
	method := "GetTrafficStat"
	u.l.InfoContext(ctx,
		method,
		slog.Any("time_range", tr),
		slog.Uint64("count", uint64(count)),
	)
	stat, err := u.mdb.GetTrafficStat(ctx, tr, count)
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
