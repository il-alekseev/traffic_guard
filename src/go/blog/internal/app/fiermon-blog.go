package app

import (
	"fiermon-blog/config"
	"fiermon-blog/internal/controllers/http/v1"
	"fiermon-blog/internal/models"
	"fiermon-blog/internal/repository/postgresRepo"
	"fiermon-blog/internal/service"
	"fiermon-blog/pkg/fslog"
	"fiermon-blog/pkg/fslog/wsl"
	"fiermon-blog/pkg/pgorm"
	"fmt"
	"gorm.io/gorm"
	"log/slog"
)

// StartApp - Инициализация хранилища, логики приложения, api handlers
func StartApp(cfg *config.Config) {
	log := fslog.NewFSlogger(cfg.Level)
	log.Info("Starting App")

	// инициализируем объект базы данных *gorm.DB
	db, err := connectDb(cfg.PG, log, models.BusinessLog{})
	if err != nil {
		log.Error("Failed to connect to database")
		return
	}

	// передаем объект *gorm.DB для реализации методов
	repo := postgresRepo.NewPostgresDB(db)

	//repo, err := postgresGorm.NewPostgres(cfg.PG, log)
	//if err != nil {
	//	log.Error("Failed to connect to postgres", wsl.Err(err))
	//	return
	//}

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

func connectDb(cfg config.PG, log *slog.Logger, model ...interface{}) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.Port,
		cfg.SSLMode,
	)

	options := []pgorm.Option{
		pgorm.AutoMigrate(true),
		pgorm.Models(model...),
	}

	postgresDB, err := pgorm.NewPostgres(dsn, log, options...)
	if err != nil {
		return nil, err
	}

	return postgresDB.GetDB(), nil
}
