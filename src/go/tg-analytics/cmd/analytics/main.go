package main

import (
	"flag"
	"log/slog"
	"tg-an/config"
	an "tg-an/internal/app/analytics"
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
	}
	cfg.App.DevVersion = "0.1.1-dev.21"
	an.Run(cfg)
}
