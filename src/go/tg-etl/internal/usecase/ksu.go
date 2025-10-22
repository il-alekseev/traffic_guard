package usecase

import (
	"context"
	"tg-etl/internal/models"
	"tg-etl/pkg/slogger/wsl"
)

func (uc *UseCase) GetLogs(ctx context.Context, start *models.IdsLog, count uint) ([]models.IdsLog, error) {
	// Получаем coun логов начиная с указанного (включая его)
	logs, err := uc.ksuDB.GetLogs(ctx, start, count)
	if err != nil {
		uc.l.ErrorContext(ctx, "get logs", wsl.Err(err))
		return nil, nil
	}
	// Если лог был задан, значит он уже учтен в ETL, поэтому его необходиомо убрать из выборки
	if start != nil {
		// Если вернулся только сама запись, значит, новых данных нет - возвращаем пустой слайс
		if len(logs) < 1 {
			return nil, nil
		}
		logs = logs[1 : len(logs)-1]
	}
	//uc.l.DebugContext(ctx, "success get logs", wsl.Int("count", len(logs)))
	return logs, nil
}
