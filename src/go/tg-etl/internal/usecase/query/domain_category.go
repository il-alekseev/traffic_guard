package usecase

import (
	"context"
	"fmt"
	"tg-etl/internal/models"
	"tg-etl/pkg/slogger/wsl"
)

func (uc *QueryUseCase) GetNegativeCategoriesByDomainID(ctx context.Context, id uint) (string, error) {
	cacheKey := fmt.Sprintf("negative_categories:domain_id:%d", id)

	// Пытаемся получить из кеша
	if cached, exists := uc.c.Get(cacheKey); exists {
		return cached.(string), nil
	}

	// Если нет в кеше, ищем в БД
	result, err := uc.etlDB.GetNegativeCategoriesByDomainID(ctx, id)
	if err != nil {
		uc.l.ErrorContext(ctx, "get negative categories by domain id", wsl.Err(err))
		return "", err
	}

	// Добавляем в кеш
	uc.c.Set(cacheKey, result)
	return result, nil
}

func (uc *QueryUseCase) GetNegativeCategoriesStatByDomainID(ctx context.Context, id uint) (map[string]float32, error) {
	cacheKey := fmt.Sprintf("negative_categories_stat:domain_id:%d", id)

	// Пытаемся получить из кеша
	if cached, exists := uc.c.Get(cacheKey); exists {
		return cached.(map[string]float32), nil
	}

	// Если нет в кеше, ищем в БД
	result, err := uc.etlDB.GetNegativeCategoriesStatByDomainID(ctx, id)
	if err != nil {
		uc.l.ErrorContext(ctx, "get negative categories stat by domain id", wsl.Err(err))
		return nil, err
	}

	// Добавляем в кеш
	uc.c.Set(cacheKey, result)
	return result, nil
}

func (uc *QueryUseCase) GetMostNegativeCategoryByDomainID(ctx context.Context, id uint) (*models.Category, float32, error) {
	cacheKey := fmt.Sprintf("most_negative_category:domain_id:%d", id)

	// Пытаемся получить из кеша
	if cached, exists := uc.c.Get(cacheKey); exists {
		if cachedData, ok := cached.(struct {
			Category *models.Category
			Percent  float32
		}); ok {
			return cachedData.Category, cachedData.Percent, nil
		}
	}

	// Если нет в кеше, ищем в БД
	cat, percent, err := uc.etlDB.GetMostNegativeCategoryByDomainID(ctx, id)
	if err != nil {
		uc.l.ErrorContext(ctx, "get most negative category by domain id", wsl.Err(err))
		return nil, 0, err
	}

	// Если категория не найдена, возвращаем nil
	if cat == nil {
		return nil, 0, nil
	}

	// Сохраняем в кеш
	uc.c.Set(cacheKey, struct {
		Category *models.Category
		Percent  float32
	}{
		Category: cat,
		Percent:  percent,
	})

	return cat, percent, nil
}

// AddDomainCategory добавляет частоту встречаемости категории у домена или добавляет запись в БД
func (uc *QueryUseCase) AddDomainCategory(ctx context.Context, c models.Category, id uint) error {
	err := uc.etlDB.AddDomainCategory(ctx, c, id)
	if err != nil {
		uc.l.ErrorContext(ctx, "add domain category", wsl.Err(err))
		return err
	}

	// Инвалидируем все кеши, связанные с категориями этого домена
	uc.c.Delete(fmt.Sprintf("negative_categories:domain_id:%d", id))
	uc.c.Delete(fmt.Sprintf("negative_categories_stat:domain_id:%d", id))
	uc.c.Delete(fmt.Sprintf("most_negative_category:domain_id:%d", id))
	uc.c.Delete(fmt.Sprintf("negative_categories_total:domain_id:%d", id))

	return nil
}

func (uc *QueryUseCase) GetNegativeCategoriesTotalByDomainID(ctx context.Context, id uint) (int, error) {
	cacheKey := fmt.Sprintf("negative_categories_total:domain_id:%d", id)

	// Пытаемся получить из кеша
	if cached, exists := uc.c.Get(cacheKey); exists {
		return cached.(int), nil
	}

	// Если нет в кеше, ищем в БД
	result, err := uc.etlDB.GetNegativeCategoriesTotalByDomainID(ctx, id)
	if err != nil {
		uc.l.ErrorContext(ctx, "get negative categories total by domain id", wsl.Err(err))
		return 0, err
	}

	// Добавляем в кеш
	uc.c.Set(cacheKey, result)
	return result, nil
}

// GetNegativeCatWithPercByDomainID - получает название категории домена и ее вероятность
// TODO: Возможно стоит сделать наиболее вероятную категорию как отдельную структуру
func (uc *QueryUseCase) GetNegativeCatWithPercByDomainID(ctx context.Context, id uint) (string, float32, error) {
	cats, err := uc.GetNegativeCategoriesStatByDomainID(ctx, id)
	if err != nil {
		uc.l.ErrorContext(ctx, "failed to get negative categories stat", wsl.Err(err))
		return "", 0, err
	}
	// Если категорий нет - возвращаем пустую строку и 0
	if len(cats) == 0 {
		return "", 0, nil
	} else {
		var name = ""
		var perc = float32(0.0)
		for n, p := range cats {
			if p > perc {
				name = n
				perc = p
			}
		}
		return name, perc, nil
	}
}
