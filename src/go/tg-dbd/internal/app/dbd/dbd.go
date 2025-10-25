package dbd

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"tg-dbd/config"
	httpserver "tg-dbd/internal/controllers/http/v1"
	"tg-dbd/internal/models"
	postresql "tg-dbd/internal/repo/postgresql"
	"tg-dbd/internal/usecase"
	"tg-dbd/pkg/pgorm/pgorm"
	"tg-dbd/pkg/slogger"
	"tg-dbd/pkg/slogger/wsl"
	"time"
)

func getDSN(host, user, pass, dbname, port, sslmode string) string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, pass, dbname, port, sslmode,
	)
}

func createConnection(ctx context.Context,
	dsn string,
	maxPoolSize int,
	l slog.Logger,
	autoMigrate bool,
	models ...interface{}) (*postresql.RepoPG, error) {
	db, err := pgorm.New(
		dsn,
		pgorm.MaxPoolSize(maxPoolSize),
		pgorm.AutoMigrate(autoMigrate),
		pgorm.Models(models...),
	)
	if err != nil {
		l.ErrorContext(ctx, "Error with creation connetcion", wsl.Err(err))
		return nil, err
	}
	l.InfoContext(ctx, "connection created successfully")
	if err := db.HealthCheck(); err != nil {
		l.ErrorContext(ctx, "database health check failed", wsl.Err(err))
		return nil, err
	} else {
		l.InfoContext(ctx, "database connection is healthy")
	}
	return postresql.New(db, l), nil
}

// TODO: красиво переписать RUN
func Run(cfg *config.Config) {
	ctx := context.Background()

	logLevel := slog.LevelDebug
	if cfg.Log.Level != "debug" {
		logLevel = slog.LevelInfo
	}
	slogger.InitLogging(logLevel)
	logger := *slog.Default()

	logger.InfoContext(ctx, "dashboard service", wsl.Info("Creating DB connection"))

	dsnETL := getDSN(cfg.PG.Host, cfg.PG.User, cfg.PG.Pass,
		cfg.PG.DBName, cfg.PG.Port, cfg.PG.SSLMode)
	db, err := createConnection(ctx, dsnETL, cfg.PG.PoolMax, logger, false, models.Session{})
	if err != nil {
		logger.ErrorContext(ctx, "dashboard service", wsl.String("create db connection error", err.Error()))
		return
	}
	dsnMetrics := getDSN(cfg.Metrics.Host, cfg.Metrics.User, cfg.Metrics.Pass,
		cfg.Metrics.DBName, cfg.Metrics.Port, cfg.Metrics.SSLMode)
	mdb, err := createConnection(ctx, dsnMetrics, cfg.Metrics.PoolMax, logger, false, models.StatsJSON{})
	if err != nil {
		logger.ErrorContext(ctx, "dashboard service", wsl.String("create mdb connection error", err.Error()))
		return
	}
	u := usecase.New(db, mdb, logger)

	server := httpserver.New(cfg, logger, u)

	go func() {
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			logger.ErrorContext(ctx, "dashboard service", wsl.Err(err))
		}
	}()
	// Создаем канал для получения сигналов
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt) // Подписываемся на сигнал прерывания (Ctrl+C)

	// Ожидаем сигнала
	<-signalChan
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Stop(ctx); err != nil {
		logger.ErrorContext(ctx, "dashboard service", wsl.Err(err))
	}
	server.Stop(ctx)
	logger.Info("dashboard http server stopped")

}
