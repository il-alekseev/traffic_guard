package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config -.
	Config struct {
		App           `yaml:"app"`
		HTTP          `yaml:"http"`
		Swagger       `yaml:"swagger"`
		Log           `yaml:"logger"`
		KSU           `yaml:"ksu_postgres"`
		ETL           `yaml:"etl_postgres"`
		EtlController `yaml:"etl_controller"`
		Kafka         `yaml:"kafka"`
	}
	// App -.
	App struct {
		Name       string `env-required:"true" yaml:"name"    env:"ETL_APP_NAME"`
		Version    string `env-required:"true" yaml:"version" env:"ETL_APP_VERSION"`
		DevVersion string `yaml:"dev_version"`
	}
	// HTTP -.
	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"ETL_HTTP_PORT"`
		Host string `env-required:"true" yaml:"host" env:"ETL_HTTP_HOST"`
	}
	Swagger struct {
		Host string `env-required:"true" yaml:"host" env:"ETL_SWAGGER"`
	}
	// Log -.
	Log struct {
		Level string `env-required:"true" yaml:"log_level"   env:"ETL_LOG_LEVEL"`
	}

	// KSU PG -.
	KSU struct {
		PoolMax int    `env-required:"true" yaml:"pool_max" env:"ETL_KSU_PG_POOL_MAX"`
		Host    string `env-required:"true" yaml:"host" env:"ETL_KSU_PG_HOST"`
		Port    string `env-required:"true" yaml:"port" env:"ETL_KSU_PG_PORT"`
		User    string `env-required:"true" yaml:"user" env:"ETL_KSU_USER"`
		Pass    string `env-required:"true" yaml:"pass" env:"ETL_KSU_PASS"`
		DBName  string `env-required:"true" yaml:"dbname" env:"ETL_KSU_DBNAME"`
		SSLMode string `env-required:"true" yaml:"sslmode" env:"ETL_KSU_SSLMODE"`
	}

	// ETL PG -.
	ETL struct {
		PoolMax int    `env-required:"true" yaml:"pool_max" env:"ETL_PG_POOL_MAX"`
		Host    string `env-required:"true" yaml:"host" env:"ETL_PG_HOST"`
		Port    string `env-required:"true" yaml:"port" env:"ETL_PG_PORT"`
		User    string `env-required:"true" yaml:"user" env:"ETL_USER"`
		Pass    string `env-required:"true" yaml:"pass" env:"ETL_PASS"`
		DBName  string `env-required:"true" yaml:"dbname" env:"ETL_DBNAME"`
		SSLMode string `env-required:"true" yaml:"sslmode" env:"ETL_SSLMODE"`
	}

	EtlController struct {
		CacheTTL  int  `env-required:"true" yaml:"cache_ttl" env:"ETL_CACHE_TTL"`
		Refresh   uint `env-required:"true" yaml:"refresh" env:"ETL_REFRESH"`
		BatchSize uint `env-required:"true" yaml:"batchsize" env:"ETL_BATCHSIZE"`
		MLAttemps uint `env-required:"true" yaml:"ml_attemps" env:"ETL_ML_ATTEMPS"`
	}

	Kafka struct {
		Host          string `env-required:"true" yaml:"host" env:"ETL_KAFKA_HOST"`
		Port          string `env-required:"true" yaml:"port" env:"ETL_KAFKA_PORT"`
		URLTopic      string `yaml:"url_topic" env:"ETL_KAFKA_URL_TOPIC"`
		MetadataTopic string `yaml:"metadata_topic" env:"ETL_KAFKA_METADATA_TOPIC"`
		MLTopic       string `yaml:"ml_topic" env:"ETL_KAFKA_ML_TOPIC"`
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
