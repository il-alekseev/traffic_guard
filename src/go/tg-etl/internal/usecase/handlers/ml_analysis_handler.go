package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"tg-etl/internal/models"
	"tg-etl/internal/usecase"
	"tg-etl/pkg/slogger/wsl"
	"time"

	"github.com/google/uuid"
)

// MLAnalysisHandler обрабатывает сообщения из ML топика
type MLAnalysisHandler struct {
	mlAttemps uint
	q         usecase.QueryUsecase
	l         *slog.Logger
}

// NewMLAnalysisHandler создает новый обработчик ML анализа
func NewMLAnalysisHandler(mlAttemps uint, queryUsecase usecase.QueryUsecase, logger *slog.Logger) *MLAnalysisHandler {
	return &MLAnalysisHandler{
		mlAttemps: mlAttemps,
		q:         queryUsecase,
		l:         logger,
	}
}

// HandleMLAnalysis обрабатывает результат ML анализа
func (h *MLAnalysisHandler) HandleMLAnalysis(ctx context.Context, result models.MLAnalysisResult) {
	// парсим RequestID
	requestID, err := uuid.Parse(result.RequestID)
	if err != nil {
		h.l.ErrorContext(ctx, "failed to parse request_id to uuid", wsl.Err(err))
		return
	}
	// Обновляем время получения метаданных для URL по requestID
	url := models.URL{
		GetCategoryAt: time.Now(),
	}
	if err := h.q.UpdateURLByRequestID(ctx, requestID, url); err != nil {
		h.l.ErrorContext(ctx, "failed to handle URL metadata", wsl.Err(err))
	}

	// Находим категорию
	category, err := h.q.GetCategoryByName(ctx, result.RecognisedClass)
	if err != nil {
		h.l.ErrorContext(ctx, "failed to handle ML data", wsl.Err(err))
		return
	}
	// Если вернулся nil, значит есть несовпадение между категориями ETL и ML
	if category == nil {
		h.l.ErrorContext(ctx, "category not recognized", wsl.String("name", result.RecognisedClass))
		return
	}
	// Находим домен по RequestID
	domain, err := h.q.GetDomainByRequestID(ctx, requestID)
	if err != nil {
		h.l.ErrorContext(ctx, "failed to handle ML data", wsl.Err(err))
		return
	}
	if domain == nil {
		h.l.WarnContext(ctx, "domain not found", wsl.String("request_id", requestID.String()))
		return
	}
	// Обрабатываем случай, когда нашлась и категория, и домен
	// Проверяем, находится ли домен в одном из списков
	list, err := h.q.GetListByDomainID(ctx, domain.ID)
	if err != nil {
		h.l.ErrorContext(ctx, "failed to get list for domain", wsl.Err(err))
		return
	}
	if list != nil {
		h.l.DebugContext(ctx, "ML handler: domain already in list",
			wsl.String("list", *list), wsl.String("domain", domain.Path))
		return
	}
	if category.Type != models.CategoryTypeNeutral {
		h.l.InfoContext(ctx, "Processing ML analysis result",
			slog.String("request_id", result.RequestID),
			slog.String("domain", domain.Path),
			slog.String("category", category.Name),
		)
	}
	// Домен не лежит в списках
	// Если категория контента негативная - вносим его в черный список
	if category.Type == models.CategoryTypeNegative {
		if err := h.q.AddDomainToList(ctx, *domain, models.Blacklist); err != nil {
			h.l.ErrorContext(ctx, "AddDomainToList", wsl.Err(err))
		}
		h.l.DebugContext(ctx, "got domain category",
			wsl.String("domain", domain.Path),
			wsl.String("category", category.Name),
		)
		// Обновляем категорию домена
		domain.CategoryID = int(category.ID)
		domain.CategorizedAt = time.Now()
	} else if category.Type == models.CategoryTypePositive {
		// Если категория положительная, инкрементируем число проверок домена для получения положительного статуса на 1
		// проверяем, не достигло ли число успешых проверок на положительный контент константе
		domain.AnalysisCount++
		// Если достигло, то вносим домен в белый список
		if domain.AnalysisCount == h.mlAttemps {
			if err := h.q.AddDomainToList(ctx, *domain, models.Whitelist); err != nil {
				h.l.ErrorContext(ctx, "AddDomainToList", wsl.Err(err))
			}
			h.l.DebugContext(ctx, "domain is positive",
				wsl.String("domain", domain.Path),
			)
		}
	}
	if err := h.q.UpdateDomain(ctx, *domain); err != nil {
		h.l.ErrorContext(ctx, "UpdateDomain", wsl.Err(err))
		return
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
