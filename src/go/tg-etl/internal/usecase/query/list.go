package usecase

import (
	"context"
	"fmt"
	"tg-etl/pkg/slogger/wsl"
)

// GetListByDomainID получает список по ID домена с кешированием
func (uc *QueryUseCase) GetListByDomainID(ctx context.Context, id uint) (*string, error) {
	cacheKey := fmt.Sprintf("list:domain_id:%d", id)

	// Пытаемся получить из кеша
	if cached, exists := uc.c.Get(cacheKey); exists {
		return cached.(*string), nil
	}

	// Если нет в кеше, ищем в БД
	list, err := uc.etlDB.GetListByDomainID(ctx, id)
	if err != nil {
		uc.l.ErrorContext(ctx, "get list by domain id from db",
			wsl.Int("id", int(id)), wsl.Err(err))
		return nil, err
	}

	// Если список есть в базе, добавляем в кеш
	if list != nil {
		uc.c.Set(cacheKey, list)
	}

	return list, nil
}

// invalidateListCache инвалидирует кеш для списков
func (uc *QueryUseCase) InvalidateListCache(id uint) {
	cacheKey := fmt.Sprintf("list:domain_id:%d", id)
	uc.c.Delete(cacheKey)
}
