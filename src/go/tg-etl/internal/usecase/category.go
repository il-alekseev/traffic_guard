package usecase

import (
	"context"
	"fmt"
	"tg-etl/internal/repo/category"
	"tg-etl/pkg/slogger/wsl"
)

func (uc *UseCase) CreateCategories(ctx context.Context, categories []string) error {
	err := uc.etlDB.CreateCategories(ctx, categories)
	if err != nil {
		uc.l.ErrorContext(ctx, "create categories", wsl.Err((err)))
		return err
	}
	return nil
}

func (uc *UseCase) GetCategoryByID(ctx context.Context, id uint) (*category.ContentCategory, error) {
	cacheKey := fmt.Sprintf("category:id:%d", id)
	// Пытаемся получить из кеша
	if cached, exists := uc.c.Get(cacheKey); exists {
		return cached.(*category.ContentCategory), nil
	}
	// Если нет в кеше, ищем в БД
	c, err := uc.etlDB.GetCategoryByID(ctx, id)
	if err != nil {
		uc.l.ErrorContext(ctx, "get category by id from db", wsl.Err((err)))
	}
	cat, err := category.ParseContentCategory(c.Name)
	if err != nil {
		uc.l.ErrorContext(ctx, "get category by id from db", wsl.Err((err)))
		return nil, err
	}
	// Добавляем в кеш
	uc.c.Set(cacheKey, &cat)

	return &cat, nil
}

func (uc *UseCase) GetCategoryByName(ctx context.Context, name string) (*category.ContentCategory, error) {
	cacheKey := fmt.Sprintf("category:name:%s", name)

	// Пытаемся получить из кеша
	if cached, exists := uc.c.Get(cacheKey); exists {
		return cached.(*category.ContentCategory), nil
	}

	// Если нет в кеше, ищем в БД
	c, err := uc.etlDB.GetCategoryByName(ctx, name)
	if err != nil {
		uc.l.ErrorContext(ctx, "get category by name from db", wsl.Err(err))
		return nil, err
	}

	// Парсим категорию из строки
	cat, err := category.ParseContentCategory(c.Name)
	if err != nil {
		uc.l.ErrorContext(ctx, "parse category by name", wsl.Err(err))
		return nil, err
	}

	// Добавляем в кеш
	uc.c.Set(cacheKey, &cat)

	return &cat, nil
}

func (uc *UseCase) GetCategories(ctx context.Context) ([]category.ContentCategory, error) {
	cacheKey := "categories:all"
	// Пытаемся получить из кеша
	if cached, exists := uc.c.Get(cacheKey); exists {
		return cached.([]category.ContentCategory), nil
	}
	// Получаем категории из БД
	categories, err := uc.etlDB.GetCategories(ctx)
	if err != nil {
		uc.l.ErrorContext(ctx, "get categories", wsl.Err(err))
		return nil, err
	}
	// Конвертируем категории из БД в ContentCategory
	result := make([]category.ContentCategory, 0, len(categories))
	for _, cat := range categories {
		contentCat, err := category.ParseContentCategory(cat.Name)
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
