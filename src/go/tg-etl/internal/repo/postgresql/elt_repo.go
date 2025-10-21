package postgresql

import (
	"cmd/etl/internal/models"
	"cmd/etl/pkg/pgorm"
	"cmd/etl/pkg/slogger"
	"context"
	"fmt"
	"log/slog"
	"time"

	"gorm.io/gorm"
)

type ELTRepoPG struct {
	db pgorm.Interface
	l  slog.Logger
}

func NewELTRepoPG(db pgorm.Interface, l slog.Logger) *ELTRepoPG {
	return &ELTRepoPG{
		db: db,
		l:  l,
	}
}

func (r *ELTRepoPG) CreateDevice(ctx context.Context, d models.Device) error {
	err := r.db.WithTx(ctx, func(tx *gorm.DB) error {
		if err := tx.Create(&d).Error; err != nil {
			return fmt.Errorf("failed to create device: %w", err)
		}
		return nil
	})
	if err != nil {
		return slogger.WrapError(ctx, err)
	}
	return nil
}

func (r *ELTRepoPG) GetDeviceByID(ctx context.Context, id int) (*models.Device, error) {
	var d models.Device
	err := r.db.GetDB().WithContext(ctx).Where("id = ?", id).First(&d).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		err = fmt.Errorf("failed to get device: %w", err)
		return nil, slogger.WrapError(ctx, err)
	}
	return &d, nil
}

func (r *ELTRepoPG) GetDevices(ctx context.Context) ([]models.Device, error) {
	var devices []models.Device
	err := r.db.GetDB().WithContext(ctx).Find(&devices).Error
	if err != nil {
		err = fmt.Errorf("failed to get devices: %w", err)
		return nil, slogger.WrapError(ctx, err)
	}
	return devices, nil
}

func (r *ELTRepoPG) GetSourceByAddr(ctx context.Context, ip string) (*models.Source, error) {
	var source models.Source
	err := r.db.GetDB().WithContext(ctx).Where("ip = ?", ip).First(&source).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		err = fmt.Errorf("failed to get source by addr %s: %w", ip, err)
		return nil, slogger.WrapError(ctx, err)
	}
	return &source, nil
}

func (r *ELTRepoPG) GetSourceByID(ctx context.Context, id uint) (*models.Source, error) {
	var source models.Source
	err := r.db.GetDB().WithContext(ctx).First(&source, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		err = fmt.Errorf("failed to get source by id %d: %w", id, err)
		return nil, slogger.WrapError(ctx, err)
	}
	return &source, nil
}

func (r *ELTRepoPG) CreateSource(ctx context.Context, source models.Source) error {
	err := r.db.WithTx(ctx, func(tx *gorm.DB) error {
		if err := tx.Create(&source).Error; err != nil {
			return fmt.Errorf("failed to create source %s: %w", source.IP, err)
		}
		return nil
	})
	if err != nil {
		return slogger.WrapError(ctx, err)
	}
	return nil
}

func (r *ELTRepoPG) UpdateSource(ctx context.Context, source models.Source) error {
	err := r.db.WithTx(ctx, func(tx *gorm.DB) error {
		result := tx.Save(&source)
		if result.Error != nil {
			return fmt.Errorf("failed to update source %d: %w", source.ID, result.Error)
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("source with id %d not found", source.ID)
		}
		return nil
	})
	if err != nil {
		return slogger.WrapError(ctx, err)
	}
	return nil
}

func (r *ELTRepoPG) GetSources(ctx context.Context) ([]models.Source, error) {
	var sources []models.Source
	err := r.db.GetDB().WithContext(ctx).Find(&sources).Error
	if err != nil {
		err = fmt.Errorf("failed to get sources: %w", err)
		return nil, slogger.WrapError(ctx, err)
	}
	return sources, nil
}

func (r *ELTRepoPG) DeleteSource(ctx context.Context, id uint) error {
	err := r.db.WithTx(ctx, func(tx *gorm.DB) error {
		result := tx.Delete(&models.Source{}, id)
		if result.Error != nil {
			return fmt.Errorf("failed to delete source %d: %w", id, result.Error)
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("source with id %d not found", id)
		}
		return nil
	})
	if err != nil {
		return slogger.WrapError(ctx, err)
	}
	return nil
}

func (r *ELTRepoPG) GetDomainByID(ctx context.Context, id uint) (*models.Domain, error) {
	var domain models.Domain
	err := r.db.GetDB().WithContext(ctx).First(&domain, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		err = fmt.Errorf("failed to get domain by id %d: %w", id, err)
		return nil, slogger.WrapError(ctx, err)
	}
	return &domain, nil
}

func (r *ELTRepoPG) GetDomainByAddr(ctx context.Context, ip string, port int) (*models.Domain, error) {
	var domain models.Domain
	err := r.db.GetDB().WithContext(ctx).Where("ip = ? AND port = ?", ip, port).First(&domain).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		err = fmt.Errorf("failed to get domain by addr %s:%d: %w", ip, port, err)
		return nil, slogger.WrapError(ctx, err)
	}
	return &domain, nil
}

func (r *ELTRepoPG) GetDomainByPath(ctx context.Context, path string) (*models.Domain, error) {
	var domain models.Domain
	err := r.db.GetDB().WithContext(ctx).Where("path = ?", path).First(&domain).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		err = fmt.Errorf("failed to get domain by path %s: %w", path, err)
		return nil, slogger.WrapError(ctx, err)
	}
	return &domain, nil
}

func (r *ELTRepoPG) CreateDomain(ctx context.Context, domain models.Domain) error {
	err := r.db.WithTx(ctx, func(tx *gorm.DB) error {
		if err := tx.Create(&domain).Error; err != nil {
			return fmt.Errorf("failed to create domain %s:%d: %w", domain.IP, domain.Port, err)
		}
		return nil
	})
	if err != nil {
		return slogger.WrapError(ctx, err)
	}
	return nil
}

func (r *ELTRepoPG) UpdateDomain(ctx context.Context, domain models.Domain) error {
	err := r.db.WithTx(ctx, func(tx *gorm.DB) error {
		result := tx.Save(&domain)
		if result.Error != nil {
			return fmt.Errorf("failed to update domain %d: %w", domain.ID, result.Error)
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("domain with id %d not found", domain.ID)
		}
		return nil
	})
	if err != nil {
		return slogger.WrapError(ctx, err)
	}
	return nil
}

func (r *ELTRepoPG) GetDomains(ctx context.Context) ([]models.Domain, error) {
	var domains []models.Domain
	err := r.db.GetDB().WithContext(ctx).Find(&domains).Error
	if err != nil {
		err = fmt.Errorf("failed to get domains: %w", err)
		return nil, slogger.WrapError(ctx, err)
	}
	return domains, nil
}

func (r *ELTRepoPG) IncrementDomainAccessCount(ctx context.Context, domainID uint) error {
	err := r.db.WithTx(ctx, func(tx *gorm.DB) error {
		result := tx.Model(&models.Domain{}).
			Where("id = ?", domainID).
			Update("access_count", gorm.Expr("access_count + ?", 1))
		if result.Error != nil {
			return fmt.Errorf("failed to increment access count for domain %d: %w", domainID, result.Error)
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("domain with id %d not found", domainID)
		}
		return nil
	})
	if err != nil {
		return slogger.WrapError(ctx, err)
	}
	return nil
}

func (r *ELTRepoPG) UpdateDomainLastAccess(ctx context.Context, domainID uint, datetime time.Time) error {
	err := r.db.WithTx(ctx, func(tx *gorm.DB) error {
		result := tx.Model(&models.Domain{}).
			Where("id = ?", domainID).
			Update("last_access_datetime", datetime)
		if result.Error != nil {
			return fmt.Errorf("failed to update last access datetime for domain %d: %w", domainID, result.Error)
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("domain with id %d not found", domainID)
		}
		return nil
	})
	if err != nil {
		return slogger.WrapError(ctx, err)
	}
	return nil
}

func (r *ELTRepoPG) DeleteDomain(ctx context.Context, id uint) error {
	err := r.db.WithTx(ctx, func(tx *gorm.DB) error {
		result := tx.Delete(&models.Domain{}, id)
		if result.Error != nil {
			return fmt.Errorf("failed to delete domain %d: %w", id, result.Error)
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("domain with id %d not found", id)
		}
		return nil
	})
	if err != nil {
		return slogger.WrapError(ctx, err)
	}
	return nil
}

// Дополнительные методы для Domain

func (r *ELTRepoPG) IncrementDomainAnalysisAttempts(ctx context.Context, domainID uint) error {
	err := r.db.WithTx(ctx, func(tx *gorm.DB) error {
		result := tx.Model(&models.Domain{}).
			Where("id = ?", domainID).
			Update("analysis_attemps_count", gorm.Expr("analysis_attemps_count + ?", 1))
		if result.Error != nil {
			return fmt.Errorf("failed to increment analysis attempts for domain %d: %w", domainID, result.Error)
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("domain with id %d not found", domainID)
		}
		return nil
	})
	if err != nil {
		return slogger.WrapError(ctx, err)
	}
	return nil
}

func (r *ELTRepoPG) IncrementDomainContentAnalysis(ctx context.Context, domainID uint) error {
	err := r.db.WithTx(ctx, func(tx *gorm.DB) error {
		result := tx.Model(&models.Domain{}).
			Where("id = ?", domainID).
			Update("content_analysis_counter", gorm.Expr("content_analysis_counter + ?", 1))
		if result.Error != nil {
			return fmt.Errorf("failed to increment content analysis for domain %d: %w", domainID, result.Error)
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("domain with id %d not found", domainID)
		}
		return nil
	})
	if err != nil {
		return slogger.WrapError(ctx, err)
	}
	return nil
}

func (r *ELTRepoPG) UpdateDomainDecision(ctx context.Context, domainID uint, decisionID uint, decisionTime time.Time) error {
	err := r.db.WithTx(ctx, func(tx *gorm.DB) error {
		result := tx.Model(&models.Domain{}).
			Where("id = ?", domainID).
			Updates(map[string]interface{}{
				"decision_id":       decisionID,
				"decision_datetime": decisionTime,
			})
		if result.Error != nil {
			return fmt.Errorf("failed to update decision for domain %d: %w", domainID, result.Error)
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("domain with id %d not found", domainID)
		}
		return nil
	})
	if err != nil {
		return slogger.WrapError(ctx, err)
	}
	return nil
}

func (r *ELTRepoPG) CreateSession(ctx context.Context, session models.Session) error {
	err := r.db.WithTx(ctx, func(tx *gorm.DB) error {
		if err := tx.Create(&session).Error; err != nil {
			return fmt.Errorf("failed to create session %s: %w", session.IP, err)
		}
		return nil
	})
	if err != nil {
		return slogger.WrapError(ctx, err)
	}
	return nil
}
func (r *ELTRepoPG) GetSessionByID(ctx context.Context, id uint) (*models.Session, error) {
	var session models.Session
	err := r.db.GetDB().WithContext(ctx).First(&session, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		err = fmt.Errorf("failed to get session by id %d: %w", id, err)
		return nil, slogger.WrapError(ctx, err)
	}
	return &session, nil
}
func (r *ELTRepoPG) GetSessions(ctx context.Context) ([]models.Session, error) {
	var sessions []models.Session
	err := r.db.GetDB().WithContext(ctx).Find(&sessions).Error
	if err != nil {
		err = fmt.Errorf("failed to get sessions: %w", err)
		return nil, slogger.WrapError(ctx, err)
	}
	return sessions, nil
}
