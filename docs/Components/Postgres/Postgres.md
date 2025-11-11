# Описание компонентов Postgres

## Общее описание

PostgreSQL — мощная объектно-реляционная система баз данных с открытым исходным кодом. Инструмент написан на языке C.

## Взаимосвязь с другими компонетами

### Postgres хранящий логи

![](Components/Postgres/img/1.png)

В рамках системы выступает внешним хранилищем для логов KSU и впоследствии источником данных для создания дашбордов в компоненте Grafana.

### Postgres хранящий конфигурации


![](Components/Postgres/img/2.png)  
  

В рамках системы выступает хранилищем для конфигураций Keycloak и Grafana.

# Описание настройки

Инструменты будет развернут в системе в виде Docker- контейнера в системе управления контейнерами Docker Compose.

## Настройка контейнера в системе

### Postgres хранящий логи

Перед запуском контейнера необходимо указать переменные окружения в .env файле.

`.env `:

  

```
POSTGRES_LOG_VOLUME_PORT=5432 # Версия образа

POSTGRES_LOG_PORT=5432 # Порт хоста

POSTGRES_LOG_USER=log_user # Имя администратора

POSTGRES_LOG_PASSWORD=log_password # Пароль администратора

POSTGRES_KSU_USER=ksu_user # Имя пользователя базы логов

POSTGRES_KSU_PASSWORD=ksu_password # Пароль пользователя базы логов

POSTGRES_KSU_DATABASE=ksu # Имя базы логов
```


`docker-compose.yml`:

  
```yaml
logsdb:
	image: postgres:${POSTGRES_TAG}
	ports:
		- ${POSTGRES_LOG_VOLUME_PORT}:${POSTGRES_LOG_PORT} 
	environment:
		- POSTGRES_USER=system_admin
		- POSTGRES_PASSWORD=adminpassword
		- POSTGRES_KSU_USER=${POSTGRES_KSU_USER}
		- POSTGRES_KSU_PASSWORD=${POSTGRES_KSU_PASSWORD}
		- POSTGRES_KSU_DATABASE=${POSTGRES_KSU_DATABASE}
		- POSTGRES_KSU_MON_DATABASE=${POSTGRES_KSU_MON_DATABASE}
	volumes:
		- ./postgres_log/data:/var/lib/postgresql/data
	restart: always
	healthcheck:
		test: ["CMD-SHELL", "pg_isready -U $POSTGRES_LOG_USER"]
		interval: 5s
		timeout: 5s
		retries: 10
		command: [
			"sh", "-c",
			"echo \"CREATE USER $POSTGRES_KSU_USER WITH PASSWORD '$POSTGRES_KSU_PASSWORD';\n\" > /docker-entrypoint-initdb.d/init.sql &&
			echo \"CREATE DATABASE $POSTGRES_KSU_DATABASE OWNER $POSTGRES_KSU_USER ;\n\" >> /docker-entrypoint-initdb.d/init.sql &&
			echo \"CREATE DATABASE $POSTGRES_KSU_MON_DATABASE OWNER $POSTGRES_KSU_USER ;\n\" >> /docker-entrypoint-initdb.d/init.sql &&
			docker-entrypoint.sh postgres "
			]
	networks:
		- database
```
  

  

### Postgres хранящий конфигурации
    
Перед запуском контейнера необходимо указать переменные окружения в .env файле.

`./deploy/.env `:

```
POSTGRES_TAG=16

POSTGRES_CONFIG_VOLUME_PORT=5433

POSTGRES_CONFIG_PORT=5432

POSTGRES_USER=config_user

POSTGRES_PASSWORD=config_password

POSTGRES_KEYCLOAK_USER=keycloak_user

POSTGRES_KEYCLOAK_PASSWORD=keycloak_password

POSTGRES_KEYCLOAK_DATABASE=keycloak

POSTGRES_GRAFANA_USER=grafana_user

POSTGRES_GRAFANA_PASSWORD=grafana_password

POSTGRES_GRAFANA_DATABASE=grafana

POSTGRES_BLOG_USER=blog_user

POSTGRES_BLOG_PASSWORD=blog_password

POSTGRES_BLOG_DATABASE=blog

POSTGRES_ETL_USER=etl_user

POSTGRES_ETL_PASSWORD=etl_password

POSTGRES_ETL_DATABASE=etl

POSTGRES_WEB_USER=web_user

POSTGRES_WEB_PASSWORD=web_password

POSTGRES_WEB_DATABASE=web

POSTGRES_ML_USER=ml_user

POSTGRES_ML_PASSWORD=ml_password

POSTGRES_ML_DATABASE=ml
```


`docker-compose.yml`:

  ```yaml
    configdb:
    image: postgres:${POSTGRES_TAG}
    ports:
      - ${POSTGRES_CONFIG_VOLUME_PORT}:${POSTGRES_CONFIG_PORT}
    env_file: ./postgres_config/.env
    volumes:
      - ./postgres_config/data:/var/lib/postgresql/data
    restart: always
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U $POSTGRES_USER"]
      interval: 5s
      timeout: 5s
      retries: 10
    command: [
      "sh", "-c",
      "echo \"CREATE USER $POSTGRES_KEYCLOAK_USER WITH PASSWORD '$POSTGRES_KEYCLOAK_PASSWORD';\n\" > /docker-entrypoint-initdb.d/init.sql &&
      echo \"CREATE USER $POSTGRES_GRAFANA_USER WITH PASSWORD '$POSTGRES_GRAFANA_PASSWORD';\n\" >> /docker-entrypoint-initdb.d/init.sql &&
      echo \"CREATE USER $POSTGRES_BLOG_USER WITH PASSWORD '$POSTGRES_BLOG_PASSWORD';\n\" >> /docker-entrypoint-initdb.d/init.sql &&
      echo \"CREATE USER $POSTGRES_ETL_USER WITH PASSWORD '$POSTGRES_ETL_PASSWORD';\n\" >> /docker-entrypoint-initdb.d/init.sql &&
      echo \"CREATE USER $POSTGRES_WEB_USER WITH PASSWORD '$POSTGRES_WEB_PASSWORD';\n\" >> /docker-entrypoint-initdb.d/init.sql &&
      echo \"CREATE USER $POSTGRES_ML_USER WITH PASSWORD '$POSTGRES_ML_PASSWORD';\n\" >> /docker-entrypoint-initdb.d/init.sql &&
      echo \"CREATE DATABASE $POSTGRES_KEYCLOAK_DATABASE OWNER $POSTGRES_KEYCLOAK_USER ;\n\" >> /docker-entrypoint-initdb.d/init.sql &&
      echo \"CREATE DATABASE $POSTGRES_GRAFANA_DATABASE OWNER $POSTGRES_GRAFANA_USER ;\n\" >> /docker-entrypoint-initdb.d/init.sql &&
      echo \"CREATE DATABASE $POSTGRES_BLOG_DATABASE OWNER $POSTGRES_BLOG_USER ;\n\" >> /docker-entrypoint-initdb.d/init.sql &&
      echo \"CREATE DATABASE $POSTGRES_ETL_DATABASE OWNER $POSTGRES_ETL_USER ;\n\" >> /docker-entrypoint-initdb.d/init.sql &&
      echo \"CREATE DATABASE $POSTGRES_WEB_DATABASE OWNER $POSTGRES_WEB_USER ;\n\" >> /docker-entrypoint-initdb.d/init.sql &&
      echo \"CREATE DATABASE $POSTGRES_ML_DATABASE OWNER $POSTGRES_ML_USER ;\n\" >> /docker-entrypoint-initdb.d/init.sql &&
      docker-entrypoint.sh postgres "
    ]
    networks:
      - continent
  ```

