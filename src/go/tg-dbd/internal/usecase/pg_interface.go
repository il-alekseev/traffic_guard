package usecase

import (
	"context"
	"tg-dbd/internal/controllers/http/v1/dto"
	"tg-dbd/internal/models"
	"tg-dbd/pkg/trparser"
)

type RepoPGInterface interface {
	GetSessions(ctx context.Context, tr *trparser.TimeRange, f models.SessionFilter, search string, p models.Pagination, s models.Sorting) ([]dto.Session, int64, error)

	GetTopCategories(ctx context.Context, tr *trparser.TimeRange, f models.CategoryFilter, count int) ([]dto.Category, error)
	GetTopDetections(ctx context.Context, tr *trparser.TimeRange, f models.DetectionFilter, p models.Pagination) ([]dto.Detection, int64, error)
	GetDetectionStat(ctx context.Context, tr *trparser.TimeRange, f models.DetectionFilter) (dto.DetectionStat, error)
}
