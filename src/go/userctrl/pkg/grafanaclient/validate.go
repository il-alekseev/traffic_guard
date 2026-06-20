package grafanaclient

import (
	"userctrl/pkg/slogger/wsl"
	"fmt"
)

type ValidateInputDataError struct {
	msg string
}

func (e *ValidateInputDataError) Error() string {
	return e.msg
}

func (gc GrafanaClient) validateUserLogin(login string) error {
	log := gc.logger.With(wsl.Label("method", "validateUserLogin"))

	// Проверка на совпадение логина пользвателя с зарезервированным логином
	gcData, err := gc.getGrafanaClientInfo()
	if err != nil {
		log.Error("getGrafanaClientInfo", wsl.Err(err))
		return err
	}
	// Если логин совпадает с зарезервированным, генерируем ошибку
	if gcData["Login"] == login {
		validateError := &ValidateInputDataError{msg: fmt.Sprintf("Login %s reserved", login)}
		return validateError
	}
	return nil
}

func (gc GrafanaClient) validateUserID(id int64) error {
	log := gc.logger.With(wsl.Label("method", "validateUserID"))

	// Проверка на совпадение ID пользователя с зарезервированным ID
	gcData, err := gc.getGrafanaClientInfo()
	if err != nil {
		log.Error("getGrafanaClientInfo", wsl.Err(err))
		return err
	}
	// Если ID пользователя совпадает с зарезервированным, генерируем ошибку
	if gcData["ID"] == id {
		validateError := &ValidateInputDataError{msg: fmt.Sprintf("User id %d reserved", id)}
		return validateError
	}
	return nil
}
