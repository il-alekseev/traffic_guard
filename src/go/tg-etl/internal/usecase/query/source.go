package usecase

import (
	"context"
	"fmt"
	"tg-etl/internal/models"
	"tg-etl/pkg/slogger/wsl"
)

// GetSourceByAddr получает источник по IP с кешированием
func (uc *QueryUseCase) GetSourceByAddr(ctx context.Context, ip string) (*models.Source, error) {
	cacheKey := fmt.Sprintf("source:addr:%s", ip)

	// Пытаемся получить из кеша
	if cached, exists := uc.c.Get(cacheKey); exists {
		return cached.(*models.Source), nil
	}

	// Если нет в кеше, ищем в БД
	source, err := uc.etlDB.GetSourceByAddr(ctx, ip)
	if err != nil {
		return nil, err
	}

	// Если источник есть в базе, добавляем в кеш
	if source != nil {
		// Сохраняем в кеш по адресу
		uc.c.Set(cacheKey, source)

		// Также сохраняем в кеш по ID для консистентности
		if source.ID != 0 {
			idCacheKey := fmt.Sprintf("source:id:%d", source.ID)
			uc.c.Set(idCacheKey, source)
		}
	}

	return source, nil
}

// CreateSource создает новый источник с инвалидацией кеша
func (uc *QueryUseCase) CreateSource(ctx context.Context, source models.Source) error {
	err := uc.etlDB.CreateSource(ctx, source)
	if err != nil {
		return err
	}
	// Инвалидируем возможные кеши
	uc.invalidateSourceCache(&source)
	return nil
}

// GetSources получает все источники с кешированием
func (uc *QueryUseCase) GetSources(ctx context.Context) ([]models.Source, error) {
	cacheKey := "sources:all"

	if cached, exists := uc.c.Get(cacheKey); exists {
		return cached.([]models.Source), nil
	}
	sources, err := uc.etlDB.GetSources(ctx)
	if err != nil {
		uc.l.ErrorContext(ctx, "get sources", wsl.Err(err))
		return nil, err
	}

	uc.c.Set(cacheKey, sources)
	return sources, nil
}

// UpdateSource обновляет источник с инвалидацией кеша
func (uc *QueryUseCase) UpdateSource(ctx context.Context, source models.Source) error {
	// Получаем старые данные для инвалидации кеша по старому адресу
	oldSource, err := uc.etlDB.GetSourceByID(ctx, source.ID)
	if err != nil {
		uc.l.ErrorContext(ctx, "get old source for cache invalidation",
			wsl.Int("id", int(source.ID)), wsl.Err(err))
		// Продолжаем выполнение, так как это не критическая ошибка
	} else if oldSource != nil {
		// Инвалидируем кеш по старому адресу
		oldAddrKey := fmt.Sprintf("source:addr:%s", oldSource.IP)
		uc.c.Delete(oldAddrKey)
	}

	// Обновляем источник в БД
	err = uc.etlDB.UpdateSource(ctx, source)
	if err != nil {
		uc.l.ErrorContext(ctx, "update source",
			wsl.Int("id", int(source.ID)), wsl.Err(err))
		return err
	}

	// Инвалидируем кеши
	uc.invalidateSourceCache(&source)
	return nil
}

// invalidateSourceCache инвалидирует все кеши связанные с источником
func (uc *QueryUseCase) invalidateSourceCache(source *models.Source) {
	// Инвалидируем кеш по ID
	uc.c.Delete(fmt.Sprintf("source:id:%d", source.ID))

	// Инвалидируем кеш по адресу
	uc.c.Delete(fmt.Sprintf("source:addr:%s", source.IP))

	// Инвалидируем кеш всех источников
	uc.c.Delete("sources:all")
}

// GetSourceByID получает источник по ID с кешированием
func (uc *QueryUseCase) GetSourceByID(ctx context.Context, id uint) (*models.Source, error) {
	cacheKey := fmt.Sprintf("source:id:%d", id)
	// Пытаемся получить из кеша
	if cached, exists := uc.c.Get(cacheKey); exists {
		return cached.(*models.Source), nil
	}
	// Если нет в кеше, ищем в БД
	source, err := uc.etlDB.GetSourceByID(ctx, id)
	if err != nil {
		uc.l.ErrorContext(ctx, "get source by id from db",
			wsl.Int("id", int(id)), wsl.Err(err))
		return nil, err
	}
	// Если источник есть в базе, добавляем в кеш
	if source != nil {
		uc.c.Set(cacheKey, source)
		// Также сохраняем в кеш по адресу для консистентности
		addrCacheKey := fmt.Sprintf("source:addr:%s", source.IP)
		uc.c.Set(addrCacheKey, source)
	}

	return source, nil
}

// DeleteSource удаляет источник с инвалидацией кеша (дополнительный метод)
func (uc *QueryUseCase) DeleteSource(ctx context.Context, id uint) error {
	// Сначала получаем источник для инвалидации кеша по адресу
	source, err := uc.etlDB.GetSourceByID(ctx, id)
	if err != nil {
		uc.l.ErrorContext(ctx, "get source for cache invalidation",
			wsl.Int("id", int(id)), wsl.Err(err))
		return err
	}
	// Удаляем источник из БД
	err = uc.etlDB.DeleteSource(ctx, id)
	if err != nil {
		uc.l.ErrorContext(ctx, "delete source",
			wsl.Int("id", int(id)), wsl.Err(err))
		return err
	}

	// Инвалидируем кеши
	if source != nil {
		uc.invalidateSourceCache(source)
	}
	return nil
}
