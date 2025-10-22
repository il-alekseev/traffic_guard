package usecase

import (
	"context"
	"tg-etl/internal/models"
)

type KSURepoPGInterface interface {
	GetLogs(ctx context.Context, start *models.IdsLog, maxCount uint) ([]models.IdsLog, error)
}
