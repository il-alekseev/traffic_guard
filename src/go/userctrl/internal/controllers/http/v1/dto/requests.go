// Package dto - Пакет содержит структуры для передачи данных между обработчиками API и клиентом
package dto

import "fmt"

// Cтруктуры для передачи данных пользователя в API

// LogUser — данные для аутентификации (логин и пароль)
type LogUser struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserUpdate — данные для обновления информации о пользователе (логин и email)
type UserUpdate struct {
	Login string `json:"login"`
	Email string `json:"email" binding:"email"`
}

// UserCreateRequest — данные для создания нового пользователя, включая ФИО, логин, email, роль и пароль с валидацией
type UserCreateRequest struct {
	Login      string `json:"login" binding:"required,min=4,max=50" example:"johndoe"`
	Email      string `json:"email" example:"john.doe@example.com"`
	FirstName  string `json:"first_name" binding:"required,min=1,max=50" example:"John"`
	LastName   string `json:"last_name" binding:"required,min=1,max=50" example:"Doe"`
	Patronymic string `json:"patronymic,omitempty" example:"Ivanovich"`
	Role       string `json:"role" binding:"required" example:"SA"`
	Password   string `json:"password" binding:"required,min=4,max=50"`
}

// ValidateUserCreateRequest - валидатор для запроса создания пользователя
func (req *UserCreateRequest) ValidateUserCreateRequest() error {
	if len(req.Password) < 4 {
		return fmt.Errorf("password must be at least 4 characters long")
	}

	if len(req.Role) == 0 {
		return fmt.Errorf("at least one role must be specified")
	}

	return nil
}

// UserUpdateData - содержит данные для обновления пользователя
type UserUpdateData struct {
	FirstName  string `json:"first_name,omitempty"`
	LastName   string `json:"last_name,omitempty"`
	Patronymic string `json:"patronymic,omitempty"`
	Email      string `json:"email,omitempty"`
	Role       string `json:"role,omitempty"`
}

// UserUpdatePass - данные для обновления пароля
type UserUpdatePass struct {
	OldPass string `json:"old_password" binding:"required"`
	NewPass string `json:"new_password" binding:"required"`
}

// UserResetPass - данные для сброса пароля
type UserResetPass struct {
	NewPass string `json:"new_password" binding:"required"`
}
