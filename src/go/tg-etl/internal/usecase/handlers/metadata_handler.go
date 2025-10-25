package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"tg-etl/internal/models"
	"tg-etl/internal/usecase"
	"tg-etl/pkg/slogger/wsl"
)

// MetadataHandler обрабатывает сообщения из metadata топика
type MetadataHandler struct {
	queryUsecase usecase.QueryUsecase
	l            *slog.Logger
}

// NewMetadataHandler создает новый обработчик метаданных
func NewMetadataHandler(queryUsecase usecase.QueryUsecase, l *slog.Logger) *MetadataHandler {
	return &MetadataHandler{
		queryUsecase: queryUsecase,
		l:            l,
	}
}

func (h *MetadataHandler) HandleURLMetadata(ctx context.Context, result models.URLMetadataResult) {
	h.l.InfoContext(ctx, "Processing URL metadata result",
		slog.String("request_id", result.RequestID),
		slog.String("url", result.URL))
}

// HandleRawMessage обрабатывает сырое сообщение из Kafka
func (h *MetadataHandler) HandleRawMessage(ctx context.Context, message []byte) error {
	var result models.URLMetadataResult
	if err := json.Unmarshal(message, &result); err != nil {
		h.l.ErrorContext(ctx, "failed to unmarshal metadata message",
			wsl.Err(err),
			slog.String("raw_message", string(message)),
		)
		return err
	}

	h.HandleURLMetadata(ctx, result)
	return nil
}
