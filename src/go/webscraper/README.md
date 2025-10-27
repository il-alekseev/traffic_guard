# Scrapper

## Описание
Scrapper — сервис на Go для массового получения контента с веб реусрсов и обогащения доменной информации. Сервис читает задачи из Kafka-топика `url-processing-requests` (Смотри формат общения в proto-via-kafka.md), параллельно обрабатывает их и публикует результаты в топики `url-metadata-results` (метаданные) и `url-content-analysis` (контент).

## Используемый стек
- Go 1.25
- github.com/markusmobius/go-trafilatura — извлечение основного контента
- github.com/oschwald/maxminddb-golang — GeoIP (локальный `data/ipinfo_lite.mmdb`)
- golang.org/x/net/html, html/charset — парсинг и декодирование HTML
- Локальный JSONL (`data/attempts.jsonl`) и файл логов (`logs/scraper.log`)

## Структура
- `cmd/scraper` — CLI-вход, загрузка конфига, запуск runtime
- `config` — YAML-конфиг и загрузчик
- `internal/models` — доменные структуры (`ScrapeResult`, `FailureReason` и др.)
- `internal/usecase` — воркер-пул, применение стратегий, диагностика
- `internal/app` — оболочка над usecase
- `internal/metrics` — счётчики in-flight/processed/failed
- `internal/processor` — вспомогательные функции нормализации и формирования сообщений
- `internal/server/http` — REST-сервер (healthcheck + тестовый `get_content`)
- `internal/infrastructure` — адаптеры (`content`, `httpclient`, `dns`, `geo`, `parser`, `quality`, `storage`, `logging`, `kafka`)
- `deploy` — Dockerfile и docker-compose для локального запуска
- `docs` — `swagger.yaml` с описанием REST API

## Конфигурация (`config/config.yaml`)
- `db_path` — путь к MaxMind базе
- `workers` — количество воркеров (≤0 → `runtime.NumCPU()*2`)
- `request_timeout` — общий таймаут (используется и для воркера, и для HTTP-клиента)
- `content.strategies` — порядок стратегий (`trafilatura`, `http`, ...)
- `content.min_text_length` — минимальная длина очищенного текста
- `content.max_content_length` — максимальная длина текста (обрезается при превышении)
- `content.user_agents_path` — путь к списку User-Agent строк (по одной на строку)
- `logging.level` — уровень логирования (`debug`, `info`, `warn`, `error`)
- `kafka.brokers` — список брокеров Kafka
- `kafka.group_id`, `kafka.client_id` — параметры consumer group
- `kafka.input_topic` — топик входящих задач (`url-processing-requests`)
- `kafka.metadata_topic` — топик метаданных (`url-metadata-results`)
- `kafka.content_topic` — топик контента (`url-content-analysis`)
- `kafka.commit_interval`, `kafka.poll_timeout` — параметры чтения/коммитов
- `kafka.auth` — блок SASL/TLS авторизации (при необходимости)
- `http.address`, `http.read_timeout`, `http.write_timeout`, `http.shutdown_timeout` — параметры REST-сервера

## Быстрый старт
```bash
# зависимости
go mod tidy

# сборка бинаря
go build -o bin/webscraper ./cmd/scraper

# запуск со стандартным конфигом
./bin/webscraper --config config/config.yaml

# REST (по умолчанию http://localhost:8010):
#   GET  /helth_check
#   POST /get_content    — тестовая обработка JSON как из Kafka
# Swagger: http://localhost:8010/swagger.yaml
```
- `/helth_check` — статус Kafka, счётчики (in_flight, processed, failed) и перечень текущих `request_id`.
- `/get_content` — выполняет обработку тела запроса и возвращает JSON, аналогичный сообщениям `url-metadata-results` и `url-content-analysis`.

### Docker

```bash
# сборка локального образа
make docker-build

docker run --rm -p 8010:8010  --add-host kafka:192.168.130.112   -e WEBSCRAPER_KAFKA_BROKERS=kafka:9092   webscraper:local
```

### Быстрая проверка REST

```bash
curl --noproxy "*" http://localhost:8010/helth_check | jq

curl --noproxy "*" -X POST http://localhost:8010/get_content \
  -H 'Content-Type: application/json' \
  -d @sample-request.json | jq
```

## Формат сообщений
- `url-metadata-results` — `models.MetadataMessage` с диагностикой, доменной информацией, статусом, стратегией, предупреждениями и `content_id` для привязки контента.
- `url-content-analysis` — `models.ContentMessage` с очищенным текстом, стратегией, метрикой, `content_id`, `user_agent` и `request_id`. Публикуется только если контент получен.

## Логирование
- Структурированные логи `slog` в stdout (уровни `DEBUG/INFO/WARN/ERROR`) с атрибутами: стадия обработки, `request_id`, `content_id`, `worker_id` и т.д.
