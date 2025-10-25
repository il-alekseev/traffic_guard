package kafkaio

import (
	"context"
	"errors"
	"time"

	"github.com/segmentio/kafka-go"

	"scrapper/config"
)

const (
	defaultBatchTimeout = 100 * time.Millisecond
	defaultBatchSize    = 100
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(topic string, cfg config.KafkaConfig) (*Producer, error) {
	if len(cfg.Brokers) == 0 {
		return nil, errors.New("kafka brokers not configured")
	}
	if topic == "" {
		return nil, errors.New("kafka topic not provided")
	}

	dialer, err := NewDialer(cfg)
	if err != nil {
		return nil, err
	}

	transport := &kafka.Transport{
		SASL:        dialer.SASLMechanism,
		TLS:         dialer.TLS,
		DialTimeout: dialer.Timeout,
		ClientID:    dialer.ClientID,
	}

	writer := &kafka.Writer{
		Addr:                   kafka.TCP(cfg.Brokers...),
		Topic:                  topic,
		Balancer:               &kafka.Hash{},
		AllowAutoTopicCreation: false,
		RequiredAcks:           kafka.RequireAll,
		BatchTimeout:           defaultBatchTimeout,
		BatchSize:              defaultBatchSize,
		Transport:              transport,
	}

	return &Producer{writer: writer}, nil
}

func (p *Producer) Publish(ctx context.Context, key, value []byte) error {
	if p == nil || p.writer == nil {
		return errors.New("kafka producer not initialized")
	}

	msg := kafka.Message{
		Key:     key,
		Value:   value,
		Time:    time.Now().UTC(),
		Headers: nil,
	}

	return p.writer.WriteMessages(ctx, msg)
}

func (p *Producer) Close() error {
	if p == nil || p.writer == nil {
		return nil
	}
	return p.writer.Close()
}
