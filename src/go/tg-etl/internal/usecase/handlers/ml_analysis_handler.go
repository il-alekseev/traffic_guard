package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"tg-etl/internal/models"
	"tg-etl/internal/usecase"
	"tg-etl/pkg/slogger/wsl"
)

// MLAnalysisHandler обрабатывает сообщения из ML топика
type MLAnalysisHandler struct {
	q usecase.QueryUsecase
	l *slog.Logger
}

// NewMLAnalysisHandler создает новый обработчик ML анализа
func NewMLAnalysisHandler(queryUsecase usecase.QueryUsecase, logger *slog.Logger) *MLAnalysisHandler {
	return &MLAnalysisHandler{
		q: queryUsecase,
		l: logger,
	}
}

// HandleMLAnalysis обрабатывает результат ML анализа
func (h *MLAnalysisHandler) HandleMLAnalysis(ctx context.Context, result models.MLAnalysisResult) {
	h.l.InfoContext(ctx, "Processing ML analysis result",
		slog.String("request_id", result.RequestID),
		slog.String("RecognisedClass", result.RecognisedClass),
	)
	domain, err := h.q.GetDomainByRequestID(ctx, result.RequestID)
	if err != nil {
		h.l.ErrorContext(ctx, "failed to handle ML data", wsl.Err(err))
		return
	}
	categoryID, err := h.q.GetCategoryID(ctx, result.RecognisedClass)
	if err != nil {
		h.l.ErrorContext(ctx, "failed to handle ML data", wsl.Err(err))
		return
	}
	if categoryID == 0 {
		h.l.ErrorContext(ctx, "category not recognized", wsl.String("name", result.RecognisedClass))
		return
	}
	if domain != nil {
		domain.CategoryID = int(categoryID)
		h.q.UpdateDomain(ctx, *domain)
	} else {
		h.l.ErrorContext(ctx, "domain not found", wsl.Err(err))
	}
}

// HandleRawMessage обрабатывает сырое сообщение из Kafka
func (h *MLAnalysisHandler) HandleRawMessage(ctx context.Context, message []byte) error {
	var result models.MLAnalysisResult
	if err := json.Unmarshal(message, &result); err != nil {
		h.l.ErrorContext(ctx, "failed to unmarshal ML analysis message",
			wsl.Err(err),
			slog.String("raw_message", string(message)),
		)
		return err
	}

	h.HandleMLAnalysis(ctx, result)
	return nil
}
