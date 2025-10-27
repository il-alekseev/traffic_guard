package app

import (
	"fiermon-blog/config"
	"fiermon-blog/internal/repository/postgresGorm"
	"fiermon-blog/internal/server"
	"fiermon-blog/internal/service"
	"fiermon-blog/pkg/fslog"
	"fiermon-blog/pkg/fslog/wsl"
)

// StartApp - Инициализация хранилища, логики приложения, api handlers
func StartApp(cfg *config.Config) {
	log := fslog.NewFSlogger(cfg.Level)
	log.Info("Starting App")

	readRepo, err := postgresGorm.NewPostgres(cfg.PG, log)
	if err != nil {
		log.Error("Failed to connect to postgres", wsl.Err(err))
		return
	}

	serv := service.NewService(readRepo, log)

	httpServer := server.NewServer(cfg.HTTP, log, serv)

	go func() {
		err = httpServer.Start()
		if err != nil {
			log.Error("Failed to start http server", wsl.Err(err))
		}
	}()
}
