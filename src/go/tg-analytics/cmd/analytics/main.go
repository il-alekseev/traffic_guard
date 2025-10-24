package main

import (
	"log/slog"
	"tg-an/config"
	an "tg-an/internal/app/analytics"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		slog.Error(err.Error())
	}
	cfg.App.DevVersion = "0.1.1-dev.10"
	an.Run(cfg)
}
