package postgresGorm

import (
	"context"
	"fiermon-blog/internal/models"
	"fiermon-blog/pkg/fslog"
)

// InsertRecord - метод вставки бизнес лога в БД
func (p *PostgresDB) InsertRecord(ctx context.Context, businessLog *models.BusinessLog) error {
	err := p.db.Create(businessLog).Error
	if err != nil {
		return fslog.WrapError(ctx, err)
	}
	return nil
}
