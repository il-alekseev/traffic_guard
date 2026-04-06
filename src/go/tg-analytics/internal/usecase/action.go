package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"tg-an/internal/controllers/http/v1/values"
	"tg-an/internal/models"
	"tg-an/pkg/blog/operations"
	pkg "tg-an/pkg/models"
	"tg-an/pkg/slogger/wsl"

	"github.com/google/uuid"
)

func (u *Usecase) Act(ctx context.Context, userMeta *models.UserMeta, action, path string) error {
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
	uuidStr := uuid.New().String()

	newValue := make(map[string]any)
	newValue["action"] = action
	newValue["domain"] = path

	// Формируем информативную строку описания
	act_type, err := models.ParseActionType(action)
	if err != nil {
		err = fmt.Errorf("%s: failed parse act type: %s", method, action)
		u.l.ErrorContext(ctx, "Database operation failed",
			wsl.String("method", method),
			wsl.String("error", err.Error()),
		)
		return err
	}

	act_str := ""
	if act_type == models.ActionTypeAllow {
		act_str = "положительное"
	} else {
		act_str = "отрицательное"
	}

	record := pkg.DtoBusinessLog{
		Description: fmt.Sprintf("%s решение по домену %s", act_str, path),
		Entity:      values.DomainEntity,
		EntityID:    path,
		NewValue:    newValue,
		OldValue:    oldValue,
		EventType:   operation,
		Context:     userMeta.ContextID,
		UserName:    userMeta.Username,
		UserRole:    userMeta.ClientRole,
	}

	_, err = u.blclient.Operations.PostAPIV1Add(&operations.PostAPIV1AddParams{
		XCallerService: values.ThisServiceName,
		XRequestID:     uuidStr,
		Record:         &record,
		Context:        ctx,
	},
	)
	if err != nil {
		err = fmt.Errorf("%s: failed to blog event: %w", method, err)
		u.l.ErrorContext(ctx, "failed to blog event",
			wsl.String("method", method),
			wsl.String("error", err.Error()),
		)
		return nil
	}
	//

	u.l.InfoContext(ctx, "Act success",
		slog.String("method", method),
		slog.String("blog uuid", uuidStr),
		slog.String("action", action),
		slog.String("path", path),
		slog.String("user", userMeta.Username),
	)
	return nil
}
