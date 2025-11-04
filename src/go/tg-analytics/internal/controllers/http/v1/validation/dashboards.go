package validation

import (
	"fmt"
	"strings"
)

type GetTopCategoriesRequest struct {
	From     string `form:"from" binding:"omitempty"`
	To       string `form:"to" binding:"omitempty"`
	HostName string `form:"hostname" binding:"omitempty,max=100"`
	Type     string `form:"type" binding:"omitempty,max=20"`
	Count    int    `form:"count" binding:"omitempty,min=1"`
}

// Validate выполняет валидацию всех полей запроса
func (r *GetTopCategoriesRequest) Validate() error {
	// Валидация количества категорий
	if r.Count < 1 {
		return fmt.Errorf("count must be greater than or equal to 1")
	}

	// Валидация временного диапазона
	if err := validateTimeRange(r.From, r.To); err != nil {
		return err
	}

	// Валидация типа сессии
	allowedTypes := []string{"Заблокирован", "Запрещен", "Ожидает", "Разрешен", ""}
	if r.Type != "" && !contains(allowedTypes, r.Type) {
		return fmt.Errorf("invalid session type")
	}

	return nil
}

// Normalize нормализует значения запроса
func (r *GetTopCategoriesRequest) Normalize() {
	// Тримим строковые поля
	r.HostName = strings.TrimSpace(r.HostName)
	r.Type = strings.TrimSpace(r.Type)

	// Устанавливаем значения по умолчанию
	if r.From == "" {
		r.From = "now-24h"
	}
	if r.To == "" {
		r.To = "now"
	}
	if r.Count == 0 {
		r.Count = 5
	}
}

// ValidateAndNormalize выполняет нормализацию и валидацию
func (r *GetTopCategoriesRequest) ValidateAndNormalize() error {
	r.Normalize()
	return r.Validate()
}

// GetTrafficStatRequest представляет запрос для получения статистики трафика
type GetTrafficStatRequest struct {
	From     string `form:"from" binding:"omitempty"`
	To       string `form:"to" binding:"omitempty"`
	Count    uint   `form:"count" binding:"omitempty,min=1"`
	HostName string `form:"hostname" binding:"omitempty,max=100"`
}

// Normalize нормализует значения запроса
func (r *GetTrafficStatRequest) Normalize() {
	// Тримим строковые поля
	r.From = strings.TrimSpace(r.From)
	r.To = strings.TrimSpace(r.To)

	// Устанавливаем значения по умолчанию
	if r.From == "" {
		r.From = "now-10m"
	}
	if r.To == "" {
		r.To = "now"
	}
	if r.Count == 0 {
		r.Count = 20
	}

}

// Validate выполняет валидацию всех полей запроса
func (r *GetTrafficStatRequest) Validate() error {
	// Валидация количества точек
	if r.Count < 1 {
		return fmt.Errorf("count must be greater than or equal to 1")
	}

	// Валидация временного диапазона
	if err := validateTimeRange(r.From, r.To); err != nil {
		return err
	}

	return nil
}

// ValidateAndNormalize выполняет нормализацию и валидацию
func (r *GetTrafficStatRequest) ValidateAndNormalize() error {
	r.Normalize()
	return r.Validate()
}
