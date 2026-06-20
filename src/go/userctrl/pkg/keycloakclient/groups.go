package keycloakclient

import (
	"context"
	"userctrl/pkg/slogger/wsl"
	"fmt"

	"github.com/Nerzal/gocloak/v13"
)

// CreateGroup создает новую группу и, при необходимости, дочерние группы.
//
// Параметры:
//
//	groupName - название основной группы
//	childGroupsNames - массив названий дочерних групп
//
// Возвращаемые значения:
//
//	string - ID созданной группы
//	error - ошибка при выполнении операции
func (kc *KeycloakClient) CreateGroup(groupName string, childGoupsNames []string) (string, error) {
	log := kc.logger.With(wsl.Label("method", "CreateGroup"))

	group := gocloak.Group{Name: &groupName}
	id, err := kc.client.CreateGroup(kc.ctx, kc.GetToken().AccessToken, kc.realm, group)
	if err != nil {
		log.Error(err.Error())
		return "", err
	}
	for _, childGroupName := range childGoupsNames {
		childGroup := gocloak.Group{Name: gocloak.StringP(childGroupName)}
		childID, err := kc.client.CreateChildGroup(kc.ctx, kc.GetToken().AccessToken, kc.realm, id, childGroup)
		if err != nil {
			log.Error("unable CreateChildGroup", wsl.Err(err))
			return "", err
		}
		log.Debug("CreateGroup: success creating child group", wsl.Info(childID))
	}
	log.Debug("CreateGroup: success creating", wsl.Label("group", id))
	return id, nil
}

// GetGroups получает все группы из Keycloak и возвращает их в виде карты.
//
// Возвращаемые значения:
//
//	map[string]string - карта групп, где ключ - ID группы, значение - название группы
//	error - ошибка при выполнении операции
func (kc *KeycloakClient) GetGroups() (map[string]string, error) {
	log := kc.logger.With(wsl.Label("method", "GetGroups"))

	params := gocloak.GetGroupsParams{}
	groupsData, err := kc.client.GetGroups(kc.ctx, kc.GetToken().AccessToken, kc.realm, params)
	groups := make(map[string]string)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	for _, group := range groupsData {
		groups[*group.ID] = *group.Name
	}
	log.Debug(fmt.Sprintf("success getting %d groups", len(groups)))
	return groups, nil
}

func (kc *KeycloakClient) GetGroupByPath(ctx context.Context, path string) (*gocloak.Group, error) {
	return kc.client.GetGroupByPath(ctx, kc.GetToken().AccessToken, kc.realm, path)
}

// GetGroupByID получает группу по её ID.
//
// Параметры:
//
//	id - ID группы
//
// Возвращаемые значения:
//
//	string - название группы
//	error - ошибка при выполнении операции
func (kc *KeycloakClient) GetGroupByID(id string) (string, error) {
	log := kc.logger.With(wsl.Label("method", "GetGroupByID"))
	group, err := kc.client.GetGroup(kc.ctx, kc.GetToken().AccessToken, kc.realm, id)
	if err != nil {
		log.Error("unable GetGroups", wsl.Err(err))
		return "", err
	}
	log.Debug("success getting", wsl.Label("group", *group.Name))
	return *group.Name, nil
}

// UpdateGroup обновляет название группы.
//
// Параметры:
//
//	id - ID группы
//	groupName - новое название группы
//
// Возвращаемое значение:
//
//	error - ошибка при выполнении операции
func (kc *KeycloakClient) UpdateGroup(id, groupName string) error {
	log := kc.logger.With(wsl.Label("method", "UpdateGroup"))
	updatedGroup := gocloak.Group{ID: &id, Name: &groupName}
	err := kc.client.UpdateGroup(kc.ctx, kc.GetToken().AccessToken, kc.realm, updatedGroup)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	log.Debug("success update", wsl.Label("group", id))
	return nil
}

// DeleteGroupByPath удаляет группу по её пути.
//
// Параметры:
//
//	groupPath - путь к группе
//
// Возвращаемое значение:
//
//	error - ошибка при выполнении операции
func (kc *KeycloakClient) DeleteGroupByPath(groupPath string) error {
	log := kc.logger.With(wsl.Label("method", "DeleteGroupByPath"))

	group, err := kc.client.GetGroupByPath(kc.ctx, kc.GetToken().AccessToken, kc.realm, groupPath)
	if err != nil {
		log.Error("unable GetGroupByPath: %v", wsl.Err(err))
		return err
	}
	err = kc.DeleteGroupByID(*group.ID)
	if err != nil {
		log.Error("unable DeleteGroupByID", wsl.Err(err))
		return err
	}
	return nil
}

// DeleteGroupByID удаляет группу по её ID.
//
// Параметры:
//
//	id - ID группы
//
// Возвращаемое значение:
//
//	error - ошибка при выполнении операции
func (kc *KeycloakClient) DeleteGroupByID(id string) error {
	log := kc.logger.With(wsl.Label("method", "DeleteGroupByID"))

	err := kc.client.DeleteGroup(kc.ctx, kc.GetToken().AccessToken, kc.realm, id)
	if err != nil {
		log.Error("unable DeleteGroup", wsl.Err(err))
		return err
	}

	log.Debug("success deleting group", wsl.Label("id", id))
	return nil
}

// AddRoleToGroup добавляет роль к группе по её ID.
//
// Параметры:
//
//	groupID - ID группы
//	roleName - название роли
//
// Возвращаемое значение:
//
//	error - ошибка при выполнении операции
func (kc *KeycloakClient) AddRoleToGroup(groupID, roleName string) error {
	log := kc.logger.With(wsl.Label("method", "AddRoleToGroup"))

	role, err := kc.client.GetClientRole(kc.ctx, kc.GetToken().AccessToken, kc.realm, kc.clientID, roleName)
	if err != nil {
		log.Error("unable GetClientRole", wsl.Err(err))
		return err
	}
	roles := []gocloak.Role{*role}
	err = kc.client.AddClientRolesToGroup(kc.ctx, kc.GetToken().AccessToken, kc.realm, kc.clientID, groupID, roles)
	if err != nil {
		log.Error("unable AddClientRolesToGroup", wsl.Err(err))
		return err
	}
	log.Debug("success adding role to group", wsl.Label("id", groupID))
	return nil
}

// AddRoleToGroupByPath добавляет роль к группе по её пути.
//
// Параметры:
//
//	groupPath - путь к группе
//	roleName - название роли
//
// Возвращаемое значение:
//
//	error - ошибка при выполнении операции
func (kc *KeycloakClient) AddRoleToGroupByPath(groupPath, roleName string) error {
	log := kc.logger.With(wsl.Label("method", "AddRoleToGroupByPath"))

	group, err := kc.client.GetGroupByPath(kc.ctx, kc.GetToken().AccessToken, kc.realm, groupPath)
	if err != nil {
		log.Error("unable GetGroupByPath", wsl.Err(err))
		return err
	}
	err = kc.AddRoleToGroup(*group.ID, roleName)
	if err != nil {
		log.Error("unable AddRoleToGroup", wsl.Err(err))
		return err
	}
	log.Debug("success adding role to group", wsl.Label("path", groupPath))
	return nil
}

// AddUserToGroup - создает новую группу в Keycloak и, при необходимости, добавляет в нее дочерние группы.
//
// Параметры:
//
//	groupName - название создаваемой группы
//	childGroupsNames - массив названий дочерних групп (может быть пустым)
//
// Возвращаемые значения:
//
//	string - ID созданной группы
//	error - ошибка при выполнении операции
func (kc *KeycloakClient) AddUserToGroup(userID, groupID string) error {
	log := kc.logger.With(wsl.Label("method", "AddUserToGroup"))

	err := kc.client.AddUserToGroup(kc.ctx, kc.GetToken().AccessToken, kc.realm, userID, groupID)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	log.Debug(fmt.Sprintf("AddUserToGroup: success adding user %s to group  %s", userID, groupID))
	return nil
}

// UpdateUserGroups — обновляет группы пользователя в Keycloak: сначала удаляет из всех текущих групп, затем добавляет в указанные новые
func (kc *KeycloakClient) UpdateUserGroups(ctx context.Context, id string, newGroups []string) error {
	log := kc.logger.With(wsl.Label("method", "UpdateUserGroups"))

	token := kc.GetToken().AccessToken
	realm := kc.realm

	// Получить все группы пользователя
	currentGroups, err := kc.client.GetUserGroups(ctx, token, realm, id, gocloak.GetGroupsParams{})
	if err != nil {
		return err
	}

	// Удалить пользователя из всех текущих групп
	for _, group := range currentGroups {
		if group.ID != nil {
			err = kc.client.DeleteUserFromGroup(ctx, token, realm, id, *group.ID)
			if err != nil {
				log.ErrorContext(ctx, err.Error())
			}
		}
	}

	// Добавить пользователя в новые группы
	for _, path := range newGroups {
		g, err := kc.client.GetGroupByPath(ctx, token, realm, path)
		if err != nil {
			return fmt.Errorf("GetGroupByPath: %w", err)
		}
		err = kc.client.AddUserToGroup(ctx, token, realm, id, *g.ID)
		if err != nil {
			return fmt.Errorf("AddUserToGroup: %w", err)
		}

	}

	return nil
}
