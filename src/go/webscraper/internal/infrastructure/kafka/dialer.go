package kafkaio

import (
	"crypto/tls"
	"fmt"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/segmentio/kafka-go/sasl/scram"

	"scrapper/config"
)

const defaultDialTimeout = 10 * time.Second

func NewDialer(cfg config.KafkaConfig) (*kafka.Dialer, error) {
	mech, err := buildSASLMechanism(cfg.Auth)
	if err != nil {
		return nil, err
	}

	var tlsConfig *tls.Config
	if cfg.Auth.TLSEnabled {
		tlsConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	return &kafka.Dialer{
		Timeout:       defaultDialTimeout,
		DualStack:     true,
		TLS:           tlsConfig,
		SASLMechanism: mech,
		ClientID:      cfg.ClientID,
	}, nil
}

func buildSASLMechanism(auth config.KafkaAuthConfig) (sasl.Mechanism, error) {
	if auth.SASLMechanism == "" {
		return nil, nil
	}

	switch strings.ToUpper(auth.SASLMechanism) {
	case "PLAIN":
		return plain.Mechanism{
			Username: auth.Username,
			Password: auth.Password,
		}, nil
	case "SCRAM-SHA-256", "SCRAM256", "SCRAM_SHA_256":
		return scram.Mechanism(scram.SHA256, auth.Username, auth.Password)
	case "SCRAM-SHA-512", "SCRAM512", "SCRAM_SHA_512":
		return scram.Mechanism(scram.SHA512, auth.Username, auth.Password)
	case "NONE":
		return nil, nil
	default:
		return nil, fmt.Errorf("unsupported SASL mechanism %q", auth.SASLMechanism)
	}
}
