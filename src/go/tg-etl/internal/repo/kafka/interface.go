package kafka

import (
	"context"
	"tg-etl/internal/models"

	"github.com/segmentio/kafka-go"
)

// Client интерфейс для работы с Kafka
type Client interface {
	// SendAnalysisRequest отправляет запрос на анализ в URL topic
	SendAnalysisRequest(ctx context.Context, req models.AnalysisRequest) error

	// StartConsumer запускает потребителей для чтения результатов из metadata и ML топиков
	StartConsumer(ctx context.Context, handlers ConsumerHandlers)

	// Wait блокирует до завершения всех потребителей
	Wait()

	// Close закрывает все соединения с Kafka
	Close(ctx context.Context) error
}

// ConsumerHandlers обработчики для разных типов сообщений
type ConsumerHandlers struct {
	URLMetadataHandler func(ctx context.Context, result models.URLMetadataResult)
	MLAnalysisHandler  func(ctx context.Context, result models.MLAnalysisResult)
	ErrorHandler       func(ctx context.Context, err error, message kafka.Message)
}
