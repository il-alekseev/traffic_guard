package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config -.
	Config struct {
		App     `yaml:"app"`
		Log     `yaml:"logger"`
		ETL     `yaml:"etl"`
		HTTP    `yaml:"http"`
		Swagger `yaml:"swagger"`
		Metrics `yaml:"metrics"`
	}
	// App -.
	App struct {
		Name       string `env-required:"true" yaml:"name"    env:"AN_APP_NAME"`
		Version    string `env-required:"true" yaml:"version" env:"AN_APP_VERSION"`
		DevVersion string `yaml:"dev_version"`
	}
	// Log -.
	Log struct {
		Level string `env-required:"true" yaml:"log_level"   env:"AN_LOG_LEVEL"`
	}

	// PG -.
	ETL struct {
		PoolMax int    `env-required:"true" yaml:"pool_max" env:"AN_ETL_POOL_MAX"`
		Host    string `env-required:"true" yaml:"host" env:"AN_ETL_HOST"`
		Port    string `env-required:"true" yaml:"port" env:"AN_ETL_PORT"`
		User    string `env-required:"true" yaml:"user" env:"AN_ETL_USER"`
		Pass    string `env-required:"true" yaml:"pass" env:"AN_ETL_PASS"`
		DBName  string `env-required:"true" yaml:"dbname" env:"AN_ETL_DBNAME"`
		SSLMode string `env-required:"true" yaml:"sslmode" env:"AN_ETL_SSLMODE"`
	}

	// HTTP -.
	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"AN_HTTP_PORT"`
		Host string `env-required:"true" yaml:"host" env:"AN_HTTP_HOST"`
	}

	Swagger struct {
		Host string `env-required:"true" yaml:"host" env:"AN_SWAGGER"`
	}
	Metrics struct {
		PoolMax int    `env-required:"true" yaml:"pool_max" env:"AN_METRICS_POOL_MAX"`
		Host    string `env-required:"true" yaml:"host" env:"AN_METRICS_HOST"`
		Port    string `env-required:"true" yaml:"port" env:"AN_METRICS_PORT"`
		User    string `env-required:"true" yaml:"user" env:"AN_METRICS_USER"`
		Pass    string `env-required:"true" yaml:"pass" env:"AN_METRICS_PASS"`
		DBName  string `env-required:"true" yaml:"dbname" env:"AN_METRICS_DBNAME"`
		SSLMode string `env-required:"true" yaml:"sslmode" env:"AN_METRICS_SSLMODE"`
	}
)

// NewConfig returns app config.
func NewConfig(cfgPath string) (*Config, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Ошибка получения текущей директории: %v\n", err)
		return nil, err
	}
	fmt.Printf("Текущая директория: %s\n", currentDir)

	cfg := &Config{}
	err = cleanenv.ReadConfig(cfgPath, cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
