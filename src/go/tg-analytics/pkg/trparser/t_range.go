package trparser

import (
	"fmt"
	"time"
)

// TimeRange представляет временной диапазон для фильтрации
type TimeRange struct {
	From time.Time
	To   time.Time
}

// String возвращает строковое представление временного диапазона
func (tr TimeRange) String() string {
	return fmt.Sprintf("From: %s, To: %s", tr.From.Format(time.RFC3339), tr.To.Format(time.RFC3339))
}

// IsValid проверяет валидность временного диапазона
func (tr TimeRange) IsValid() bool {
	return !tr.From.IsZero() && !tr.To.IsZero() && tr.From.Before(tr.To)
}

// Duration возвращает продолжительность диапазона
func (tr TimeRange) Duration() time.Duration {
	return tr.To.Sub(tr.From)
}
