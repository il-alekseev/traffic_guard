package postgresql

import (
	"cmd/etl/internal/models"
	"cmd/etl/pkg/pgorm"
	"cmd/etl/pkg/slogger"
	"context"
	"fmt"
	"log/slog"
)

type KSURepoPG struct {
	db pgorm.Interface
	l  slog.Logger
}

func NewKSURepoPG(db pgorm.Interface, l slog.Logger) *KSURepoPG {
	return &KSURepoPG{
		db: db,
		l:  l,
	}
}

func (r KSURepoPG) GetLogs(ctx context.Context, start *models.IdsLog, maxCount uint) ([]models.IdsLog, error) {
	var logs []models.IdsLog
	query := r.db.GetDB().WithContext(ctx).Model(&models.IdsLog{})
	if start != nil {
		query = query.Where("id > ?", start.ID).Where("timestamp >= ?", start.Timestamp)
	}
	query = query.Order("id ASC").Limit(int(maxCount))
	err := query.Find(&logs).Error
	if err != nil {
		err = fmt.Errorf("failed to get devices: %w", err)
		return nil, slogger.WrapError(ctx, err)
	}
	return logs, nil
}
