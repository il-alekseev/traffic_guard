package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"tg-etl/config"
	"tg-etl/internal/models"
	"tg-etl/internal/repo/kafka"
	"tg-etl/pkg/blog"
	"tg-etl/pkg/usercontrol"
)

type UseCase struct {
	q QueryUsecase
	// TODO: добавить функционал очистки кешей по TTL
	batchSize      uint
	lastLog        *models.LastLog
	kc             kafka.Client
	processingLock sync.Mutex
	l              slog.Logger
	userCl         *usercontrol.Usercontrol
	blclient       *blog.Blog
}

func New(cfg *config.Config, q QueryUsecase, kc kafka.Client, l slog.Logger, userCl *usercontrol.Usercontrol, blclient *blog.Blog) (*UseCase, error) {
	ctx := context.Background()
	// Заполняем вспомогательные таблицы для ETL
	// Проверяем, пуста ли таблица с категориями, если пуста, то добавляем категории
	categories, err := q.GetCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %v", err)
	}
	if len(categories) == 0 {
		err = q.CreatePredefinedCategories(context.Background(), models.PredefinedCategories)
		if err != nil {
			return nil, fmt.Errorf("failed to create categories: %v", err)
		}
	}
	// Загружаем запись о последнем обработанном логе (в случае, если сервис останавливался)
	lastLog, err := q.GetLastLog(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get last log: %v", err)
	}
	uc := UseCase{
		q:         q,
		batchSize: uint(cfg.BatchSize),
		lastLog:   lastLog,
		kc:        kc,
		l:         l,
		userCl:    userCl,
		blclient:  blclient,
	}
	return &uc, nil
}
