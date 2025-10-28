package validation

import (
	"fmt"
	"tg-an/pkg/trparser"
	"time"
)

// validateTimeRange проверяет временной диапазон с использованием trparser
func validateTimeRange(from, to string) error {
	// Используем trparser для парсинга временного диапазона
	parser := &trparser.TimeRangeParser{}
	now := time.Now()

	// Формируем строку для парсера в формате "from=X&to=Y"
	timeRangeStr := fmt.Sprintf("from=%s&to=%s", from, to)

	timeRange, err := parser.Parse(timeRangeStr, now)
	if err != nil {
		return fmt.Errorf("invalid time range: %v", err)
	}

	// Дополнительная проверка на слишком большой диапазон
	if timeRange.Duration() > 30*24*time.Hour {
		return fmt.Errorf("time range cannot exceed 30 days")
	}

	return nil
}
