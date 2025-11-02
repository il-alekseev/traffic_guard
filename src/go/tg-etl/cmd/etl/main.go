package main

import (
	"flag"
	"log/slog"
	"tg-etl/config"
	"tg-etl/internal/app/etl"
)

func main() {
	cfgPath := "./config/config.yaml" // Путь к конфигу в контейнере (по умолчанию)

	runType := flag.String("type", "local", "Type of running service")
	// Парсинг флагов
	flag.Parse()
	if *runType == "local" {
		// Используется локальная БД для аналитики
		cfgPath = "./deploy/local/config.yaml"
	} else if *runType == "remote" {
		// Используется удаленная БД для аналитики
		cfgPath = "./deploy/remote/config.yaml"
	} else {
		slog.Error("Invalid run type argument value",
			"value", *runType,
			"allowed", "local, remote")
		return
	}

	cfg, err := config.NewConfig(cfgPath)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	cfg.App.DevVersion = "0.1.1-dev.14"
	etl.Run(cfg)
}
