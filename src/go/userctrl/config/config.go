package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config - объединение всех переменных
	Config struct {
		App      `yaml:"app"`
		HTTP     `yaml:"http"`
		Log      `yaml:"logger"`
		Swagger  `yaml:"swagger"`
		KeyCloak `yaml:"keycloak"`
		Grafana  `yaml:"grafana"`
		BlogCl   `yaml:"blog_cl"`
	}

	// App - основные метаданные приложения
	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"USER_CTRL_APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"USER_CTRL_APP_VERSION"`
	}

	// HTTP - переменные для настройки http сервера
	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"USER_CTRL_HTTP_PORT"`
		Host string `env-required:"true" yaml:"host" env:"USER_CTRL_HTTP_HOST"`
	}

	// Log - настройка уровня логирования
	Log struct {
		Level string `env-required:"true" yaml:"log_level"   env:"USER_CTRL_LOG_LEVEL"`
	}

	// Swagger - хост для swagger
	Swagger struct {
		Host string `env-required:"true" yaml:"host" env:"USER_CTRL_SWAGGER_HOST"`
	}

	// KeyCloak - подключение к KeyCloak
	KeyCloak struct {
		Proto               string `env-required:"true" yaml:"proto" env:"USER_CTRL_KEYCLOAK_PROTO"`
		Port                string `env-required:"true" yaml:"port" env:"USER_CTRL_KEYCLOAK_PORT"`
		Host                string `env-required:"true" yaml:"host" env:"USER_CTRL_KEYCLOAK_HOST"`
		Login               string `env-required:"true" yaml:"login" env:"USER_CTRL_KEYCLOAK_LOGIN"`
		Pass                string `env-required:"true" yaml:"password" env:"USER_CTRL_KEYCLOAK_PASS"`
		Realm               string `env-required:"true" yaml:"realm" env:"USER_CTRL_KEYCLOAK_REALM"`
		SystemRealm         string `env-required:"true" yaml:"system_realm" env:"USER_CTRL_KEYCLOAK_SYSTEM_REALM"`
		GrafanaClientID     string `env-required:"true" yaml:"grafana_client_id" env:"USER_CTRL_KEYCLOAK_GRAFANA_CLIENT_ID"`
		GrafanaClientUUID   string `env-required:"true" yaml:"grafana_client_uuid" env:"USER_CTRL_KEYCLOAK_GRAFANA_CLIENT_UUID"`
		GrafanaClientSecret string `env-required:"true" yaml:"grafana_client_secret" env:"USER_CTRL_KEYCLOAK_GRAFANA_CLIENT_SECRET"`
	}

	// Grafana - для автоматизации выставления кук графаны при авторизации
	Grafana struct {
		Proto    string `env-required:"true" yaml:"proto" env:"USER_CTRL_GRAFANA_PROTO"`
		Port     string `env-required:"true" yaml:"port" env:"USER_CTRL_GRAFANA_PORT"`
		Host     string `env-required:"true" yaml:"host" env:"USER_CTRL_GRAFANA_HOST"`
		AuthPath string `env-required:"true" yaml:"auth_path" env:"USER_CTRL_GRAFANA_AUTH_PATH"`
		Username string `env-required:"true" yaml:"username" env:"USER_CTRL_GRAFANA_USERNAME"`
		Pass     string `env-required:"true" yaml:"pass" env:"USER_CTRL_GRAFANA_PASS"`
	}

	// TODO: везде клиентов других сервисов переписать под такой стандарт
	// BlogCL - клиент для создания бизнес-логов
	BlogCl struct {
		Proto string `env-required:"true" yaml:"proto" env:"USER_CTRL_BLOG_PROTO"`
		Host  string `env-required:"true" yaml:"host" env:"USER_CTRL_BLOG_HOST"`
		Port  string `env-required:"true" yaml:"port" env:"USER_CTRL_BLOG_PORT"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}

	// Происходит чтение конфига, после чтение переменных окружения
	// Приоритет в пользу переменных окружения
	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}
