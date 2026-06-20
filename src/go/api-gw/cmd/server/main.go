package main

import (
	"api-gateway/config"
	"api-gateway/internal/app/apigw"
	"fmt"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		fmt.Printf("read config: %s", err.Error())
		return
	}
	cfg.App.Version = "1.0.2-dev"

	apigw.Run(cfg)
}
