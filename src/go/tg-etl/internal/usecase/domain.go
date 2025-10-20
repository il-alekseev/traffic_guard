package usecase

import (
	"cmd/etl/internal/models"
	"cmd/etl/pkg/slogger/wsl"
	"context"
	"fmt"
	"log/slog"
	"time"
)

// GetDomainByID получает домен по ID с кешированием
func (uc *UseCase) GetDomainByID(ctx context.Context, id uint) (*models.Domain, error) {
	cacheKey := fmt.Sprintf("domain:id:%d", id)

	// Пытаемся получить из кеша
	if cached, exists := uc.c.Get(cacheKey); exists {
		uc.l.DebugContext(ctx, "cache hit for domain by id", wsl.Int("id", int(id)))
		return cached.(*models.Domain), nil
	}

	uc.l.DebugContext(ctx, "cache miss for domain by id", slog.Int("id", int(id)))

	// Если нет в кеше, ищем в БД
	domain, err := uc.etlDB.GetDomainByID(ctx, id)
	if err != nil {
		uc.l.ErrorContext(ctx, "get domain by id from db",
			wsl.Int("id", int(id)), wsl.Err(err))
		return nil, err
	}

	// Если домен есть в базе, добавляем в кеш
	if domain != nil {
		uc.l.DebugContext(ctx, "success got domain by id", slog.Int("id", int(id)))
		uc.c.Set(cacheKey, domain)
	}

	return domain, nil
}

// GetDomainByAddr получает домен по IP и порту с кешированием
func (uc *UseCase) GetDomainByAddr(ctx context.Context, ip string, port int) (*models.Domain, error) {
	cacheKey := fmt.Sprintf("domain:addr:%s:%d", ip, port)

	// Пытаемся получить из кеша
	if cached, exists := uc.c.Get(cacheKey); exists {
		uc.l.DebugContext(ctx, "cache hit for domain by addr",
			wsl.String("ip", ip), wsl.Int("port", port))
		return cached.(*models.Domain), nil
	}

	uc.l.DebugContext(ctx, "cache miss for domain by addr",
		slog.String("ip", ip), slog.Int("port", port))

	// Если нет в кеше, ищем в БД
	domain, err := uc.etlDB.GetDomainByAddr(ctx, ip, port)
	if err != nil {
		uc.l.ErrorContext(ctx, "get domain by addr from db",
			wsl.String("ip", ip), wsl.Int("port", port), wsl.Err(err))
		return nil, err
	}

	// Если домен есть в базе, добавляем в кеш
	if domain != nil {
		uc.l.DebugContext(ctx, "success got domain by addr",
			slog.String("ip", ip), slog.Int("port", port))
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
func (uc *UseCase) CreateDomain(ctx context.Context, domain models.Domain) error {
	err := uc.etlDB.CreateDomain(ctx, domain)
	if err != nil {
		uc.l.ErrorContext(ctx, "create domain",
			wsl.String("ip", domain.IP), wsl.Int("port", domain.Port), wsl.Err(err))
		return err
	}

	// Инвалидируем возможные кеши
	uc.invalidateDomainCache(&domain)

	uc.l.DebugContext(ctx, "success create domain",
		slog.String("ip", domain.IP), slog.Int("port", domain.Port))
	return nil
}

// UpdateDomain обновляет домен с инвалидацией кеша
func (uc *UseCase) UpdateDomain(ctx context.Context, domain models.Domain) error {
	// Получаем старые данные для инвалидации кеша по старому адресу
	oldDomain, err := uc.etlDB.GetDomainByID(ctx, domain.ID)
	if err != nil {
		uc.l.ErrorContext(ctx, "get old domain for cache invalidation",
			wsl.Int("id", int(domain.ID)), wsl.Err(err))
		// Продолжаем выполнение, так как это не критическая ошибка
	} else if oldDomain != nil {
		// Инвалидируем кеш по старому адресу
		oldAddrKey := fmt.Sprintf("domain:addr:%s:%d", oldDomain.IP, oldDomain.Port)
		uc.c.Delete(oldAddrKey)
	}

	// Обновляем домен в БД
	err = uc.etlDB.UpdateDomain(ctx, domain)
	if err != nil {
		uc.l.ErrorContext(ctx, "update domain",
			wsl.Int("id", int(domain.ID)), wsl.Err(err))
		return err
	}

	// Инвалидируем кеши
	uc.invalidateDomainCache(&domain)

	uc.l.DebugContext(ctx, "success update domain",
		slog.Int("id", int(domain.ID)),
		slog.String("ip", domain.IP),
		slog.Int("port", domain.Port))
	return nil
}

// GetDomains получает все домены с кешированием
func (uc *UseCase) GetDomains(ctx context.Context) ([]models.Domain, error) {
	cacheKey := "domains:all"

	if cached, exists := uc.c.Get(cacheKey); exists {
		uc.l.DebugContext(ctx, "cache hit for all domains")
		return cached.([]models.Domain), nil
	}

	uc.l.DebugContext(ctx, "cache miss for all domains")

	domains, err := uc.etlDB.GetDomains(ctx)
	if err != nil {
		uc.l.ErrorContext(ctx, "get domains", wsl.Err(err))
		return nil, err
	}

	uc.c.Set(cacheKey, domains)
	uc.l.DebugContext(ctx, "success got domains", slog.Int("count", len(domains)))
	return domains, nil
}

// IncrementAccessCount инкрементирует счетчик доступа с инвалидацией кеша
func (uc *UseCase) IncrementAccessCount(ctx context.Context, domainID uint) error {
	err := uc.etlDB.IncrementDomainAccessCount(ctx, domainID)
	if err != nil {
		uc.l.ErrorContext(ctx, "increment domain access count",
			wsl.Int("id", int(domainID)), wsl.Err(err))
		return err
	}

	// Инвалидируем кеш для этого домена
	uc.c.Delete(fmt.Sprintf("domain:id:%d", domainID))

	uc.l.DebugContext(ctx, "success increment domain access count", slog.Int("id", int(domainID)))
	return nil
}

// UpdateLastAccessDatetime обновляет время последнего доступа с инвалидацией кеша
func (uc *UseCase) UpdateDomainLastAccess(ctx context.Context, domainID uint, datetime time.Time) error {
	err := uc.etlDB.UpdateDomainLastAccess(ctx, domainID, datetime)
	if err != nil {
		uc.l.ErrorContext(ctx, "update domain last access datetime",
			wsl.Int("id", int(domainID)), wsl.Err(err))
		return err
	}

	// Инвалидируем кеш для этого домена
	uc.c.Delete(fmt.Sprintf("domain:id:%d", domainID))

	uc.l.DebugContext(ctx, "success update domain last access datetime",
		slog.Int("id", int(domainID)), slog.Time("datetime", datetime))
	return nil
}

// invalidateDomainCache инвалидирует все кеши связанные с доменом
func (uc *UseCase) invalidateDomainCache(domain *models.Domain) {
	// Инвалидируем кеш по ID
	uc.c.Delete(fmt.Sprintf("domain:id:%d", domain.ID))

	// Инвалидируем кеш по адресу
	uc.c.Delete(fmt.Sprintf("domain:addr:%s:%d", domain.IP, domain.Port))

	// Инвалидируем кеш всех доменов
	uc.c.Delete("domains:all")

	uc.l.Debug("invalidated domain cache",
		slog.Int("id", int(domain.ID)),
		slog.String("ip", domain.IP),
		slog.Int("port", domain.Port),
	)
}
