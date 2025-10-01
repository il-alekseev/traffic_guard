package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

const configPathEnvName = "BLOG_CONFIG_PATH"

var ErrNoFile = errors.New("config file not exist")

type (
	// Config -.
	Config struct {
		App     `yaml:"app"`
		HTTP    `yaml:"http"`
		Log     `yaml:"logger"`
		PG      `yaml:"postgres"`
		Swagger `yaml:"swagger"`
	}

	// App -.
	App struct {
		Name    string `env-required:"false" yaml:"name"    env:"BLOG_APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"BLOG_APP_VERSION"`
	}

	// HTTP -.
	HTTP struct {
		Host        string        `env-required:"true" yaml:"host" env:"BLOG_HTTP_HOST"`
		Port        string        `env-required:"true" yaml:"port" env:"BLOG_HTTP_PORT"`
		Timeout     time.Duration `env-required:"true" yaml:"timeout" env:"BLOG_HTTP_TIMEOUT"`
		IdleTimeout time.Duration `env-required:"true" yaml:"idle-timeout" env:"BLOG_HTTP_IDLE_TIMEOUT"`
	}

	// Log -.
	Log struct {
		Level string `env-required:"true" yaml:"log_level"   env:"LOG_LEVEL"`
	}

	// PG -.
	PG struct {
		PoolMax  int    `env-required:"true" yaml:"pool_max" env:"BLOG_PG_POOL_MAX"`
		Host     string `env-required:"true" yaml:"host" env:"BLOG_PG_HOST"`
		Port     string `env-required:"true" yaml:"port" env:"BLOG_PG_PORT"`
		User     string `env-required:"true" yaml:"user" env:"BLOG_PG_USER"`
		Password string `env-required:"true" yaml:"password" env:"BLOG_PG_PASSWORD"`
		DBName   string `env-required:"true" yaml:"dbname" env:"BLOG_PG_DBNAME"`
		SSLMode  string `env-default:"disable" yaml:"sslmode" env:"BLOG_PG_SSLMODE"`
	}

	Swagger struct {
		Host string `env-required:"true" yaml:"host" env:"BLOG_SWAGGER_HOST"`
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
