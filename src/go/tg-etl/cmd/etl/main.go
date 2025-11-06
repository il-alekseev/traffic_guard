package main

import (
	"flag"
	"log/slog"
	"tg-etl/config"
	"tg-etl/internal/app/etl"
)

func main() {
	cfgPath := "./config/config.yaml" // Путь к конфигу в контейнере (по умолчанию)

	runType := flag.String("type", "", "Type of running service")
	// Парсинг флагов
	flag.Parse()
	if *runType == "local" {
		// Используется локальная БД для аналитики
		cfgPath = "./deploy/local/config.yaml"
	} else if *runType == "remote" {
		// Используется удаленная БД для аналитики
		cfgPath = "./deploy/remote/config.yaml"
	}

	cfg, err := config.NewConfig(cfgPath)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	cfg.App.DevVersion = "0.1.1-dev.17"
	etl.Run(cfg)
}
