package kafkaio

import (
	"context"
	"errors"
	"time"

	"github.com/segmentio/kafka-go"

	"scrapper/config"
)

const (
	defaultMinBytes    = 1
	defaultMaxBytes    = 10 << 20 // 10 MiB
	defaultPollTimeout = 2 * time.Second
	defaultCommitEvery = 5 * time.Second
)

var ErrEmptyPoll = errors.New("kafka poll timeout")

type Consumer struct {
	reader      *kafka.Reader
	pollTimeout time.Duration
}

func NewConsumer(cfg config.KafkaConfig) (*Consumer, error) {
	if len(cfg.Brokers) == 0 {
		return nil, errors.New("kafka brokers not configured")
	}
	if cfg.InputTopic == "" {
		return nil, errors.New("kafka input topic not configured")
	}
	if cfg.GroupID == "" {
		return nil, errors.New("kafka group id not configured")
	}

	dialer, err := NewDialer(cfg)
	if err != nil {
		return nil, err
	}

	readerConfig := kafka.ReaderConfig{
		Brokers:        cfg.Brokers,
		GroupID:        cfg.GroupID,
		Topic:          cfg.InputTopic,
		MinBytes:       defaultMinBytes,
		MaxBytes:       defaultMaxBytes,
		CommitInterval: resolveCommitInterval(cfg.CommitInterval.Duration),
		Dialer:         dialer,
	}

	reader := kafka.NewReader(readerConfig)

	return &Consumer{
		reader:      reader,
		pollTimeout: resolvePollTimeout(cfg.PollTimeout.Duration),
	}, nil
}

func (c *Consumer) Fetch(ctx context.Context) (kafka.Message, error) {
	ctx, cancel := c.withTimeout(ctx)
	defer cancel()

	msg, err := c.reader.FetchMessage(ctx)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return kafka.Message{}, ErrEmptyPoll
		}
		return kafka.Message{}, err
	}

	return msg, nil
}

func (c *Consumer) Commit(ctx context.Context, msg kafka.Message) error {
	return c.reader.CommitMessages(ctx, msg)
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}

func (c *Consumer) withTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	timeout := c.pollTimeout
	if timeout <= 0 {
		timeout = defaultPollTimeout
	}
	return context.WithTimeout(ctx, timeout)
}

func resolveCommitInterval(value time.Duration) time.Duration {
	if value <= 0 {
		return defaultCommitEvery
	}
	return value
}

func resolvePollTimeout(value time.Duration) time.Duration {
	if value <= 0 {
		return defaultPollTimeout
	}
	return value
}
