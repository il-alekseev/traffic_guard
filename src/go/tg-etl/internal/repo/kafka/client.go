package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"sync"
	"tg-etl/config"
	"tg-etl/internal/models"
	"tg-etl/pkg/slogger/wsl"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaClient struct {
	producer  *kafka.Writer
	consumers []*kafka.Reader
	l         *slog.Logger
	cfg       config.Kafka
	wg        sync.WaitGroup
}

// New создает новый клиент Kafka
func New(ctx context.Context, cfg config.Kafka, logger *slog.Logger) (*KafkaClient, error) {
	// Формируем адрес брокера
	brokerAddr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	brokers := []string{brokerAddr}

	// Валидация конфигурации
	if cfg.URLTopic == "" {
		return nil, fmt.Errorf("URL topic is required")
	}

	if cfg.MetadataTopic == "" {
		return nil, fmt.Errorf("metadata topic is required")
	}

	if cfg.MLTopic == "" {
		return nil, fmt.Errorf("ML topic is required")
	}

	// Проверяем, что сервер Kafka доступен
	err := checkKafkaConnection(ctx, brokers[0])
	if err != nil {
		logger.ErrorContext(ctx, "new kafka client", wsl.Err(err))
		return nil, err
	}

	// Простой кастомный Dialer для чтения, который использует IP напрямую
	// В противном случае, клиент пытается найти сервер Kafka по DNS
	// (Особенности работы библиотеки kafka-go)
	dialer := &kafka.Dialer{
		Timeout:   30 * time.Second,
		DualStack: true,
		DialFunc: func(ctx context.Context, network, addr string) (net.Conn, error) {
			var d net.Dialer
			conn, err := d.DialContext(ctx, network, brokerAddr)
			if err != nil {
				return nil, fmt.Errorf("failed to connect to %s (requested %s): %w", brokerAddr, addr, err)
			}
			return conn, nil
		},
	}

	// Producer с кастомным Transport для записи, который использует IP напрямую
	// В противном случае, клиент пытается найти сервер Kafka по DNS
	// (Особенности работы библиотеки kafka-go)
	transport := &kafka.Transport{
		Dial: func(ctx context.Context, network, addr string) (net.Conn, error) {
			// Игнорируем переданный addr и используем наш IP напрямую
			var d net.Dialer
			return d.DialContext(ctx, network, brokerAddr)
		},
	}

	// Producer для отправки запросов в URL topic
	producer := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        cfg.URLTopic,
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
		RequiredAcks: kafka.RequireOne,
		//Logger:       kafka.LoggerFunc(logger.Info),
		Transport: transport,
		ErrorLogger: kafka.LoggerFunc(func(s string, i ...interface{}) {
			logger.ErrorContext(ctx, "Kafka writer error", "msg", s, "args", i)
		}),
	}

	// Consumers для чтения из metadata и ML топиков
	readTopics := []string{cfg.MetadataTopic, cfg.MLTopic}
	consumers := make([]*kafka.Reader, 0, len(readTopics))

	for _, topic := range readTopics {
		consumer := kafka.NewReader(kafka.ReaderConfig{
			Brokers:        brokers,
			Topic:          topic,
			MinBytes:       10e3, // 10KB
			MaxBytes:       10e6, // 10MB
			MaxWait:        time.Second,
			GroupID:        "tg-etl-consumer-group",
			StartOffset:    kafka.FirstOffset,
			CommitInterval: 1 * time.Second,
			Dialer:         dialer,
		})
		consumers = append(consumers, consumer)
	}

	client := &KafkaClient{
		producer:  producer,
		consumers: consumers,
		l:         logger,
		cfg:       cfg,
	}

	logger.DebugContext(ctx, "Kafka client initialized",
		"broker", brokerAddr,
		"url_topic", cfg.URLTopic,
		"metadata_topic", cfg.MetadataTopic,
		"ml_topic", cfg.MLTopic)
	return client, nil
}

func checkKafkaConnection(ctx context.Context, broker string) error {
	conn, err := kafka.DialContext(ctx, "tcp", broker)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.ApiVersions()
	return err
}

// SendAnalysisRequest отправляет запрос на анализ в URL topic
func (kc *KafkaClient) SendAnalysisRequest(ctx context.Context, req models.AnalysisRequest) error {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal analysis request: %w", err)
	}

	message := kafka.Message{
		Key:   []byte(req.RequestID.String()), // Используем RequestID как ключ
		Value: jsonData,
		Time:  time.Now(),
	}

	err = kc.producer.WriteMessages(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to write message to topic %s: %w", kc.cfg.URLTopic, err)
	}

	//kc.l.DebugContext(ctx, "Analysis request sent",
	//	"request_id", req.RequestID,
	//	"dst_type", req.Dst.Type,
	//	"dst_resource", req.Dst.Resource,
	//	"src_ip", req.Src.IP,
	//	"topic", kc.cfg.URLTopic)

	return nil
}

// StartConsumer запускает потребителей для чтения результатов из metadata и ML топиков
func (kc *KafkaClient) StartConsumer(ctx context.Context, handlers ConsumerHandlers) {
	kc.l.DebugContext(ctx, "Starting Kafka consumers",
		"metadata_topic", kc.cfg.MetadataTopic,
		"ml_topic", kc.cfg.MLTopic)

	for i, consumer := range kc.consumers {
		kc.wg.Add(1)
		go kc.consumeTopic(ctx, consumer, handlers, i)
	}
}

// consumeTopic читает сообщения из конкретного топика
func (kc *KafkaClient) consumeTopic(ctx context.Context, consumer *kafka.Reader, handlers ConsumerHandlers, consumerID int) {
	defer kc.wg.Done()

	topic := consumer.Config().Topic
	kc.l.DebugContext(ctx, "Starting consumer", "topic", topic, "consumer_id", consumerID)

	for {
		select {
		case <-ctx.Done():
			kc.l.DebugContext(ctx, "Stopping consumer", "topic", topic, "consumer_id", consumerID)
			return
		default:
			msg, err := consumer.ReadMessage(ctx)
			if err != nil {
				kc.l.ErrorContext(ctx, "Failed to read message",
					"topic", topic,
					"consumer_id", consumerID,
					"error", err)
				continue
			}

			kc.handleMessage(ctx, msg, handlers, consumerID)
		}
	}
}

func (kc *KafkaClient) handleMessage(ctx context.Context, msg kafka.Message, handlers ConsumerHandlers, consumerID int) {
	//kc.l.DebugContext(ctx, "Received message",
	//	"topic", msg.Topic,
	//	"partition", msg.Partition,
	//	"offset", msg.Offset,
	//	"consumer_id", consumerID)

	switch msg.Topic {
	case kc.cfg.MetadataTopic:
		var result models.URLMetadataResult
		if err := json.Unmarshal(msg.Value, &result); err != nil {
			kc.l.ErrorContext(ctx, "Failed to unmarshal URL metadata result",
				"error", err,
				"consumer_id", consumerID,
				"topic", kc.cfg.MetadataTopic)
			if handlers.ErrorHandler != nil {
				handlers.ErrorHandler(ctx, err, msg)
			}
			return
		}

		if handlers.URLMetadataHandler != nil {
			handlers.URLMetadataHandler(ctx, result)
		} else {
			kc.l.WarnContext(ctx, "No handler for URL metadata result",
				"consumer_id", consumerID,
				"topic", kc.cfg.MetadataTopic)
		}

	case kc.cfg.MLTopic:
		var result models.MLAnalysisResult
		if err := json.Unmarshal(msg.Value, &result); err != nil {
			kc.l.ErrorContext(ctx, "Failed to unmarshal ML analysis result",
				"error", err,
				"consumer_id", consumerID,
				"topic", kc.cfg.MLTopic)
			if handlers.ErrorHandler != nil {
				handlers.ErrorHandler(ctx, err, msg)
			}
			return
		}

		if handlers.MLAnalysisHandler != nil {
			handlers.MLAnalysisHandler(ctx, result)
		} else {
			kc.l.WarnContext(ctx, "No handler for ML analysis result",
				"consumer_id", consumerID,
				"topic", kc.cfg.MLTopic)
		}

	default:
		kc.l.WarnContext(ctx, "Unknown topic",
			"topic", msg.Topic,
			"consumer_id", consumerID)
	}
}

// Wait блокирует до завершения всех потребителей
func (kc *KafkaClient) Wait() {
	kc.wg.Wait()
}

// Close закрывает все соединения с Kafka
func (kc *KafkaClient) Close(ctx context.Context) error {
	var errs []error

	if err := kc.producer.Close(); err != nil {
		errs = append(errs, fmt.Errorf("failed to close producer: %w", err))
	}

	for i, consumer := range kc.consumers {
		if err := consumer.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close consumer %d: %w", i, err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing Kafka client: %v", errs)
	}

	kc.l.DebugContext(ctx, "Kafka client closed successfully")
	return nil
}
