RUN_TYPE ?= local
AN_URL="http://127.0.0.1:7000/swagger/doc.json"
BLOG_URL="http://127.0.0.1:8004/swagger/doc.json"

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

# Генерация клиента для сервиса аналитики
.PHONY: gen-an
gen-an:
	swagger generate client -c analytics -f $(AN_URL) -A analytics -t pkg

# Генерация клиента для бизнес-логов
.PHONY: gen-blog
gen-blog:
#	swagger generate client -c blog -f $(BLOG_URL) -A blog -t pkg
	swagger generate client -c blog -f ../blog/docs/swagger.json -A blog -t pkg
