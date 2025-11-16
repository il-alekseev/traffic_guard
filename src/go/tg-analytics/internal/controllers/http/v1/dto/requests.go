package dto

// SearchRequest представляет параметры фильтрации для запросов
// Используется для поиска по частичному совпадению строк
type SearchRequest struct {
	Like string `form:"like" binding:"omitempty"` // Строка для поиска по шаблону (LIKE)
}

type DetectionActRequest struct {
	Action string `form:"action"`
	Path   string `form:"path"`
}
