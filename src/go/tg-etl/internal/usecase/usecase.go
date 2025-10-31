package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"tg-etl/config"
	"tg-etl/internal/models"
	"tg-etl/internal/repo/kafka"
)

type UseCase struct {
	q QueryUsecase
	// TODO: добавить функционал очистки кешей по TTL
	batchSize      uint
	lastLog        *models.IdsLog
	kc             kafka.Client
	processingLock sync.Mutex
	l              slog.Logger
}

func New(cfg *config.Config, q QueryUsecase, kc kafka.Client, l slog.Logger) (*UseCase, error) {
	uc := UseCase{
		q:         q,
		batchSize: uint(cfg.BatchSize),
		lastLog:   nil,
		kc:        kc,
		l:         l,
	}
	// Заполняем вспомогательные таблицы для ETL
	// Проверяем, пуста ли таблица с категориями, если пуста, то добавляем категории
	categories, err := q.GetCategories(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %v", err)
	}
	if len(categories) == 0 {
		err = q.CreatePredefinedCategories(context.Background(), models.PredefinedCategories)
		if err != nil {
			return nil, fmt.Errorf("failed to create categories: %v", err)
		}
	}
	return &uc, nil
}
