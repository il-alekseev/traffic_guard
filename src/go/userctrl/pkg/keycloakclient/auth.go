package keycloakclient

import (
	"context"
	"userctrl/pkg/slogger"
	"userctrl/pkg/slogger/wsl"
	"fmt"
	"strings"

	"github.com/Nerzal/gocloak/v13"
)

// Login — метод клиента Keycloak для аутентификации пользователя по логину и паролю
// При успешном входе возвращает JWT-токен
func (kc *KeycloakClient) Login(ctx context.Context,
	username string,
	password string) (*gocloak.JWT, error) {
	token, err := kc.client.Login(ctx, kc.clinetName, kc.clientSecret, kc.realm, username, password)
	log := kc.logger.With(wsl.Label("method", "Login"))

	ctx = slogger.WithLog(ctx, "user", username)
	ctx = slogger.WithLog(ctx, "realm", kc.realm)

	log.DebugContext(ctx, "login attempt")

	if err != nil {
		// Проверяем тип ошибки
		// Обрабатывает специфические ошибки Keycloak:
		// - временный пароль (400 с сообщением "temporary password") — переадресует на LoginTempPass
		// - незавершённая настройка аккаунта (400, "Account is not fully set up") — также переходит к смене пароля
		// - неверные учётные данные (401) — возвращает ошибку "invalid credentials"
		if apiErr, ok := err.(*gocloak.APIError); ok {
			switch {
			case apiErr.Code == 400 && strings.Contains(apiErr.Message, "temporary password"):
				log.InfoContext(ctx, "Temporary password detected")
				token, err := kc.LoginTempPass(ctx, username, password)
				return token, slogger.WrapError(ctx, err)
			case apiErr.Code == 400 && strings.Contains(apiErr.Message, "Account is not fully set up"):
				log.InfoContext(ctx, fmt.Sprintf("Account is not fully set up: %v", apiErr))
				token, err := kc.LoginTempPass(ctx, username, password)
				return token, slogger.WrapError(ctx, err)
			case apiErr.Code == 401 && strings.Contains(apiErr.Message, "Invalid user credentials"):
				log.InfoContext(ctx, "Invalid credentials")
				return nil, slogger.WrapError(ctx, fmt.Errorf("invalid credentials %w", err))
			default:
				log.InfoContext(ctx, fmt.Sprintf("Keycloak API error: %v", apiErr))
			}
		} else {
			log.InfoContext(ctx, "Login", wsl.Err(err))
		}
		return nil, slogger.WrapError(ctx, err)
	}
	log.InfoContext(ctx, fmt.Sprintf("token: %+v\n", *token))
	return token, nil
}

// Ультра мега костыль
// Чтобы залогинить пользователя с временным паролем
// 1. Меняем пароль на тот же - он перестаёт быть временным
// 2. Логинимся
// 1. Сбрасываем пароль на тот же - он снова временный
func (kc *KeycloakClient) LoginTempPass(ctx context.Context, username, password string) (*gocloak.JWT, error) {
	log := kc.logger.With(wsl.Label("method", "LoginTempPass"))

	params := gocloak.GetUsersParams{
		Username: &username,
	}
	users, err := kc.client.GetUsers(ctx, kc.GetToken().AccessToken, kc.realm, params)
	if err != nil {
		return nil, slogger.WrapError(ctx, err)
	}
	log.DebugContext(ctx, fmt.Sprintf("GetUsers: %+v", users))
	if len(users) < 1 {
		return nil, ErrUserNotFound
	}
	user := users[0]
	err = kc.setСonstantPassword(ctx, *user.ID, password)
	if err != nil {
		log.DebugContext(ctx, "kc.ChangePass", wsl.Err(err))
		return nil, slogger.WrapError(ctx, err)
	}
	token, err := kc.client.Login(ctx, kc.clinetName, kc.clientSecret, kc.realm, username, password)
	if err != nil {
		log.DebugContext(ctx, "kc.client.Login", wsl.Err(err))
		return nil, slogger.WrapError(ctx, err)
	}
	err = kc.ResetPass(ctx, *user.ID, password)
	if err != nil {
		log.DebugContext(ctx, "kc.ResetPass", wsl.Err(err))
		return nil, slogger.WrapError(ctx, err)
	}
	return token, nil
}

// Logout — производит выход пользователя по переданному refresh token
func (kc *KeycloakClient) Logout(ctx context.Context,
	refreshToken string) error {
	return kc.client.Logout(ctx, kc.clinetName, kc.clientSecret, kc.realm, refreshToken)
}

// RefreshToken — обновляет access и refresh токены с использованием текущего refresh token
func (kc *KeycloakClient) RefreshToken(ctx context.Context,
	refreshToken string) (*gocloak.JWT, error) {
	return kc.client.RefreshToken(ctx, refreshToken, kc.clinetName, kc.clientSecret, kc.realm)
}

// ValidationPassword — проверяет валидность пароля пользователя через попытку входа в Keycloak
// Возвращает true, если вход успешен или требуется смена временного пароля (учётная запись активна)
// Возвращает false при неверных учётных данных или других ошибках аутентификации
func (kc *KeycloakClient) ValidationPassword(ctx context.Context, username, password string) bool {
	log := kc.logger.With(wsl.Label("method", "validatePassword"))

	ctx = slogger.WithLog(ctx, "username", username)

	log.DebugContext(ctx, "ValidationPassword", wsl.Label("username", username))
	token, err := kc.client.Login(ctx, kc.clinetName, kc.clientSecret, kc.realm, username, password)
	log.Debug(fmt.Sprintf("login attempt for user: %s in realm: %s", username, kc.realm))
	if err != nil {
		// Проверяем тип ошибки
		if apiErr, ok := err.(*gocloak.APIError); ok {
			switch {
			case apiErr.Code == 400 && strings.Contains(apiErr.Message, "temporary password"):
				log.InfoContext(ctx, "Temporary password detected")
				return true
			case apiErr.Code == 400 && strings.Contains(apiErr.Message, "Account is not fully set up"):
				log.InfoContext(ctx, fmt.Sprintf("Account is not fully set up: %v", apiErr))
				return true
			case apiErr.Code == 401 && strings.Contains(apiErr.Message, "Invalid user credentials"):
				log.InfoContext(ctx, "Invalid credentials")
				return false
			default:
				log.InfoContext(ctx, fmt.Sprintf("Keycloak API error: %v", apiErr))
			}
		} else {
			log.InfoContext(ctx, "Login", wsl.Err(err))
		}
		return false
	}

	// После успешной проверки токен сразу аннулируется (Logout), чтобы не оставлять активную сессию
	err = kc.client.Logout(ctx, kc.clinetName, kc.clientSecret, kc.realm, token.RefreshToken)
	if err != nil {
		log.ErrorContext(ctx, "unable logout", wsl.Err(err))
	}
	return true
}
