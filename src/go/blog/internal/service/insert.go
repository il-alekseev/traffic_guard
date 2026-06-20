package service

import (
	"context"
	"fiermon-blog/internal/models"
	"time"
)

// AddRecordToRepo - добавить бизнес лог в БД
func (s *Service) AddRecordToRepo(ctx context.Context, businessLog *models.BusinessLog) error {
	businessLog.Timestamp = time.Now()
	return s.repoW.InsertRecord(ctx, businessLog)
}
