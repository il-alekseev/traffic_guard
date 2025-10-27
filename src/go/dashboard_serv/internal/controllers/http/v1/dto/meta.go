// dto/models.go
package dto

import "time"

// TimeRangeFilter - фильтр по временному диапазону
// swagger:model
type TimeRangeFilter struct {
	// Начало периода
	// Required: true
	// Example: 2024-01-15T10:00:00Z
	Start time.Time `json:"start"`

	// Конец периода
	// Required: true
	// Example: 2024-01-15T11:00:00Z
	End time.Time `json:"end"`

	// Фильтр по NGFW (опционально)
	// Example: ["ngfw-1", "ngfw-2"]
	NGFWIDs []string `json:"ngfw_ids,omitempty"`

	// Фильтр по категориям (опционально)
	// Example: ["category-1", "category-2"]
	CategoryIDs []string `json:"category_ids,omitempty"`
}

// RequestStatsResponse - ответ со статистикой запросов
// swagger:model
type RequestStatsResponse struct {
	// Общая статистика
	Total RequestStats `json:"total"`

	// Статистика по NGFW
	ByNGFW map[string]RequestStats `json:"by_ngfw,omitempty"`

	// Период данных
	Period TimeRange `json:"period"`
}

// TrafficStatsResponse - ответ со статистикой трафика
// swagger:model
type TrafficStatsResponse struct {
	// Общая статистика
	Total TrafficStats `json:"total"`

	// Статистика по NGFW
	ByNGFW map[string]TrafficStats `json:"by_ngfw,omitempty"`

	// Период данных
	Period TimeRange `json:"period"`
}

// TopResourcesResponse - ответ с топом ресурсов
// swagger:model
type TopResourcesResponse struct {
	// Топ 25 ресурсов по доменным именам
	TopDomains []ResourceItem `json:"top_domains"`

	// Топ 25 категорий
	TopCategories []CategoryItem `json:"top_categories"`

	// Период данных
	Period TimeRange `json:"period"`
}

// BlockedCategoriesResponse - ответ с запрещенными категориями
// swagger:model
type BlockedCategoriesResponse struct {
	// Сработавшие запрещающие категории
	BlockedCategories []BlockedCategory `json:"blocked_categories"`

	// Период данных
	Period TimeRange `json:"period"`
}

// BlockedResourcesResponse - ответ с ресурсами запрещенных категорий
// swagger:model
type BlockedResourcesResponse struct {
	// Ресурсы запрещенных категорий
	Resources []BlockedResource `json:"resources"`

	// Общее количество
	Total int `json:"total"`

	// Период данных
	Period TimeRange `json:"period"`
}

// Вспомогательные структуры

// RequestStats - статистика запросов
// swagger:model
type RequestStats struct {
	Total   int64   `json:"total"`
	Allowed int64   `json:"allowed"`
	Blocked int64   `json:"blocked"`
	RPS     float64 `json:"requests_per_second"`
}

// TrafficStats - статистика трафика
// swagger:model
type TrafficStats struct {
	TotalBytes int64   `json:"total_bytes"`
	Incoming   int64   `json:"incoming"`
	Outgoing   int64   `json:"outgoing"`
	BPS        float64 `json:"bytes_per_second"`
}

// ResourceItem - элемент ресурса
// swagger:model
type ResourceItem struct {
	Domain       string    `json:"domain"`
	RequestCount int64     `json:"request_count"`
	LastAccess   time.Time `json:"last_access"`
	NGFWID       string    `json:"ngfw_id"`
}

// CategoryItem - элемент категории
// swagger:model
type CategoryItem struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	RequestCount int64  `json:"request_count"`
}

// BlockedCategory - запрещенная категория
// swagger:model
type BlockedCategory struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	BlockCount  int64  `json:"block_count"`
	Description string `json:"description"`
}

// BlockedResource - ресурс запрещенной категории
// swagger:model
type BlockedResource struct {
	IP           string    `json:"ip"`
	Domain       string    `json:"domain"`
	SourceIP     string    `json:"source_ip"`
	URL          string    `json:"url"`
	RequestCount int64     `json:"request_count"`
	LastAccess   time.Time `json:"last_access"`
	NGFWID       string    `json:"ngfw_id"`
	CategoryID   string    `json:"category_id"`
	CategoryName string    `json:"category_name"`
}

// TimeRange - временной диапазон
// swagger:model
type TimeRange struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}
