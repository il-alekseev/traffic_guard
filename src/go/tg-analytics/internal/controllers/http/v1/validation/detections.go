package validation

import (
	"fmt"
	"strings"
	"tg-an/pkg/models"
)

type GetTopDetectionsRequest struct {
	From     string `form:"from" binding:"omitempty"`
	To       string `form:"to" binding:"omitempty"`
	HostName string `form:"hostname" binding:"omitempty,max=100"`
	Category string `form:"category" binding:"omitempty,max=50"`
	Page     int    `form:"page" binding:"omitempty,min=1"`
	Limit    int    `form:"limit" binding:"omitempty,min=1"`
}

// Normalize нормализует значения запроса
func (r *GetTopDetectionsRequest) Normalize() {
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
	if r.Page == 0 {
		r.Page = 1
	}
	if r.Limit == 0 {
		r.Limit = 10
	}
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
