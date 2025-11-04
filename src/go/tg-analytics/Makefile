RUN_TYPE ?= local
AN_URL="http://127.0.0.1:7000/swagger/doc.json"

# Запуск tg-analytics микросервиса
.PHONY: run
run: 
	go run ./cmd/analytics/main.go -type=$(RUN_TYPE)

.PHONY: run-local
run-local:
	$(MAKE) run RUN_TYPE=local

.PHONY: run-remote
run-remote:
	$(MAKE) run RUN_TYPE=remote

# Генерация Swagger документации
.PHONY: swagger
swagger:
	swag init -g .\internal\controllers\http\v1\server.go -o ./docs

# Генерация клиента для сервиса дашбордов
.PHONY: gen
gen:
	swagger generate client -c analytics -f $(AN_URL) -A analytics -t pkg