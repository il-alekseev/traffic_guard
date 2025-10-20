package dto

type SuccessResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"` // Описание ошибки
}

// PaginationMeta содержит метаданные пагинации
// Используется вместе со списками данных для навигации по страницам
type PaginationMeta struct {
	Page  int   `json:"page"`  // Текущая страница (начинается с 1)
	Limit int   `json:"limit"` // Количество элементов на странице
	Total int64 `json:"total"` // Общее количество элементов
	Pages int   `json:"pages"` // Общее количество страниц
}

type ListResponse struct {
	Data interface{}    `json:"data"` // Список однотипных данных
	Meta PaginationMeta `json:"meta"` // Метаданные пагинации
}

type DataPointsResponse struct {
	Type  string `json:"type"`  // Тип данных
	Data  []uint `json:"data"`  // Список однотипных данных
	Count uint   `json:"count"` // Число точек с данными
}
