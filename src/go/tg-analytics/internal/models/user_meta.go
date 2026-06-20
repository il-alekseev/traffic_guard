package models

// UserMeta содержит метаданные пользователя, извлеченные из JWT токена.
// Используется для передачи информации о пользователе между слоями приложения.
type UserMeta struct {
	UUID       string `json:"uuid"`
	Username   string `json:"username"`
	ClientRole string `json:"client_role"`
	ShortRole  string `json:"short_role"`
	ContextID  string `json:"context_id"`
}
