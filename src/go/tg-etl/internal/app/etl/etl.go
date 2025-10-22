package etl

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"tg-etl/config"
	"tg-etl/internal/controllers/etl"
	v1 "tg-etl/internal/controllers/http/v1"
	"tg-etl/internal/models"
	"tg-etl/internal/repo/postgresql"
	"tg-etl/internal/usecase"
	"tg-etl/pkg/pgorm"
	"tg-etl/pkg/slogger"
	"tg-etl/pkg/slogger/wsl"
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
	models ...interface{}) (*pgorm.Postgres, error) {
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
	return db, nil
}

func createKSUConnection(ctx context.Context, dsn string, maxPoolSize int, l slog.Logger, models ...interface{}) (*postgresql.KSURepoPG, error) {
	db, err := createConnection(ctx, dsn, maxPoolSize, l, false, models...)
	if err != nil {
		return nil, err
	}
	return postgresql.NewKSURepoPG(db, l), nil
}

func createETLConnection(ctx context.Context, dsn string, maxPoolSize int, l slog.Logger, models ...interface{}) (*postgresql.ELTRepoPG, error) {
	db, err := createConnection(ctx, dsn, maxPoolSize, l, true, models...)
	if err != nil {
		return nil, err
	}
	return postgresql.NewELTRepoPG(db, l), nil
}

func Run(cfg *config.Config) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	logLevel := slog.LevelDebug
	if cfg.Log.Level != "debug" {
		logLevel = slog.LevelInfo
	}
	slogger.InitLogging(logLevel)
	logger := *slog.Default()

	logger.InfoContext(ctx, "ETL service", wsl.Info("Creating DB connections"))

	dsnKSU := getDSN(cfg.KSU.Host, cfg.KSU.User, cfg.KSU.Pass,
		cfg.KSU.DBName, cfg.KSU.Port, cfg.KSU.SSLMode)
	dbKSU, err := createKSUConnection(ctx, dsnKSU, cfg.KSU.PoolMax, logger, models.IdsLog{})
	if err != nil {
		logger.ErrorContext(ctx, "ETL service", wsl.String("create KSU db connection error", err.Error()))
		return
	}

	dsnETL := getDSN(cfg.ETL.Host, cfg.ETL.User, cfg.ETL.Pass,
		cfg.ETL.DBName, cfg.ETL.Port, cfg.ETL.SSLMode)
	dbETL, err := createETLConnection(ctx, dsnETL, cfg.ETL.PoolMax, logger,
		// Автомиграция только базовых таблиц
		models.Category{},
		models.CategoryDomain{},
		models.Decision{},
		models.Detection{},
		models.Device{},
		models.Domain{},
		models.Session{},
		models.Source{},
		models.Status{},
	)
	if err != nil {
		logger.ErrorContext(ctx, "ETL service", wsl.String("create ETL db connection error", err.Error()))
		return
	}

	uc := usecase.New(cfg, dbKSU, dbETL, logger)

	etlController := etl.New(
		*cfg,
		uc,
		logger,
	)

	// HTTP сервер
	server := v1.New(cfg, logger)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Запуск HTTP сервера
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			logger.ErrorContext(ctx, "etl service", wsl.Err(err))
		}
	}()
	// Запуск  ETL процессора
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.InfoContext(ctx, "starting ETL service")
		etlController.Start(ctx)
	}()

	<-sigChan
	logger.InfoContext(ctx, "received shutdown signal")
	server.Stop(ctx)
	logger.InfoContext(ctx, "ETL http server stopped")
	cancel()
	wg.Wait()
	logger.InfoContext(ctx, "ETL service stopped completely")
}
