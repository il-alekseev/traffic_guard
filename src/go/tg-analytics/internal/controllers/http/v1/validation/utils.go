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

	_, err := parser.Parse(timeRangeStr, now)
	if err != nil {
		return fmt.Errorf("invalid time range: %v", err)
	}

	return nil
}

// contains проверяет наличие строки в слайсе
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
