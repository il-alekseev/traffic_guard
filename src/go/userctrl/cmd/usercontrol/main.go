package main

import (
	"fmt"
	"userctrl/config"
	"userctrl/internal/app/usercontrol"
)

func main() {
	// Получение переменных из конфигурационного файла или переменного окружения
	cfg, err := config.NewConfig()
	if err != nil {
		fmt.Printf("read config: %s", err.Error())
		return
	}

	// Запуск основной логики
	usercontrol.Run(cfg)
}
