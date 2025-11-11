# Описание компонента Keycloak

## Общее описание

Keycloak это инструмент для управление идентификацией и доступом к сервисам с открытым исходным кодом. Пользователи аутентифицируются с помощью Keycloak, а не отдельных приложений. Это означает, приложениям не нужно иметь дело с формами входа, аутентификацией пользователей и сохранением пользователей. После входа в систему благодаря Keycloak пользователям не нужно повторно входить в систему, чтобы получить доступ к другому приложению. Инструмент написан на языке Java. 

## Взаимосвязь с другими компонетами

  
  
![](Components/Keycloak/img/1.png)  
  

В рамках системы Keycloak, как инструмент для управление идентификацией и доступом к сервису Grafana.



# Описание настройки

Инструмент будет развернут в системе в виде Docker- контейнера в системе управления контейнерами Docker Compose.

## Настройка контейнера в системе

Перед запуском контейнера необходимо указать переменные окружения в двух .env файлах. Первый отвечает за саму работу компонента, а второй за парметры запуска контейнера.
  

`./deploy/keycloak/.env` :

  
```
KEYCLOAK_ADMIN=system_admin

KEYCLOAK_ADMIN_PASSWORD=admin

KC_PROXY_ADDRESS_FORWARDING=true

KC_HOSTNAME_STRICT=false

KC_HOSTNAME=192.168.130.114

KC_HOSTNAME_PORT=8443

KC_HTTPS_CERTIFICATE_FILE=/certs/fullchain.crt

KC_HTTPS_CERTIFICATE_KEY_FILE=/certs/keycloak.key

KC_HTTP_ENABLED=false

KC_HTTPS_ENABLED=true

KC_DB=postgres

KC_DB_URL=jdbc:postgresql://192.168.130.114:5433/keycloak

KC_DB_USERNAME=keycloak_user

KC_DB_PASSWORD=keycloak_password

KEYCLOAK_GRAFANA_ID=grafana-sso

KEYCLOAK_GRAFANA_SECRET=uAgEvDTaFYvzqr287p2qZFIi7habElf8

KEYCLOAK_GRAFANA_REDIRECT_ONE=https://192.168.130.114:3000/login/generic_oauth

KEYCLOAK_GRAFANA_REDIRECT_TWO=https://localhost:3000/login/generic_oauth

KEYCLOAK_GRAFANA_REDIRECT_THREE=https://grafana:3000/login/generic_oauth

KEYCLOAK_GRAFANA_PROTOCOL=openid-connect

KEYCLOAK_GLOBAL_PASSWORD=admin

WEBHOOK_HTTP_BASE_PATH=192.168.130.114
```


`./deploy/.env`:

```
...
KEYCLOAK_TAG=24.0.3

KEYCLOAK_PORT=8443

KEYCLOAK_VOLUME_PORT=8443
...
```

`docker-compose.yml`:

```yaml
keycloak:
	image: quay.io/keycloak/keycloak:${KEYCLOAK_TAG}
	# build:
		# context: ../patch-keycloak
		# dockerfile: Dockerfile
		# image: keycloak-webhook
	ports:
		- ${KEYCLOAK_VOLUME_PORT}:${KEYCLOAK_PORT}
	env_file: ./keycloak/.env
	volumes:
		- ./keycloak/realm-import.json:/opt/keycloak/data/import/realm-import.json
		- ./keycloak:/certs
		- ./keycloak/data:/opt/keycloak/data
	restart: always
	healthcheck:
		test: ["CMD-SHELL", "echo -n > /dev/tcp/localhost/${KEYCLOAK_PORT}"]
		interval: 10s
		timeout: 5s
		retries: 6
		start_period: 20s
	command: start --https-certificate-file=/certs/fullchain.crt --https-certificate-key-file=/certs/keycloak.key --import-realm --verbose
	depends_on:	
		configdb:
			condition: service_healthy
	networks:
		- continent
```


## Особенности настройки

Помимо настроек запуска компонента необходимо указать импортируемые настройки для работы с Grafana.

  
  ```json
  {
  "realm": "grafanasso",
  "enabled": true,
  "roles": {
    "client": {
      "grafana-sso": [
        {
          "name": "SA",
          "description": "",
          "composite": false,
          "clientRole": true,
          "attributes": {}
        }
      ]
    }
  },
  "clients": [
    {
      "clientId": "${KEYCLOAK_GRAFANA_ID}",
      "enabled": true,
      "secret": "${KEYCLOAK_GRAFANA_SECRET}",
      "redirectUris": [
        "${KEYCLOAK_GRAFANA_REDIRECT_ONE}",
        "${KEYCLOAK_GRAFANA_REDIRECT_TWO}",
        "${KEYCLOAK_GRAFANA_REDIRECT_THREE}"
      ],
      "standardFlowEnabled": true,
      "directAccessGrantsEnabled": true,
      "publicClient": false,
      "serviceAccountsEnabled": false,
      "protocol": "openid-connect",
      "protocolMappers": [
        {
          "name": "client-roles",
          "protocol": "${KEYCLOAK_GRAFANA_PROTOCOL}",
          "protocolMapper": "oidc-usermodel-client-role-mapper",
          "consentRequired": false,
          "config": {
            "multivalued": "true",
            "userinfo.token.claim": "true",
            "id.token.claim": "true",
            "access.token.claim": "true",
            "claim.name": "roles",
            "jsonType.label": "String"
          }
        },
        {
            "name": "username",
            "protocol": "openid-connect",
            "protocolMapper": "oidc-usermodel-property-mapper",
            "consentRequired": false,
            "config": {
                "userinfo.token.claim": "true",
                "user.attribute": "username",
                "id.token.claim": "true",
                "access.token.claim": "true",
                "claim.name": "preferred_username",
                "jsonType.label": "String"
            }
        },
        {
            "name": "email",
            "protocol": "openid-connect",
            "protocolMapper": "oidc-usermodel-property-mapper",
            "consentRequired": false,
            "config": {
                "userinfo.token.claim": "true",
                "user.attribute": "email",
                "id.token.claim": "true",
                "access.token.claim": "true",
                "claim.name": "email",
                "jsonType.label": "String"
            }
        },
        {
            "name": "groups",
            "protocol": "openid-connect",
            "protocolMapper": "oidc-group-membership-mapper",
            "consentRequired": false,
            "config": {
                "full.path": "true",
                "id.token.claim": "true",
                "access.token.claim": "true",
                "claim.name": "groups",
                "userinfo.token.claim": "true"
            }
        }
      ],
      "defaultClientScopes": [
        "profile",
        "roles",
        "email"
      ],
      "optionalClientScopes": []
    }
  ],
  
  "components": {
    "org.keycloak.userprofile.UserProfileProvider": [
      {
        "providerId": "declarative-user-profile",
        "subComponents": {},
        "config": {
          "kc.user.profile.config": [
            "{\"attributes\":[{\"name\":\"username\",\"displayName\":\"${username}\",\"validations\":{\"length\":{\"min\":3,\"max\":255},\"username-prohibited-characters\":{},\"up-username-not-idn-homograph\":{}},\"permissions\":{\"view\":[\"admin\",\"user\"],\"edit\":[\"admin\",\"user\"]},\"multivalued\":false},{\"name\":\"email\",\"displayName\":\"${email}\",\"validations\":{\"email\":{},\"length\":{\"max\":255}},\"annotations\":{},\"permissions\":{\"view\":[\"admin\",\"user\"],\"edit\":[\"admin\",\"user\"]},\"multivalued\":false},{\"name\":\"firstName\",\"displayName\":\"${firstName}\",\"validations\":{\"length\":{\"max\":255},\"person-name-prohibited-characters\":{}},\"required\":{\"roles\":[\"user\"]},\"permissions\":{\"view\":[\"admin\",\"user\"],\"edit\":[\"admin\",\"user\"]},\"multivalued\":false},{\"name\":\"lastName\",\"displayName\":\"${lastName}\",\"validations\":{\"length\":{\"max\":255},\"person-name-prohibited-characters\":{}},\"required\":{\"roles\":[\"user\"]},\"permissions\":{\"view\":[\"admin\",\"user\"],\"edit\":[\"admin\",\"user\"]},\"multivalued\":false},{\"name\":\"patronymic\",\"displayName\":\"Patronymic\",\"validations\":{},\"annotations\":{},\"permissions\":{\"view\":[\"admin\",\"user\"],\"edit\":[\"admin\",\"user\"]},\"multivalued\":false}],\"groups\":[{\"name\":\"user-metadata\",\"displayHeader\":\"User metadata\",\"displayDescription\":\"Attributes, which refer to user metadata\"}]}"
          ]
        }
      }
    ]
  },
  "users": [
    {
      "username": "admin",
      "enabled": true,
      "email": "admin@targos.com",
      "firstName": "Admin",
      "lastName": "Admin",
      "groups": ["global-SA"], 
      "credentials": [
          {
              "type": "password",
              "value": "${KEYCLOAK_GLOBAL_PASSWORD}",
              "temporary": false
          }
      ]
    }
  ],
  "groups": [
    {
      "name": "global-SA",
      "path": "/global-SA",
      "subGroups": [],
      "attributes": {},
      "realmRoles": [],
      "clientRoles": {
        "realm-management": [
          "realm-admin"
        ],
        "grafana-sso": [
          "SA"
        ]
      }
    }
  ]

}
  ```