# Makefile for deploying and managing the complex

include ./deploy/.env
export

COMPOSE_FILE = ./deploy/docker-compose.yml
CERT_GEN_SCRIPT = ./deploy/cert_gen.sh
COMPOSE_PROJECT_NAME ?= deploy
GRAFANA_CERT_DIR = ./deploy/grafana
KEYCLOAK_CERT_DIR = ./deploy/keycloak


.PHONY: all
all: cert-gen deploy


.PHONY: cert-gen
cert-gen:
	@echo "Generating certificates..."
	$(CERT_GEN_SCRIPT)

.PHONY: deploy
deploy:
	@echo "Create Data Directories..."
	if [ ! -d "./deploy/postgres_log/data" ]; then \
		mkdir ./deploy/postgres_log/data; \
		chmod 777 ./deploy/postgres_log/data; \
	fi
	if [ ! -d "./deploy/postgres_config/data" ]; then \
		mkdir ./deploy/postgres_config/data; \
		chmod 777 ./deploy/postgres_config/data; \
	fi
	if [ ! -d "./deploy/keycloak/data" ]; then \
		mkdir ./deploy/keycloak/data; \
		chmod 777 ./deploy/keycloak/data; \
	fi
	if [ ! -d "./deploy/grafana/data" ]; then \
		mkdir ./deploy/grafana/data; \
		chmod 777 ./deploy/grafana/data; \
	fi
	if [ ! -d "./deploy/prometheus/data" ]; then \
		mkdir ./deploy/prometheus/data; \
		chmod 777 ./deploy/prometheus/data; \
	fi
	@echo "Starting Docker containers..."
	docker compose up -d 

.PHONY: stop
stop:
	@echo "Stopping Docker containers..."
	docker compose -f $(COMPOSE_FILE) --project-name $(COMPOSE_PROJECT_NAME) stop

.PHONY: clean
clean:
	@echo "Stopping and removing Docker containers and volumes..."
	docker compose -f $(COMPOSE_FILE) --project-name $(COMPOSE_PROJECT_NAME) down -v --rmi local

	@echo "Remove Data Directories"
	if [ -d "./deploy/postgres_log/data" ]; then \
		sudo rm -rf ./deploy/postgres_log/data; \
	fi
	if [ -d "./deploy/postgres_config/data" ]; then \
		sudo rm -rf ./deploy/postgres_config/data; \
	fi
	if [ -d "./deploy/keycloak/data" ]; then \
		sudo rm -rf ./deploy/keycloak/data; \
	fi
	if [ -d "./deploy/grafana/data" ]; then \
		sudo rm -rf ./deploy/grafana/data; \
	fi
	if [ -d "./deploy/prometheus/data" ]; then \
		sudo rm -rf ./deploy/prometheus/data; \
	fi


.PHONY: status
status:
	@echo "Checking status of Docker containers..."
	docker compose ps

.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make all       - Set up environment and deploy containers"
	@echo "  make cert-gen  - Generate certificates"
	@echo "  make deploy    - Start Docker containers"
	@echo "  make stop      - Stop Docker containers"
	@echo "  make clean     - Stop and remove containers, volumes, and PAM users"
	@echo "  make status    - Check status of Docker containers"
	@echo "  make help      - Display this help message"
