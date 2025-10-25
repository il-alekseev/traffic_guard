package usecase

import (
	"context"
	"tg-etl/internal/models"
	"time"
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

	GetDomainByID(ctx context.Context, id uint) (*models.Domain, error)
	GetDomainByAddr(ctx context.Context, ip string, port int) (*models.Domain, error)
	GetDomainByPath(ctx context.Context, path string) (*models.Domain, error)
	CreateDomain(ctx context.Context, domain models.Domain) error
	UpdateDomain(ctx context.Context, domain models.Domain) error
	GetDomains(ctx context.Context) ([]models.Domain, error)
	IncrementDomainAccessCount(ctx context.Context, domainID uint) error
	UpdateDomainLastAccess(ctx context.Context, domainID uint, datetime time.Time) error
	DeleteDomain(ctx context.Context, id uint) error

	CreateSession(ctx context.Context, session models.Session) error
	GetSessionByID(ctx context.Context, id uint) (*models.Session, error)
	GetSessions(ctx context.Context) ([]models.Session, error)

	CreateCategories(ctx context.Context, categories []string) error
	GetCategoryByID(ctx context.Context, id uint) (*models.Category, error)
	GetCategoryByName(ctx context.Context, name string) (*models.Category, error)
	GetCategories(ctx context.Context) ([]models.Category, error)
}
