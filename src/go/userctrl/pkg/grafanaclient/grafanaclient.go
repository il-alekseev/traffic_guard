package grafanaclient

import (
	"userctrl/pkg/slogger/wsl"
	"log/slog"

	"github.com/go-openapi/strfmt"
	goapi "github.com/grafana/grafana-openapi-client-go/client"
)

// GrafanaClientInterface — интерфейс для управления пользователями в Grafana: добавление администратора, удаление пользователей
type GrafanaClientInterface interface {
	AddAdminUserToAllOrgs(login string) error
	DeleteUsers() error
	DeleteUserById(username int64) error
}

// GrafanaClient — реализация клиента с поддержкой HTTP-запросов к API Grafana
type GrafanaClient struct {
	client   goapi.GrafanaHTTPAPI
	config   goapi.TransportConfig
	provider string
	logger   slog.Logger
}

// NewGrafanaClient — конструктор, создающий экземпляр клиента с указанной конфигурацией подключения и провайдером
func NewGrafanaClient(config goapi.TransportConfig, provider string, l slog.Logger) *GrafanaClient {
	client := goapi.NewHTTPClientWithConfig(strfmt.Default, &config)
	return &GrafanaClient{
		config:   config,
		client:   *client,
		provider: provider,
		logger:   l,
	}
}

// getGrafanaClientInfo — получает информацию о текущем авторизованном пользователе (логин, ID, email, организация)
// и возвращает её в виде карты; при ошибках логирует их и возвращает nil и ошибку
func (gc *GrafanaClient) getGrafanaClientInfo() (map[string]any, error) {
	log := gc.logger.With(wsl.Label("method", "getGrafanaClientInfo"))
	// example:
	// login -> admin
	// userID -> 1
	// email -> admin@localhost
	// orgName -> Main Org.
	// orgNameID -> 1
	api := gc.client.SignedInUser
	response, err := api.GetSignedInUser()
	if err != nil {
		log.Error("unable GetSignedInUser", wsl.Err(err))
		return nil, err
	}
	user_info := make(map[string]any)
	user_info["Login"] = response.Payload.Login
	user_info["ID"] = response.Payload.ID
	user_info["Email"] = response.Payload.Email
	user_info["OrgID"] = response.Payload.OrgID
	resp, err := gc.client.Org.GetCurrentOrg()
	if err != nil {
		log.Error("unable GetCurrentOrg", wsl.Err(err))
		return nil, err
	}
	user_info["OrgName"] = resp.Payload.Name
	return user_info, nil
}
