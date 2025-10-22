package usecase

import (
	"context"
	"log/slog"
	"tg-etl/internal/models"
	k "tg-etl/internal/repo/kafka"

	"github.com/segmentio/kafka-go"
)

// registerKafkaHandlers регистрирует обработчики для сообщений из Kafka
func (uc *UseCase) registerKafkaHandlers() {
	handlers := k.ConsumerHandlers{
		URLMetadataHandler: uc.handleURLMetadata,
		MLAnalysisHandler:  uc.handleMLAnalysis,
		ErrorHandler:       uc.handleKafkaError,
	}

	// Запускаем потребителей в отдельной горутине
	go func() {
		ctx := context.Background()
		uc.kc.StartConsumer(ctx, handlers)
		uc.kc.Wait()
	}()
}

// handleURLMetadata обрабатывает результаты метаданных URL
func (uc *UseCase) handleURLMetadata(ctx context.Context, result models.URLMetadataResult) {
	uc.l.InfoContext(ctx, "Processing URL metadata result",
		slog.String("request_id", result.RequestID),
		slog.String("url", result.URL))
}

// handleMLAnalysis обрабатывает результаты ML анализа
func (uc *UseCase) handleMLAnalysis(ctx context.Context, result models.MLAnalysisResult) {
	uc.l.InfoContext(ctx, "Processing ML analysis result",
		slog.String("request_id", result.RequestID),
		slog.String("url", result.URL),
		slog.String("category", result.Category),
		slog.Float64("confidence", result.Confidence))
}

// handleKafkaError обрабатывает ошибки из Kafka
func (uc *UseCase) handleKafkaError(ctx context.Context, err error, message kafka.Message) {
	uc.l.ErrorContext(ctx, "Kafka message processing error",
		slog.String("topic", message.Topic),
		slog.String("error", err.Error()))
}
