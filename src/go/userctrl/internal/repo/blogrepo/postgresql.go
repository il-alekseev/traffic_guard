package blogrepo

import (
	"context"
	"fmt"
	"log/slog"

	"userctrl/pkg/bizlogger"
	"userctrl/pkg/pgorm"
	"userctrl/pkg/slogger"

	"gorm.io/gorm"
)

type BlogRepoPG struct {
	db pgorm.Interface
}

// New — инициализация слоя репозитория
func New(db pgorm.Interface) *BlogRepoPG {
	repo := BlogRepoPG{
		db: db,
	}
	return &repo
}

// InsertLogInDB — сохраняет запись бизнес-лога (bizlogger.BusinessLog) в БД в рамках транзакции;
// при ошибке логирует её и возвращает обёрнутую ошибку с контекстом
func (r *BlogRepoPG) InsertLogInDB(ctx context.Context, bizlog *bizlogger.BusinessLog) error {
	return r.db.WithTx(ctx, func(tx *gorm.DB) error {
		if err := tx.Create(&bizlog).Error; err != nil {
			err = fmt.Errorf("failed to write bizlog: %w", err)
			slog.ErrorContext(ctx, err.Error())
			return slogger.WrapError(ctx, err)
		}
		return nil
	})
}
