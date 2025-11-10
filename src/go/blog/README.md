# Fiermon-blog - микросервис для отображения логов из БД Postgres c пагинацией и фильтрацией
**Поддерживаемая структура таблицы в БД**
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

## Сборка и эксплуатация
Для сборки приложения используется multistage docker build

Требования:
- Хост с установленым docker
- Доступ в интернет для корректной установки базовых образов и загрузки зависимостей
- Установленный make

### Сборка приложения 
```bash
sudo make build
```
для проверки используйте `docker images`, доджна быть строка `myproject/bislog              latest`

### Эксплуатация
Для настройки контейнера требуется конфигурационный файл. Пример конфигурационного файла расположен в дирректории `config/config.yml`
Путь к кастомному файлу указывается в переменной `CONFIG_PATH` в Makefile

Запуск контейнера:
```bash
sudo make start
```

Остановка контейнера:
```bash
sudo make stop
```

Поддерживаемые переменные среды:
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
- BLOG_PG_HOST - string ("172.17.134.91")
- BLOG_PG_PORT - string ("5432")
- BLOG_PG_USER - string ("blog")
- BLOG_PG_PASSWORD - string ("blogpassword")
- BLOG_PG_DBNAME - string ("blog")
- BLOG_PG_SSLMODE - string ("disable")

**Swagger**
- BLOG_SWAGGER_HOST - string ("127.0.0.1")

## Особенности запросов
Для каждого запроса необходима Bearer Authorization, токен выдается микросервисом `fiermon-api-gw`.
В данном токене есть уровни доступа которые влияют на выдачу логов:
- SA - доступны все логи
- CA - доступны только логи с ролью CA
- CO - доступны только логи с ролью CO

Если роль не определена - ошибка `there is no suitable role`