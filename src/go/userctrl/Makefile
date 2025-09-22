.PHONY: swagger

# Установка swag если не установлен
install_swag:
	go install github.com/swaggo/swag/cmd/swag@latest

# Генерация Swagger документации
swagger:
	swag init -g .\internal\controllers\http\v1\server.go -o ./docs

# Запуск сервера с Swagger UI
run: 
	go run ./cmd/usercontrol/main.go

.PHONY: vet
vet:
	go vet ./...

.PHONY: lint
lint:
	golint ./...

.PHONY: docker-up
docker-up:
	docker compose --file build/docker-compose.local.yml --env-file build/local.env up -d