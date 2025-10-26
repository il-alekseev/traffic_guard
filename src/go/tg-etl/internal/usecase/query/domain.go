package usecase

import (
	"context"
	"fmt"
	"tg-etl/internal/models"
	"tg-etl/pkg/slogger/wsl"
)

func (uc *QueryUseCase) GetDomainByRequestID(ctx context.Context, id string) (*models.Domain, error) {
	return uc.etlDB.GetDomainByRequestID(ctx, id)
}

// GetDomainByID получает домен по ID с кешированием
func (uc *QueryUseCase) GetDomainByID(ctx context.Context, id uint) (*models.Domain, error) {
	cacheKey := fmt.Sprintf("domain:id:%d", id)
	// Пытаемся получить из кеша
	if cached, exists := uc.c.Get(cacheKey); exists {
		return cached.(*models.Domain), nil
	}
	// Если нет в кеше, ищем в БД
	domain, err := uc.etlDB.GetDomainByID(ctx, id)
	if err != nil {
		uc.l.ErrorContext(ctx, "get domain by id from db",
			wsl.Int("id", int(id)), wsl.Err(err))
		return nil, err
	}
	// Если домен есть в базе, добавляем в кеш
	if domain != nil {
		uc.c.Set(cacheKey, domain)
	}

	return domain, nil
}
func (uc *QueryUseCase) GetDomainByPath(ctx context.Context, path string) (*models.Domain, error) {
	cacheKey := fmt.Sprintf("domain:path:%s", path)
	// Пытаемся получить из кеша
	if cached, exists := uc.c.Get(cacheKey); exists {
		return cached.(*models.Domain), nil
	}
	// Если нет в кеше, ищем в БД
	domain, err := uc.etlDB.GetDomainByPath(ctx, path)
	if err != nil {
		uc.l.ErrorContext(ctx, "get domain by path from db",
			wsl.String("path", path), wsl.Err(err))
		return nil, err
	}
	// Если домен есть в базе, добавляем в кеш
	if domain != nil {
		uc.c.Set(cacheKey, domain)
	}
	return domain, nil
}

// GetDomainByAddr получает домен по IP и порту с кешированием
func (uc *QueryUseCase) GetDomainByAddr(ctx context.Context, ip string, port int) (*models.Domain, error) {
	cacheKey := fmt.Sprintf("domain:addr:%s:%d", ip, port)
	// Пытаемся получить из кеша
	if cached, exists := uc.c.Get(cacheKey); exists {
		return cached.(*models.Domain), nil
	}
	// Если нет в кеше, ищем в БД
	domain, err := uc.etlDB.GetDomainByAddr(ctx, ip, port)
	if err != nil {
		uc.l.ErrorContext(ctx, "get domain by addr from db",
			wsl.String("ip", ip), wsl.Int("port", port), wsl.Err(err))
		return nil, err
	}
	// Если домен есть в базе, добавляем в кеш
	if domain != nil {
		uc.c.Set(cacheKey, domain)
		// Также сохраняем в кеш по ID для консистентности
		if domain.ID != 0 {
			idCacheKey := fmt.Sprintf("domain:id:%d", domain.ID)
			uc.c.Set(idCacheKey, domain)
		}
	}

	return domain, nil
}

// CreateDomain создает новый домен с инвалидацией кеша
func (uc *QueryUseCase) CreateDomain(ctx context.Context, domain models.Domain) error {
	err := uc.etlDB.CreateDomain(ctx, domain)
	if err != nil {
		uc.l.ErrorContext(ctx, "create domain",
			wsl.String("ip", domain.IP), wsl.Int("port", domain.Port), wsl.Err(err))
		return err
	}
	// Инвалидируем возможные кеши
	uc.invalidateDomainCache(&domain)
	return nil
}

// UpdateDomain обновляет домен с инвалидацией кеша
// Обновление возможно по ID, доменному имени или IP домена
func (uc *QueryUseCase) UpdateDomain(ctx context.Context, domain models.Domain) error {
	id := domain.ID
	var oldDomain *models.Domain
	var err error
	// Получаем старые данные для инвалидации кеша по старому адресу
	if id != 0 {
		oldDomain, err = uc.etlDB.GetDomainByID(ctx, domain.ID)
		if err != nil {
			uc.l.ErrorContext(ctx, "get old domain for cache invalidation",
				wsl.Int("id", int(domain.ID)), wsl.Err(err))
		}
	} else if domain.Path != "" {
		oldDomain, err := uc.etlDB.GetDomainByPath(ctx, domain.Path)
		if err != nil {
			uc.l.ErrorContext(ctx, "get old domain for cache invalidation",
				wsl.String("path", domain.Path), wsl.Err(err))
		} else if oldDomain != nil {
			id = oldDomain.ID
		}
	} else if domain.IP != "" {
		oldDomain, err := uc.etlDB.GetDomainByAddr(ctx, domain.IP, domain.Port)
		if err != nil {
			uc.l.ErrorContext(ctx, "get old domain for cache invalidation",
				wsl.String("IP", domain.IP), wsl.Err(err))
		} else if oldDomain != nil {
			id = oldDomain.ID
		}
	} else {
		err = fmt.Errorf("empty domain")
		uc.l.ErrorContext(ctx, "update domain", wsl.Err(err))
		return err
	}
	if oldDomain != nil {
		// Инвалидируем все кеши
		uc.invalidateDomainCache(oldDomain)
	}
	// Обновляем домен в БД
	err = uc.etlDB.UpdateDomainByID(ctx, domain, int(id))
	if err != nil {
		uc.l.ErrorContext(ctx, "update domain",
			wsl.Int("id", int(domain.ID)), wsl.Err(err))
		return err
	}
	// Инвалидируем кеши
	uc.invalidateDomainCache(&domain)
	return nil
}

// GetDomains получает все домены с кешированием
func (uc *QueryUseCase) GetDomains(ctx context.Context) ([]models.Domain, error) {
	cacheKey := "domains:all"

	if cached, exists := uc.c.Get(cacheKey); exists {
		return cached.([]models.Domain), nil
	}
	domains, err := uc.etlDB.GetDomains(ctx)
	if err != nil {
		uc.l.ErrorContext(ctx, "get domains", wsl.Err(err))
		return nil, err
	}

	uc.c.Set(cacheKey, domains)
	return domains, nil
}

// invalidateDomainCache инвалидирует все кеши связанные с доменом
func (uc *QueryUseCase) invalidateDomainCache(domain *models.Domain) {
	// Инвалидируем кеш по ID
	uc.c.Delete(fmt.Sprintf("domain:id:%d", domain.ID))

	// Инвалидируем кеш по адресу
	uc.c.Delete(fmt.Sprintf("domain:addr:%s:%d", domain.IP, domain.Port))

	// Инвалидируем кеш по доменному имени
	uc.c.Delete(fmt.Sprintf("domain:path:%s", domain.Path))

	// Инвалидируем кеш всех доменов
	uc.c.Delete("domains:all")
}
