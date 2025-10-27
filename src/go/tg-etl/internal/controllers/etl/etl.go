package etl

import (
	"context"
	"log/slog"
	"tg-etl/config"
	"tg-etl/internal/repo/category"
	"tg-etl/pkg/slogger/wsl"
	"time"
)

type EtlController struct {
	u       UsecaseInterface
	refresh time.Duration
	l       slog.Logger
}

func New(cfg config.Config,
	u UsecaseInterface,
	l slog.Logger,
) *EtlController {
	return &EtlController{
		u:       u,
		refresh: time.Second * time.Duration(cfg.Refresh),
		l:       l,
	}
}

func (e *EtlController) Start(ctx context.Context) {
	e.l.InfoContext(ctx, "starting IDS log processor service",
		wsl.Int("refresh interval", int(e.refresh)))
	// Проверяем, пуста ли таблица с категориями, если пуста, то добавляем категории
	categories, err := e.u.GetCategories(ctx)
	if err != nil {
		e.l.ErrorContext(ctx, "failed to get categories", wsl.Err(err))
		return
	}
	if len(categories) == 0 {
		err = e.u.CreateCategories(ctx, category.AllCategoryStrings())
		if err != nil {
			e.l.ErrorContext(ctx, "failed to create categories", wsl.Err(err))
			return
		}
	}
	// Создаем таймер для соблюдения интервала
	ticker := time.NewTicker(e.refresh)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			e.l.InfoContext(ctx, "stopping etl processor controller")
			return

		case <-ticker.C:
			if err := e.u.ProcessNewLogs(ctx); err != nil {
				e.l.ErrorContext(ctx, "failed to process logs", wsl.Err(err))
			}
		}
	}
}
