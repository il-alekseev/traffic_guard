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
- `content.archive_strategies` — резервные стратегии (например, `wayback`, `archive.today`) которые пробуются только если обнаружены анти-бот страницы
- `content.min_text_length` — минимальная длина очищенного текста
- `content.max_content_length` — максимальная длина текста (обрезается при превышении)
- `content.user_agents_path` — путь к списку User-Agent строк (по одной на строку)
- `content.anti_bot_phrases_path` — словарь анти-бот/капча фраз (по одной на строку), используется для быстрого обнаружения заглушек
- `content.strip_json_fragments` — если `true`, после очистки HTML удаляются крупные JSON-фрагменты (гидрация front-end)
- `logging.level` — уровень логирования (`debug`, `info`, `warn`, `error`)
- `kafka.brokers` — список брокеров Kafka
- `kafka.group_id`, `kafka.client_id` — параметры consumer group
- `kafka.input_topic` — топик входящих задач (`url-processing-requests`)
- `kafka.metadata_topic` — топик метаданных (`url-metadata-results`)
- `kafka.content_topic` — топик контента (`url-content-analysis`)
- `kafka.commit_interval`, `kafka.poll_timeout` — параметры чтения/коммитов
- `kafka.auth` — блок SASL/TLS авторизации (при необходимости)
- `http.address`, `http.read_timeout`, `http.write_timeout`, `http.shutdown_timeout` — параметры REST-сервера
- Переменные окружения:
  - `WEBSCRAPER_KAFKA_BROKERS` — перекрывает список брокеров из YAML (через запятую)
  - `WEBSCRAPER_HTTP_ADDRESS` — задаёт адрес REST-сервера
  - `WEBSCRAPER_HOST_ALIASES` — карта `hostname=ip` (через запятую/точку с запятой/перевод строки). Например `kafka=192.168.130.112` заставит сервис подключаться к брокеру по IP, не требуя `--add-host`.
  - Пример заполнения лежит в корне репозитория: `cp .env.example .env && nano .env`.

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

cat <<'EOF' > .env
WEBSCRAPER_KAFKA_BROKERS=kafka:9092
WEBSCRAPER_HOST_ALIASES=kafka=192.168.130.112
WEBSCRAPER_HTTP_ADDRESS=:8010
EOF

# -p <host_port>:<container_port>. В контейнере сервис слушает WEBSCRAPER_HTTP_ADDRESS, снаружи можно пробросить любой порт.
docker run --rm -p 8010:8010 --env-file .env webscraper:local
```

### Быстрая проверка REST

```bash
curl --noproxy "*" http://localhost:8010/helth_check | jq

curl --noproxy "*" -X POST http://localhost:8010/get_content \
  -H 'Content-Type: application/json' \
  -d @sample-request.json | jq
```

### Инструменты

**Kafka publisher** (`tools/url_publisher`)

```bash
go run ./tools/url_publisher \
  -input urls.txt \
  -env .env \
  -config config/config.yaml \
  -src-ip 10.0.0.5 \
  -dst-type url
```

- читает URL из файла и публикует `ProcessingRequest` в Kafka;
- при `-dry-run` только печатает `request_id`;
- `-env` позволяет переиспользовать готовый `.env` с брокерами.

**REST tester** (`tools/http_tester`)

```bash
go run ./tools/http_tester \
  -input urls.txt \
  -endpoint http://localhost:8010/get_content \
  -src-ip 10.0.0.5 \
  -dst-type url \
  -delay 500ms \
  -retries 3 \
  -retry-delay 250ms
```

- собирает тот же JSON, что приходит из Kafka, и дергает `/get_content`;
- назначает `request_id` (или использует `-request-id`);
- печатает статус и полный JSON-ответ, поддерживает автоповторы при сетевых ошибках.

### Обход анти-бот страниц

- В `data/anti_bot_phrases.txt` хранится список фраз, по которым определяется, что сайт показал заглушку/капчу (например, «Вы не робот?», SmartCaptcha, отключённый JavaScript). Если очищенный текст короче 500 символов и содержит хотя бы одну фразу, стратегия немедленно отклоняется, предупреждение попадает в метаданные, а контент больше не публикуется в `url-content-analysis`.
- При срабатывании детектора сначала выполняются до трёх повторных попыток с разными User-Agent из `user_agents.txt` (порядок выбирается случайно). Если даже это не помогает, автоматически подключаются архивные стратегии (`wayback`, `archive.today`), которые пытаются вытянуть текст из Wayback Machine и archive.today. Эти стратегии перечислены в `content.archive_strategies` и запускаются только когда действительно нужны.
- Любой очищенный текст короче ~150 символов считается недостоверным: сервис меняет User-Agent до трёх раз и, при неудаче, переключается на архивные источники.
- HTTP статусы 401/403/429/451 сразу вызывают перебор User-Agent + постановку Wayback/Archive.today в очередь.
- В топик `url-content-analysis` теперь попадают только полностью успешные (`status=ok`) результаты — частичные ответы с Failure причинами отбрасываются.
- Чтобы пополнить словарь фраз, достаточно добавить строки в `data/anti_bot_phrases.txt` (формат одна фраза на строку, комментарии можно начинать с `#`).

## Стратегии получения контента

1. **Trafilatura** — внешняя библиотека, которая сама ходит в сеть, отбрасывает рекламу и возвращает чистый текст. Используем в первую очередь.
2. **HTTP** — наш универсальный клиент: скачивает HTML, прогоняет через cleaner, запускает эвристики.
3. **Wayback** — архивный снимок из Wayback Machine. Берём «сырой» `/web/0id_…`, декодируем оригинальную кодировку (включая `X-Archive-Orig-Content-Encoding`).
4. **Archive.today** — зеркало archive.today, если Wayback не помог.

Каждая стратегия:
- перебирает до 3 произвольных User-Agent из `data/user_agents.txt`;
- при 401/403/429/451 переключается на новый User-Agent и добавляет архивы;
- после очистки режет встроенные JSON-фрагменты и проверяет длину текста (минимум ~150 символов);
- в случае капчи или короткого текста автоматически ставит в очередь архивные стратегии.

Только результаты со статусом `ok` попадают в `url-content-analysis`.

## GeoIP база (MaxMind/IPinfo)

Мы используем `data/ipinfo_lite.mmdb` (совместимо с MaxMind). Свежий файл можно бесплатно получить у [ipinfo.io](https://ipinfo.io/products/free-ip-database):
1. Зарегистрироваться и запросить бесплатный доступ к Lite-базе.
2. Скачать `ipinfo-lite.mmdb`.
3. Положить файл в `data/ipinfo_lite.mmdb` либо обновить `db_path` в `config/config.yaml`.

## Формат сообщений
- `url-metadata-results` — `models.MetadataMessage` с диагностикой, доменной информацией, статусом, стратегией, предупреждениями и `content_id` для привязки контента.
- `url-content-analysis` — `models.ContentMessage` с очищенным текстом, стратегией, метрикой, `content_id`, `user_agent` и `request_id`. Публикуется только если контент получен.

## Логирование
- Структурированные логи `slog` в stdout (уровни `DEBUG/INFO/WARN/ERROR`) с атрибутами: стадия обработки, `request_id`, `content_id`, `worker_id` и т.д.
