APP_NAME = fiermon-blog
DOCKER_IMAGE = $(APP_NAME)
VERSION = 0.0.1
CONFIG_PATH = config/config.yml

.PHONY: all build run stop

all: build

build:
	docker build -t $(DOCKER_IMAGE):$(VERSION) .

.PHONY: run
run:
	go run ./cmd/main.go

.PHONY: docker-up
docker-up:
	docker compose --file build/docker-compose.yaml --env-file build/.env up -d

stop:
	docker stop $(APP_NAME)

.PHONY: swagger
swagger:
	swag init -g ./internal/controllers/http/v1/server.go -o ./docs