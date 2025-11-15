package service

import (
	"context"
	"fiermon-blog/internal/models"
)

// GetRecords - получение логов из базы данных
func (s *Service) GetRecords(ctx context.Context, meta *models.UserMeta, page, limit int, role, contextID, search string) ([]models.BusinessLog, models.Meta, error) {
	return s.repoR.GetRecord(
		ctx,
		meta,
		page,
		limit,
		role,
		contextID,
		search,
	)
}
