package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"tg-an/internal/controllers/http/v1/dto"
	"tg-an/internal/controllers/http/v1/values"
	"tg-an/internal/models"
	"tg-an/pkg/blog/operations"
	pkg "tg-an/pkg/models"
	"tg-an/pkg/slogger/wsl"
	"tg-an/pkg/trparser"

	"github.com/go-openapi/runtime"
)

func (u *Usecase) GetTopDetections(ctx context.Context, userMeta *models.UserMeta, tr *trparser.TimeRange, f models.DetectionFilter, a string, p models.Pagination) ([]dto.Detection, int64, error) {
	method := "GetTopDetections"
	u.l.InfoContext(ctx,
		method,
		slog.Any("time_range", tr),
		slog.Any("filter", f),
		slog.Any("pagination", p),
		slog.String("user", userMeta.Username),
	)

	// Проверяем роль пользователя
	// И если она CA, то фильтруем по хосту
	if userMeta.ShortRole == values.ContextAdmin {
		f.HostName = userMeta.ContextID
	}

	detections, total, err := u.db.GetTopDetections(ctx, tr, f, a, p)
	if err != nil {
		err = fmt.Errorf("%s: failed to get top detections: %w", method, err)
		u.l.ErrorContext(ctx, "Database operation failed",
			wsl.String("method", method),
			wsl.String("error", err.Error()),
		)
		return nil, 0, err
	}

	u.l.InfoContext(ctx, "Top detections retrieved",
		slog.String("method", method),
		slog.Int("count", len(detections)),
		slog.Int("total", int(total)),
		slog.String("user", userMeta.Username),
	)
	return detections, total, nil
}

// GetDetectionStat возвращает статистику обнаружений за указанный период
func (u *Usecase) GetDetectionStat(ctx context.Context, userMeta *models.UserMeta, tr *trparser.TimeRange, f models.DetectionFilter) (dto.DetectionStat, error) {
	method := "GetDetectionStat"
	u.l.InfoContext(ctx,
		method,
		slog.Any("time_range", tr),
		slog.Any("filter", f),
		slog.String("user", userMeta.Username),
	)

	// Проверяем роль пользователя
	// И если она CA, то фильтруем по хосту
	if userMeta.ShortRole == values.ContextAdmin {
		f.HostName = userMeta.ContextID
	}

	stat, err := u.db.GetDetectionStat(ctx, tr, f)
	if err != nil {
		err = fmt.Errorf("%s: failed to get detection stat: %w", method, err)
		u.l.ErrorContext(ctx, "Database operation failed",
			wsl.String("method", method),
			wsl.String("error", err.Error()),
		)
		return dto.DetectionStat{}, err
	}

	u.l.InfoContext(ctx, "detection stat retrieved",
		slog.String("method", method),
		slog.Int64("Detected", stat.Detected),
		slog.Int64("Allowed", stat.Allowed),
		slog.Int64("Denied", stat.Denied),
		slog.Int64("Unresolved", stat.Unresolved),
		slog.String("user", userMeta.Username),
	)
	return stat, nil
}

func (u *Usecase) Act(ctx context.Context, userMeta *models.UserMeta, authInfo runtime.ClientAuthInfoWriter, action, path string) error {
	method := "Act"
	u.l.InfoContext(ctx,
		method,
		wsl.String("action", action),
		slog.Any("path", path),
		slog.String("user", userMeta.Username),
	)

	// Проверяем права доступа для выполнения действий
	if userMeta.ShortRole != values.SystemAdmin && userMeta.ShortRole != values.ContextAdmin {
		err := fmt.Errorf("user %s does not have permission to perform actions", userMeta.Username)
		u.l.ErrorContext(ctx, "Permission denied",
			wsl.String("method", method),
			wsl.String("error", err.Error()),
		)
		return err
	}

	// Проверяем, является ли действие над доменом первым, если нет, то получаем значение предыдущего действия
	lastActionStr, err := u.db.GetDomainAction(ctx, path)
	if err != nil {
		err = fmt.Errorf("%s: failed get domain action: %w", method, err)
		u.l.ErrorContext(ctx, "Database operation failed",
			wsl.String("method", method),
			wsl.String("error", err.Error()),
		)
		return err
	}

	operation := "CREATE"
	oldValue := make(map[string]any)
	if lastActionStr != "" {
		lastAction, err := models.ParseActionType(lastActionStr)
		if err != nil {
			err = fmt.Errorf("%s: failed parse domain action: %w", method, err)
			return err
		}
		operation = "UPDATE"
		oldValue["path"] = path
		oldValue["action"] = lastAction.String()

	}

	err = u.db.Act(ctx, userMeta.Username, action, path)
	if err != nil {
		err = fmt.Errorf("%s: failed act with domain: %w", method, err)
		u.l.ErrorContext(ctx, "Database operation failed",
			wsl.String("method", method),
			wsl.String("error", err.Error()),
		)
		return err
	}

	// Запись события в бизнес-лог

	newValue := make(map[string]any)
	newValue["action"] = action
	newValue["domain"] = path

	record := pkg.DtoBusinessLog{
		Description: "Решение по домену",
		Entity:      "Domain",
		EntityID:    path,
		NewValue:    newValue,
		OldValue:    oldValue,
		EventType:   operation,
		Context:     userMeta.ContextID,
		UserName:    userMeta.Username,
		UserRole:    userMeta.ClientRole,
	}

	_, err = u.blclient.Operations.PostAPIV1Add(&operations.PostAPIV1AddParams{
		Record:  &record,
		Context: ctx,
	},
		authInfo,
	)
	if err != nil {
		err = fmt.Errorf("%s: failed to blog event: %w", method, err)
		u.l.ErrorContext(ctx, "failed to blog event",
			wsl.String("method", method),
			wsl.String("error", err.Error()),
		)
		return err
	}
	//

	u.l.InfoContext(ctx, "Act success",
		slog.String("method", method),
		slog.String("action", action),
		slog.String("path", path),
		slog.String("user", userMeta.Username),
	)
	return nil
}
