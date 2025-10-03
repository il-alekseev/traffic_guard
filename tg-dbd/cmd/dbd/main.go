package main

import (
	"log/slog"
	"tg-dbd/config"
	"tg-dbd/internal/app/dbd"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		slog.Error(err.Error())
	}
	dbd.Run(cfg)
}
