package usecase

import (
	"log/slog"
	"userctrl/config"
	"userctrl/pkg/bizlogger"
	"userctrl/pkg/grafanaclient"
	"userctrl/pkg/grafcookier"
	"userctrl/pkg/keycloakclient"
)

type RepoPGInterface interface {
}

// UseCase — центральный компонент слоя бизнес-логики, координирующий взаимодействие между внешними сервисами:
// Keycloak (управление пользователями и ролями), Clixon (NACM-политики), SSH, Grafana (аутентификация и сессии),
// а также логирование бизнес-событий и техническое логирование
type UseCase struct {
	kc             keycloakclient.KeycloakClientInterface
	grafanaCookier grafcookier.GrafCookierInterface
	grafanaCl      *grafanaclient.GrafanaClient
	blog           bizlogger.LoggerInterface
	l              slog.Logger
}

// New — конструктор UseCase, принимающий конфигурацию и зависимости, необходимые для работы бизнес-логики
func New(
	cfg *config.Config,
	kc keycloakclient.KeycloakClientInterface,
	grafanaCookier grafcookier.GrafCookierInterface,
	grafanaCl *grafanaclient.GrafanaClient,
	blog bizlogger.LoggerInterface,
	l slog.Logger,
) *UseCase {
	uc := UseCase{
		kc:             kc,
		grafanaCookier: grafanaCookier,
		grafanaCl:      grafanaCl,
		blog:           blog,
		l:              l,
	}
	return &uc
}
