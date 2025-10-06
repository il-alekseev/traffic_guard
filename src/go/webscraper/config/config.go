package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalYAML(value *yaml.Node) error {
	var str string
	if err := value.Decode(&str); err == nil {
		if str == "" {
			d.Duration = 0
			return nil
		}

		dur, err := time.ParseDuration(str)
		if err != nil {
			return err
		}
		d.Duration = dur
		return nil
	}

	var num int64
	if err := value.Decode(&num); err == nil {
		d.Duration = time.Duration(num)
		return nil
	}

	return fmt.Errorf("invalid duration")
}

type ContentConfig struct {
    Strategies     []string `yaml:"strategies"`
    MinTextLength  int      `yaml:"min_text_length"`
    RepositoryPath string   `yaml:"repository_path"`
}

type LoggingConfig struct {
    File string `yaml:"file"`
}

type Config struct {
    InputPath      string        `yaml:"input_path"`
    DBPath         string        `yaml:"db_path"`
    Workers        int           `yaml:"workers"`
    RequestTimeout Duration      `yaml:"request_timeout"`
    HTTPTimeout    Duration      `yaml:"http_timeout"`
    Content        ContentConfig `yaml:"content"`
    Logging        LoggingConfig `yaml:"logging"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	return &cfg, nil
}
