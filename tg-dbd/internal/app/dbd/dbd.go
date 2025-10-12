package dbd

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"tg-dbd/config"
	httpserver "tg-dbd/internal/controllers/http/v1"
	"tg-dbd/internal/usecase"
	"tg-dbd/pkg/logger"
	"time"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	u := usecase.New(l)

	server := httpserver.New(&cfg.HTTP, l, u)

	go func() {
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			l.Fatal("Failed to run server: %v", err)
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
		l.Error("Server forced to shutdown: %v", err)
	}
	server.Stop(ctx)
	l.Info("dashboard http server stopped")

}
