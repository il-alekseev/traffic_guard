package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"fiermon-blog/config"
	"fiermon-blog/docs"
	"fiermon-blog/internal/app"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	cfg.App.Version = "1.0.2-dev"

	docs.SwaggerInfo.Host = cfg.Swagger.Host + ":" + cfg.HTTP.Port

	app.StartApp(cfg)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done
}
