package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"tg-dbd/pkg/slogger/wsl"
)

func (u *Usecase) GetDevices(ctx context.Context) ([]string, error) {
	method := "GetDevices"
	u.l.InfoContext(ctx,
		method,
	)
	devices, err := u.db.GetDevices(ctx)
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
