package etl

import (
	"context"
	"log/slog"
	"tg-etl/config"
	u "tg-etl/internal/usecase"
	"tg-etl/pkg/slogger/wsl"
	"time"
)

type EtlController struct {
	u       u.UsecaseInterface
	refresh time.Duration
	l       slog.Logger
	cfg     config.EtlController
}

func New(cfg config.Config,
	u u.UsecaseInterface,
	l slog.Logger,
) *EtlController {
	return &EtlController{
		u:       u,
		refresh: time.Second * time.Duration(cfg.Refresh),
		l:       l,
		cfg:     cfg.EtlController,
	}
}

func (e *EtlController) Start(ctx context.Context) {
	e.l.InfoContext(ctx, "starting IDS log processor service",
		wsl.Int("refresh interval", int(e.refresh)))
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
