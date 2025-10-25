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
	etl.Run(cfg)
}
