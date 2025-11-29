package app

import (
	"fiermon-blog/config"
	"fiermon-blog/internal/controllers/http/v1"
	"fiermon-blog/internal/repository/postgresGorm"
	"fiermon-blog/internal/service"
	"fiermon-blog/pkg/fslog"
	"fiermon-blog/pkg/fslog/wsl"
)

// StartApp - Инициализация хранилища, логики приложения, api handlers
func StartApp(cfg *config.Config) {
	log := fslog.NewFSlogger(cfg.Level)
	log.Info("Starting App")

	repo, err := postgresGorm.NewPostgres(cfg.PG, log)
	if err != nil {
		log.Error("Failed to connect to postgres", wsl.Err(err))
		return
	}

	serv := service.NewService(repo, repo, log)

	httpServer := v1.NewServer(cfg.HTTP, log, serv, cfg.App.Version)

	go func() {
		log.Info("Starting HTTP server on address", wsl.Label("Host", cfg.HTTP.Host+":"+cfg.HTTP.Port))
		err = httpServer.Start()
		if err != nil {
			log.Error("Failed to start http controllers", wsl.Err(err))
		}
	}()
}
