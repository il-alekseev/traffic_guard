package usecase

import (
	"context"
	"tg-etl/internal/models"
	"tg-etl/pkg/slogger/wsl"
)

func (uc *QueryUseCase) GetLogs(ctx context.Context, start *models.LastLog, count uint) ([]models.IdsLog, error) {
	// Получаем coun логов начиная с указанного (включая его)
	logs, err := uc.ksuDB.GetLogs(ctx, start, count)
	if err != nil {
		uc.l.ErrorContext(ctx, "get logs", wsl.Err(err))
		return nil, nil
	}
	return logs, nil
}
