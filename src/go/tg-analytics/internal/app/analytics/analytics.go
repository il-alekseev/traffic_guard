package an

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"tg-an/config"
	httpserver "tg-an/internal/controllers/http/v1"
	"tg-an/internal/models"
	postresql "tg-an/internal/repo/postgresql"
	"tg-an/internal/usecase"
	"tg-an/pkg/blog"
	"tg-an/pkg/cslogger"
	"tg-an/pkg/pgorm/pgorm"
	"tg-an/pkg/slogger/wsl"
	"time"

	httptransport "github.com/go-openapi/runtime/client"
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

	logger := *cslogger.NewColorLogger()
	//logger := *slog.Default()

	logger.InfoContext(ctx, "analytics service", wsl.Info("Creating DB connection"))

	dsnETL := getDSN(cfg.ETL.Host, cfg.ETL.User, cfg.ETL.Pass,
		cfg.ETL.DBName, cfg.ETL.Port, cfg.ETL.SSLMode)
	db, err := createConnection(ctx, dsnETL, cfg.ETL.PoolMax, logger, false, models.Session{})
	if err != nil {
		logger.ErrorContext(ctx, "analytics service", wsl.String("create db connection error", err.Error()))
		return
	}
	dsnMetrics := getDSN(cfg.Metrics.Host, cfg.Metrics.User, cfg.Metrics.Pass,
		cfg.Metrics.DBName, cfg.Metrics.Port, cfg.Metrics.SSLMode)
	mdb, err := createConnection(ctx, dsnMetrics, cfg.Metrics.PoolMax, logger, false, models.StatsJSON{})
	if err != nil {
		logger.ErrorContext(ctx, "analytics service", wsl.String("create mdb connection error", err.Error()))
		return
	}
	// Инициализация клиента для бизнес-логов
	blclient := blog.New(
		httptransport.New(
			cfg.BlogServ.Host+":"+cfg.BlogServ.Port,
			"/",
			[]string{cfg.BlogServ.Proto},
		),
		nil,
	)

	u := usecase.New(db, mdb, blclient, logger)

	server := httpserver.New(cfg, logger, u)

	go func() {
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			logger.ErrorContext(ctx, "analytics service", wsl.Err(err))
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
		logger.ErrorContext(ctx, "analytics service", wsl.Err(err))
	}
	server.Stop(ctx)
	logger.Info("analytics http server stopped")

}
