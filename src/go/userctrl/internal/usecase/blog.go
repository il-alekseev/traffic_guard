package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"userctrl/internal/controllers/http/v1/values"
	"userctrl/internal/models"
	"userctrl/pkg/blog/blog/operations"
	blog "userctrl/pkg/blog/models"

	"github.com/google/uuid"
)

// blogCreateUser — логирует событие создания пользователя с деталями новой учётной записи
func (uc *UseCase) blogCreateUser(ctx context.Context, userMeta models.UserMeta, newUser models.User) {
	method := "blogCreateUser"

	type BlogCreateUser struct {
		ID         string `json:"user_id"`
		Login      string `json:"login"`
		Email      string `json:"email"`
		FirstName  string `json:"first_name"`
		LastName   string `json:"last_name"`
		Patronymic string `json:"patronymic"`
		Role       string `json:"role"`
	}

	uuidStr := uuid.New().String()

	description := fmt.Sprintf("пользователь %s (%s) создал пользователя %s (%s)",
		userMeta.Username,
		userMeta.ClientRole,
		newUser.Login,
		newUser.Role)

	record := blog.DtoBusinessLog{
		Description: description,
		Entity:      values.UserEntity,
		EntityID:    userMeta.Username,
		NewValue: BlogCreateUser{
			ID:         newUser.ID,
			Login:      newUser.Login,
			Email:      newUser.Email,
			FirstName:  newUser.FirstName,
			LastName:   newUser.LastName,
			Patronymic: newUser.Patronymic,
			Role:       newUser.Role,
		},
		OldValue:  "",
		EventType: values.CreateOperation,
		Context:   userMeta.ContextID,
		UserName:  userMeta.Username,
		UserRole:  userMeta.ClientRole,
	}

	_, err := uc.blogCl.Operations.PostAPIV1Add(&operations.PostAPIV1AddParams{
		XCallerService: values.ThisServiceName,
		XRequestID:     uuidStr,
		Record:         &record,
		Context:        ctx,
	},
	)
	if err != nil {
		err = fmt.Errorf("%s: failed to blog event: %w", method, err)
		uc.l.ErrorContext(ctx, "failed to blog event",
			slog.String("method", method),
			slog.String("error", err.Error()),
		)
	}
}

// blogUpdateUser — логирует событие обновления данных пользователя
func (uc *UseCase) blogUpdateUser(ctx context.Context, userMeta models.UserMeta, oldUser models.User, newUser models.User) {
	method := "blogUpdateUser"

	type BlogUpdateUser struct {
		Email      string `json:"email"`
		FirstName  string `json:"first_name"`
		LastName   string `json:"last_name"`
		Patronymic string `json:"patronymic"`
		Role       string `json:"role"`
	}

	uuidStr := uuid.New().String()

	description := fmt.Sprintf("пользователь %s (%s) обновил данные пользователя %s",
		userMeta.Username,
		userMeta.ClientRole,
		newUser.Login)

	record := blog.DtoBusinessLog{
		Description: description,
		Entity:      values.UserEntity,
		EntityID:    newUser.Login,
		NewValue: BlogUpdateUser{
			Email:      newUser.Email,
			FirstName:  newUser.FirstName,
			LastName:   newUser.LastName,
			Patronymic: newUser.Patronymic,
			Role:       newUser.Role,
		},
		OldValue: BlogUpdateUser{
			Email:      oldUser.Email,
			FirstName:  oldUser.FirstName,
			LastName:   oldUser.LastName,
			Patronymic: oldUser.Patronymic,
			Role:       oldUser.Role,
		},
		EventType: values.UpdateOperation,
		Context:   userMeta.ContextID,
		UserName:  userMeta.Username,
		UserRole:  userMeta.ClientRole,
	}

	_, err := uc.blogCl.Operations.PostAPIV1Add(&operations.PostAPIV1AddParams{
		XCallerService: values.ThisServiceName,
		XRequestID:     uuidStr,
		Record:         &record,
		Context:        ctx,
	})
	if err != nil {
		err = fmt.Errorf("%s: failed to blog event: %w", method, err)
		uc.l.ErrorContext(ctx, "failed to blog event",
			slog.String("method", method),
			slog.String("error", err.Error()),
		)
	}
}

// blogDeleteUser — логирует событие удаления пользователя с сохранением данных о нём
func (uc *UseCase) blogDeleteUser(ctx context.Context, userMeta models.UserMeta, oldUser models.User) {
	method := "blogDeleteUser"

	type BlogDeleteUser struct {
		ID         string `json:"user_id"`
		Login      string `json:"login"`
		Email      string `json:"email"`
		FirstName  string `json:"first_name"`
		LastName   string `json:"last_name"`
		Patronymic string `json:"patronymic"`
		Role       string `json:"role"`
	}

	uuidStr := uuid.New().String()

	description := fmt.Sprintf("пользователь %s (%s) удалил пользователя %s",
		userMeta.Username,
		userMeta.ClientRole,
		oldUser.Login)

	record := blog.DtoBusinessLog{
		Description: description,
		Entity:      values.UserEntity,
		EntityID:    oldUser.Login,
		NewValue:    "",
		OldValue: BlogDeleteUser{
			ID:         oldUser.ID,
			Login:      oldUser.Login,
			Email:      oldUser.Email,
			FirstName:  oldUser.FirstName,
			LastName:   oldUser.LastName,
			Patronymic: oldUser.Patronymic,
			Role:       oldUser.Role,
		},
		EventType: values.DeleteOperation,
		Context:   userMeta.ContextID,
		UserName:  userMeta.Username,
		UserRole:  userMeta.ClientRole,
	}

	_, err := uc.blogCl.Operations.PostAPIV1Add(&operations.PostAPIV1AddParams{
		XCallerService: values.ThisServiceName,
		XRequestID:     uuidStr,
		Record:         &record,
		Context:        ctx,
	})
	if err != nil {
		err = fmt.Errorf("%s: failed to blog event: %w", method, err)
		uc.l.ErrorContext(ctx, "failed to blog event",
			slog.String("method", method),
			slog.String("error", err.Error()),
		)
	}
}

// blogResetUserPass — логирует сброс пароля пользователя
func (uc *UseCase) blogResetUserPass(ctx context.Context, userMeta models.UserMeta, user models.User) {
	method := "blogResetUserPass"

	uuidStr := uuid.New().String()

	description := fmt.Sprintf("пользователь %s (%s) сбросил пароль пользователю %s",
		userMeta.Username,
		userMeta.ClientRole,
		user.Login)

	record := blog.DtoBusinessLog{
		Description: description,
		Entity:      values.UserEntity,
		EntityID:    user.Login,
		NewValue:    "",
		OldValue:    "",
		EventType:   values.UpdateOperation,
		Context:     userMeta.ContextID,
		UserName:    userMeta.Username,
		UserRole:    userMeta.ClientRole,
	}

	_, err := uc.blogCl.Operations.PostAPIV1Add(&operations.PostAPIV1AddParams{
		XCallerService: values.ThisServiceName,
		XRequestID:     uuidStr,
		Record:         &record,
		Context:        ctx,
	})
	if err != nil {
		err = fmt.Errorf("%s: failed to blog event: %w", method, err)
		uc.l.ErrorContext(ctx, "failed to blog event",
			slog.String("method", method),
			slog.String("error", err.Error()),
		)
	}
}
