package usecase

import (
	"cmd/etl/internal/models"
	"context"
)

type KSURepoPGInterface interface {
	GetLogs(ctx context.Context, start *models.IdsLog, maxCount uint) ([]models.IdsLog, error)
}
