package postgresql

import (
	"context"
	"fmt"
	"tg-etl/internal/models"
	"tg-etl/pkg/slogger"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

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

func (r *ELTRepoPG) UpdateDomainByID(ctx context.Context, domain models.Domain, id int) error {
	err := r.db.WithTx(ctx, func(tx *gorm.DB) error {
		// Создаем map для обновления только переданных полей
		updates := make(map[string]interface{})

		if domain.IP != "" {
			updates["ip"] = domain.IP
		}
		if domain.Country != "" {
			updates["country"] = domain.Country
		}
		if domain.Path != "" {
			updates["path"] = domain.Path
		}
		if domain.Port != 0 {
			updates["port"] = domain.Port
		}
		if domain.CategoryID != 0 {
			updates["category_id"] = domain.CategoryID
		}
		if domain.ActionID != 0 {
			updates["action_id"] = domain.ActionID
		}
		if domain.AnalysisCount != 0 {
			updates["analysis_count"] = domain.AnalysisCount
		}
		if domain.NegRate != 0 {
			updates["neg_rate"] = domain.NegRate
		}
		if !domain.CategorizedAt.IsZero() {
			updates["categorized_at"] = domain.CategorizedAt
		}

		// Если нет полей для обновления - выходим
		if len(updates) == 0 {
			return nil
		}

		result := tx.Model(&models.Domain{}).Where("id = ?", id).Updates(updates)
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

func (r ELTRepoPG) GetDomainByRequestID(ctx context.Context, requestID uuid.UUID) (*models.Domain, error) {
	var domain models.Domain

	// Используем JOIN для связи таблиц URL и Domain через domain_id
	err := r.db.GetDB().WithContext(ctx).
		Table("domains").
		Joins("JOIN urls ON domains.id = urls.domain_id").
		Where("urls.request_id = ?", requestID).
		First(&domain).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		err = fmt.Errorf("failed to get domain by request ID %s: %w", requestID, err)
		return nil, slogger.WrapError(ctx, err)
	}

	return &domain, nil
}
