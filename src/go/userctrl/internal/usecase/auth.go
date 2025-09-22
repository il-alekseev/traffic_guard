package usecase

import (
	"context"
	"userctrl/internal/controllers/http/v1/values"
	"userctrl/internal/models"
	"userctrl/pkg/grafcookier"
	"userctrl/pkg/slogger"
	"userctrl/pkg/slogger/wsl"
	"fmt"

	"github.com/Nerzal/gocloak/v13"
)

// SignIn выполняет аутентификацию пользователя
func (uc *UseCase) SignIn(ctx context.Context, username, password string) (models.Token, error) {
	log := uc.l.With(wsl.Label("method", "SignIn"))
	token := models.Token{}

	ctx = slogger.WithLog(ctx, "username", username)
	// получение данных о пользователе из Keycloak
	user, count, err := uc.kc.GetUsersByUsername(ctx, username)
	if err != nil {
		err = fmt.Errorf("uc.kc.GetUsersByUsername, %w", err)
		log.ErrorContext(ctx, err.Error())
		return token, slogger.WrapError(ctx, err)
	}

	// проверка наличия записей
	if count < 1 {
		err = fmt.Errorf("SignIn %w", models.ErrUserNotFound)
		log.ErrorContext(ctx, err.Error())
		return token, slogger.WrapError(ctx, err)
	}

	tempPass := false
	u := user[0]
	if u.IsNeedToChangePassword {
		tempPass = true
		if err := uc.kc.ChangePass(ctx, u.ID, password, password); err != nil {
			log.ErrorContext(ctx, "uc.kc.ChangePass", wsl.Err(err))
		}
	}

	// Каналы для получения результатов из горутин
	type jwtResult struct {
		jwt *gocloak.JWT
		err error
	}
	jwtChan := make(chan jwtResult)

	type cookieResult struct {
		cookie *grafcookier.GrafCookies
		err    error
	}
	cookieChan := make(chan cookieResult)

	// Запускаем горутины для параллельного выполнения
	go func() {
		jwt, err := uc.kc.Login(ctx, username, password)
		jwtChan <- jwtResult{jwt, err}
	}()

	go func() {
		cookie, err := uc.grafanaCookier.GetCookies(username, password)
		cookieChan <- cookieResult{cookie, err}
	}()

	// Ожидаем результаты из обеих горутин
	jwtRes := <-jwtChan
	if jwtRes.err != nil {
		err = fmt.Errorf("keycloak client login: %w", jwtRes.err)
		log.ErrorContext(ctx, err.Error())
		return token, slogger.WrapError(ctx, err)
	}
	// заполнение структуры токена
	token.AccessToken = jwtRes.jwt.AccessToken
	token.RefreshToken = jwtRes.jwt.RefreshToken
	token.ExpiresIn = jwtRes.jwt.ExpiresIn

	cookieRes := <-cookieChan
	if cookieRes.err != nil {
		err = fmt.Errorf("uc.grafanaCookier.GetCookies, username=%s: %w", username, cookieRes.err)
		log.ErrorContext(ctx, err.Error())
		return token, slogger.WrapError(ctx, err)
	}

	if u.Role == values.SystemAdmin {
		if err := uc.grafanaCl.AddAdminUserToAllOrgs(username); err != nil {
			err = fmt.Errorf("uc.grafanaCl.AddAdminUserToAllOrgs, username=%s: %w", username, err)
			log.ErrorContext(ctx, err.Error())
		}
	}

	if tempPass {
		u := user[0]
		if err := uc.kc.ResetPass(ctx, u.ID, password); err != nil {
			log.ErrorContext(ctx, "uc.kc.ResetPass, username=%s: %w", username, err)
		}
	}

	token.GrafanaSession = cookieRes.cookie.GrafanaSession
	token.GrafanaSessionExpiry = cookieRes.cookie.GrafanaSessionExpiry

	return token, nil
}

// SignOut - выполняет выход пользователя из системы
func (uc *UseCase) SignOut(ctx context.Context, refreshToken string) error {
	err := uc.kc.Logout(ctx, refreshToken)
	return err
}

// RefreshToken - выполняет выход пользователя из системы
func (uc *UseCase) RefreshToken(ctx context.Context, refreshToken string) (models.Token, error) {
	token := models.Token{}

	// запрос в Keycloak для обновления токена
	jwt, err := uc.kc.RefreshToken(ctx, refreshToken)
	if err != nil {
		return token, fmt.Errorf("keycloak client RefreshToken: %w", err)
	}
	token.AccessToken = jwt.AccessToken
	token.RefreshToken = jwt.RefreshToken
	token.ExpiresIn = jwt.ExpiresIn
	return token, nil
}

// ValidateToken проверяет валидность токена
func (uc *UseCase) ValidateToken(ctx context.Context, accessToken string) (int, error) {
	res, err := uc.kc.ValidateToken(ctx, accessToken)
	if err != nil || !res {
		return 1, err
	}
	return 0, nil
}
