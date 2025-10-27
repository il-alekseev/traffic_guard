.PHONY: run build clean

# Имя исполняемого файла
BINARY_NAME=api-gateway.exe

USERCONTROL_URL="http://127.0.0.1:8002/swagger/doc.json"
CTXCONTROL_URL="http://127.0.0.1:8003/swagger/doc.json"
BLOG_SERV_URL="http://127.0.0.1:8004/swagger/doc.json"

# Цель по умолчанию
all: generate build

# Генерация кода из OpenAPI спецификации
.PHONY: gen-u 
gen-u:
	swagger generate client -c usercontrol -f $(USERCONTROL_URL) -A usercontrol -t pkg

.PHONY: gen-c 
gen-c:
	swagger generate client -c ctxcontrol -f $(CTXCONTROL_URL) -A ctxcontrol -t pkg

.PHONY: gen-b 
gen-b:
	swagger generate client -c blogserv -f $(BLOG_SERV_URL) -A blogserv -t pkg

run:
	go run cmd/server/main.go

build:
	go build -o bin/$(BINARY_NAME) cmd/server/main.go

clean:
	rm -rf bin/ internal/generated/

.PHONY: swagger
swagger:
	swag init -g .\internal\controllers\http\v1\server.go -o ./docs

.PHONY: vet
vet:
	go vet ./...

.PHONY: lint
lint:
	golint ./...


.PHONY: docker-up
docker-up:
	docker compose --file build/docker-compose.local.yml --env-file build/local.env up -d