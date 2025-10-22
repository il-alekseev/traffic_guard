package usecase

import (
	"context"
	"fmt"
	"tg-etl/internal/models"
	"tg-etl/pkg/slogger/wsl"
)

// GetSourceByAddr получает источник по IP с кешированием
func (uc *UseCase) GetSourceByAddr(ctx context.Context, ip string) (*models.Source, error) {
	cacheKey := fmt.Sprintf("source:addr:%s", ip)

	// Пытаемся получить из кеша
	if cached, exists := uc.c.Get(cacheKey); exists {
		//uc.l.DebugContext(ctx, "cache hit for source by addr",
		//	wsl.String("ip", ip), wsl.Int("port", port))
		return cached.(*models.Source), nil
	}

	//uc.l.DebugContext(ctx, "cache miss for source by addr",
	//	slog.String("ip", ip), slog.Int("port", port))

	// Если нет в кеше, ищем в БД
	source, err := uc.etlDB.GetSourceByAddr(ctx, ip)
	if err != nil {
		//uc.l.ErrorContext(ctx, "get source by addr from db",
		//	wsl.String("ip", ip), wsl.Int("port", port), wsl.Err(err))
		return nil, err
	}

	// Если источник есть в базе, добавляем в кеш
	if source != nil {
		//uc.l.DebugContext(ctx, "success got source by addr",
		//	slog.String("ip", ip), slog.Int("port", port))

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
func (uc *UseCase) CreateSource(ctx context.Context, source models.Source) error {
	err := uc.etlDB.CreateSource(ctx, source)
	if err != nil {
		//uc.l.ErrorContext(ctx, "create source",
		//	wsl.String("ip", source.IP), wsl.Int("port", source.Port), wsl.Err(err))
		return err
	}

	// Инвалидируем возможные кеши
	uc.invalidateSourceCache(&source)

	//uc.l.DebugContext(ctx, "success create source",
	//	slog.String("ip", source.IP), slog.Int("port", source.Port))
	return nil
}

// GetSources получает все источники с кешированием
func (uc *UseCase) GetSources(ctx context.Context) ([]models.Source, error) {
	cacheKey := "sources:all"

	if cached, exists := uc.c.Get(cacheKey); exists {
		//uc.l.DebugContext(ctx, "cache hit for all sources")
		return cached.([]models.Source), nil
	}

	//uc.l.DebugContext(ctx, "cache miss for all sources")

	sources, err := uc.etlDB.GetSources(ctx)
	if err != nil {
		uc.l.ErrorContext(ctx, "get sources", wsl.Err(err))
		return nil, err
	}

	uc.c.Set(cacheKey, sources)
	//uc.l.DebugContext(ctx, "success got sources", slog.Int("count", len(sources)))
	return sources, nil
}

// UpdateSource обновляет источник с инвалидацией кеша
func (uc *UseCase) UpdateSource(ctx context.Context, source models.Source) error {
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

	//uc.l.DebugContext(ctx, "success update source",
	//	slog.Int("id", int(source.ID)),
	//	slog.String("ip", source.IP),
	//	slog.Int("port", source.Port))
	return nil
}

// invalidateSourceCache инвалидирует все кеши связанные с источником
func (uc *UseCase) invalidateSourceCache(source *models.Source) {
	// Инвалидируем кеш по ID
	uc.c.Delete(fmt.Sprintf("source:id:%d", source.ID))

	// Инвалидируем кеш по адресу
	uc.c.Delete(fmt.Sprintf("source:addr:%s", source.IP))

	// Инвалидируем кеш всех источников
	uc.c.Delete("sources:all")

	//uc.l.Debug("invalidated source cache",
	//	slog.Int("id", int(source.ID)),
	//	slog.String("ip", source.IP),
	//	slog.Int("port", source.Port))
}

// GetSourceByID получает источник по ID с кешированием
func (uc *UseCase) GetSourceByID(ctx context.Context, id uint) (*models.Source, error) {
	cacheKey := fmt.Sprintf("source:id:%d", id)

	// Пытаемся получить из кеша
	if cached, exists := uc.c.Get(cacheKey); exists {
		//uc.l.DebugContext(ctx, "cache hit for source by id", wsl.Int("id", int(id)))
		return cached.(*models.Source), nil
	}

	//uc.l.DebugContext(ctx, "cache miss for source by id", slog.Int("id", int(id)))

	// Если нет в кеше, ищем в БД
	source, err := uc.etlDB.GetSourceByID(ctx, id)
	if err != nil {
		uc.l.ErrorContext(ctx, "get source by id from db",
			wsl.Int("id", int(id)), wsl.Err(err))
		return nil, err
	}

	// Если источник есть в базе, добавляем в кеш
	if source != nil {
		//uc.l.DebugContext(ctx, "success got source by id", slog.Int("id", int(id)))
		uc.c.Set(cacheKey, source)

		// Также сохраняем в кеш по адресу для консистентности
		addrCacheKey := fmt.Sprintf("source:addr:%s", source.IP)
		uc.c.Set(addrCacheKey, source)
	}

	return source, nil
}

// DeleteSource удаляет источник с инвалидацией кеша (дополнительный метод)
func (uc *UseCase) DeleteSource(ctx context.Context, id uint) error {
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

	//uc.l.DebugContext(ctx, "success delete source", slog.Int("id", int(id)))
	return nil
}
