package postgresql

import (
	"context"
	"errors"
	"fmt"
	"tg-etl/internal/models"
	"tg-etl/pkg/slogger"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GetURLByPath получает URL по пути
func (r *ELTRepoPG) GetURLByPath(ctx context.Context, path string) (*models.URL, error) {
	var url models.URL

	err := r.db.GetDB().WithContext(ctx).
		Where("path = ?", path).
		First(&url).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		err = fmt.Errorf("failed to get URL by path %s: %w", path, err)
		return nil, slogger.WrapError(ctx, err)
	}

	return &url, nil
}

func (r *ELTRepoPG) GetURLByPathDomain(ctx context.Context, path string, id uint) (*models.URL, error) {
	var url models.URL

	err := r.db.GetDB().WithContext(ctx).
		Where("path = ? and domain_id = ?", path, id).
		First(&url).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		err = fmt.Errorf("failed to get URL by path %s: %w", path, err)
		return nil, slogger.WrapError(ctx, err)
	}

	return &url, nil
}

// CreateURL создает новый URL, если он еще не существует
func (r *ELTRepoPG) CreateURL(ctx context.Context, url models.URL) error {
	err := r.db.GetDB().WithContext(ctx).Create(&url).Error
	if err != nil {
		err = fmt.Errorf("failed to create URL: %w", err)
		return slogger.WrapError(ctx, err)
	}

	return nil
}

func (r *ELTRepoPG) GetURLByRequestID(ctx context.Context, requestID uuid.UUID) (*models.URL, error) {
	var url models.URL
	err := r.db.GetDB().WithContext(ctx).
		Where("request_id = ?", requestID).
		First(&url).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get URL by request ID %s: %w", requestID, err)
	}

	return &url, nil
}

func (r *ELTRepoPG) UpdateURLByRequestID(ctx context.Context, requestID uuid.UUID, newURL models.URL) error {
	err := r.db.WithTx(ctx, func(tx *gorm.DB) error {
		// Создаем map для обновления только переданных полей
		updates := make(map[string]interface{})
		// Обновляем только те поля, которые были явно заданы (не zero value)
		if !newURL.GetCategoryAt.IsZero() {
			updates["get_category_at"] = newURL.GetCategoryAt
		}
		if !newURL.GetMetaDataAt.IsZero() {
			updates["get_metadata_at"] = newURL.GetMetaDataAt
		}
		if !newURL.PutKafkaAt.IsZero() {
			updates["put_kafka_at"] = newURL.PutKafkaAt
		}
		if !newURL.IDSLogsAt.IsZero() {
			updates["ids_logs_at"] = newURL.IDSLogsAt
		}

		// Если нет полей для обновления - выходим
		if len(updates) == 0 {
			return nil
		}

		result := tx.Model(&models.URL{}).Where("request_id = ?", requestID).Updates(updates)
		if result.Error != nil {
			return fmt.Errorf("failed to update URL %d: %w", requestID, result.Error)
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("URL with request_id %d not found", requestID)
		}
		return nil
	})
	if err != nil {
		return slogger.WrapError(ctx, err)
	}
	return nil
}
