package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"tg-an/internal/controllers/http/v1/values"
	"tg-an/internal/models"
	"tg-an/pkg/slogger/wsl"
)

func (u *Usecase) GetDevices(ctx context.Context, userMeta *models.UserMeta) ([]string, error) {
	method := "GetDevices"
	u.l.InfoContext(ctx,
		method,
	)

	var hostname = ""
	// Проверяем роль пользователя
	// И если она CA, то фильтруем по хосту
	if userMeta.ShortRole == values.ContextAdmin {
		hostname = userMeta.ContextID
	}

	devices, err := u.db.GetDevices(ctx, hostname)
	if err != nil {
		err = fmt.Errorf("%s: failed to get devices: %w", method, err)
		u.l.ErrorContext(ctx, "Database operation failed",
			wsl.String("method", method),
			wsl.String("error", err.Error()),
		)
		return nil, err
	}
	u.l.InfoContext(ctx, "Devices retrieved",
		slog.String("method", method),
		slog.Int("devices_count", len(devices)),
	)
	return devices, nil
}

func (u *Usecase) GetContentCategories(ctx context.Context) ([]string, error) {
	method := "GetContentCategories"
	u.l.InfoContext(ctx,
		method,
	)
	categories, err := u.db.GetContentCategories(ctx)
	if err != nil {
		err = fmt.Errorf("%s: failed to get categories: %w", method, err)
		u.l.ErrorContext(ctx, "Database operation failed",
			wsl.String("method", method),
			wsl.String("error", err.Error()),
		)
		return nil, err
	}
	u.l.InfoContext(ctx, "Categories retrieved",
		slog.String("method", method),
		slog.Int("categories_count", len(categories)),
	)
	return categories, nil
}
