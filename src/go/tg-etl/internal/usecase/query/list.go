package usecase

import (
	"context"
	"fmt"
	"tg-etl/internal/models"
	"tg-etl/pkg/slogger/wsl"
)

func (uc *QueryUseCase) AddDomainToList(ctx context.Context, domain models.Domain, list models.ListType) error {
	existedList, err := uc.GetListByDomainID(ctx, domain.ID)
	if err != nil {
		uc.l.ErrorContext(ctx, "get list for domain", wsl.Err((err)))
		return err
	}
	// Домен может быть одновременно только в одном списке!
	if existedList != nil {
		err = fmt.Errorf("domain already in list %s", *existedList)
		uc.l.ErrorContext(ctx, "get list for domain", wsl.Err((err)))
		return err
	}
	if err := uc.etlDB.AddDomainToList(ctx, domain, list); err != nil {
		uc.l.ErrorContext(ctx, "add domain to list", wsl.Err((err)))
		return err
	}
	// Инвалидируем возможные кеши
	uc.c.Delete(fmt.Sprintf("list:domain_id:%d", domain.ID))
	return nil
}

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
