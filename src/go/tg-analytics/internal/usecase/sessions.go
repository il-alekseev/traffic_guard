package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"tg-an/internal/controllers/http/v1/dto"
	"tg-an/internal/controllers/http/v1/values"
	"tg-an/internal/models"
	"tg-an/pkg/slogger/wsl"
	"tg-an/pkg/trparser"
)

// GetSessions возвращает список сессий с пагинацией
func (u *Usecase) GetSessions(ctx context.Context, userMeta *models.UserMeta, tr *trparser.TimeRange, f models.SessionFilter, search string, p models.Pagination, s models.Sorting) ([]dto.Session, int64, error) {
	method := "GetSessions"
	u.l.InfoContext(ctx,
		method,
		slog.Any("time_range", tr),
		slog.Any("filter", f),
		slog.String("search", search),
		slog.Any("pagination", p),
		slog.Any("sorting", s),
	)

	// Проверяем роль пользователя
	// И если она CA, то фильтруем по хосту
	if userMeta.ShortRole == values.ContextAdmin {
		f.HostName = userMeta.ContextID
	}

	sessions, total, err := u.db.GetSessions(ctx, tr, f, search, p, s)
	if err != nil {
		err = fmt.Errorf("%s: failed to get sessions: %w", method, err)
		u.l.ErrorContext(ctx, "Database operation failed",
			wsl.String("method", method),
			wsl.String("error", err.Error()),
		)
		return nil, 0, err
	}

	u.l.InfoContext(ctx, "Sessions retrieved",
		slog.String("method", method),
		slog.Int("count", len(sessions)),
		slog.Int64("total", total),
	)
	return sessions, total, nil
}
