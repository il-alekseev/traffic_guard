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
// TODO:  вынести настройки списков в отдельную структуру-конфиг
type MLAnalysisHandler struct {
	mlAttemps   uint
	maxNegCount uint
	q           usecase.QueryUsecase
	l           *slog.Logger
}

// NewMLAnalysisHandler создает новый обработчик ML анализа
func NewMLAnalysisHandler(mlAttemps uint, maxNegCount uint, queryUsecase usecase.QueryUsecase, logger *slog.Logger) *MLAnalysisHandler {
	return &MLAnalysisHandler{
		mlAttemps:   mlAttemps,
		maxNegCount: maxNegCount,
		q:           queryUsecase,
		l:           logger,
	}
}

// HandleMLAnalysis обрабатывает результат ML анализа
// TODO: Оптимизировать функцию
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
		h.l.WarnContext(ctx, "failed to handle URL metadata", wsl.Err(err))
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
		h.l.ErrorContext(ctx, "failed to get domain by request_id", wsl.Err(err))
		return
	}
	if domain == nil {
		h.l.WarnContext(ctx, "domain not found", wsl.String("request_id", requestID.String()))
		return
	}
	h.l.DebugContext(ctx, "reseived message from ML: ",
		wsl.String("domain", domain.Path),
		wsl.String("category", result.RecognisedClass),
	)
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
	// Добавляем или увеличиваем на 1 число встречаемости категории в таблице DomainCategory
	if err := h.q.AddDomainCategory(ctx, *category, domain.ID); err != nil {
		h.l.ErrorContext(ctx, "failed to add domain category in db", wsl.Err(err))
		return
	}

	// Обновляем запись о домене
	if category.Type == models.CategoryTypePositive {
		// Если категория положительная, инкрементируем число проверок
		domain.AnalysisCount++

		// Если достигли или превысили порог, вносим в белый список
		if domain.AnalysisCount >= h.mlAttemps {
			if err := h.q.AddDomainToList(ctx, *domain, models.Whitelist); err != nil {
				h.l.ErrorContext(ctx, "failed to add domain to whitelist", wsl.Err(err))
			} else {
				h.l.DebugContext(ctx, "domain added to whitelist",
					wsl.String("domain", domain.Path))
			}
		}
	}

	// Обновляем наиболее негативную категорию домена
	newCat, perc, err := h.q.GetMostNegativeCategoryByDomainID(ctx, domain.ID)
	if err != nil {
		h.l.ErrorContext(ctx, "failed to get most negative category", wsl.Err(err))
		// Не возвращаемся, продолжаем с текущими значениями
	}
	// Если домен ни разу не определялся как негативный и категория новая, изменяем его положительную категорию
	if newCat == nil && domain.CategoryID != int(category.ID) {
		domain.CategoryID = int(category.ID)
		domain.CategorizedAt = time.Now()
	}
	// Если домен хотя бы раз определялся как негативный, то обновляем его рейтинг негативных категорий и категорию тоже
	if newCat != nil {
		// Обновляем вероятность
		domain.NegRate = perc
		// Если ID категорий не совпадает, обновляем
		if domain.CategoryID != int(newCat.ID) {
			domain.CategoryID = int(newCat.ID)
			domain.CategorizedAt = time.Now()
		}
	} else {
		h.l.WarnContext(ctx, "no negative categories found for domain", wsl.Int("domain_id", int(domain.ID)))
	}

	// Проверяем, нужно ли включать домен в черный список
	if category.Type == models.CategoryTypeNegative {
		negCount, err := h.q.GetNegativeCategoriesTotalByDomainID(ctx, domain.ID)
		if err != nil {
			h.l.ErrorContext(ctx, "failed to get negative categories total", wsl.Err(err))
		} else if negCount >= int(h.maxNegCount) {
			if err := h.q.AddDomainToList(ctx, *domain, models.Blacklist); err != nil {
				h.l.ErrorContext(ctx, "failed to add domain to blacklist", wsl.Err(err))
			} else {
				h.l.DebugContext(ctx, "domain added to blacklist",
					wsl.String("domain", domain.Path),
					wsl.Int("neg_count", negCount))
			}
		}

		h.l.DebugContext(ctx, "got domain category",
			wsl.String("domain", domain.Path),
			wsl.String("category", category.Name),
		)
	}

	// Обновляем сам домен
	if err := h.q.UpdateDomain(ctx, *domain); err != nil {
		h.l.WarnContext(ctx, "failed to update domain", wsl.Err(err))
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
