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
	"tg-etl/internal/repo/kafka"
	ksuRepo "tg-etl/internal/repo/postgresql"
	etlRepo "tg-etl/internal/repo/postgresql/etl"
	"tg-etl/internal/usecase"
	"tg-etl/internal/usecase/handlers"
	q "tg-etl/internal/usecase/query"
	"tg-etl/pkg/blog"
	"tg-etl/pkg/cslogger"
	"tg-etl/pkg/pgorm"
	"tg-etl/pkg/slogger/wsl"
	"tg-etl/pkg/usercontrol"

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

func createKSUConnection(ctx context.Context, dsn string, maxPoolSize int, l slog.Logger, models ...interface{}) (*ksuRepo.KSURepoPG, error) {
	db, err := createConnection(ctx, dsn, maxPoolSize, l, false, models...)
	if err != nil {
		return nil, err
	}
	return ksuRepo.NewKSURepoPG(db, l), nil
}

func createETLConnection(ctx context.Context, dsn string, maxPoolSize int, l slog.Logger, models ...interface{}) (*etlRepo.ELTRepoPG, error) {
	db, err := createConnection(ctx, dsn, maxPoolSize, l, true, models...)
	if err != nil {
		return nil, err
	}
	return etlRepo.NewELTRepoPG(db, l), nil
}

func Run(cfg *config.Config) {
	// Создаем контекст
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	logger := *cslogger.NewColorLogger()

	logger.InfoContext(ctx, "ETL service",
		wsl.Info("Creating ksu db connection"),
		wsl.String("host", cfg.KSU.Host),
		wsl.String("db", cfg.KSU.DBName),
		wsl.String("user", cfg.KSU.User),
	)
	dsnKSU := getDSN(cfg.KSU.Host, cfg.KSU.User, cfg.KSU.Pass,
		cfg.KSU.DBName, cfg.KSU.Port, cfg.KSU.SSLMode)
	dbKSU, err := createKSUConnection(ctx, dsnKSU, cfg.KSU.PoolMax, logger, models.IdsLog{})
	if err != nil {
		logger.ErrorContext(ctx, "ETL service", wsl.String("create KSU db connection error", err.Error()))
		return
	}

	logger.InfoContext(ctx, "ETL service",
		wsl.Info("Creating etl db connection"),
		wsl.String("host", cfg.ETL.Host),
		wsl.String("db", cfg.ETL.DBName),
		wsl.String("user", cfg.ETL.User),
	)
	dsnETL := getDSN(cfg.ETL.Host, cfg.ETL.User, cfg.ETL.Pass,
		cfg.ETL.DBName, cfg.ETL.Port, cfg.ETL.SSLMode)
	dbETL, err := createETLConnection(ctx, dsnETL, cfg.ETL.PoolMax, logger,
		// Автомиграция только базовых таблиц
		models.Device{},
		models.Session{},
		models.Source{},
		models.Domain{},
		models.Action{},
		models.Category{},
		models.DomainControlLists{},
		models.URL{},
		models.LastLog{},
		models.DomainCategory{},
	)
	if err != nil {
		logger.ErrorContext(ctx, "ETL service", wsl.String("create ETL db connection error", err.Error()))
		return
	}

	// Инициализация Kafka клиента
	logger.InfoContext(ctx, "ETL service", wsl.Info("Creating Kafka client"))
	kc, err := kafka.New(ctx, cfg.Kafka, &logger)
	if err != nil {
		logger.ErrorContext(ctx, "ETL service", wsl.String("create Kafka client error", err.Error()))
		return
	}
	defer kc.Close(ctx)

	// Инициализация QueryUsecase для работы с базами с поддержкой m-cashe
	q := q.New(cfg, dbKSU, dbETL, logger)

	// Инициализация клиента usercontrol
	userCl := usercontrol.New(
		httptransport.New(
			cfg.UserControl.Host+":"+cfg.UserControl.Port,
			"/",
			[]string{cfg.UserControl.Proto},
		),
		nil,
	)

	// Инициализация клиента для бизнес-логов
	blclient := blog.New(
		httptransport.New(
			cfg.BlogServ.Host+":"+cfg.BlogServ.Port,
			"/",
			[]string{cfg.BlogServ.Proto},
		),
		nil,
	)

	//Инициализация ProcessorUsecase для обработки данных
	p, err := usecase.New(cfg, q, kc, logger, userCl, blclient)
	if err != nil {
		logger.ErrorContext(ctx, "ETL service", wsl.String("failed to create usecase", err.Error()))
		return
	}

	// Инициализация обработчиков Kafka сообщений
	logger.InfoContext(ctx, "ETL service", wsl.Info("Initializing Kafka handlers"))
	metadataHandler := handlers.NewMetadataHandler(q, &logger)
	mlAnalysisHandler := handlers.NewMLAnalysisHandler(cfg.MLAttemps, cfg.BlacklistTriesCount, q, &logger)
	consumerHandlers := handlers.NewCompositeHandler(metadataHandler, mlAnalysisHandler, &logger)

	// Запуск обработчиков Kafka сообщений
	logger.InfoContext(ctx, "ETL service", wsl.Info("Starting Kafka consumers"))
	kc.StartConsumer(ctx, consumerHandlers)

	// Инициализация контроллеров
	logger.InfoContext(ctx, "ETL service", wsl.Info("Initializing controllers"))

	// Создаем полный UsecaseInterface для обратной совместимости
	uc := struct {
		usecase.QueryUsecase
		usecase.ProcessorUseCase
	}{
		QueryUsecase:     q,
		ProcessorUseCase: p,
	}

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
		logger.InfoContext(ctx, "Starting HTTP server")
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			logger.ErrorContext(ctx, "HTTP server error", wsl.Err(err))
		}
	}()

	// Запуск  ETL процессора
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.InfoContext(ctx, "starting ETL service")
		etlController.Start(ctx)
	}()

	// Ожидание сигналов завершения
	<-sigChan
	logger.InfoContext(ctx, "Received shutdown signal")

	// Graceful shutdown
	logger.InfoContext(ctx, "Stopping HTTP server")
	server.Stop(ctx)

	logger.InfoContext(ctx, "Stopping ETL processor")
	cancel()

	// Ждем завершения Kafka consumers
	logger.InfoContext(ctx, "Waiting for Kafka consumers to finish")
	kc.Wait()

	logger.InfoContext(ctx, "Waiting for goroutines to finish")
	wg.Wait()

	logger.InfoContext(ctx, "ETL service stopped completely")
}
