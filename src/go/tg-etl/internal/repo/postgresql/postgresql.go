package postgresql

import (
	"cmd/etl/internal/models"
	"cmd/etl/pkg/logger"
	"cmd/etl/pkg/pgorm"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type RepoPG struct {
	db pgorm.Interface
	// TODO: переделать под slog
	logger logger.Interface
}

func New(db pgorm.Interface, l logger.Interface) *RepoPG {
	return &RepoPG{
		db:     db,
		logger: l,
	}
}

func (r *RepoPG) GetSession(ctx *context.Context) (*models.IdsLog, error) {
	return nil, nil
}

func (r *RepoPG) CreateLog(ctx *context.Context, l *models.IdsLog) error {
	err := r.db.WithTx(*ctx, func(tx *gorm.DB) error {
		if err := tx.Create(l); err != nil {
			return fmt.Errorf("failed to create log: %v", err)
		}
		return nil
	})
	return err
}
