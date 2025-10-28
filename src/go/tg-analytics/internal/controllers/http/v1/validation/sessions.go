package validation

import (
	"fmt"
	"strings"
)

type GetSessionsRequest struct {
	From     string `form:"from" binding:"omitempty"`
	To       string `form:"to" binding:"omitempty"`
	HostName string `form:"hostname" binding:"omitempty,max=100"`
	Category string `form:"category" binding:"omitempty,max=50"`
	Type     string `form:"type" binding:"omitempty,max=20"`
	Search   string `form:"search" binding:"omitempty,max=200"`
	Page     int    `form:"page" binding:"omitempty,min=1"`
	Limit    int    `form:"limit" binding:"omitempty,min=1"`
	OrderBy  string `form:"order_by" binding:"omitempty"`
	OrderDir string `form:"order_dir" binding:"omitempty,oneof=asc desc"`
}

// Normalize нормализует значения запроса
func (r *GetSessionsRequest) Normalize() {
	// Тримим строковые поля
	r.HostName = strings.TrimSpace(r.HostName)
	r.Category = strings.TrimSpace(r.Category)
	r.Type = strings.TrimSpace(r.Type)
	r.Search = strings.TrimSpace(r.Search)
	r.OrderBy = strings.TrimSpace(r.OrderBy)
	r.OrderDir = strings.TrimSpace(r.OrderDir)

	// Устанавливаем значения по умолчанию
	if r.Page == 0 {
		r.Page = 1
	}
	if r.Limit == 0 {
		r.Limit = 10
	}
	if r.OrderBy == "" {
		r.OrderBy = "datetime_utc"
	}
	if r.OrderDir == "" {
		r.OrderDir = "desc"
	}
	if r.From == "" {
		r.From = "now-10m"
	}
	if r.To == "" {
		r.To = "now"
	}

	// Приводим order_dir к нижнему регистру
	r.OrderDir = strings.ToLower(r.OrderDir)
}

// Validate выполняет валидацию всех полей запроса
func (r *GetSessionsRequest) Validate() error {
	// Валидация пагинации
	if r.Page < 1 {
		return fmt.Errorf("page must be greater than or equal to 1")
	}
	if r.Limit < 1 {
		return fmt.Errorf("limit must be greater than or equal to 1")
	}

	// Валидация временного диапазона
	if err := validateTimeRange(r.From, r.To); err != nil {
		return err
	}

	// Валидация поиска
	if r.Search != "" && len(r.Search) < 2 {
		return fmt.Errorf("search query must be at least 2 characters long")
	}

	// Валидация сортировки
	allowedOrderFields := []string{"id", "datetime_utc", "type", "status", "url", "proto", "host_name", "src_ip", "src_country", "username", "dst_ip", "dst_port", "dst_country", "category", ""}
	if !contains(allowedOrderFields, r.OrderBy) {
		return fmt.Errorf("invalid order_by field")
	}

	// Валидация типа сессии
	allowedTypes := []string{"Заблокирован", "Запрещен", "Ожидает", "Разрешен", ""}
	if r.Type != "" && !contains(allowedTypes, r.Type) {
		return fmt.Errorf("invalid session type")
	}

	return nil
}

// ValidateAndNormalize выполняет нормализацию и валидацию
func (r *GetSessionsRequest) ValidateAndNormalize() error {
	r.Normalize()
	return r.Validate()
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
