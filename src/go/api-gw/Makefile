.PHONY: run build clean

# Имя исполняемого файла
BINARY_NAME=api-gateway.exe

ANALYTICS_URL="http://127.0.0.1:7000/swagger/doc.json"
USERCONTROL_URL="http://127.0.0.1:8002/swagger/doc.json"
BLOG_SERV_URL="http://127.0.0.1:8004/swagger/doc.json"

# Цель по умолчанию
all: generate build

# Генерация кода из OpenAPI спецификации
.PHONY: gen-a
gen-a:
	swagger generate client -c analytics -f $(ANALYTICS_URL) -A analytics -t pkg
	# swagger generate client -c analytics -f ../tg-analytics/docs/swagger.json -A analytics -t pkg

.PHONY: gen-u 
gen-u:
	swagger generate client -c usercontrol -f $(USERCONTROL_URL) -A usercontrol -t pkg
	#swagger generate client -c usercontrol -f ../userctrl/docs/swagger.json -A usercontrol -t pkg

.PHONY: gen-b 
gen-b:
	swagger generate client -c blogserv -f $(BLOG_SERV_URL) -A blogserv -t pkg
	#swagger generate client -c blogserv -f ../blog/docs/swagger.json -A blogserv -t pkg

run:
	go run cmd/server/main.go

build:
	go build -o bin/$(BINARY_NAME) cmd/server/main.go

clean:
	rm -rf bin/ internal/generated/

.PHONY: swagger
swagger:
	swag init -g ./internal/controllers/http/v1/server.go -o ./docs
	# убирает поле omitempty которое добавляется при кодогенерации клиентов
	#find pkg -type f -name '*.go' -exec sed -i '' 's/,omitempty//' {} + # для macos
    #find pkg -type f -name '*.go' -exec sed -i 's/,omitempty//' {} + # для Linux

.PHONY: vet
vet:
	go vet ./...

.PHONY: lint
lint:
	golint ./...


.PHONY: docker-up
docker-up:
	docker compose --file build/docker-compose.local.yml --env-file build/local.env up -d