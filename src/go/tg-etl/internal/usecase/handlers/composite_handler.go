package handlers

import (
	"context"
	"log/slog"
	k "tg-etl/internal/repo/kafka"
	"tg-etl/pkg/slogger/wsl"

	"github.com/segmentio/kafka-go"
)

// CompositeHandler объединяет все обработчики Kafka
type CompositeHandler struct {
}

// NewCompositeHandler создает композитный обработчик
func NewCompositeHandler(
	metadataHandler *MetadataHandler,
	mlAnalysisHandler *MLAnalysisHandler,
	logger *slog.Logger,
) k.ConsumerHandlers {
	return k.ConsumerHandlers{
		URLMetadataHandler: metadataHandler.HandleURLMetadata,
		MLAnalysisHandler:  mlAnalysisHandler.HandleMLAnalysis,
		ErrorHandler:       handleError,
	}
}

// handleError обрабатывает ошибки Kafka
func handleError(ctx context.Context, err error, message kafka.Message) {
	slog.ErrorContext(ctx, "Kafka message processing error",
		wsl.Err(err),
		slog.String("topic", message.Topic),
		slog.Int("partition", message.Partition),
		slog.Int64("offset", message.Offset),
	)
}

// RawMessageHandler для обработки сырых сообщений (опционально)
type RawMessageHandler struct {
	metadataHandler   *MetadataHandler
	mlAnalysisHandler *MLAnalysisHandler
	logger            *slog.Logger
}

// NewRawMessageHandler создает обработчик сырых сообщений
func NewRawMessageHandler(
	metadataHandler *MetadataHandler,
	mlAnalysisHandler *MLAnalysisHandler,
	logger *slog.Logger,
) *RawMessageHandler {
	return &RawMessageHandler{
		metadataHandler:   metadataHandler,
		mlAnalysisHandler: mlAnalysisHandler,
		logger:            logger,
	}
}

// HandleMessage определяет тип сообщения и направляет его в соответствующий обработчик
func (h *RawMessageHandler) HandleMessage(ctx context.Context, topic string, message []byte) error {
	switch topic {
	case "url-metadata-results": // Замените на актуальный топик из конфига
		return h.metadataHandler.HandleRawMessage(ctx, message)
	case "ml-analysis-results": // Замените на актуальный топик из конфига
		return h.metadataHandler.HandleRawMessage(ctx, message)
	default:
		h.logger.WarnContext(ctx, "unknown topic for message",
			slog.String("topic", topic),
			slog.String("message", string(message)),
		)
		return nil
	}
}
