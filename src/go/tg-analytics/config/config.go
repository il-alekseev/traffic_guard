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
		PG      `yaml:"postgres"`
		HTTP    `yaml:"http"`
		Swagger `yaml:"swagger"`
		Metrics `yaml:"metrics"`
	}
	// App -.
	App struct {
		Name       string `env-required:"true" yaml:"name"    env:"DBD_APP_NAME"`
		Version    string `env-required:"true" yaml:"version" env:"DBD_APP_VERSION"`
		DevVersion string `yaml:"dev_version"`
	}
	// Log -.
	Log struct {
		Level string `env-required:"true" yaml:"log_level"   env:"DBD_LOG_LEVEL"`
	}

	// PG -.
	PG struct {
		PoolMax int    `env-required:"true" yaml:"pool_max" env:"DBD_PG_POOL_MAX"`
		Host    string `env-required:"true" yaml:"host" env:"DBD_PG_HOST"`
		Port    string `env-required:"true" yaml:"port" env:"DBD_PG_PORT"`
		User    string `env-required:"true" yaml:"user" env:"DBD_USER"`
		Pass    string `env-required:"true" yaml:"pass" env:"DBD_PASS"`
		DBName  string `env-required:"true" yaml:"dbname" env:"DBD_DBNAME"`
		SSLMode string `env-required:"true" yaml:"sslmode" env:"DBD_SSLMODE"`
	}

	// HTTP -.
	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"DBD_HTTP_PORT"`
		Host string `env-required:"true" yaml:"host" env:"DBD_HTTP_HOST"`
	}

	Swagger struct {
		Host string `env-required:"true" yaml:"host" env:"DBD_SWAGGER_HOST"`
	}
	Metrics struct {
		PoolMax int    `env-required:"true" yaml:"pool_max" env:"DBD_METRICS_POOL_MAX"`
		Host    string `env-required:"true" yaml:"host" env:"DBD_METRICS_HOST"`
		Port    string `env-required:"true" yaml:"port" env:"DBD_METRICS_PORT"`
		User    string `env-required:"true" yaml:"user" env:"DBD_METRICS_USER"`
		Pass    string `env-required:"true" yaml:"pass" env:"DBD_METRICS_PASS"`
		DBName  string `env-required:"true" yaml:"dbname" env:"DBD_METRICS_DBNAME"`
		SSLMode string `env-required:"true" yaml:"sslmode" env:"DBD_METRICS_SSLMODE"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Ошибка получения текущей директории: %v\n", err)
		return nil, err
	}
	fmt.Printf("Текущая директория: %s\n", currentDir)

	cfgPath := "./config/config.yaml"

	if _, err := os.Stat(cfgPath); err == nil {
		fmt.Printf("Файл %s существует\n", cfgPath)
	} else if os.IsNotExist(err) {
		fmt.Printf("Файл %s не существует\n", cfgPath)
		return nil, err
	} else {
		fmt.Printf("Ошибка при проверке файла: %v\n", err)
		return nil, err
	}

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
