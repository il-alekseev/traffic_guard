package service

import (
	"context"
	"fiermon-blog/internal/models"
	"fiermon-blog/pkg/fslog/wsl"
	"log/slog"
)

// Read - реализует методы базы данных
type Read interface {
	GetRecord(ctx context.Context, meta *models.UserMeta, page, limit int, role, contextID, search string) ([]models.BusinessLog, models.Meta, error)
}
type Write interface {
	InsertRecord(ctx context.Context, businessLog *models.BusinessLog) error
}
type Service struct {
	log   *slog.Logger
	repoR Read
	repoW Write
}

func NewService(repoR Read, repoW Write, log *slog.Logger) *Service {
	log = log.With(wsl.Label("layer", "service"))
	return &Service{
		repoR: repoR,
		repoW: repoW,
		log:   log,
	}
}
