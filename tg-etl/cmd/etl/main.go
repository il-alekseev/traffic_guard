package main

import (
	"cmd/etl/config"
	"cmd/etl/internal/app/etl"
	"log/slog"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		slog.Error(err.Error())
	}
	etl.Run(cfg)
}
