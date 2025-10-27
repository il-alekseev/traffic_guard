package models

import "errors"

// Стандартные ошибки приложения, используемые в бизнес-логике и обработчиках
// Определяют типовые ситуации: пользователь или роль не найдены, неверные учётные данные,
// уже существующий пользователь, недействительный токен, отсутствие прав доступа
var (
	ErrUserNotFound       = errors.New("user not found")
	ErrRoleNotFound       = errors.New("role not found")
	ErrGroupNotFound      = errors.New("group not found")
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrPermissionDenied   = errors.New("permission denied")
)
