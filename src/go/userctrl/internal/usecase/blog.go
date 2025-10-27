package usecase

import (
	"userctrl/internal/models"
	"userctrl/pkg/bizlogger"
)

// blogCreateUser — логирует событие создания пользователя с деталями новой учётной записи
func (uc *UseCase) blogCreateUser(userMeta models.UserMeta, newUser models.User) {
	type BlogCreateUser struct {
		ID         string `json:"user_id"`
		Login      string `json:"login"`
		Email      string `json:"email"`
		FirstName  string `json:"first_name"`
		LastName   string `json:"last_name"`
		Patronymic string `json:"patronymic"`
		Role       string `json:"role"`
	}

	description := "the user has been created"
	uc.blog.LogCreate(
		bizlogger.EntityUser,
		userMeta.Username,
		userMeta.ShortRole,
		userMeta.ContextID,
		newUser.Login,
		BlogCreateUser{
			ID:         newUser.ID,
			Login:      newUser.Login,
			Email:      newUser.Email,
			FirstName:  newUser.FirstName,
			LastName:   newUser.LastName,
			Patronymic: newUser.Patronymic,
			Role:       newUser.Role,
		},
		&description,
	)
}

// blogUpdateUser — логирует событие обновления данных пользователя (в текущей реализации старые и новые данные передаются одинаково)
func (uc *UseCase) blogUpdateUser(userMeta models.UserMeta, oldUser models.User, newUser models.User) {
	type BlogUpdateUser struct {
		Email      string `json:"email"`
		FirstName  string `json:"first_name"`
		LastName   string `json:"last_name"`
		Patronymic string `json:"patronymic"`
		Role       string `json:"role"`
	}

	description := "updated user information"
	uc.blog.LogUpdate(
		bizlogger.EntityUser,
		userMeta.Username,
		userMeta.ShortRole,
		userMeta.ContextID,
		newUser.Login,
		BlogUpdateUser{
			Email:      newUser.Email,
			FirstName:  newUser.FirstName,
			LastName:   newUser.LastName,
			Patronymic: newUser.Patronymic,
			Role:       newUser.Role,
		},
		BlogUpdateUser{
			Email:      newUser.Email,
			FirstName:  newUser.FirstName,
			LastName:   newUser.LastName,
			Patronymic: newUser.Patronymic,
			Role:       newUser.Role,
		},
		&description,
	)
}

// blogDeleteUser — логирует событие удаления пользователя с сохранением данных о нём
func (uc *UseCase) blogDeleteUser(userMeta models.UserMeta, oldUser models.User) {
	type BlogDeleteUser struct {
		ID         string `json:"user_id"`
		Login      string `json:"login"`
		Email      string `json:"email"`
		FirstName  string `json:"first_name"`
		LastName   string `json:"last_name"`
		Patronymic string `json:"patronymic"`
		Role       string `json:"role"`
	}

	description := "the user has been deleted"
	uc.blog.LogDelete(
		bizlogger.EntityUser,
		userMeta.Username,
		userMeta.ShortRole,
		userMeta.ContextID,
		oldUser.Login,
		BlogDeleteUser{
			ID:         oldUser.ID,
			Login:      oldUser.Login,
			Email:      oldUser.Email,
			FirstName:  oldUser.FirstName,
			LastName:   oldUser.LastName,
			Patronymic: oldUser.Patronymic,
			Role:       oldUser.Role,
		},
		&description,
	)
}

// blogResertUserPass — логирует сброс пароля пользователя (без передачи данных в payloa
func (uc *UseCase) blogResertUserPass(userMeta models.UserMeta, user models.User) {
	description := "password reset"
	uc.blog.LogUpdate(
		bizlogger.EntityUser,
		userMeta.Username,
		userMeta.ShortRole,
		userMeta.ContextID,
		user.Login,
		nil,
		nil,
		&description,
	)
}
