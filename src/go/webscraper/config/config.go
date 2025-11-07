package config

import (
	"fmt"
	"net"
	"os"
	"strings"
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
	Strategies         []string `yaml:"strategies"`
	ArchiveStrategies  []string `yaml:"archive_strategies"`
	MinTextLength      int      `yaml:"min_text_length"`
	RepositoryPath     string   `yaml:"repository_path"`
	MaxContentLength   int      `yaml:"max_content_length"`
	UserAgentsPath     string   `yaml:"user_agents_path"`
	AntiBotPhrasesPath string   `yaml:"anti_bot_phrases_path"`
	StripJSONFragments bool     `yaml:"strip_json_fragments"`
}

type LoggingConfig struct {
	Level string `yaml:"level"`
}

type KafkaAuthConfig struct {
	Username      string `yaml:"username"`
	Password      string `yaml:"password"`
	SASLMechanism string `yaml:"sasl_mechanism"`
	TLSEnabled    bool   `yaml:"tls_enabled"`
}

type KafkaConfig struct {
	Brokers        []string        `yaml:"brokers"`
	GroupID        string          `yaml:"group_id"`
	ClientID       string          `yaml:"client_id"`
	InputTopic     string          `yaml:"input_topic"`
	MetadataTopic  string          `yaml:"metadata_topic"`
	ContentTopic   string          `yaml:"content_topic"`
	CommitInterval Duration        `yaml:"commit_interval"`
	PollTimeout    Duration        `yaml:"poll_timeout"`
	Auth           KafkaAuthConfig `yaml:"auth"`
}

type HTTPConfig struct {
	Address         string   `yaml:"address"`
	ReadTimeout     Duration `yaml:"read_timeout"`
	WriteTimeout    Duration `yaml:"write_timeout"`
	ShutdownTimeout Duration `yaml:"shutdown_timeout"`
}

type Config struct {
	DBPath         string        `yaml:"db_path"`
	Workers        int           `yaml:"workers"`
	RequestTimeout Duration      `yaml:"request_timeout"`
	Content        ContentConfig `yaml:"content"`
	Logging        LoggingConfig `yaml:"logging"`
	Kafka          KafkaConfig   `yaml:"kafka"`
	HTTP           HTTPConfig    `yaml:"http"`
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

func ApplyEnvOverrides(cfg *Config) {
	if cfg == nil {
		return
	}

	if brokers := os.Getenv("WEBSCRAPER_KAFKA_BROKERS"); brokers != "" {
		values := strings.Split(brokers, ",")
		sanitized := make([]string, 0, len(values))
		for _, b := range values {
			trimmed := strings.TrimSpace(b)
			if trimmed != "" {
				sanitized = append(sanitized, trimmed)
			}
		}
		if len(sanitized) > 0 {
			cfg.Kafka.Brokers = sanitized
		}
	}

	if addr := os.Getenv("WEBSCRAPER_HTTP_ADDRESS"); addr != "" {
		cfg.HTTP.Address = strings.TrimSpace(addr)
	}

	if archive := os.Getenv("WEBSCRAPER_ARCHIVE_STRATEGIES"); archive != "" {
		values := strings.Split(archive, ",")
		normalized := make([]string, 0, len(values))
		for _, v := range values {
			trimmed := strings.TrimSpace(v)
			if trimmed != "" {
				normalized = append(normalized, trimmed)
			}
		}
		if len(normalized) > 0 {
			cfg.Content.ArchiveStrategies = normalized
		}
	}

	if raw := os.Getenv("WEBSCRAPER_HOST_ALIASES"); raw != "" {
		aliases := parseHostAliases(raw)
		if len(aliases) > 0 && len(cfg.Kafka.Brokers) > 0 {
			cfg.Kafka.Brokers = rewriteBrokerHosts(cfg.Kafka.Brokers, aliases)
		}
	}

	if path := os.Getenv("WEBSCRAPER_ANTIBOT_PHRASES_PATH"); path != "" {
		cfg.Content.AntiBotPhrasesPath = strings.TrimSpace(path)
	}

	if strip := os.Getenv("WEBSCRAPER_STRIP_JSON"); strip != "" {
		cfg.Content.StripJSONFragments = strings.EqualFold(strip, "1") ||
			strings.EqualFold(strip, "true") ||
			strings.EqualFold(strip, "yes")
	}
}

func parseHostAliases(raw string) map[string]string {
	result := make(map[string]string)
	if strings.TrimSpace(raw) == "" {
		return result
	}

	separators := func(r rune) bool {
		switch r {
		case ',', ';', '\n':
			return true
		default:
			return false
		}
	}

	pairs := strings.FieldsFunc(raw, separators)
	for _, pair := range pairs {
		if pair == "" {
			continue
		}
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			continue
		}
		host := strings.ToLower(strings.TrimSpace(parts[0]))
		addr := strings.TrimSpace(parts[1])
		if host == "" || addr == "" {
			continue
		}
		result[host] = addr
	}

	return result
}

func rewriteBrokerHosts(brokers []string, aliases map[string]string) []string {
	if len(brokers) == 0 || len(aliases) == 0 {
		return brokers
	}

	rewritten := make([]string, 0, len(brokers))
	for _, broker := range brokers {
		trimmed := strings.TrimSpace(broker)
		if trimmed == "" {
			continue
		}

		host, port, err := net.SplitHostPort(trimmed)
		if err != nil {
			rewritten = append(rewritten, trimmed)
			continue
		}

		if mapped, ok := aliases[strings.ToLower(host)]; ok && mapped != "" {
			rewritten = append(rewritten, net.JoinHostPort(mapped, port))
		} else {
			rewritten = append(rewritten, trimmed)
		}
	}

	return rewritten
}
