package usecase

import (
	"context"
	"fmt"
	"tg-etl/internal/models"
	"tg-etl/pkg/slogger/wsl"
)

func (uc *QueryUseCase) CreatePredefinedCategories(ctx context.Context, categories []models.Category) error {
	err := uc.etlDB.CreateCategories(ctx, categories)
	if err != nil {
		uc.l.ErrorContext(ctx, "create categories", wsl.Err((err)))
		return err
	}
	return nil
}

// TODO: добавить кеширование
func (uc *QueryUseCase) GetCategoryIDByName(ctx context.Context, name string) (uint, error) {
	c, err := uc.etlDB.GetCategoryByName(ctx, name)
	if err != nil {
		uc.l.ErrorContext(ctx, "get category by name from db", wsl.Err(err))
		return 0, err
	}
	if c == nil {
		return 0, nil
	} else {
		return c.ID, nil
	}
}

func (uc *QueryUseCase) GetCategoryByID(ctx context.Context, id uint) (*models.Category, error) {
	cacheKey := fmt.Sprintf("category:id:%d", id)
	// Пытаемся получить из кеша
	if cached, exists := uc.c.Get(cacheKey); exists {
		return cached.(*models.Category), nil
	}
	// Если нет в кеше, ищем в БД
	c, err := uc.etlDB.GetCategoryByID(ctx, id)
	if err != nil {
		uc.l.ErrorContext(ctx, "get category by id from db", wsl.Err((err)))
	}
	cat, err := models.ParseContentCategory(c.Name)
	if err != nil {
		uc.l.ErrorContext(ctx, "get category by id from db", wsl.Err((err)))
		return nil, err
	}
	// Добавляем в кеш
	uc.c.Set(cacheKey, &cat)

	return &cat, nil
}

func (uc *QueryUseCase) GetCategoryByName(ctx context.Context, name string) (*models.Category, error) {
	cacheKey := fmt.Sprintf("category:name:%s", name)

	// Пытаемся получить из кеша
	if cached, exists := uc.c.Get(cacheKey); exists {
		return cached.(*models.Category), nil
	}

	// Если нет в кеше, ищем в БД
	c, err := uc.etlDB.GetCategoryByName(ctx, name)
	if err != nil {
		uc.l.ErrorContext(ctx, "get category by name from db", wsl.Err(err))
		return nil, err
	}

	// Парсим категорию из строки
	cat, err := models.ParseContentCategory(c.Name)
	if err != nil {
		uc.l.ErrorContext(ctx, "parse category by name", wsl.Err(err))
		return nil, err
	}

	// Добавляем в кеш
	uc.c.Set(cacheKey, &cat)

	return &cat, nil
}

func (uc *QueryUseCase) GetCategories(ctx context.Context) ([]models.Category, error) {
	cacheKey := "categories:all"
	// Пытаемся получить из кеша
	if cached, exists := uc.c.Get(cacheKey); exists {
		return cached.([]models.Category), nil
	}
	// Получаем категории из БД
	categories, err := uc.etlDB.GetCategories(ctx)
	if err != nil {
		uc.l.ErrorContext(ctx, "get categories", wsl.Err(err))
		return nil, err
	}
	// Конвертируем категории из БД в ContentCategory
	result := make([]models.Category, 0, len(categories))
	for _, cat := range categories {
		contentCat, err := models.ParseContentCategory(cat.Name)
		if err != nil {
			uc.l.ErrorContext(ctx, "parse category", wsl.Err(err))
			continue // Пропускаем некорректные категории, но продолжаем обработку
		}
		result = append(result, contentCat)
	}
	// Добавляем в кеш
	uc.c.Set(cacheKey, result)
	return result, nil
}
