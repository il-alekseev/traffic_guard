package usecase

import (
	"context"
	"tg-etl/internal/models"

	"github.com/google/uuid"
)

type ETLRepoPGInterface interface {
	CreateDevice(ctx context.Context, device models.Device) error
	GetDeviceByID(ctx context.Context, id int) (*models.Device, error)
	GetDevices(ctx context.Context) ([]models.Device, error)

	GetSourceByAddr(ctx context.Context, ip string) (*models.Source, error)
	GetSourceByID(ctx context.Context, id uint) (*models.Source, error)
	CreateSource(ctx context.Context, source models.Source) error
	UpdateSource(ctx context.Context, source models.Source) error
	GetSources(ctx context.Context) ([]models.Source, error)
	DeleteSource(ctx context.Context, id uint) error

	GetDomainByRequestID(ctx context.Context, requestID uuid.UUID) (*models.Domain, error)
	GetDomainByID(ctx context.Context, id uint) (*models.Domain, error)
	GetDomainByAddr(ctx context.Context, ip string, port int) (*models.Domain, error)
	GetDomainByPath(ctx context.Context, path string) (*models.Domain, error)
	CreateDomain(ctx context.Context, domain models.Domain) error
	UpdateDomainByID(ctx context.Context, domain models.Domain, id int) error
	GetDomains(ctx context.Context) ([]models.Domain, error)
	DeleteDomain(ctx context.Context, id uint) error

	CreateSession(ctx context.Context, session models.Session) error
	GetSessionByID(ctx context.Context, id uint) (*models.Session, error)
	GetSessions(ctx context.Context) ([]models.Session, error)

	CreateCategories(ctx context.Context, categories []models.Category) error
	GetCategoryByID(ctx context.Context, id uint) (*models.Category, error)
	GetCategoryByName(ctx context.Context, name string) (*models.Category, error)
	GetCategories(ctx context.Context) ([]models.Category, error)

	AddDomainToList(ctx context.Context, domain models.Domain, list models.ListType) error
	GetListByDomainID(ctx context.Context, id uint) (*string, error)

	GetURLByPath(ctx context.Context, path string) (*models.URL, error)
	GetURLByPathDomain(ctx context.Context, path string, id uint) (*models.URL, error)
	GetURLByRequestID(ctx context.Context, requestID uuid.UUID) (*models.URL, error)
	CreateURL(ctx context.Context, url models.URL) error
	UpdateURLByRequestID(ctx context.Context, requestID uuid.UUID, newURL models.URL) error

	GetLastLog(ctx context.Context) (*models.LastLog, error)
	CreateOrUpdateLastLog(ctx context.Context, log models.LastLog) error
}

type KSURepoPGInterface interface {
	GetLogs(ctx context.Context, start *models.LastLog, maxCount uint) ([]models.IdsLog, error)
}
