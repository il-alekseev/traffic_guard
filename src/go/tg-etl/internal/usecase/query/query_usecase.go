package usecase

import (
	"log/slog"
	"tg-etl/config"
	"tg-etl/internal/models"
	"tg-etl/internal/repo/m_cache"
	"time"
)

type QueryUseCase struct {
	MLAttemps uint
	ksuDB     KSURepoPGInterface
	etlDB     ETLRepoPGInterface
	// TODO: добавить функционал очистки кешей по TTL
	c       m_cache.MemoryCache
	lastLog *models.IdsLog
	l       slog.Logger
}

func New(cfg *config.Config, ksuDB KSURepoPGInterface, etlDB ETLRepoPGInterface, l slog.Logger) *QueryUseCase {
	uc := QueryUseCase{
		MLAttemps: cfg.MLAttemps,
		ksuDB:     ksuDB,
		etlDB:     etlDB,
		c:         *m_cache.New(time.Duration(cfg.TTL) * time.Minute),
		lastLog:   nil,
		l:         l,
	}
	return &uc
}
