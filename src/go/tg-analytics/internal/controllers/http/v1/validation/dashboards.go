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
	allowedTypes := []string{"Разрешен", "Запрещен", "VPN", ""}
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

// GetRequestStatRequest представляет запрос для получения статистики запросов
type GetRequestStatRequest struct {
	From        string `form:"from" binding:"omitempty"`
	To          string `form:"to" binding:"omitempty"`
	Count       uint   `form:"count" binding:"omitempty,min=1"`
	HostName    string `form:"hostname" binding:"omitempty,max=100"`
	RequestType string `form:"request_type" binding:"omitempty,max=100"`
}

// Normalize нормализует значения запроса
func (r *GetRequestStatRequest) Normalize() {
	// Тримим строковые поля
	r.From = strings.TrimSpace(r.From)
	r.To = strings.TrimSpace(r.To)
	r.HostName = strings.TrimSpace(r.HostName)
	r.RequestType = strings.TrimSpace(r.RequestType)

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
func (r *GetRequestStatRequest) Validate() error {
	// Валидация количества точек
	if r.Count < 1 {
		return fmt.Errorf("count must be greater than or equal to 1")
	}

	// Валидация временного диапазона
	if err := validateTimeRange(r.From, r.To); err != nil {
		return err
	}

	// Валидация типа запроса (если нужно добавить конкретные допустимые значения)
	allowedRequestTypes := []string{"allowed", "blocked", "before_block", "pending"}
	if r.RequestType != "" && !contains(allowedRequestTypes, r.RequestType) {
		return fmt.Errorf("invalid request type")
	}

	return nil
}

// ValidateAndNormalize выполняет нормализацию и валидацию
func (r *GetRequestStatRequest) ValidateAndNormalize() error {
	r.Normalize()
	return r.Validate()
}

// GetDeviceStatRequest представляет запрос для получения состояния устройств
type GetDeviceStatRequest struct {
	From  string `form:"from" binding:"omitempty"`
	To    string `form:"to" binding:"omitempty"`
	Count uint   `form:"count" binding:"omitempty,min=1"`
}

// Normalize нормализует значения запроса
func (r *GetDeviceStatRequest) Normalize() {
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
func (r *GetDeviceStatRequest) Validate() error {
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
func (r *GetDeviceStatRequest) ValidateAndNormalize() error {
	r.Normalize()
	return r.Validate()
}

// GetAnomaliesRequest представляет запрос для получения состояния устройств
type GetAnomaliesRequest struct {
	From     string `form:"from" binding:"omitempty"`
	To       string `form:"to" binding:"omitempty"`
	HostName string `form:"hostname" binding:"omitempty,max=100"`
}

// Normalize нормализует значения запроса
func (r *GetAnomaliesRequest) Normalize() {
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
}

// Validate выполняет валидацию всех полей запроса
func (r *GetAnomaliesRequest) Validate() error {
	// Валидация временного диапазона
	if err := validateTimeRange(r.From, r.To); err != nil {
		return err
	}

	return nil
}

// ValidateAndNormalize выполняет нормализацию и валидацию
func (r *GetAnomaliesRequest) ValidateAndNormalize() error {
	r.Normalize()
	return r.Validate()
}

// GetAnomaliesRequest представляет запрос для получения состояния устройств
type ActRequest struct {
	Action string `form:"action"`
	Path   string `form:"path"`
}

// Normalize нормализует значения запроса
func (r *ActRequest) Normalize() {
	// Тримим строковые поля
	r.Action = strings.TrimSpace(r.Action)
	r.Path = strings.TrimSpace(r.Path)
}

// Validate выполняет валидацию всех полей запроса
func (r *ActRequest) Validate() error {
	// Валидация типа действия
	allowedTypes := []string{"alllow", "deny"}
	if r.Action != "" && !contains(allowedTypes, r.Action) {
		return fmt.Errorf("invalid action")
	}
	return nil
}

// ValidateAndNormalize выполняет нормализацию и валидацию
func (r *ActRequest) ValidateAndNormalize() error {
	r.Normalize()
	return r.Validate()
}

// GetAnomaliesRequest представляет запрос для получения состояния устройств
type GetDeviceReportRequest struct {
	From     string `form:"from" binding:"omitempty"`
	To       string `form:"to" binding:"omitempty"`
	HostName string `form:"hostname" binding:"omitempty,max=100"`
}

// Normalize нормализует значения запроса
func (r *GetDeviceReportRequest) Normalize() {
	// Тримим строковые поля
	r.From = strings.TrimSpace(r.From)
	r.To = strings.TrimSpace(r.To)
	r.HostName = strings.TrimSpace(r.HostName)
	// Устанавливаем значения по умолчанию
	if r.From == "" {
		r.From = "now-10m"
	}
	if r.To == "" {
		r.To = "now"
	}
}

// Validate выполняет валидацию всех полей запроса
func (r *GetDeviceReportRequest) Validate() error {
	// Валидация временного диапазона
	if err := validateTimeRange(r.From, r.To); err != nil {
		return err
	}

	return nil
}

// ValidateAndNormalize выполняет нормализацию и валидацию
func (r *GetDeviceReportRequest) ValidateAndNormalize() error {
	r.Normalize()
	return r.Validate()
}
