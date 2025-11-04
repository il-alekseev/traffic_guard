package usecase

import (
	"context"
	"fmt"
	"tg-etl/internal/models"
	"tg-etl/pkg/slogger/wsl"

	"github.com/google/uuid"
)

func (uc *QueryUseCase) CreateURL(ctx context.Context, url models.URL) error {
	err := uc.etlDB.CreateURL(ctx, url)
	if err != nil {
		return err
	}
	// Инвалидируем возможные кеши
	uc.invalidateUrlCache(&url)
	return nil
}

func (uc *QueryUseCase) GetURLByPath(ctx context.Context, path string) (*models.URL, error) {
	cacheKey := fmt.Sprintf("url:path:%s", path)

	// Пытаемся получить из кеша
	if cached, exists := uc.c.Get(cacheKey); exists {
		return cached.(*models.URL), nil
	}
	// Если нет в кеше, ищем в БД
	url, err := uc.etlDB.GetURLByPath(ctx, path)
	if err != nil {
		uc.l.ErrorContext(ctx, "get url by path from db",
			wsl.String("path", path), wsl.Err(err))
		return nil, err
	}

	// Если URL есть в базе, добавляем в кеш
	if url != nil {
		uc.c.Set(cacheKey, url)

		// Также сохраняем в кеш по ID для консистентности
		if url.ID != 0 {
			idCacheKey := fmt.Sprintf("url:id:%d", url.ID)
			uc.c.Set(idCacheKey, url)
		}
	}

	return url, nil
}

func (uc *QueryUseCase) GetURLByPathDomain(ctx context.Context, path string, id uint) (*models.URL, error) {
	cacheKey := fmt.Sprintf("url:path:id%s:%d", path, id)

	// Пытаемся получить из кеша
	if cached, exists := uc.c.Get(cacheKey); exists {
		return cached.(*models.URL), nil
	}
	// Если нет в кеше, ищем в БД
	url, err := uc.etlDB.GetURLByPathDomain(ctx, path, id)
	if err != nil {
		uc.l.ErrorContext(ctx, "get url by path from db",
			wsl.String("path", path), wsl.Err(err))
		return nil, err
	}

	// Если URL есть в базе, добавляем в кеш
	if url != nil {
		uc.c.Set(cacheKey, url)

		// Также сохраняем в кеш по ID для консистентности
		if url.ID != 0 {
			idCacheKey := fmt.Sprintf("url:id:%d", url.ID)
			uc.c.Set(idCacheKey, url)
		}
	}
	return url, nil
}

// TODO:  добавить кеширование
func (uc *QueryUseCase) GetURLByRequestID(ctx context.Context, requestID uuid.UUID) (*models.URL, error) {
	return uc.etlDB.GetURLByRequestID(ctx, requestID)
}

// invalidateSourceCache инвалидирует все кеши связанные с источником
func (uc *QueryUseCase) invalidateUrlCache(url *models.URL) {
	// Инвалидируем кеш по ID
	uc.c.Delete(fmt.Sprintf("url:id:%d", url.ID))

	// Инвалидируем кеш по пути
	uc.c.Delete(fmt.Sprintf("url:path:%s", url.Path))

	// Инвалидируем кеш по пути и ID домена
	uc.c.Delete(fmt.Sprintf("url:path:id%s:%d", url.Path, url.DomainID))
	// Инвалидируем кеш всех источников

	uc.c.Delete("url:all")
}
