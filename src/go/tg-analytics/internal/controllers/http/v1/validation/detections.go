package validation

import (
	"fmt"
	"strings"
	"tg-an/pkg/models"
)

type GetTopDetectionsRequest struct {
	From     string `form:"from" binding:"omitempty"`
	To       string `form:"to" binding:"omitempty"`
	Action   string `form:"action" binding:"omitempty"`
	HostName string `form:"hostname" binding:"omitempty,max=100"`
	Category string `form:"category" binding:"omitempty,max=50"`
	Page     int    `form:"page" binding:"omitempty,min=1"`
	Limit    int    `form:"limit" binding:"omitempty,min=1"`
	Search   string `form:"search" binding:"omitempty,max=200"`
	OrderBy  string `form:"order_by" binding:"omitempty"`
	OrderDir string `form:"order_dir" binding:"omitempty,oneof=asc desc"`
}

// Normalize нормализует значения запроса
func (r *GetTopDetectionsRequest) Normalize() {
	// Тримим строковые поля
	r.HostName = strings.TrimSpace(r.HostName)
	r.Category = strings.TrimSpace(r.Category)
	r.Search = strings.TrimSpace(r.Search)
	r.OrderBy = strings.TrimSpace(r.OrderBy)
	r.OrderDir = strings.TrimSpace(r.OrderDir)

	// Устанавливаем значения по умолчанию
	if r.From == "" {
		r.From = "now-10m"
	}
	if r.To == "" {
		r.To = "now"
	}
	if r.Page == 0 {
		r.Page = 1
	}
	if r.Limit == 0 {
		r.Limit = 10
	}

	// Приводим order_dir к нижнему регистру
	r.OrderDir = strings.ToLower(r.OrderDir)
}

// Validate выполняет валидацию всех полей запроса
func (r *GetTopDetectionsRequest) Validate() error {
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

	// Валидация категории
	allowedCategories := []string{}
	for _, c := range models.PredefinedCategories {
		allowedCategories = append(allowedCategories, c.Name)
	}

	if r.Category != "" && !contains(allowedCategories, r.Category) {
		return fmt.Errorf("invalid category")
	}
	// Валидация типа действия
	allowedActions := []string{}
	for _, a := range models.AllActions {
		allowedActions = append(allowedActions, a.String())
	}

	if r.Action != "" && !contains(allowedActions, r.Action) {
		return fmt.Errorf("invalid action")
	}

	// Валидация поиска
	if r.Search != "" && len(r.Search) < 2 {
		return fmt.Errorf("search query must be at least 2 characters long")
	}
	// Валидация сортировки
	allowedOrderFields := []string{
		"",
		"domain",
		"request_count",
		"categorized_at",
	}
	if !contains(allowedOrderFields, r.OrderBy) {
		return fmt.Errorf("invalid order_by field")
	}
	// Валидация направления сортировки
	allowedOrderDirs := []string{"desc", "asc", ""}
	if !contains(allowedOrderDirs, r.OrderDir) {
		return fmt.Errorf("invalid order_by field")
	}

	return nil
}

// ValidateAndNormalize выполняет нормализацию и валидацию
func (r *GetTopDetectionsRequest) ValidateAndNormalize() error {
	r.Normalize()
	return r.Validate()
}

type GetDetectionStatRequest struct {
	From     string `form:"from" binding:"omitempty"`
	To       string `form:"to" binding:"omitempty"`
	HostName string `form:"hostname" binding:"omitempty,max=100"`
	Category string `form:"category" binding:"omitempty,max=50"`
}

// Normalize нормализует значения запроса
func (r *GetDetectionStatRequest) Normalize() {
	// Тримим строковые поля
	r.HostName = strings.TrimSpace(r.HostName)
	r.Category = strings.TrimSpace(r.Category)

	// Устанавливаем значения по умолчанию
	if r.From == "" {
		r.From = "now-10m"
	}
	if r.To == "" {
		r.To = "now"
	}
}

// Validate выполняет валидацию всех полей запроса
func (r *GetDetectionStatRequest) Validate() error {
	// Валидация временного диапазона
	if err := validateTimeRange(r.From, r.To); err != nil {
		return err
	}

	// Валидация категории
	allowedCategories := []string{}
	for _, c := range models.PredefinedCategories {
		allowedCategories = append(allowedCategories, c.Name)
	}

	if r.Category != "" && !contains(allowedCategories, r.Category) {
		return fmt.Errorf("invalid category")
	}

	return nil
}

func (r *GetDetectionStatRequest) ValidateAndNormalize() error {
	r.Normalize()
	return r.Validate()
}
