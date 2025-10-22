package usecase

import (
	"log/slog"
	"tg-etl/config"
	"tg-etl/internal/models"
	"tg-etl/internal/repo/m_cache"
	"time"
)

type UseCase struct {
	ksuDB KSURepoPGInterface
	etlDB ETLRepoPGInterface
	// TODO: добавить функционал очистки кешей по TTL
	c        m_cache.MemoryCache
	maxCount uint
	lastLog  *models.IdsLog
	l        slog.Logger
}

func New(cfg *config.Config, ksuDB KSURepoPGInterface, etlDB ETLRepoPGInterface, l slog.Logger) *UseCase {
	uc := UseCase{
		ksuDB:    ksuDB,
		etlDB:    etlDB,
		c:        *m_cache.New(time.Duration(cfg.TTL) * time.Minute),
		maxCount: uint(cfg.MaxCount),
		lastLog:  nil,
		l:        l,
	}
	return &uc
}
