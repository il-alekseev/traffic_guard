package keycloakclient

import (
	"userctrl/pkg/slogger/wsl"
	"fmt"
	"strings"

	"github.com/Nerzal/gocloak/v13"
)

// CreateRole создает новую роль в Keycloak.
//
// Параметры:
//
//	name - название создаваемой роли
//
// Возвращаемые значения:
//
//	string - ID созданной роли
//	error - ошибка при выполнении операции
func (kc *KeycloakClient) CreateRole(name string) (string, error) {
	log := kc.logger.With("method", "CreateRole")

	role := gocloak.Role{Name: gocloak.StringP(name)}
	id, err := kc.client.CreateClientRole(kc.ctx, kc.GetToken().AccessToken, kc.realm, kc.clientID, role)
	if err != nil {
		log.Error("unable CreateClientRole", wsl.Err(err))
		return "", err
	}
	log.Debug("success creating", wsl.Label("role", id))
	return id, nil
}

// GetRoles получает все роли из Keycloak и возвращает их в виде карты.
//
// Возвращаемые значения:
//
//	map[string]string - карта ролей, где ключ - ID роли, значение - название роли
//	error - ошибка при выполнении операции
func (kc *KeycloakClient) GetRoles() (map[string]string, error) {
	log := kc.logger.With("method", "GetRoles")

	params := gocloak.GetRoleParams{}
	roles_data, err := kc.client.GetClientRoles(kc.ctx, kc.GetToken().AccessToken, kc.realm, kc.clientID, params)
	if err != nil {
		log.Error("unable CreateClientRole", wsl.Err(err))
		return nil, err
	}
	roles := make(map[string]string)
	for _, role := range roles_data {
		roles[*role.ID] = *role.Name
	}
	log.Debug(fmt.Sprintf("success getting %d roles", len(roles)))
	return roles, nil
}

// GetRoleNameByID получает название роли по её ID.
//
// Параметры:
//
//	id - ID роли
//
// Возвращаемые значения:
//
//	string - название роли
//	error - ошибка при выполнении операции
func (kc *KeycloakClient) GetRoleNameByID(id string) (string, error) {
	log := kc.logger.With("method", "GetRoleNameByID")

	role, err := kc.client.GetClientRoleByID(kc.ctx, kc.GetToken().AccessToken, kc.realm, id)
	if err != nil {
		log.Error("GetClientRoleByID", wsl.Err(err))
		return "", err
	}
	log.Debug("success getting role", wsl.Label("id", *role.Name))
	return *role.Name, nil
}

// GetRoleIDByName получает ID роли по её названию.
//
// Параметры:
//
//	name - название роли
//
// Возвращаемые значения:
//
//	string - ID роли
//	error - ошибка при выполнении операции
func (kc *KeycloakClient) GetRoleIDByName(name string) (string, error) {
	log := kc.logger.With("method", "GetRoleIDByName")

	role, err := kc.client.GetClientRole(kc.ctx, kc.GetToken().AccessToken, kc.realm, kc.clientID, name)
	if err != nil {
		log.Error("unable GetClientRole", wsl.Err(err))
		return "", err
	}
	log.Debug("success getting role", wsl.Label("name", *role.ID))
	return *role.ID, nil
}

// UpdateRole обновляет название роли.
//
// Параметры:
//
//	id - ID роли
//	name - новое название роли
//
// Возвращаемые значения:
//
//	string - обновленное название роли
//	error - ошибка при выполнении операции
func (kc *KeycloakClient) UpdateRole(id, name string) (string, error) {
	log := kc.logger.With("method", "UpdateRole")

	role, err := kc.client.GetClientRoleByID(kc.ctx, kc.GetToken().AccessToken, kc.realm, id)
	if err != nil {
		log.Error("unable GetClientRoleByID", wsl.Err(err))
		return "", err
	}
	role.Name = &name
	err = kc.client.UpdateRole(kc.ctx, kc.GetToken().AccessToken, kc.realm, kc.clientID, *role)
	if err != nil {
		log.Error("UpdateRole", wsl.Err(err))
		return "", err
	}
	log.Debug("success updating role")
	return *role.Name, nil
}

// DeleteRoleByID удаляет роль по её ID.
//
// Параметры:
//
//	id - ID роли
//
// Возвращаемое значение:
//
//	error - ошибка при выполнении операции
func (kc *KeycloakClient) DeleteRoleByID(id string) error {
	log := kc.logger.With("method", "DeleteRoleByID")

	name, err := kc.GetRoleNameByID(id)
	if err != nil {
		log.Error("unable GetRoleNameByID", wsl.Err(err))
		return err
	}
	err = kc.DeleteRoleByName(name)
	if err != nil {
		log.Error("unable DeleteRoleByName", wsl.Err(err))
		return err
	}
	log.Debug("success getting role by id")
	return nil
}

// DeleteRoleByName удаляет роль по её названию.
//
// Параметры:
//
//	name - название роли
//
// Возвращаемое значение:
//
//	error - ошибка при выполнении операции
func (kc *KeycloakClient) DeleteRoleByName(name string) error {
	log := kc.logger.With("method", "DeleteRoleByName")

	err := kc.client.DeleteClientRole(kc.ctx, kc.GetToken().AccessToken, kc.realm, kc.clientID, name)
	if err != nil {
		log.Error("DeleteClientRole", wsl.Err(err))
		return err
	}
	log.Debug("success deleting role by name")
	return nil
}

func (kc *KeycloakClient) getGroupsByRole(role string) *[]string {
	switch {
	case role == "SA":
		return &[]string{"/global-SA"}
	case strings.HasPrefix(role, "CA-"):
		context := strings.TrimPrefix(role, "CA-")
		return &[]string{fmt.Sprintf("/%s/CA", context)}
	case strings.HasPrefix(role, "CO-"):
		context := strings.TrimPrefix(role, "CO-")
		return &[]string{fmt.Sprintf("/%s/CO", context)}
	default:
		return nil
	}
}
