ETL_URL="http://127.0.0.1:7000/swagger/doc.json"
# Запуск ETL  микросервиса
.PHONY: run
run: 
	go run ./cmd/etl/main.go

# Генерация Swagger документации
.PHONY: swagger
swagger:
	swag init -g .\internal\controllers\http\v1\server.go -o ./docs
