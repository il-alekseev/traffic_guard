package dto

// SearchRequest представляет параметры фильтрации для запросов
// Используется для поиска по частичному совпадению строк
type SearchRequest struct {
	Like string `form:"like" binding:"omitempty"` // Строка для поиска по шаблону (LIKE)
}
