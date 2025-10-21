package etl

import (
	"cmd/etl/internal/models"
	"context"
	"time"
)

// ETLUseCase интерфейс для всех use cases ETL системы
type UsecaseInterface interface {
	SourceUseCase
	DomainUseCase
	DeviceUseCase
	SessionUseCase
	IdsUseCase
	ProcessUsecase
}

// SourceUseCase определяет методы для работы с источниками
type SourceUseCase interface {
	// GetSourceByAddr получает источник по IP и порту с кешированием
	GetSourceByAddr(ctx context.Context, ip string) (*models.Source, error)
	// GetSourceByID получает источник по ID с кешированием
	GetSourceByID(ctx context.Context, id uint) (*models.Source, error)
	// CreateSource создает новый источник с инвалидацией кеша
	CreateSource(ctx context.Context, source models.Source) error
	// UpdateSource обновляет источник с инвалидацией кеша
	UpdateSource(ctx context.Context, source models.Source) error
	// GetSources получает все источники с кешированием
	GetSources(ctx context.Context) ([]models.Source, error)
	// DeleteSource удаляет источник с инвалидацией кеша
	DeleteSource(ctx context.Context, id uint) error
}

// DomainUseCase определяет методы для работы с доменами
type DomainUseCase interface {
	// GetDomainByID получает домен по ID с кешированием
	GetDomainByID(ctx context.Context, id uint) (*models.Domain, error)
	// GetDomainByAddr получает домен по IP и порту с кешированием
	GetDomainByAddr(ctx context.Context, ip string, port int) (*models.Domain, error)
	// CreateDomain создает новый домен с инвалидацией кеша
	CreateDomain(ctx context.Context, domain models.Domain) error
	// UpdateDomain обновляет домен с инвалидацией кеша
	UpdateDomain(ctx context.Context, domain models.Domain) error
	// GetDomains получает все домены с кешированием
	GetDomains(ctx context.Context) ([]models.Domain, error)
	// IncrementAccessCount инкрементирует счетчик доступа с инвалидацией кеша
	IncrementAccessCount(ctx context.Context, domainID uint) error
	// UpdateDomainLastAccess обновляет время последнего доступа
	UpdateDomainLastAccess(ctx context.Context, domainID uint, datetime time.Time) error
}

// DeviceUseCase определяет методы для работы с устройствами
type DeviceUseCase interface {
	// GetDeviceByID получает устройство по ID с кешированием
	GetDeviceByID(ctx context.Context, id int) (*models.Device, error)
	// CreateDevice создает новое устройство с инвалидацией кеша
	CreateDevice(ctx context.Context, device models.Device) error
	// GetDevices получает все устройства с кешированием
	GetDevices(ctx context.Context) ([]models.Device, error)
}

// SessionUseCase определяет методы для работы с сессиями
type SessionUseCase interface {
	// CreateSession создает новую сессию
	CreateSession(ctx context.Context, session models.Session) error
	// GetSessionByID получает сессию по ID
	GetSessionByID(ctx context.Context, id uint) (*models.Session, error)
	// GetSessions получает все сессии с пагинацией
	GetSessions(ctx context.Context, limit, offset int) ([]models.Session, error)
}

type IdsUseCase interface {
	GetLogs(ctx context.Context, start *models.IdsLog, count uint) ([]models.IdsLog, error)
}

type ProcessUsecase interface {
	// Основные методы обработки логов
	ProcessNewLogs(ctx context.Context) error
	ProcessLog(ctx context.Context, log models.IdsLog) error
	// Методы обработки сущностей
	ProcessSource(ctx context.Context, log models.IdsLog) (*models.Source, error)
	ProcessDomain(ctx context.Context, log models.IdsLog) (*models.Domain, error)
	ProcessDevice(ctx context.Context, log models.IdsLog) (*models.Device, error)
	ProcessSession(ctx context.Context, log models.IdsLog, source *models.Source, domain *models.Domain, device *models.Device) error
}
