DBD_URL="http://127.0.0.1:7000/swagger/doc.json"
# Запуск dashboard микросервиса
.PHONY: run
run: 
	go run ./cmd/dbd/main.go

# Генерация Swagger документации
.PHONY: swagger
swagger:
	swag init -g .\internal\controllers\http\v1\server.go -o ./docs

# Генерация клиента для сервиса дашбордов
.PHONY: gen
gen:
	swagger generate client -c dbd -f $(DBD_URL) -A dbd -t pkg
