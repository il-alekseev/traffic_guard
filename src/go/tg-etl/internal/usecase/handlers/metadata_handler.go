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
	q usecase.QueryUsecase
	l *slog.Logger
}

// NewMetadataHandler создает новый обработчик метаданных
func NewMetadataHandler(queryUsecase usecase.QueryUsecase, l *slog.Logger) *MetadataHandler {
	return &MetadataHandler{
		q: queryUsecase,
		l: l,
	}
}

func (h *MetadataHandler) HandleURLMetadata(ctx context.Context, result models.URLMetadataResult) {
	//h.l.DebugContext(ctx, "Processing URL metadata result",
	//	slog.String("request_id", result.RequestID),
	//	slog.String("url", result.URL))
	// Обновляем поля домена полученной информацией
	// TODO: Продумать кейсы с различными данными (domain, ip, url)
	domain := models.Domain{
		IP:      result.Domain.IP,
		Country: result.Domain.Geo.CountryCode,
		Path:    result.Domain.Name,
	}
	if err := h.q.UpdateDomain(ctx, domain); err != nil {
		h.l.ErrorContext(ctx, "failed to handle URL metadata", wsl.Err(err))
	}
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
