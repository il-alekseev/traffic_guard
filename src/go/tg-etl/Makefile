RUN_TYPE ?= local
ETL_URL="http://127.0.0.1:7000/swagger/doc.json"
USERCONTROL_URL="http://192.168.130.112:8002/swagger/doc.json"
BLOG_URL="http://127.0.0.1:8004/swagger/doc.json"

# Запуск ETL  микросервиса
.PHONY: run
run:
	go run ./cmd/etl/main.go -type=$(RUN_TYPE)

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

# Генерация клиента UserControl
.PHONY: gen-u 
gen-u:
	swagger generate client -c usercontrol -f $(USERCONTROL_URL) -A usercontrol -t pkg

# Генерация клиента для бизнес-логов
.PHONY: gen-blog
gen-blog:
	swagger generate client -c blog -f $(BLOG_URL) -A blog -t pkg