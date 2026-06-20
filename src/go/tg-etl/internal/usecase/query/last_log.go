package usecase

import (
	"context"
	"tg-etl/internal/models"
)

func (uc *QueryUseCase) GetLastLog(ctx context.Context) (*models.LastLog, error) {
	return uc.etlDB.GetLastLog(ctx)
}

func (uc *QueryUseCase) CreateOrUpdateLastLog(ctx context.Context, log models.LastLog) error {
	return uc.etlDB.CreateOrUpdateLastLog(ctx, log)
}
