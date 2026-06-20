package usercontrol

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"
	"userctrl/config"
	"userctrl/docs"
	httpserver_v1 "userctrl/internal/controllers/http/v1"
	"userctrl/internal/usecase"
	blog "userctrl/pkg/blog/blog"
	"userctrl/pkg/grafanaclient"
	"userctrl/pkg/grafcookier"
	"userctrl/pkg/keycloakclient"
	"userctrl/pkg/pgorm"
	"userctrl/pkg/slogger"
	"userctrl/pkg/slogger/wsl"

	"github.com/gin-gonic/gin"
	goapi "github.com/grafana/grafana-openapi-client-go/client"

	httptransport "github.com/go-openapi/runtime/client"
)

func makeDSN(host, user, pass, dbname, port, sslmode string) string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, pass, dbname, port, sslmode,
	)
}

// getDB - обертка для подключения к БД и получения экземпляра *pgorm.Postgres
func getDB(dsn string, l slog.Logger, models ...interface{}) (*pgorm.Postgres, error) {
	db, err := pgorm.New(dsn,
		pgorm.MaxPoolSize(10),
		pgorm.AutoMigrate(true),
		pgorm.Models(models...),
	)
	if err != nil {
		l.Error("new pg repo", wsl.Err(err))
		return nil, fmt.Errorf("new pg repo: %w", err)
	}
	l.Debug("Connect OK")
	if err := db.HealthCheck(); err != nil {
		err := fmt.Errorf("database health check failed: %s", err.Error())
		l.Error(err.Error())
		return nil, err
	} else {
		l.Info("Database connection is healthy")
	}
	l.Debug("Repo new")
	return db, nil
}

// Run - инициализация всех зависимостей, запуск компонентов приложения
func Run(cfg *config.Config) error {
	ctx := context.Background()

	gin.SetMode(gin.ReleaseMode)

	// Инициализация
	logLevel := slog.LevelDebug
	if cfg.Log.Level != "debug" {
		logLevel = slog.LevelInfo
	}
	slogger.InitLogging(logLevel)

	// set swagger host
	docs.SwaggerInfo.Host = cfg.Swagger.Host + ":" + cfg.HTTP.Port

	// keycloak client init
	keycloakURL := fmt.Sprintf(
		"%s://%s:%s",
		cfg.KeyCloak.Proto,
		cfg.KeyCloak.Host,
		cfg.KeyCloak.Port,
	)
	logger := *slog.Default()

	// инициализация клиента для работы с Keycloak
	kc, err := keycloakclient.NewKeycloakClient(
		keycloakURL,
		cfg.KeyCloak.Login,
		cfg.KeyCloak.Pass,
		cfg.KeyCloak.SystemRealm,
		cfg.KeyCloak.Realm,
		cfg.KeyCloak.GrafanaClientID,
		cfg.KeyCloak.GrafanaClientUUID,
		cfg.KeyCloak.GrafanaClientSecret,
		logger,
	)
	if err != nil {
		slog.ErrorContext(ctx, "keycloakclient.NewKeycloakClient", wsl.Err(err))
		return err
	}

	// инициализация клиента grafana
	grafanaURL := fmt.Sprintf(
		"%s://%s:%s",
		cfg.Grafana.Proto,
		cfg.Grafana.Host,
		cfg.Grafana.Port,
	)
	grafCookier, err := grafcookier.New(
		grafanaURL,
		cfg.Grafana.AuthPath,
		logger,
	)
	if err != nil {
		err = fmt.Errorf("grafcookier.New, url=%s: error=%w", grafanaURL, err)
		logger.ErrorContext(ctx, err.Error())
		return err
	}

	tr := goapi.TransportConfig{
		Host:     cfg.Grafana.Host + ":" + cfg.Grafana.Port,
		BasePath: "/api",
		Schemes:  []string{cfg.Grafana.Proto},
		BasicAuth: url.UserPassword(
			cfg.Grafana.Username,
			cfg.Grafana.Pass,
		),
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		NumRetries:       3,
		RetryTimeout:     1,
		RetryStatusCodes: []string{"420", "5xx"},
		HTTPHeaders:      map[string]string{},
	}

	grafanaCl := grafanaclient.NewGrafanaClient(tr, "generic_oauth", logger)

	blogCl := blog.New(
		httptransport.New(
			cfg.BlogCl.Host+":"+cfg.BlogCl.Port,
			"/",
			[]string{cfg.BlogCl.Proto},
		),
		nil,
	)
	// инициализация слоя бизнес логики
	uc := usecase.New(
		cfg,
		kc,
		grafCookier,
		grafanaCl,
		blogCl,
		logger,
	)

	// инициализация слоя http
	server := httpserver_v1.New(cfg, uc, logger)

	// запуск http server
	go func() {
		if err := server.Run(); err != nil && err != http.ErrServerClosed {
			logger.Error("Failed to run server", wsl.Err(err))
		}
	}()

	// реализация grace full shutdown
	// Создаем канал для получения сигналов
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM) // Подписываемся на сигнал прерывания (Ctrl+C)

	// Ожидаем сигнала
	<-signalChan
	log.Println("Received shutdown signal, stopping server...")

	// Останавливаем сервер
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Stop(ctx); err != nil {
		logger.Error("Server forced to shutdown", wsl.Err(err))
	}

	logger.Info("Server stopped gracefully")

	return nil
}
