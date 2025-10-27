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
type Service struct {
	log  *slog.Logger
	repo Read
}

func NewService(repo Read, log *slog.Logger) *Service {
	log = log.With(wsl.Label("layer", "service"))
	return &Service{
		repo: repo,
		log:  log,
	}
}
