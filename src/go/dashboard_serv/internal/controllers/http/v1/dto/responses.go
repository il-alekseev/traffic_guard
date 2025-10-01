// Package dto - Пакет содержит структуры для передачи данных между обработчиками API и клиентом
package dto

// SuccessResponse — успешный ответ с сообщением
type SuccessResponse struct {
	Message string `json:"message"`
}

// ErrorResponse — ответ с описанием ошибки
type ErrorResponse struct {
	Error string `json:"error"`
}
