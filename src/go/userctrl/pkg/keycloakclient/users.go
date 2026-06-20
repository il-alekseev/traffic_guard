package keycloakclient

import (
	"context"
	"errors"
	"userctrl/internal/models"
	"userctrl/pkg/slogger"
	"userctrl/pkg/slogger/wsl"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/Nerzal/gocloak/v13"
)

// convertKeyCloakUserToUser — преобразует объект пользователя из Keycloak (gocloak.User) в локальную модель User
// Извлекает отчество из атрибутов пользователя, если доступно
func (kc *KeycloakClient) convertKeyCloakUserToUser(user *gocloak.User, clientRole string) *User {
	patronymic := ""
	if user.Attributes != nil {
		data, ok := (*user.Attributes)["patronymic"]
		if ok {
			patronymic = data[0]
		} else {
			kc.logger.Warn("Attribute {patronymic} not found")
		}
	}
	isNeedToChangePassword := false
	if len(*user.RequiredActions) > 0 && (*user.RequiredActions)[0] == "UPDATE_PASSWORD" {
		isNeedToChangePassword = true
	}

	isSuperAdmin := false
	if clientRole == "SA" {
		isSuperAdmin = true
	}
	u := User{
		ID:                     *user.ID,
		Username:               gocloak.PString(user.Username),
		Email:                  gocloak.PString(user.Email),
		FirstName:              gocloak.PString(user.FirstName),
		LastName:               gocloak.PString(user.LastName),
		CreatedAt:              strconv.Itoa(int(gocloak.PInt64(user.CreatedTimestamp))),
		Patronymic:             patronymic,
		IsNeedToChangePassword: isNeedToChangePassword,
		Role:                   clientRole,
		IsSuperAdmin:           isSuperAdmin,
	}
	return &u
}

// checkAttributeExists проверяет существование указанного атрибута в профиле пользователя Keycloak.
//
// Параметры:
//
//	attrName - имя атрибута, существование которого необходимо проверить
//
// Возвращаемые значения:
//
//	bool - true, если атрибут существует, false - если не существует
//	error - ошибка, если произошла ошибка при выполнении запроса
func (kc *KeycloakClient) checkAttributeExists(attrName string) (bool, error) {
	log := kc.logger.With(wsl.Label("method", "checkAttributeExists"))
	// Создание HTTP запроса
	// "https://localhost:8443/admin/realms/tsum/users/profile
	url := fmt.Sprintf("%s/admin/realms/%s/users/profile", kc.basePath, kc.realm)
	req, err := http.NewRequestWithContext(kc.ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Error("NewRequestWithContext", wsl.Err(err))
		return false, err
	}
	// Добавление токена авторизации
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", kc.GetToken().AccessToken))
	req.Header.Set("Content-Type", "application/json")
	// Отправка запроса
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("client.Do", wsl.Err(err))
		return false, err
	}
	defer resp.Body.Close()
	// Проверка статуса ответа
	if resp.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("Request Error: %s", resp.Status)
		log.Error(msg)
		return false, &KeycloakClientError{msg: msg}
	}
	// Чтение тела ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("io.ReadAll", wsl.Err(err))
		return false, err
	}
	stringBody := string(body)
	log.Debug("Server resp", wsl.Info(stringBody))
	// Проверка на наличие атрибута
	if strings.Contains(stringBody, attrName) {
		return true, nil
	} else {
		return false, &KeycloakClientError{msg: fmt.Sprintf("Atribute {%s} not found", attrName)}
	}
}

// CreateUser создает нового пользователя в Keycloak.
//
// Параметры:
//
//	u - структура с данными пользователя (User)
//	groupPath - указатель на путь к группе (*string)
//	password - пароль пользователя (string)
//
// Возвращаемые значения:
//
//	string - ID созданного пользователя
//	error - ошибка при выполнении операции
func (kc *KeycloakClient) CreateUser(u User, groupPath *string, password string) (*User, error) {
	log := kc.logger.With(wsl.Label("method", "CreateUser"))
	// _, err := kc.checkAttributeExists("patronymic")
	// if err != nil {
	// 	kc.logger.Error("checkAttributeExists: %v", err)
	// 	return "", err
	// }
	credentials := []gocloak.CredentialRepresentation{
		{Type: gocloak.StringP("password"), Value: &password, Temporary: gocloak.BoolP(true)},
	}
	attributes := make(map[string][]string)
	if u.Patronymic != "" {
		attributes = map[string][]string{
			"patronymic": {u.Patronymic},
		}
	}

	groups := []string{}
	if groupPath != nil {
		groups = []string{*groupPath}
	}
	user := gocloak.User{
		Username:    &u.Username,
		FirstName:   &u.FirstName,
		LastName:    &u.LastName,
		Enabled:     gocloak.BoolP(true),
		Email:       &u.Email,
		Credentials: &credentials,
		Groups:      &groups,
		Attributes:  &attributes,
	}
	id, err := kc.client.CreateUser(kc.ctx, kc.GetToken().AccessToken, kc.realm, user)
	if err != nil {
		log.Error("unable CreateUser", wsl.Err(err))
		return nil, err
	}
	// Добавляем роль пользователю
	if u.Role != "" {
		err = kc.AddRoleToUser(id, u.Role)
		if err != nil {
			log.Error("unable AddRoleToUser", wsl.Err(err))
			return nil, err
		}
	}
	log.Debug("success creating", wsl.Label("user", id))
	res, err := kc.GetUserByID(id)
	if err != nil {
		return nil, fmt.Errorf("kc.GetUserByID: %w", err)
	}
	return res, nil
}

func (kc *KeycloakClient) getUsersWithParams(ctx context.Context, params gocloak.GetUsersParams) ([]*User, int, error) {
	log := kc.logger.With(wsl.Label("method", "getUsersWithParams"))

	data, err := kc.client.GetUsers(ctx, kc.GetToken().AccessToken, kc.realm, params)
	if err != nil {
		log.ErrorContext(ctx, "unable GetUsers", wsl.Err(err))
		return nil, 0, slogger.WrapError(ctx, err)
	}
	var users []*User
	for _, user := range data {
		u, err := kc.GetUserByID(*user.ID)
		if err != nil {
			log.ErrorContext(ctx, "unable GetUserByID", wsl.Err(err))
			return nil, 0, slogger.WrapError(ctx, err)
		}
		u.ID = *user.ID
		users = append(users, u)
	}
	log.DebugContext(ctx, "success getting users", wsl.Int("value", len(users)))
	total, err := kc.client.GetUserCount(ctx, kc.GetToken().AccessToken, kc.realm, params)
	if err != nil {
		err = fmt.Errorf("kc.client.GetUserCount: %w", err)
		log.DebugContext(ctx, err.Error())
		return nil, 0, slogger.WrapError(ctx, err)
	}
	return users, total, nil
}

// GetUsers получает всех пользователей из Keycloak и возвращает их в виде карты.
//
// Возвращаемые значения:
//
//	map[string]User - карта пользователей, где ключ - ID пользователя
//	int64 - всего пользователей по данным фильтрам
//	error - ошибка при выполнении операции
func (kc *KeycloakClient) GetUsers(ctx context.Context, page, limit int, role, contextID, username, search string) ([]*User, int, error) {
	// role: Enums(SA, CA, CO, "")
	if role != "" || contextID != "" {
		// нет стандартного поиска по ролям - это решение
		return kc.getUsersWithRoleFlter(ctx, page, limit, role, contextID, search)
	}
	first := page * limit
	params := gocloak.GetUsersParams{
		First: &first,
		Max:   &limit,
	}
	if username != "" {
		params.Username = &username
	}
	if search != "" {
		params.Search = &search
	}
	return kc.getUsersWithParams(ctx, params)
}

func searchUser(u *User, search string) bool {
	if search == "" {
		return true
	}
	search = strings.ToLower(search)
	strs := []string{u.FirstName, u.LastName, u.Patronymic, u.Username, u.Email}
	for _, s := range strs {
		if strings.Contains(strings.ToLower(s), search) {
			return true
		}
	}
	return false
}

func (kc *KeycloakClient) getUsersWithClientRoleFlter(ctx context.Context, clientRole string) ([]*User, error) {
	log := kc.logger.With(wsl.Label("method", "getUsersWithClientRole"))

	groups, err := kc.client.GetGroupsByClientRole(ctx, kc.GetToken().AccessToken, kc.realm, clientRole, kc.clientID)
	if err != nil {
		err = fmt.Errorf("GetGroupsByRole: %w", err)
		return nil, slogger.WrapError(ctx, err)
	}
	if len(groups) == 0 {
		err = errors.New("GetGroupsByRole: groups not found")
		return nil, slogger.WrapError(ctx, err)
	}

	log.DebugContext(ctx, fmt.Sprintf("GetGroupsByRole: role = %s, groups = %+v", clientRole, groups))
	g := groups[0]
	param := gocloak.GetGroupsParams{}
	data, err := kc.client.GetGroupMembers(ctx, kc.GetToken().AccessToken, kc.realm, *g.ID, param)
	if err != nil {
		err = fmt.Errorf("unable GetUsersByClientRoleName: %w", err)
		log.ErrorContext(ctx, err.Error())
		return nil, slogger.WrapError(ctx, err)
	}
	log.DebugContext(ctx, fmt.Sprintf("GetGroupsByRole: group path = %s, users = %+v", *g.Path, data))
	users := make([]*User, 0, len(data))
	for _, user := range data {
		u, err := kc.GetUserByID(*user.ID)
		if err != nil {
			log.ErrorContext(ctx, err.Error())
			return nil, slogger.WrapError(ctx, err)
		}
		u.ID = *user.ID
		users = append(users, u)
	}

	return users, nil
}

func (kc *KeycloakClient) getUsersWithShortRoleFlter(ctx context.Context, shortRole string) ([]*User, error) {
	log := kc.logger.With(wsl.Label("method", "getUsersWithShortRole"))

	// shortRole - SA, CA, CO
	if shortRole == "SA" {
		return kc.getUsersWithClientRoleFlter(ctx, shortRole)
	} else if shortRole == "CA" || shortRole == "CO" {
		res := make([]*User, 0)
		search := fmt.Sprintf("%s-", shortRole)
		params := gocloak.GetRoleParams{
			Search: &search,
		}
		roles, err := kc.client.GetClientRoles(ctx, kc.GetToken().AccessToken, kc.realm, kc.clientID, params)
		if err != nil {
			err = fmt.Errorf("kc.client.GetClientRoles, %s: %w", shortRole, err)
			return nil, slogger.WrapError(ctx, err)
		}
		for _, r := range roles {
			log.DebugContext(ctx, "get users", wsl.Label("role", *r.Name))
			if strings.HasPrefix(*r.Name, shortRole+"-") {
				groups := kc.getGroupsByRole(*r.Name)
				if len(*groups) < 1 {
					continue
				}
				paths := *groups
				group, err := kc.client.GetGroupByPath(ctx, kc.GetToken().AccessToken, kc.realm, paths[0])
				if err != nil {
					err = fmt.Errorf("GetGroupByPath %s: %w", paths[0], err)
					log.ErrorContext(ctx, err.Error())
					continue
				}
				users, err := kc.client.GetGroupMembers(ctx, kc.GetToken().AccessToken, kc.realm, *group.ID, gocloak.GetGroupsParams{})
				if err != nil {
					err = fmt.Errorf("kc.client.GetUsersByClientRoleName, role=%s: %w", *r.Name, err)
					log.ErrorContext(ctx, err.Error())
				} else {
					addUsers := make([]*User, 0, len(users))
					for _, u := range users {
						addUsers = append(addUsers, kc.convertKeyCloakUserToUser(u, *r.Name))
					}
					res = append(res, addUsers...)
				}
			}
		}
		return res, nil
	} else {
		return nil, slogger.WrapError(ctx, fmt.Errorf("%w : shortRole=%s", ErrInvalidRole, shortRole))
	}
}

func (kc *KeycloakClient) getUsersWithContextIDFlter(ctx context.Context, contextID string) ([]*User, error) {
	log := kc.logger.With(wsl.Label("method", "getUsersWithContextID"))

	// shortRole - SA, CA, CO
	pathCA := fmt.Sprintf("/%s/CA", contextID)
	pathCO := fmt.Sprintf("/%s/CO", contextID)
	groupCA, err := kc.client.GetGroupByPath(ctx, kc.GetToken().AccessToken, kc.realm, pathCA)
	if err != nil {
		return nil, slogger.WrapError(ctx, fmt.Errorf("GetGroupByPath %s: %w", pathCA, err))
	}
	groupCO, err := kc.client.GetGroupByPath(ctx, kc.GetToken().AccessToken, kc.realm, pathCO)
	if err != nil {
		return nil, slogger.WrapError(ctx, fmt.Errorf("GetGroupByPath %s: %w", pathCO, err))
	}
	param := gocloak.GetGroupsParams{}
	usersCA, err := kc.client.GetGroupMembers(ctx, kc.GetToken().AccessToken, kc.realm, *groupCA.ID, param)
	if err != nil {
		return nil, slogger.WrapError(ctx, fmt.Errorf("GetGroupMembers %s: %w", pathCA, err))
	}
	usersCO, err := kc.client.GetGroupMembers(ctx, kc.GetToken().AccessToken, kc.realm, *groupCO.ID, param)
	if err != nil {
		return nil, slogger.WrapError(ctx, fmt.Errorf("GetGroupMembers %s: %w", pathCO, err))
	}
	users := make([]*User, 0, len(usersCA)+len(usersCO))
	for _, user := range usersCA {
		u, err := kc.GetUserByID(*user.ID)
		if err != nil {
			log.ErrorContext(ctx, "unable GetUserByID", wsl.Err(err))
			return nil, slogger.WrapError(ctx, err)
		}
		u.ID = *user.ID
		users = append(users, u)
	}
	for _, user := range usersCO {
		u, err := kc.GetUserByID(*user.ID)
		if err != nil {
			kc.logger.ErrorContext(ctx, "unable GetUserByID", wsl.Err(err))
			return nil, slogger.WrapError(ctx, err)
		}
		u.ID = *user.ID
		users = append(users, u)
	}
	// TODO: по идее ещё и отсортировать по username
	return users, nil
}

func (kc *KeycloakClient) getUsersWithRoleFlter(ctx context.Context, page, limit int, shortRole, contextID, search string) ([]*User, int, error) {
	var users []*User
	var err error = nil
	if shortRole != "" && contextID != "" { //
		clientRole := shortRole + "-" + contextID
		if shortRole == "SA" {
			clientRole = "SA"
		}
		users, err = kc.getUsersWithClientRoleFlter(ctx, clientRole)
		if err != nil {
			return nil, 0, fmt.Errorf("getUsersWithClientRoleFlter: %w", err)
		}
	} else if shortRole != "" {
		if shortRole == "SA" {
			users, err = kc.getUsersWithClientRoleFlter(ctx, shortRole)
			if err != nil {
				return nil, 0, fmt.Errorf("getUsersWithClientRoleFlter: %w", err)
			}
		} else {
			users, err = kc.getUsersWithShortRoleFlter(ctx, shortRole) // shortRole = CA or CO
			if err != nil {
				return nil, 0, fmt.Errorf("getUsersWithShortRoleFlter: %w", err)
			}
		}
	} else {
		// contextID != "" && shortRole == ""
		users, err = kc.getUsersWithContextIDFlter(ctx, contextID)
		if err != nil {
			return nil, 0, fmt.Errorf("getUsersWithClientRoleFlter: %w", err)
		}
	}

	var res []*User
	for _, u := range users {
		if searchUser(u, search) {
			res = append(res, u)
		}
	}
	total := len(res)
	first := page * limit
	end := first + limit
	if first > len(res)-1 {
		return nil, 0, nil
	}
	if end >= len(res) {
		return res[first:], total, nil
	}
	return res[first:end], total, nil
}

func (kc *KeycloakClient) GetAllUsers(ctx context.Context) ([]*User, int, error) {
	params := gocloak.GetUsersParams{}
	return kc.getUsersWithParams(ctx, params)
}

func (kc *KeycloakClient) GetMetaUsersInContext(ctx context.Context, contextID string) (*models.MetaDataForContext, error) {
	meta := models.MetaDataForContext{
		ContextID: contextID,
	}
	ca := fmt.Sprintf("CA-%s", contextID)
	co := fmt.Sprintf("CO-%s", contextID)
	usersCA, err := kc.getUsersWithClientRoleFlter(ctx, ca)
	if err != nil {
		return nil, fmt.Errorf("getUsersWithClientRoleFlter - role=%s: %w", ca, err)
	}
	usersCO, err := kc.getUsersWithClientRoleFlter(ctx, co)
	if err != nil {
		return nil, fmt.Errorf("getUsersWithClientRoleFlter - role=%s: %w", co, err)
	}
	meta.CountCA = len(usersCA)
	meta.CountCO = len(usersCO)
	return &meta, nil
}

func (kc *KeycloakClient) GetUsersByUsername(ctx context.Context, username string) ([]*User, int, error) {
	params := gocloak.GetUsersParams{
		Username: &username,
	}
	return kc.getUsersWithParams(ctx, params)
}

func (kc *KeycloakClient) GetUsersWithId() ([]User, error) { //
	log := kc.logger.With(wsl.Label("method", "getUsersWithId"))

	params := gocloak.GetUsersParams{}
	data, err := kc.client.GetUsers(kc.ctx, kc.GetToken().AccessToken, kc.realm, params)
	if err != nil {
		log.Error("unable GetUsers", wsl.Err(err))
		return nil, err
	}
	users := make([]User, 0, len(data))
	for _, user := range data {
		u, err := kc.GetUserByID(*user.ID)
		if err != nil {
			log.Error("unable GetUserByID", wsl.Err(err))
			return nil, err
		}
		users = append(users, *u)
	}
	log.Debug("success getting", wsl.Int("number", len(users)))
	return users, nil
}

// GetUserByID получает пользователя по его ID.
//
// Параметры:
//
//	id - ID пользователя
//
// Возвращаемые значения:
//
//	*User - указатель на структуру пользователя
//	error - ошибка при выполнении операции
func (kc *KeycloakClient) GetUserByID(id string) (*User, error) {
	log := kc.logger.With(wsl.Label("method", "getUserByID"))

	// Находим основные параметры пользователя
	user, err := kc.client.GetUserByID(kc.ctx, kc.GetToken().AccessToken, kc.realm, id)
	if err != nil {
		log.Error("unable GetUserByID", wsl.Err(err))
		return nil, err
	}
	// Находим отчество пользователя
	patronymic := ""
	if user.Attributes != nil {
		data, ok := (*user.Attributes)["patronymic"]
		if ok {
			patronymic = data[0]
		} else {
			log.Warn("Attribute {patronymic} not found")
		}
	}
	// Находим, нужно ли обновить пароль пользователя
	isNeedToChangePassword := false
	if len(*user.RequiredActions) > 0 && (*user.RequiredActions)[0] == "UPDATE_PASSWORD" {
		isNeedToChangePassword = true
	}

	// Находим роль пользователя
	//roles, err := kc.client.GetClientRolesByUserID(kc.ctx, kc.GetToken().AccessToken, kc.realm, kc.clientID, id)
	roles, err := kc.client.GetCompositeClientRolesByUserID(kc.ctx, kc.GetToken().AccessToken, kc.realm, kc.clientID, id) // Илья
	if err != nil {
		log.Error("unable GetClientRolesByUserID", wsl.Err(err))
		return nil, err
	}
	role := ""
	if len(roles) > 0 {
		role = *roles[0].Name
	}
	// Определяем, является ли пользователь суперадмином
	isSuperAdmin := false
	if role == "SA" {
		isSuperAdmin = true
	}
	u := User{
		ID:                     id,
		Username:               gocloak.PString(user.Username),
		Email:                  gocloak.PString(user.Email),
		FirstName:              gocloak.PString(user.FirstName),
		LastName:               gocloak.PString(user.LastName),
		CreatedAt:              strconv.Itoa(int(gocloak.PInt64(user.CreatedTimestamp))),
		Patronymic:             patronymic,
		IsNeedToChangePassword: isNeedToChangePassword,
		Role:                   role,
		IsSuperAdmin:           isSuperAdmin,
	}
	return &u, nil
}

// UpdateUser обновляет информацию о пользователе.
//
// Параметры:
//
//	id - ID пользователя (string)
//	u - структура с обновленными данными пользователя (User)
//
// Возвращаемое значение:
//
//	error - ошибка при выполнении операции
func (kc *KeycloakClient) UpdateUser(ctx context.Context, userId string, u User) error {
	log := kc.logger.With(wsl.Label("method", "updateUser"))

	attributes := map[string][]string{
		"patronymic": {u.Patronymic},
	}

	user := gocloak.User{
		ID:         &userId,
		FirstName:  &u.FirstName,
		LastName:   &u.LastName,
		Email:      &u.Email,
		Attributes: &attributes,
		Groups:     kc.getGroupsByRole(u.Role),
	}

	if u.Role != "" {
		log.DebugContext(ctx, fmt.Sprintf("update use groups: %+v", *user.Groups))
	}

	err := kc.client.UpdateUser(ctx, kc.GetToken().AccessToken, kc.realm, user)
	if err != nil {
		log.ErrorContext(ctx, "unable UpdateUser", wsl.Err(err))
		return err
	}
	groups := kc.getGroupsByRole(u.Role)
	if groups != nil {
		if err = kc.UpdateUserGroups(ctx, userId, *groups); err != nil {
			log.ErrorContext(ctx, "unable UpdateUserGroups", wsl.Err(err))
			return err
		}
	}

	log.DebugContext(ctx, "success updating user")
	return nil
}

// DeleteUserByID удаляет пользователя по его ID.
//
// Параметры:
//
//	id - ID пользователя
//
// Возвращаемое значение:
//
//	error - ошибка при выполнении операции
func (kc *KeycloakClient) DeleteUserByID(ctx context.Context, id string) error {
	log := kc.logger.With(wsl.Label("method", "deleteUserByID"))

	err := kc.client.DeleteUser(ctx, kc.GetToken().AccessToken, kc.realm, id)
	if err != nil {
		log.ErrorContext(ctx, "DeleteUser", wsl.Err(err))
		return err
	}
	log.DebugContext(ctx, "success deleting user")
	return nil
}

// AddRoleToUser добавляет роль пользователю.
//
// Параметры:
//
//	userID - ID пользователя
//	roleName - название роли
//
// Возвращаемое значение:
//
//	error - ошибка при выполнении операции
func (kc *KeycloakClient) AddRoleToUser(userID, roleName string) error {
	log := kc.logger.With(wsl.Label("method", "addRoleToUser"))

	role, err := kc.client.GetClientRole(kc.ctx, kc.GetToken().AccessToken, kc.realm, kc.clientID, roleName)
	if err != nil {
		log.Error("unable GetClientRole", wsl.Err(err))
		return err
	}
	roles := []gocloak.Role{*role}
	err = kc.client.AddClientRolesToUser(kc.ctx, kc.GetToken().AccessToken, kc.realm, kc.clientID, userID, roles)
	if err != nil {
		log.Error("unable AddClientRolesToUser", wsl.Err(err))
		return err
	}
	log.Debug("AddRoleToUser: success adding role to user", wsl.Label("id", userID))
	return nil
}

func (kc *KeycloakClient) GetUserCount(ctx context.Context) (int, error) {
	return kc.client.GetUserCount(ctx, kc.GetToken().AccessToken, kc.realm, gocloak.GetUsersParams{})
}

func (kc *KeycloakClient) ChangePass(ctx context.Context, userID string, oldPass, newPass string) error {
	user, err := kc.GetUserByID(userID)
	if err != nil {
		return fmt.Errorf("kc.GetUserByID, userID=%s: %w", userID, err)
	}

	if kc.ValidationPassword(ctx, user.Username, oldPass) {
		return kc.client.SetPassword(ctx, kc.GetToken().AccessToken, userID, kc.realm, newPass, false)
	}
	return ErrInvalidUserCredentials
}

func (kc *KeycloakClient) setСonstantPassword(ctx context.Context, userID string, pass string) error {
	return kc.client.SetPassword(ctx, kc.GetToken().AccessToken, userID, kc.realm, pass, false)
}

func (kc *KeycloakClient) ResetPass(ctx context.Context, userID string, newPass string) error {
	return kc.client.SetPassword(ctx, kc.GetToken().AccessToken, userID, kc.realm, newPass, true)
}
