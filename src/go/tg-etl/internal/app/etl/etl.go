package etl

import (
	"cmd/etl/config"
	"cmd/etl/internal/controllers/etl"
	"cmd/etl/internal/models"
	"cmd/etl/internal/repo/postgresql"
	"cmd/etl/internal/usecase"
	"cmd/etl/pkg/pgorm"
	"cmd/etl/pkg/slogger"
	"cmd/etl/pkg/slogger/wsl"
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
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
	ctx := context.Background()
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		logger.InfoContext(ctx, "received shutdown signal")
		cancel()
	}()

	// Запуск процессора
	logger.InfoContext(ctx, "starting ETL service")
	etlController.Start(ctx)
	logger.InfoContext(ctx, "ETL service stopped")
}
