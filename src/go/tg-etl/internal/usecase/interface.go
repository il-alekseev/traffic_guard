package usecase

import (
	"context"
	"tg-etl/internal/models"

	"github.com/google/uuid"
)

type UsecaseInterface interface {
	QueryUsecase
	ProcessorUseCase
}

// UseCaseInterface определяет методы для обработки ETL данных
type ProcessorUseCase interface {
	ProcessNewLogs(ctx context.Context) (bool, error)
}

// QueryUsecase - композитный интерфейс для всех запросов
type QueryUsecase interface {
	ActionUseCase
	SourceUseCase
	DomainUseCase
	DeviceUseCase
	SessionUseCase
	IdsUseCase
	CategoryUseCase
	ListUseCase
	URLUseCase
	LastLogUseCase
}

// ActionUseCase
type ActionUseCase interface {
	GetActionByDomainID(ctx context.Context, id uint) (*models.Action, error)
}

// SourceUseCase определяет методы для работы с источниками
type SourceUseCase interface {
	GetSourceByAddr(ctx context.Context, ip string) (*models.Source, error)
	GetSourceByID(ctx context.Context, id uint) (*models.Source, error)
	CreateSource(ctx context.Context, source models.Source) error
	UpdateSource(ctx context.Context, source models.Source) error
	GetSources(ctx context.Context) ([]models.Source, error)
	DeleteSource(ctx context.Context, id uint) error
}

// DomainUseCase определяет методы для работы с доменами
type DomainUseCase interface {
	GetDomainByID(ctx context.Context, id uint) (*models.Domain, error)
	GetDomainByAddr(ctx context.Context, ip string, port int) (*models.Domain, error)
	GetDomainByPath(ctx context.Context, path string) (*models.Domain, error)
	GetDomainByRequestID(ctx context.Context, requestID uuid.UUID) (*models.Domain, error)
	CreateDomain(ctx context.Context, domain models.Domain) error
	UpdateDomain(ctx context.Context, domain models.Domain) error
	GetDomains(ctx context.Context) ([]models.Domain, error)
}

// DeviceUseCase определяет методы для работы с устройствами
type DeviceUseCase interface {
	GetDeviceByID(ctx context.Context, id int) (*models.Device, error)
	CreateDevice(ctx context.Context, device models.Device) error
	GetDevices(ctx context.Context) ([]models.Device, error)
}

// SessionUseCase определяет методы для работы с сессиями
type SessionUseCase interface {
	CreateSession(ctx context.Context, session models.Session) error
	GetSessionByID(ctx context.Context, id uint) (*models.Session, error)
	GetSessions(ctx context.Context, limit, offset int) ([]models.Session, error)
}

// IdsUseCase определяет методы для работы с IDS логами
type IdsUseCase interface {
	GetLogs(ctx context.Context, start *models.LastLog, count uint) ([]models.IdsLog, error)
}

// CategoryUseCase определяет методы для работы с категориями
type CategoryUseCase interface {
	CreatePredefinedCategories(ctx context.Context, categories []models.Category) error
	GetCategoryByID(ctx context.Context, id uint) (*models.Category, error)
	GetCategoryByName(ctx context.Context, name string) (*models.Category, error)
	GetCategoryIDByName(ctx context.Context, name string) (uint, error)
	GetCategories(ctx context.Context) ([]models.Category, error)
}

// ListUseCase определяет методы для работы со списками
type ListUseCase interface {
	AddDomainToList(ctx context.Context, domain models.Domain, list models.ListType) error
	GetListByDomainID(ctx context.Context, id uint) (*string, error)
}

// URL методы
type URLUseCase interface {
	CreateURL(ctx context.Context, url models.URL) error
	GetURLByPath(ctx context.Context, path string) (*models.URL, error)
	GetURLByPathDomain(ctx context.Context, path string, id uint) (*models.URL, error)
	UpdateURLByRequestID(ctx context.Context, requestID uuid.UUID, newURL models.URL) error
}

// Работа с последним обработанным логом
type LastLogUseCase interface {
	GetLastLog(ctx context.Context) (*models.LastLog, error)
	CreateOrUpdateLastLog(ctx context.Context, log models.LastLog) error
}
