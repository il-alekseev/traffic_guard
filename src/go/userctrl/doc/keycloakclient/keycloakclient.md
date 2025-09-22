# Клиент Keycloak

## Требования к развертыванию клиента KeyCloack

### Требования к ролям учетной записи

- realm-management realm-admin - для создания пользователей, групп, ролей клиентов ![realm-admin](client_role.png)
  
### Требования к атрибутам пользователя

- patronymic - для добавления отчества пользователя ![patronymic](attribute.png)

## Описание полей клиента

```go

type KeycloakClient struct {
    basePath string
    client   *goCloak.GoCloak
    token    goCloak.JWT
    ctx      context.Context
    realm    string
    clientID string
    logger   logger.Interface
}

```

Структура KeycloakClient представляет собой клиентскую реализацию для работы с Keycloak API через библиотеку goCloak. Она инкапсулирует все необходимые компоненты для взаимодействия с сервером аутентификации.

### Поля структуры

- basePath - базовая часть URL для подключения к Keycloak
- client - указатель на экземпляр GoCloak, который является основным клиентом для взаимодействия с API Keycloak
- token - JWT-токен, используемый для аутентификации запросов к API
- ctx - контекст выполнения операций, который может быть использован для управления временем жизни запросов
- realm - название Realm в Keycloak, с которым работает клиент
- clientID - идентификатор клиента в Keycloak
- logger - интерфейс логгера для записи сообщений и ошибок

### Пример использования

```go

var url string = "https://localhost:8443"
var login string = "admin"
var password string = "admin"
var realm string = "tsum"
var l logger.Logger = *logger.New("debug")

var testClient *KeycloakClient = NewKeycloakClient(url, login, password, realm, "grafana-sso", &l)

```
