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

func (u *Usecase) GetTopDetections(ctx context.Context, tr *trparser.TimeRange, f models.DetectionFilter, a string, p models.Pagination) ([]dto.Detection, int64, error) {
	method := "GetTopDetections"
	u.l.InfoContext(ctx,
		method,
		slog.Any("time_range", tr),
		slog.Any("filter", f),
		slog.Any("pagination", p),
	)

	detections, total, err := u.db.GetTopDetections(ctx, tr, f, a, p)
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
		slog.Int64("Detected", stat.Detected),
		slog.Int64("Allowed", stat.Allowed),
		slog.Int64("Denied", stat.Denied),
		slog.Int64("Unresolved", stat.Unresolved),
	)
	return stat, nil
}
