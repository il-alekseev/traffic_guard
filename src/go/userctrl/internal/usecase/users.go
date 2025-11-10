package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"userctrl/internal/controllers/http/v1/dto"
	"userctrl/internal/models"
	"userctrl/internal/usecase/utils"
	"userctrl/pkg/keycloakclient"
	"userctrl/pkg/slogger"
	"userctrl/pkg/slogger/wsl"
)

// convertUserFromKeyCloakClient - копирование из *keycloakclient.User в models.User
func convertUserFromKeyCloakClient(u *keycloakclient.User) models.User {
	user := models.User{
		ID:                     u.ID,
		Login:                  u.Username,
		Email:                  u.Email,
		FirstName:              u.FirstName,
		LastName:               u.LastName,
		Patronymic:             u.Patronymic,
		Role:                   u.Role,
		CreatedAt:              u.CreatedAt,
		IsNeedToChangePassword: u.IsNeedToChangePassword,
		IsSuperAdmin:           u.IsSuperAdmin,
	}
	return user
}

// GetUsers возвращает список пользователей с пагинацией
func (uc *UseCase) GetUsers(ctx context.Context, meta *models.UserMeta, page, limit int, role, contextID, search string) ([]models.User, int, error) {
	page = page - 1 // так так нумерация с 1

	var ctxFilter string
	var roleFilter string

	// Обрабатываем права доступа в зависимости от роли
	switch meta.ShortRole {
	case "SA":
		// Super Admin - использует переданные фильтры как есть
		ctxFilter = contextID
		roleFilter = role

	case "CA":
		// Context Admin - может видеть только пользователей своего контекста
		if contextID != "" && contextID != meta.ContextID {
			return nil, 0, models.ErrPermissionDenied
		}
		if role == "SA" {
			return nil, 0, models.ErrPermissionDenied
		}
		ctxFilter = meta.ContextID // Принудительно устанавливаем контекст
		roleFilter = role          // Но фильтр по роли оставляем (кроме SA)

	case "CO":
		// Common User - не имеет прав на просмотр списка пользователей
		return nil, 0, models.ErrPermissionDenied

	default:
		// Неизвестная роль - отказ в доступе
		return nil, 0, models.ErrPermissionDenied

	}

	usersKC, total, err := uc.kc.GetUsers(ctx, page, limit, roleFilter, ctxFilter, "", search)

	if err != nil {
		err = fmt.Errorf("keyclakclient get users: %w", err)
		return nil, 0, err
	}
	users := make([]models.User, 0, len(usersKC))
	for _, val := range usersKC {
		u := convertUserFromKeyCloakClient(val)
		users = append(users, u)
	}

	return users, total, nil
}

// GetByID возвращает пользователя по ID
func (uc *UseCase) GetByID(ctx context.Context, meta *models.UserMeta, userID string) (models.User, error) {
	res := models.User{}

	// Анонимная функция для получения и конвертации пользователя
	getUser := func(id string) (models.User, error) {
		userKC, err := uc.kc.GetUserByID(id)
		if err != nil {
			return models.User{}, models.ErrUserNotFound

		}
		return convertUserFromKeyCloakClient(userKC), nil
	}

	// Если запрашивающий пользователь запрашивает самого себя
	if meta.UUID == userID {
		return getUser(userID)
	}

	// Для всех остальных случаев проверяем роли
	switch meta.ShortRole {
	case "SA":
		// Super Admin имеет доступ ко всем пользователям
		return getUser(userID)

	case "CA":
		// Context Admin имеет доступ только к пользователям своего контекста
		u, err := getUser(userID)
		if err != nil {
			return res, err
		}

		_, contextID, err := utils.ParseRole(u.Role)
		if err != nil {
			return res, err
		}
		if meta.ContextID == contextID {
			return u, nil
		}
		return res, models.ErrPermissionDenied

	default:
		// Все остальные роли не имеют доступа к чужим данным
		return res, models.ErrPermissionDenied
	}
}

// GetProfile возвращает профиль текущего пользователя
func (uc *UseCase) GetProfile(ctx context.Context, meta *models.UserMeta, userID string) (models.User, error) {
	return uc.GetByID(ctx, meta, userID)
}

// Create создает нового пользователя
// SA может создавать пользователей с любой ролью
// CA может создавать только пользователей своего контекста
func (uc *UseCase) Create(ctx context.Context, meta *models.UserMeta, u dto.UserCreateRequest) (models.User, error) {
	log := uc.l.With(wsl.Label("method", "Create"))

	res := models.User{}

	// Проверка прав доступа
	switch meta.ShortRole {
	case "SA":
		break
	case "CA":
		// Context Admin - проверка соответствия контекста
		_, contextID, err := utils.ParseRole(u.Role)
		if err != nil {
			return res, fmt.Errorf("invalid role format: %w", err)
		}

		if contextID != meta.ContextID {
			return res, models.ErrPermissionDenied
		}

	default:
		return res, models.ErrPermissionDenied
	}

	users, _, err := uc.kc.GetUsersByUsername(ctx, u.Login)
	if err != nil {
		err = fmt.Errorf("keycloak GetUsers: %w", err)
		log.DebugContext(ctx, err.Error())
		return res, err
	}
	// Проверяем, существует ли пользователь с таким email/логином
	for _, us := range users {
		if u.Login == us.Username {
			err = fmt.Errorf("the username '%s' is already occupied, : %w", u.Login, models.ErrUserExists)
			return res, err
		}
	}
	// Проверка, что такая роль есть
	roles, _, err := uc.GetRoles(ctx, meta, 1, 1, u.Role)
	if err != nil {
		err = fmt.Errorf("keycloak GetRoles: %w", err)
		log.DebugContext(ctx, err.Error())
		return res, err
	}

	roleIsExist := false
	for _, r := range *roles {
		if r.Role == u.Role {
			roleIsExist = true
			break
		}
	}

	if !roleIsExist {
		err = fmt.Errorf("role '%s': %w", u.Role, models.ErrRoleNotFound)
		log.DebugContext(ctx, err.Error())
		return res, err
	}

	rolePath := utils.GetGroupsByRole(u.Role) // переобразуем роль в путь до группы
	// Проверка, что такая группа есть
	g, err := uc.kc.GetGroupByPath(ctx, rolePath)
	if err != nil {
		err = fmt.Errorf("keycloak GetGroups: %w", err)
		log.DebugContext(ctx, err.Error())
		return res, err
	}

	for g == nil || *g.Path != rolePath {
		err = fmt.Errorf("group '%s': %w", rolePath, models.ErrGroupNotFound)
		log.DebugContext(ctx, err.Error())
		return res, err
	}

	user := keycloakclient.User{
		Username:   u.Login,
		FirstName:  u.FirstName,
		LastName:   u.LastName,
		Patronymic: u.Patronymic,
		Email:      u.Email,
	}

	log.DebugContext(ctx, fmt.Sprintf("try create user: %+v", user))
	//TODO: проверка, что такой путь до группы есть
	kcuser, err := uc.kc.CreateUser(user, &rolePath, u.Password)
	if err != nil {
		err = fmt.Errorf("uc.kc.CreateUser: %w", err)
		log.DebugContext(ctx, err.Error())
		return res, err
	}

	resUser := convertUserFromKeyCloakClient(kcuser)
	uc.blogCreateUser(*meta, resUser)
	return resUser, nil
}

// Update обновляет данные пользователя
func (uc *UseCase) Update(ctx context.Context, meta *models.UserMeta, userID string, u dto.UserUpdateData) error {
	log := uc.l.With(wsl.Label("method", "Update"))

	// Проверяем базовые права доступа
	switch {
	case meta.UUID == userID:
		// Пользователь может редактировать себя, независимо от роли
		break

	case meta.ShortRole == "SA":
		// Super Admin имеет полный доступ
		break

	case meta.ShortRole == "CA":
		// Context Admin может редактировать только пользователей своего контекста
		oldUser, err := uc.kc.GetUserByID(userID)
		if err != nil {
			log.ErrorContext(ctx, "uc.kc.GetUserByID", wsl.Label("userID", userID), wsl.Err(err))
			return models.ErrUserNotFound

		}

		_, contextID, err := utils.ParseRole(oldUser.Role)
		if err != nil {
			return err
		}
		if contextID != meta.ContextID {
			return models.ErrPermissionDenied
		}

	case meta.ShortRole == "CO":
		// Common User может редактировать только себя
		return models.ErrPermissionDenied

	default:
		// Неизвестная роль
		return models.ErrPermissionDenied
	}

	// Получаем текущие данные пользователя для логирования изменений
	oldUser, err := uc.kc.GetUserByID(userID)
	if err != nil {
		log.ErrorContext(ctx, "uc.kc.GetUserByID ", wsl.Label("userID", userID), wsl.Err(err))
		return models.ErrUserNotFound
	}

	err = uc.kc.UpdateUser(ctx,
		userID,
		keycloakclient.User{
			FirstName:  u.FirstName,
			LastName:   u.LastName,
			Email:      u.Email,
			Patronymic: u.Patronymic,
			Role:       u.Role,
		},
	)
	if err != nil {
		return err
	}

	// Логируем изменения
	newUser, err := uc.kc.GetUserByID(userID)
	if err != nil {
		err = fmt.Errorf("uc.kc.GetUserByID, userID=%s: %w", userID, err)
		log.ErrorContext(ctx, err.Error())
		return err
	}
	uc.blogUpdateUser(
		*meta,
		convertUserFromKeyCloakClient(oldUser),
		convertUserFromKeyCloakClient(newUser),
	)
	return nil
}

// Delete удаляет пользователя
func (uc *UseCase) Delete(ctx context.Context, meta *models.UserMeta, userID string) error {
	log := uc.l.With(wsl.Label("method", "Delete"))

	// Проверка базовых прав доступа
	switch meta.ShortRole {
	case "SA":
		break
	case "CA":
		// Context Admin может удалять только пользователей своего контекста
		user, err := uc.kc.GetUserByID(userID)
		if err != nil {
			return models.ErrUserNotFound
		}

		_, contextID, err := utils.ParseRole(user.Role)
		if err != nil {
			return fmt.Errorf("failed to parse user role: %w", err)
		}

		if contextID != meta.ContextID {
			return models.ErrPermissionDenied
		}
	default:
		return models.ErrPermissionDenied

	}
	// TODO: Сброс ssh сессии Linux

	// Получаем данные пользователя для логгирования и очистки
	user, err := uc.kc.GetUserByID(userID)
	if err != nil {
		return models.ErrUserNotFound

	}

	ctx = slogger.WithLog(ctx, "username", user.Username)
	// Если пользователь есть в Grafana
	grafanaUserID, err := uc.grafanaCl.GetUserByLogin(user.Username)
	if err != nil {
		err = fmt.Errorf("uc.grafanaCl.GetUserByLogin, username=%s: %w", user.Username, err)
		log.ErrorContext(ctx, err.Error())
		// return err
	} else {
		if err := uc.grafanaCl.LogoutUser(grafanaUserID); err != nil {
			err = fmt.Errorf("uc.grafanaCl.LogoutUser, username=%s: %w", user.Username, err)
			log.ErrorContext(ctx, err.Error())
			//return err
		}

		if err := uc.grafanaCl.DeleteUserById(grafanaUserID); err != nil {
			err = fmt.Errorf("uc.grafanaCl.DeleteUserById, username=%s: %w", user.Username, err)
			log.ErrorContext(ctx, err.Error())
			//return err
		}
	}

	// Основная операция - удаление пользователя из Keycloak
	if err := uc.kc.DeleteUserByID(ctx, userID); err != nil {
		err = fmt.Errorf("failed to delete user from Keycloak: %w", err)
		log.ErrorContext(ctx, err.Error())
		return slogger.WrapError(ctx, err)
	}

	uc.blogDeleteUser(*meta, convertUserFromKeyCloakClient(user))

	return nil
}

// Count возвращает общее количество пользователей
func (uc *UseCase) Count(ctx context.Context, meta *models.UserMeta) (int, error) {
	if meta.ShortRole != "SA" {
		return 0, models.ErrPermissionDenied
	}
	return uc.kc.GetUserCount(ctx)
}

func (uc *UseCase) GetUsersCountForContext(ctx context.Context, meta *models.UserMeta, ctxID string) (*models.MetaDataForContext, error) {
	m, err := uc.kc.GetMetaUsersInContext(ctx, ctxID)
	if err != nil {
		return nil, fmt.Errorf("uc.kc.GetMetaUsersInContext, ctxID=%s: %w", ctxID, err)
	}
	return m, nil
}

// UpdatePassword обновляет пароль пользователя
func (uc *UseCase) UpdatePassword(ctx context.Context, meta *models.UserMeta, userID string, oldPassword, newPassword string) error {
	if meta.UUID != userID {
		return models.ErrPermissionDenied
	}
	err := uc.kc.ChangePass(ctx, userID, oldPassword, newPassword)
	return err
}

// ResetPassword сбрасывает пароль пользователя
// SA может сбрасывать пароль любому пользователю
// CA может сбрасывать пароль только пользователям своего контекста
func (uc *UseCase) ResetPassword(ctx context.Context, meta *models.UserMeta, userID string, newPassword string) error {
	log := uc.l.With(wsl.Label("method", "ResetPassword"))

	ctx = slogger.WithLog(ctx, "userID", userID)
	// Проверка прав доступа
	switch meta.ShortRole {
	case "SA":
		break
	case "CA":
		// Context Admin - проверяем принадлежность к контексту
		user, err := uc.kc.GetUserByID(userID)
		if err != nil {
			log.ErrorContext(ctx, "failed to get user for password reset", wsl.Err(err))
			return slogger.WrapError(ctx, models.ErrUserNotFound)
		}

		_, contextID, err := utils.ParseRole(user.Role)
		if err != nil {
			return slogger.WrapError(ctx, fmt.Errorf("failed to parse user role: %w", err))
		}

		if contextID != meta.ContextID {
			return slogger.WrapError(ctx, models.ErrPermissionDenied)
		}

	default:
		// Все остальные роли не имеют прав на сброс пароля
		return slogger.WrapError(ctx, models.ErrPermissionDenied)
	}

	// Выполняем сброс пароля
	if err := uc.kc.ResetPass(ctx, userID, newPassword); err != nil {
		log.ErrorContext(ctx, "failed to reset password for user %v", wsl.Err(err))
		return slogger.WrapError(ctx, fmt.Errorf("password reset failed: %w", err))
	}

	// Логируем действие
	user, err := uc.kc.GetUserByID(userID)
	if err != nil {
		log.ErrorContext(ctx, "failed to get user for auditing", wsl.Err(err))
		// Не прерываем выполнение, так как операция уже выполнена
	} else {
		uc.blogResertUserPass(*meta, convertUserFromKeyCloakClient(user))
	}

	return nil
}

// UpdateRoles обновляет роли пользователя
func (uc *UseCase) UpdateRole(ctx context.Context, meta *models.UserMeta, userID string, role string) error {
	// TODO
	return errors.New("в разработке")
}

// RemoveRoles удаляет роли у пользователя
func (uc *UseCase) RemoveRole(ctx context.Context, meta *models.UserMeta, userID string, role string) error {
	// TODO
	return errors.New("в разработке")
}

func (uc *UseCase) GetCountUsersByRole(ctx context.Context, meta *models.UserMeta) (*models.CountUserByRole, error) {
	log := uc.l.With(wsl.Label("method", "GetCountUsersByRole"))

	if meta.ShortRole != "SA" {
		return nil, models.ErrPermissionDenied
	}
	users, _, err := uc.kc.GetAllUsers(ctx)
	if err != nil {
		err = fmt.Errorf("uc.kc.GetAllUsers: %w", err)
		log.ErrorContext(ctx, err.Error())
		return nil, err
	}
	res := models.CountUserByRole{}
	for _, u := range users {
		switch {
		case u.Role == "SA":
			res.CountSA += 1
		case strings.HasPrefix(u.Role, "CA-"):
			res.CountCA += 1
		case strings.HasPrefix(u.Role, "CO-"):
			res.CountCO += 1
		}
	}
	return &res, nil
}
