# Makefile для управления docker-compose

# Имя сервиса из docker-compose.yml
SERVICE=traffic-classifier

# Сборка образа
build:
	docker compose -f deploy/docker-compose.yml build

# Запуск контейнера в фоне
up:
	docker compose -f deploy/docker-compose.yml up -d

# Запуск в интерактивном режиме (логи в консоль)
logs:
	docker compose -f deploy/docker-compose.yml logs -f $(SERVICE)

# Остановка контейнера
down:
	docker compose -f deploy/docker-compose.yml down

# Выполнить команду внутри контейнера
exec:
	docker compose -f deploy/docker-compose.yml exec $(SERVICE) /bin/bash
