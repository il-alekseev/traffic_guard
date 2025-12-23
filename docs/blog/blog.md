# blog service

## Общее описание сервиса blog

**Сервис blog** - используется как единая точка управления логами (получение, добавление). Используемая версия языка программирования Go: 1.24.2

### Основные функции
1. Документирование API (Swagger/OpenAPI) - предоставляет интерактивную информацию по адресу http://BLOG_SWAGGER_HOST/swagger/index.html
2. Получение всех логов с фильтрацией и пагинацией
3. Добавление логов

### Зависимости
Для корректной работы требуется предварительно запущенный PostgreSQL

### Таблицы
**Структура таблицы в БД**
```sql
CREATE TABLE business_logs (
  id SERIAL PRIMARY KEY,
  timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
  event_type VARCHAR(10) NOT NULL, -- CREATE, UPDATE, DELETE
  entity VARCHAR(20) NOT NULL, -- user, context
  username VARCHAR(36) NOT NULL,
  user_role VARCHAR(2) NOT NULL, -- SA, CA
  context VARCHAR(50),
  entity_id VARCHAR(50) NOT NULL, -- userID or ContextID
  old_value JSONB,
  new_value JSONB,
  description TEXT
);
```

### Технологии
- **Язык**: Go 1.24.2+
- **API**: REST (Gin), Swagger-документация
- **База данных**: PostgreSQL
- **Логирование**: slog со структурированным выводом
- **Документация**: Swagger/OpenAPI
- **Конфигурация**: cleanenv для загрузки настроек
- **Развертывание**: Docker-контейнеры

## Сборка и эксплуатация
Для сборки приложения используется multistage docker build

Требования:
- Хост с установленным docker
- Доступ в интернет для корректной установки базовых образов и загрузки зависимостей
- Установленный make

Все команды по сборке и запуску вынесены в `Makefile`
Поддерживаемые команды: 
   `all` - 
   `build` - сборка docker образа
   `run` - локальный запуск приложения с использованием компилятора Go
   `stop` - остановка docker контейнера
   `docker-up` - сборка и запуск с использованием docker compose
   `swagger` - генерация swagger документации

## Поддерживаемые переменные среды
**App**
- BLOG_APP_NAME - string ("fiermon-blog")
- BLOG_APP_VERSION - string ("1.0.0")

**Http server**
- BLOG_HTTP_HOST - string ("127.0.0.1")
- BLOG_HTTP_PORT - string ("8004")
- BLOG_HTTP_TIMEOUT - time.Duration (4s)
- BLOG_HTTP_IDLE_TIMEOUT - time.Duration (60s)

**Уровень логирования**
- LOG_LEVEL - string ("debug")

**Подключение к БД Postgres**
- BLOG_PG_POOL_MAX - int (2)
- BLOG_PG_HOST - string ("1.1.1.1")
- BLOG_PG_PORT - string ("5432")
- BLOG_PG_USER - string ("blog")
- BLOG_PG_PASSWORD - string ("blogpassword")
- BLOG_PG_DBNAME - string ("blog")
- BLOG_PG_SSLMODE - string ("disable")

**Swagger**
- BLOG_SWAGGER_HOST - string ("127.0.0.1"")

## Особенности запросов
Для каждого запроса необходима Bearer Authorization, токен выдается микросервисом `api-gw`.
В данном токене есть уровни доступа которые влияют на выдачу логов:
- SA - доступны все логи
- CA - доступны только логи с ролью CA

Если роль не определена - ошибка `there is no suitable role`