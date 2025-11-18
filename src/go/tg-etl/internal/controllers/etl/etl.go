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

	for {
		select {
		case <-ctx.Done():
			e.l.InfoContext(ctx, "stopping etl processor controller")
			return

		default:
			immediate, err := e.u.ProcessNewLogs(ctx)
			if err != nil {
				e.l.ErrorContext(ctx, "failed to process logs", wsl.Err(err))
			}

			// Если нужно немедленное выполнение, не ждем таймер
			if immediate {
				e.l.DebugContext(ctx, "immediate processing requested, continuing without delay")
				continue
			}

			// Создаем таймер для соблюдения интервала
			ticker := time.NewTicker(e.refresh)

			select {
			case <-ctx.Done():
				ticker.Stop()
				e.l.InfoContext(ctx, "stopping etl processor controller")
				return
			case <-ticker.C:
				ticker.Stop()
			}
		}
	}
}
