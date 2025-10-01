# Scrapper

## Описание
Scrapper — сервис на Go для массового снятия HTML-страниц, извлечения очищенного текста и обогащения доменной информации. На вход подаются URL (из файла, stdin или напрямую через `--url`), на выходе формируется JSON с текстом, данными о домене и диагностикой применённых стратегий.

## Используемый стек
- Go 1.25
- github.com/markusmobius/go-trafilatura — извлечение основного контента
- github.com/oschwald/maxminddb-golang — GeoIP (локальный `data/ipinfo_lite.mmdb`)
- golang.org/x/net/html, html/charset — парсинг и декодирование HTML
- Локальный JSONL (`data/attempts.jsonl`) и файл логов (`logs/scraper.log`)

## Структура
- `cmd/scraper` — CLI-вход, загрузка конфига, запуск runtime
- `config` — YAML-конфиг и загрузчик
- `internal/app` — оболочка над usecase
- `internal/domain` — доменные структуры (`ScrapeResult`, `FailureReason` и др.)
- `internal/usecase` — воркер-пул, применение стратегий, диагностика
- `internal/infrastructure` — адаптеры (`content`, `httpclient`, `dns`, `geo`, `parser`, `quality`, `storage`, `logging`)

## Конфигурация (`config/config.yaml`)
- `input_path` — файл с URL или `-` для stdin
- `db_path` — путь к MaxMind базе
- `workers` — количество воркеров (≤0 → `runtime.NumCPU()*2`)
- `request_timeout`, `http_timeout` — таймауты на обработку и HTTP
- `content.strategies` — порядок стратегий (`trafilatura`, `http`, ...)
- `content.min_text_length` — минимальная длина очищенного текста
- `content.repository_path` — файл-аудит JSONL
- `logging.file` — путь к лог-файлу (по умолчанию `logs/scraper.log`)

## Быстрый старт
```bash
# зависимости
go mod tidy

# сборка
go build -o scraper ./cmd/scraper

# запуск по конфигу
./scraper --config config/config.yaml

# отладка единичных URL (override input_path)
./scraper --url https://example.com --url https://another.site
```

## Формат вывода
На stdout выводится массив `[]ScrapeResult`:
- `url` — исходный адрес
- `domain` — name/ip/geo (континент, страна, ASN и т.д.)
- `status` — `ok`, `partial`, `error`
- `content` — очищенный текст
- `strategy` — стратегия, давшая результат
- `metric` — числовая метрика качества
- `warnings` — массив диагностических сообщений
- `failure` — код причины (`fetch_error`, `clean_error`, `evaluate_error`, `empty_content`, `requirements_not_met`, `no_successful_strategy`)
- `error` — краткое описание в зависимости от `failure`

## Логирование и аудит
- `logs/scraper.log` — worker, домен, стратегия, статус, причина
- `data/attempts.jsonl` — подробности всех попыток (URL, стратегия, ошибки, timestamp)

## TODO / roadmap
- Интеграция с Kafka (публикация результатов и диагностик)
- Docker-образ и docker-compose для развёртывания
- REST API: Swagger/OpenAPI, health/readiness endpoints
- Proxy fallback: использование прокси (по `client_external_ip` или конфигу), если прямой запрос неудачен
- Стратегия Wayback Machine для недоступных страниц
- AI-поддержка (firecrawl.dev, Qwen, «Яндекс Пересказ»)
- Подключение БД (Postgres/ClickHouse/SQLite) для долговременного хранения
- Управление прокси (глобальные настройки, per-strategy, пулы)

## Дальнейшие шаги
- Kafka + Docker
- REST слой с Swagger и healthcheck
- Реализация дополнительных стратегий (Wayback, AI)
- Вынос результатов в БД вместо файлов
- Система прокси/обход блокировок

