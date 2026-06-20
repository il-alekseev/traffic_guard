package keycloakclient

import (
	"context"
	"crypto/tls"
	"userctrl/internal/models"
	"userctrl/pkg/slogger"
	"userctrl/pkg/slogger/wsl"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/Nerzal/gocloak/v13"
)

type KeycloakClientInterface interface {
	// Auth
	Login(ctx context.Context, username string, password string) (*gocloak.JWT, error)
	Logout(ctx context.Context, refreshToken string) error
	RefreshToken(ctx context.Context, refreshToken string) (*gocloak.JWT, error)
	ValidateToken(ctx context.Context, accessToken string) (bool, error)
	GetUserCount(ctx context.Context) (int, error)
	ChangePass(ctx context.Context, userID string, oldPass, newPass string) error
	ResetPass(ctx context.Context, userID string, newPass string) error
	GetGroupByPath(ctx context.Context, path string) (*gocloak.Group, error)
	Close()
	GetMetaUsersInContext(ctx context.Context, contextID string) (*models.MetaDataForContext, error)
	// CRUD users
	CreateUser(u User, groupPath *string, password string) (*User, error)
	// Получение всех пользователей
	GetUsers(ctx context.Context, page, limit int, role, contextID, username, search string) ([]*User, int, error)
	GetAllUsers(ctx context.Context) ([]*User, int, error)
	GetUsersWithId() ([]User, error)
	GetUsersByUsername(ctx context.Context, username string) ([]*User, int, error)
	// Получение пользователя по ID
	GetUserByID(id string) (*User, error)
	// Обновление пользователя
	UpdateUser(ctx context.Context, id string, u User) error
	// Удаление пользователя по ID
	DeleteUserByID(ctx context.Context, id string) error
	// Добавление роли пользователю
	AddRoleToUser(userID, roleName string) error
	// CRUD roles
	// Создание роли
	CreateRole(name string) (string, error)
	// Получение всех ролей
	GetRoles() (map[string]string, error)
	// Получение имени роли по ID
	GetRoleNameByID(id string) (string, error)
	// Получение ID роли по имени
	GetRoleIDByName(name string) (string, error)
	// Обновление роли
	UpdateRole(id, name string) (string, error)
	// Удаление роли по ID
	DeleteRoleByID(id string) error
	// Удаление роли по имени
	DeleteRoleByName(name string) error
	// CRUD groups
	// Создание группы
	CreateGroup(groupName string, childGroupsNames []string) (string, error)
	// Получение списка всех групп
	GetGroups() (map[string]string, error)
	// Получение группы по ID
	GetGroupByID(id string) (string, error)
	// Обновление группы
	UpdateGroup(id, groupName string) error
	// Удаление группы по пути
	DeleteGroupByPath(groupPath string) error
	// Удаление группы по ID
	DeleteGroupByID(id string) error
	// Добавление роли в группу по ID
	AddRoleToGroup(groupID, roleName string) error
	// Добавление роли в группу по пути
	AddRoleToGroupByPath(groupPath, roleName string) error
	// Добавление пользователя в группу
	AddUserToGroup(userID, groupID string) error
}

type KeycloakClientError struct {
	msg string
}

func (e *KeycloakClientError) Error() string {
	return e.msg
}

type User struct {
	ID                     string `json:"user_id" example:"a71032a2-87eb-4d4e-a2a0-d86404672799"`
	Username               string
	Email                  string
	FirstName              string
	LastName               string
	Patronymic             string
	Role                   string
	CreatedAt              string
	IsNeedToChangePassword bool
	IsSuperAdmin           bool
}

type KeycloakClient struct {
	basePath     string
	client       *gocloak.GoCloak
	token        gocloak.JWT
	tokenMutex   sync.RWMutex
	ctx          context.Context
	realm        string // tsum
	systemRealm  string // master
	clientID     string // uuid
	clinetName   string // grafana-sso
	clientSecret string
	logger       slog.Logger
	stopChan     chan struct{}
}

// NewKeycloakClient создает новый клиент для работы с Keycloak.
//
// Параметры:
//
//		url - URL сервера Keycloak
//		login - логин администратора
//		password - пароль администратора
//		realm - название realm в Keycloak
//		clientID - ID клиента в Keycloak
//	 	clientSecret - серкрет клиента
//		l - интерфейс логгера для записи сообщений
//
// Возвращаемое значение:
//
//	*KeycloakClient - указатель на инициализированный клиент
func NewKeycloakClient(url, login, password, systemRealm, realm, clientID, clientUUID, clientSecret string, l slog.Logger) (*KeycloakClient, error) {
	basePath := url
	client := gocloak.NewClient(url)
	ctx := context.Background()
	rastyClient := client.RestyClient()
	// TODO: Добавить использование TLS сертификатов
	// Настройка TLS конфигурации
	tlsConfig := &tls.Config{
		RootCAs:            nil,  // Можно добавить CA сертификаты
		ClientCAs:          nil,  // Можно добавить клиентские CA
		InsecureSkipVerify: true, // НЕ рекомендуется для продакшена
	}
	// Применение конфигурации
	rastyClient.SetTLSClientConfig(tlsConfig)
	l.InfoContext(ctx, "Try connect to keycloak", wsl.Label("url", url))
	fmt.Printf("login=%s, pass=%s, systemRealm=%s", login, password, systemRealm)
	token, err := client.LoginAdmin(ctx, login, password, systemRealm)
	if err != nil {
		l.ErrorContext(ctx, "NewKeycloakClient", wsl.Err(err))
		return nil, slogger.WrapError(ctx, err)
	}
	params := gocloak.GetClientsParams{}
	clients, err := client.GetClients(ctx, token.AccessToken, realm, params)
	if err != nil {
		l.ErrorContext(ctx, "GetClients", wsl.Err(err))
		return nil, slogger.WrapError(ctx, err)
	}
	id := ""
	for _, cl := range clients {
		if *cl.ClientID == clientID {
			id = *cl.ID
			break
		}
	}
	if id == "" {
		l.ErrorContext(ctx, fmt.Sprintf("GetClients: client {%s} not found", clientID))
		return nil, slogger.WrapError(ctx, err)
	}

	kc := KeycloakClient{
		basePath:     basePath,
		client:       client,
		token:        *token,
		ctx:          ctx,
		systemRealm:  systemRealm,
		realm:        realm,
		clientSecret: clientSecret,
		clientID:     id,
		clinetName:   clientID,
		logger:       l,
		tokenMutex:   sync.RWMutex{},
		stopChan:     make(chan struct{}),
	}

	// Start the token refresh goroutine
	go kc.startTokenRefresh(1*time.Minute, login, password)

	return &kc, nil
}

// startTokenRefresh starts a goroutine that refreshes the token at regular intervals.
func (kc *KeycloakClient) startTokenRefresh(refreshInterval time.Duration, login, password string) {
	log := kc.logger.With(wsl.Label("method", "startTokenRefresh"))
	ticker := time.NewTicker(refreshInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			kc.refreshToken(login, password)
		case <-kc.stopChan:
			log.Info("Stopping token refresh goroutine")
			return
		}
	}
}

func (kc *KeycloakClient) refreshToken(login, password string) {
	log := kc.logger.With(wsl.Label("method", "refreshToken"))

	kc.tokenMutex.RLock()
	refreshToken := kc.token.RefreshToken
	kc.tokenMutex.RUnlock()

	if refreshToken == "" {
		log.Error("No refresh token available")
		return
	}
	newToken, err := kc.client.RefreshToken(kc.ctx, refreshToken, "admin-cli", "", kc.systemRealm)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to refresh token: %v", err))
		// TODO: try to re-login if refresh failed
		timeout := 2 * time.Second
		ctx, _ := context.WithTimeout(context.Background(), timeout)

		newToken, err = kc.client.LoginAdmin(ctx, login, password, kc.systemRealm)
		kc.tokenMutex.Lock()
		kc.token = *newToken
		kc.tokenMutex.Unlock()
		if err != nil {
			log.Error("Failed to re-login if", wsl.Err(err))
		}
		log.Debug("Re-login successfully")
		return
	}

	kc.tokenMutex.Lock()
	kc.token = *newToken
	kc.tokenMutex.Unlock()

	log.Debug("Token refreshed successfully")
}

func (kc *KeycloakClient) IsTokenValid() bool {
	kc.tokenMutex.RLock()
	defer kc.tokenMutex.RUnlock()

	if kc.token.AccessToken == "" {
		return false
	}

	// Simple check - you could also parse JWT to check expiry
	return time.Now().Before(time.Unix(int64(kc.token.ExpiresIn), 0))
}

func (kc *KeycloakClient) GetToken() gocloak.JWT {
	kc.tokenMutex.RLock()
	defer kc.tokenMutex.RUnlock()
	return kc.token
}

// Close stops the token refresh goroutine and cleans up resources.
func (kc *KeycloakClient) Close() {
	close(kc.stopChan)
}
