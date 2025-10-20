package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"tg-dbd/internal/controllers/http/v1/dto"
	"tg-dbd/internal/models"
	"tg-dbd/pkg/slogger/wsl"
	"tg-dbd/pkg/trparser"
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

func (u *Usecase) GetTopDetections(ctx context.Context, tr *trparser.TimeRange, f models.DetectionFilter, p models.Pagination) ([]dto.Detection, int64, error) {
	method := "GetTopDetections"
	u.l.InfoContext(ctx,
		method,
		slog.Any("time_range", tr),
		slog.Any("filter", f),
		slog.Any("pagination", p),
	)

	detections, total, err := u.db.GetTopDetections(ctx, tr, f, p)
	if err != nil {
		err = fmt.Errorf("%s: failed to get top detections: %w", method, err)
		u.l.ErrorContext(ctx, "Database operation failed",
			wsl.String("method", method),
			wsl.String("error", err.Error()),
		)
		return nil, 0, err
	}

	u.l.InfoContext(ctx, "Top detections retrieved",
		slog.String("method", method),
		slog.Int("count", len(detections)),
		slog.Int("total", int(total)),
	)
	return detections, total, nil
}

// GetDetectionStat возвращает статистику обнаружений за указанный период
func (u *Usecase) GetDetectionStat(ctx context.Context, tr *trparser.TimeRange, f models.DetectionFilter) (dto.DetectionStat, error) {
	method := "GetDetectionStat"
	u.l.InfoContext(ctx,
		method,
		slog.Any("time_range", tr),
		slog.Any("filter", f),
	)

	stat, err := u.db.GetDetectionStat(ctx, tr, f)
	if err != nil {
		err = fmt.Errorf("%s: failed to get detection stat: %w", method, err)
		u.l.ErrorContext(ctx, "Database operation failed",
			wsl.String("method", method),
			wsl.String("error", err.Error()),
		)
		return dto.DetectionStat{}, err
	}

	u.l.InfoContext(ctx, "detection stat retrieved",
		slog.String("method", method),
		slog.Int("Detected", stat.Detected),
		slog.Int("Accepted", stat.Accepted),
		slog.Int("Denied", stat.Denied),
		slog.Int("Unresolved", stat.Unresolved),
	)
	return stat, nil
}
