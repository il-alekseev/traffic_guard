package usecase

import (
	"context"
	"fmt"
	"tg-etl/internal/models"
	"tg-etl/pkg/slogger/wsl"
)

func (uc *QueryUseCase) GetActionByDomainID(ctx context.Context, id uint) (*models.Action, error) {
	cacheKey := fmt.Sprintf("action:id:%d", id)
	// Пытаемся получить из кеша
	if cached, exists := uc.c.Get(cacheKey); exists {
		return cached.(*models.Action), nil
	}

	action, err := uc.etlDB.GetActionByDomainID(ctx, id)
	if err != nil {
		uc.l.ErrorContext(ctx, "get action by id from db", wsl.Err((err)))
	}
	if action != nil {
		uc.c.Set(cacheKey, action)
	}
	return action, nil
}
