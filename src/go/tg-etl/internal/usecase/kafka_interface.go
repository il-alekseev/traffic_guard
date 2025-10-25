package usecase

import (
	"context"
	"tg-etl/internal/models"
)

// KafkaRepository интерфейс для работы с Kafka
type KafkaRepository interface {
	SendAnalysisRequest(ctx context.Context, req models.AnalysisRequest) error
}
