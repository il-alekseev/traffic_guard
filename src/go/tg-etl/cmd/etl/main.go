package main

import (
	"log/slog"
	"tg-etl/config"
	"tg-etl/internal/app/etl"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		slog.Error(err.Error())
	}
	cfg.App.DevVersion = "0.1.1-dev.8"
	etl.Run(cfg)
}
