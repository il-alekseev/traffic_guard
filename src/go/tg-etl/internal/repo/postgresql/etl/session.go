package postgresql

import (
	"context"
	"fmt"
	"tg-etl/internal/models"
	"tg-etl/pkg/slogger"

	"gorm.io/gorm"
)

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
