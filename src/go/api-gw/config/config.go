package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config -.
	Config struct {
		App         `yaml:"app"`
		HTTP        `yaml:"http"`
		Log         `yaml:"logger"`
		Swagger     `yaml:"swagger"`
		UserControl `yaml:"usercontrol"`
		
		BlogServ  `yaml:"blog_serv"`
		Analytics `yaml:"analytics"`
		KeyCloak  `yaml:"keycloak"`
	}

	// App - contains basic application metadata.
	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"API_GW_APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"API_GW_APP_VERSION"`
	}

	// HTTP - contains web server configuration.
	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"API_GW_HTTP_PORT"`
		Host string `env-required:"true" yaml:"host" env:"API_GW_HTTP_HOST"`
	}

	// Log - contains logging configuration parameters.
	Log struct {
		Level string `env-required:"true" yaml:"log_level"   env:"API_GW_LOG_LEVEL"`
	}

	// Swagger - хост, который будет использоваться для swagger
	Swagger struct {
		Host string `env-required:"true" yaml:"host" env:"API_GW_SWAGGER_HOST"`
	}

	// UserControl - параметры UserControl сервиса
	UserControl struct {
		Proto string `env-required:"true" yaml:"proto" env:"API_GW_USERCONTROL_PROTO"`
		Host  string `env-required:"true" yaml:"host" env:"API_GW_USERCONTROL_HOST"`
		Port  string `env-required:"true" yaml:"port" env:"API_GW_USERCONTROL_PORT"`
	}

	// BlogServ - параметры BlogServ сервиса
	BlogServ struct {
		Proto string `env-required:"true" yaml:"proto" env:"API_GW_BLOG_SERV_PROTO"`
		Host  string `env-required:"true" yaml:"host" env:"API_GW_BLOG_SERV_HOST"`
		Port  string `env-required:"true" yaml:"port" env:"API_GW_BLOG_SERV_PORT"`
	}

	// Analytics - параметры Analytics сервиса
	Analytics struct {
		Proto string `env-required:"true" yaml:"proto" env:"API_GW_ANALYTICS_PROTO"`
		Host  string `env-required:"true" yaml:"host" env:"API_GW_ANALYTICS_HOST"`
		Port  string `env-required:"true" yaml:"port" env:"API_GW_ANALYTICS_PORT"`
	}

	KeyCloak struct {
		// Публичный RSA (RS256) ключ для проверки подписи токенов
		// Получить в keycloak https://keycloak_host:port/{{realm}}/realm-settings/keys
		// realm - в начальной конфигурации равен tsum
		// подробнее в документации
		PemFile string `env-required:"true" yaml:"pem_fiel_path" env:"API_GW_KEYCLOAK_PEM"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}

	// Сначала читает параметры из конфига - удобно при разработке и локальном запуске.
	// Затем смотрит переменные окружения и перезаписывает параметры из конфига.
	// Таким образом параметры ENV имеют приоритет над конфигом - удобно при развертывании
	// в docker
	err := cleanenv.ReadConfig("./config/config.yml", cfg) //./config/config.yml   fiermon/fiermon-api-gw/config/config.yml
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}
