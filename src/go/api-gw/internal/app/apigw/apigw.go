package apigw

import (
	"api-gateway/config"
	"api-gateway/docs"
	v1 "api-gateway/internal/controllers/http/v1"
	"api-gateway/pkg/analytics"
	"api-gateway/pkg/blogserv"
	"api-gateway/pkg/ctxcontrol"
	"api-gateway/pkg/slogger"
	"api-gateway/pkg/usercontrol"
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	httptransport "github.com/go-openapi/runtime/client"
)

// Run - инициализация и запуск основных слоев сервиса
func Run(cfg *config.Config) {
	ctx := context.Background()

	gin.SetMode(gin.ReleaseMode)

	// Установка ip swagger документации
	docs.SwaggerInfo.Host = cfg.Swagger.Host + ":" + cfg.HTTP.Port

	// Настройка логирования
	logLevel := slog.LevelDebug
	if cfg.Log.Level != "debug" {
		logLevel = slog.LevelInfo
	}
	slogger.InitLogging(logLevel)

	// инициализация слоя usecase
	// инициализация клиента сервиса usercontrol
	userCtrlCl := usercontrol.New(
		httptransport.New(
			cfg.UserControl.Host+":"+cfg.UserControl.Port,
			"/",
			[]string{cfg.UserControl.Proto},
		),
		nil,
	)

	// инициализация клиента сервиса ctxcontrol
	ctxCtrlCl := ctxcontrol.New(
		httptransport.New(
			cfg.CtxControl.Host+":"+cfg.CtxControl.Port,
			"/",
			[]string{cfg.UserControl.Proto},
		),
		nil,
	)

	// инициализация клиента сервиса blog
	blogCl := blogserv.New(
		httptransport.New(
			cfg.BlogServ.Host+":"+cfg.BlogServ.Port,
			"/",
			[]string{cfg.BlogServ.Proto},
		),
		nil,
	)

	// инициализация клиента сервиса analytics
	analyticsCl := analytics.New(
		httptransport.New(
			cfg.Analytics.Host+":"+cfg.Analytics.Port,
			"/",
			[]string{cfg.Analytics.Proto},
		),
		nil,
	)

	// инициализация http сервера и обработчиков
	server, err := v1.New(ctx, cfg, userCtrlCl, ctxCtrlCl, blogCl, analyticsCl)
	if err != nil {
		slog.ErrorContext(slogger.ErrorCtx(ctx, err), "create http server: "+err.Error())
	}
	slog.InfoContext(ctx, "create http server")
	stopChan := make(chan struct{})

	// http server run
	go func() {
		if err := server.Run(); err != nil && err != http.ErrServerClosed {
			slog.Error("Failed to run server", "err", err)
			close(stopChan)
			return
		}
	}()

	// Создаем канал для получения сигналов
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM) // Подписываемся на сигнал прерывания (Ctrl+C)

	// Ожидаем сигнала
	select {
	case <-signalChan:
		slog.Info("Received shutdown signal, stopping server...")
	case <-stopChan:
		slog.Info("Server stopped due to internal error")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Stop(ctx); err != nil {
		slog.Error("Server forced to shutdown", "err", err)
	}

	slog.Info("Server stopped gracefully")
}
