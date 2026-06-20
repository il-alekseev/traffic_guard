// Package dto - Пакет содержит структуры для передачи данных между обработчиками API и клиентом
package dto

import "userctrl/internal/models"

// SuccessResponse — успешный ответ с сообщением
type SuccessResponse struct {
	Message string `json:"message"`
}

// ErrorResponse — ответ с описанием ошибки
type ErrorResponse struct {
	Error string `json:"error"`
}

// UsersListResponse — список пользователей с метаданными пагинации
type UsersListResponse struct {
	Data []models.User  `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

// RolesListResponse — список ролей с метаданными пагинации
type RolesListResponse struct {
	Data []models.Role  `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

// PaginationMeta — метаданные пагинации: текущая страница, лимит, общее количество, число страниц
type PaginationMeta struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
	Pages int `json:"pages"`
}

// CountResponse — ответ, содержащий числовое значение счётчика

type CountResponse struct {
	Count int `json:"count"`
}

// CountUserForContextResponse — количество пользователей по категориям (CA, CO)
type CountUserForContextResponse struct {
	CountCA int `json:"count_ca"`
	CountCO int `json:"count_co"`
}

// CountUserByRoleResponse — количество пользователей по ролям (SA, CA, CO)
type CountUserByRoleResponse struct {
	CountSA int `json:"count_sa"`
	CountCA int `json:"count_ca"`
	CountCO int `json:"count_co"`
}
