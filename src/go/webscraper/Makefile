.PHONY: build run tidy lint vet swagger docker-build docker-up docker-down

build:
	mkdir -p bin
	CGO_ENABLED=0 go build -o bin/webscraper ./cmd/scraper

run:
	go run ./cmd/scraper --config config/config.yaml

vet:
	go vet ./...

lint:
	golangci-lint run || true

swagger:
	@echo "Swagger specification located at docs/swagger.yaml"

tidy:
	go mod tidy

docker-build:
	docker build -f deploy/Dockerfile -t webscraper:local .

docker-up:
	docker compose -f deploy/docker-compose.local.yml up -d --build

docker-down:
	docker compose -f deploy/docker-compose.local.yml down
