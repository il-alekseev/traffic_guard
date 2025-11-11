# Описание компонента “Стек компонентов Kafka”

## Общее описание


Стек компонентов Kafka состоит из трех компонентов:

- **Kafka**- это распределённая платформа для обмена сообщениями и потоковой обработки данных в реальном времени, построенная на принципе публикации-подписки.

- **Zookeeper** - это централизованный сервис для координации и управления распределёнными системами, такой как Kafka, который хранит их метаданные и обеспечивает синхронизацию между узлами.

- **Kafka-UI** - это веб-интерфейс для визуального мониторинга и управления кластером Apache Kafka, который позволяет удобно просматривать топики, сообщения, потребителей.

## Взаимосвязь с другими компонетами
  
  
![](Components/Kafka/img/1.png)
  
  

В рамках системы стек отвечает за общение между компонентами системы:

- Webscraper
- Tg-ETL
- Neuro


# Описание настройки

Инструмент будет развернут в системе в виде Docker- контейнера в системе управления контейнерами Docker Compose.


## Настройка контейнера в системе

### Kafka

Перед запуском контейнера необходимо указать переменные окружения в двух .env файлах. Первый отвечает за саму работу компонента, а второй за парметры запуска контейнера.

`./deploy/kafka/.env` :  

```
KAFKA_BROKER_ID=1

KAFKA_ZOOKEEPER_CONNECT=zookeeper:2181

KAFKA_LISTENERS=PLAINTEXT://0.0.0.0:9092

KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://kafka:9092

KAFKA_LISTENER_SECURITY_PROTOCOL_MAP=PLAINTEXT:PLAINTEXT

KAFKA_INTER_BROKER_LISTENER_NAME=PLAINTEXT

KAFKA_AUTO_CREATE_TOPICS_ENABLE="true"

KAFKA_DELETE_TOPIC_ENABLE="true"

KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR=1

KAFKA_TRANSACTION_STATE_LOG_REPLICATION_FACTOR=1

KAFKA_TRANSACTION_STATE_LOG_MIN_ISR=1

KAFKA_LOG_RETENTION_HOURS=168

KAFKA_LOG_RETENTION_BYTES=1073741824

KAFKA_LOG_SEGMENT_BYTES=1073741824

KAFKA_GROUP_INITIAL_REBALANCE_DELAY_MS=0
```
  
`./deploy/.env`:

```
...
KAFKA_APP_NAME=kafka

KAFKA_HTTP_PORT=9092
...
```
  

`docker-compose.yml`:

```yaml
  kafka:
    image: confluentinc/cp-kafka:7.4.0
    container_name: ${KAFKA_APP_NAME}
    ports:
      - ${KAFKA_HTTP_PORT}:${KAFKA_HTTP_PORT}
    volumes:
      - ./kafka/data:/var/lib/kafka/data
      - ./kafka/init_topics.sh:/tmp/c
    env_file: ./kafka/.env 
    healthcheck:
      test: ["CMD", "bash", "-c", "kafka-topics --bootstrap-server localhost:${KAFKA_HTTP_PORT} --list && /tmp/init-topics.sh"]
      interval: 20s
      timeout: 15s
      retries: 5
      start_period: 40s
    restart: always
    depends_on:
      zookeeper:
        condition: service_healthy   
    networks:
      - continent
```

А также при запуске контейнера запускается скрипт создания топиков.

`./deploy/kafka/init-topics.sh`

```bash
#!/bin/bash

  

# Простой скрипт для создания топиков через healthcheck

BOOTSTRAP_SERVER="localhost:9092"

  

# Функция создания топика если не существует

create_topic_if_not_exists() {

local topic=$1

local partitions=$2

local retention=$3

local cleanup_policy=$4

if ! kafka-topics --bootstrap-server $BOOTSTRAP_SERVER --list | grep -q "$topic"; then

echo "Creating topic: $topic"

kafka-topics --bootstrap-server $BOOTSTRAP_SERVER --create \

--topic "$topic" \

--partitions "$partitions" \

--replication-factor 1 \

--config "retention.ms=$retention" \

--config "cleanup.policy=$cleanup_policy"

fi

}

  

# Создаем топики если их нет

create_topic_if_not_exists "ml-analysis-results" 3 604800000 "delete"

create_topic_if_not_exists "url-content-analysis" 3 604800000 "delete"

create_topic_if_not_exists "url-metadata-results" 3 604800000 "delete"

create_topic_if_not_exists "url-processing-requests" 3 604800000 "delete"

  

exit 0
```
  

  

### Zookeeper

Перед запуском контейнера необходимо указать переменные окружения в двух .env файлах. Первый отвечает за саму работу компонента, а второй за парметры запуска контейнера.  

`./deploy/zookeeper/.env`:

```
ZOOKEEPER_CLIENT_PORT=2181

ZOOKEEPER_TICK_TIME=2000

ZOOKEEPER_SYNC_LIMIT=2
```

`./deploy/.env`:

```
...
ZOO_APP_NAME=zookeeper

ZOO_PORT=2181
...
```

`docker-compose.yml`:

```yaml
  zookeeper:
    image: confluentinc/cp-zookeeper:7.4.0
    container_name: ${ZOO_APP_NAME}
    ports:
      - ${ZOO_PORT}:${ZOO_PORT}
    volumes:
      - ./zookeeper/data:/var/lib/zookeeper/data
      - ./zookeeper/log:/var/lib/zookeeper/log
    env_file: ./zookeeper/.env
    healthcheck:
      test: ["CMD", "bash", "-c", "echo stat | nc localhost ${ZOO_PORT}"]
      interval: 10s
      timeout: 5s
      retries: 5
    restart: always
    networks:
      - continent
```

### Kafka-UI

Перед запуском контейнера необходимо указать переменные окружения в двух .env файлах. Первый отвечает за саму работу компонента, а второй за парметры запуска контейнера.  

`./deploy/kafka_ui/.env`:

```
KAFKA_CLUSTERS_0_NAME: local-kafka

KAFKA_CLUSTERS_0_BOOTSTRAPSERVERS: kafka:9092

KAFKA_CLUSTERS_0_PROPERTIES_SECURITY_PROTOCOL: PLAINTEXT

KAFKA_CLUSTERS_0_PROPERTIES_REQUEST_TIMEOUT_MS: 30000

KAFKA_CLUSTERS_0_PROPERTIES_RETRY_BACKOFF_MS: 1000

KAFKA_CLUSTERS_0_READONLY: "false"

KAFKA_CLUSTERS_0_SHOW_INTERNAL_TOPICS: "true"

KAFKA_CLUSTERS_0_DISABLE_LOGS_DUMPING: "false"
```

`./deploy/.env`:

```
KAFKAUI_APP_NAME=kafka-ui

KAFKAUI_PORT=8080

KAFKAUI_VOLUME_PORT=8081
```

`docker-compose.yml`:

```yaml
  kafka-ui:
    image: provectuslabs/kafka-ui:latest
    container_name: ${KAFKAUI_APP_NAME}
    ports:
      - ${KAFKAUI_VOLUME_PORT}:${KAFKAUI_PORT}
    env_file: ./kafka_ui/.env
    volumes:
      - ./kafka_ui/data:/etc/kafka-ui
    healthcheck:
      test: ["CMD-SHELL", "wget --no-verbose --tries=1 --spider http://localhost:${KAFKAUI_PORT}/actuator/health || exit 1"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s
    restart: always
    depends_on:
      kafka:
        condition: service_healthy
    networks:
      - continent
```

