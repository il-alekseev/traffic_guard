package grafanaclient

import (
	"userctrl/pkg/slogger/wsl"
	"fmt"
	"sync"

	"github.com/grafana/grafana-openapi-client-go/models"
)

// GetUserByLogin получает ID пользователя по логину
//
// Параметры:
//
//	login - логин или email пользователя для поиска
//
// Возвращает:
//
//	int64 - ID пользователя
//	error - ошибку при получении данных (nil если успешно)
func (gc *GrafanaClient) GetUserByLogin(login string) (int64, error) {
	log := gc.logger.With(wsl.Label("method", "GetUserByLogin"))
	// Если пользователь пытается получить информацию о сервисном аккаунте, ничего не выводим
	err := gc.validateUserLogin(login)
	if err != nil {
		err = fmt.Errorf("warn get reserved login: %w", err)
		log.Warn(err.Error())
		return 0, err
	}

	api := gc.client.Users
	response, err := api.GetUserByLoginOrEmail(login)
	if err != nil {
		log.Error("Error getting user by loginOrEmail", wsl.Err(err))
		return 0, err
	}
	return response.Payload.ID, nil
}

// DeleteUserById удаляет пользователя из Grafana по его ID.
//
// Параметры:
//
//	id int64 - идентификатор удаляемого пользователя
//
// Возвращаемое значение:
//
//	error - ошибка при удалении (nil если успешно)
func (gc *GrafanaClient) DeleteUserById(userID int64) error {
	log := gc.logger.With(wsl.Label("method", "DeleteUserById"))

	err := gc.validateUserID(userID)
	if err != nil {
		log.Error("unable DeleteUserById: error delete reserved ID")
		return err
	}
	api := gc.client.AdminUsers
	_, err = api.AdminDeleteUser(userID)
	if err != nil {
		log.Error("Error deleting user by ID", wsl.Err(err))
		return err
	}
	log.Debug("Success deleting user")
	return nil
}

func (gc *GrafanaClient) LogoutUser(userID int64) error {
	log := gc.logger.With(wsl.Label("method", "LogoutUser"))

	if _, err := gc.client.AdminUsers.AdminLogoutUser(userID); err != nil {
		err = fmt.Errorf("gc.client.AdminUsers.AdminLogoutUser, userID=%d: %w", userID, err)
		log.Error(err.Error())
		return err
	}
	return nil
}

func (gc *GrafanaClient) AddAdminUserToAllOrgs(login string) error {
	log := gc.logger.With(wsl.Label("method", "AddAdminUserToAllOrgs"))
	orgs, err := gc.GetAllOrganizations()
	if err != nil {
		return fmt.Errorf("gc.GetAllOrganizations: %w", err)
	}

	wg := sync.WaitGroup{}
	wg.Add(len(orgs))
	for _, org := range orgs {
		go func() {
			defer wg.Done()
			if org.ID == 1 {
				return
			}
			params := &models.AddOrgUserCommand{
				LoginOrEmail: login,
				Role:         "Admin",
			}

			// Отправляем запрос
			_, err := gc.client.Orgs.AddOrgUser(org.ID, params)
			if err != nil {
				log.Error("failed to add user to organization", wsl.Err(err))
			}
		}()
	}

	wg.Wait()

	return nil
}
