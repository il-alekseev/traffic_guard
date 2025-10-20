package postgresql

import (
	"cmd/etl/internal/models"
	"cmd/etl/pkg/pgorm"
	"context"
	"fmt"
	"log/slog"

	"gorm.io/gorm"
)

type RepoPG struct {
	db pgorm.Interface
	// TODO: переделать под slog
	logger slog.Logger
}

func New(db pgorm.Interface, l slog.Logger) *RepoPG {
	return &RepoPG{
		db:     db,
		logger: l,
	}
}

func (r *RepoPG) GetSession(ctx context.Context) (*models.IdsLog, error) {
	return nil, nil
}

func (r *RepoPG) CreateLog(ctx context.Context, l *models.IdsLog) error {
	err := r.db.WithTx(ctx, func(tx *gorm.DB) error {
		return tx.Create(l).Error
	})
	return err
}

func (r *RepoPG) GetLogs(ctx context.Context) ([]models.IdsLog, error) {
	var logs []models.IdsLog
	err := r.db.GetDB().WithContext(ctx).Find(&logs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get logs: %w", err)
	}
	return logs, nil
}
