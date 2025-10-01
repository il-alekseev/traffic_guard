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
	// TODO: парсинг параметров
	// TODO: чтение конфига
	cfg, err := config.NewConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	// TODO: настройка swagger
	docs.SwaggerInfo.Host = cfg.Swagger.Host + ":" + cfg.HTTP.Port
	// TODO: запуск основной логики
	app.StartApp(cfg)

	// TODO: grace full shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done
}
