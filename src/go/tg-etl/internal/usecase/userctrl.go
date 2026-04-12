package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"tg-etl/pkg/usercontrol/roles"

	"github.com/go-openapi/runtime"
)

// hasDeviceInRoles проверяет, есть ли имя узла в списке ролей
func (uc *UseCase) hasDeviceInRoles(authInfo *runtime.ClientAuthInfoWriterFunc, hostname string) (bool, error) {
	params := roles.GetV1RolesParams{}
	roles, err := uc.userCl.Roles.GetV1Roles(&params, authInfo)
	if err != nil {
		return false, err
	}

	if roles == nil || roles.Payload == nil || roles.Payload.Data == nil {
		return false, nil
	}

	for _, role := range roles.Payload.Data {
		if role == nil {
			continue
		}
		if role.Role != nil {
			if strings.Contains(*role.Role, hostname) {
				return true, nil
			}
		}
	}
	return false, nil
}

func (uc *UseCase) addRolesToKeyCloak(
	ctx context.Context,
	authInfo *runtime.ClientAuthInfoWriterFunc,
	hostname string,
) error {
	// Создание групп и ролей в keycloak
	userClReq, err := uc.userCl.Roles.PutV1RolesContext(
		&roles.PutV1RolesContextParams{
			ContextID: &hostname,
			Context:   ctx,
		},
		authInfo,
	)
	if err != nil {
		return fmt.Errorf("user control client -PutV1RolesContext : %w", err)
	}
	if userClReq.Code() != http.StatusOK {
		err = fmt.Errorf("user control client - PutV1RolesContext- code %d : %s",
			userClReq.Code(),
			userClReq.GetPayload().Message,
		)
		return err
	}
	return nil
}

// syncRolesForDevices - проверяет наличие ролей для устройств, лежащих в БД и в случае их отсутствия - создает их
func (uc *UseCase) syncRolesForDevices(ctx context.Context) error {
	method := "checkRolesForDevices"
	// Получаем авторизацию для usercontrol
	authInfo, err := uc.getServiceAuthInfo(ctx)
	if err != nil {
		return fmt.Errorf("failed to get auth: %w", err)
	}
	// Получаем список всех устройств в БД
	devs, err := uc.q.GetDevices(ctx)
	if err != nil {
		return fmt.Errorf("%s: failed to get devices: %w", method, err)
	}
	for _, d := range devs {
		// Проверяем, создавалась ли уже роль для данного сетевого узла
		roleExists, err := uc.hasDeviceInRoles(&authInfo, d.HostName)
		if err != nil {
			return fmt.Errorf("failed to get roles for device: %w", err)
		}
		// Если не создавалась, то создаем
		if !roleExists {
			if err = uc.addRolesToKeyCloak(ctx, &authInfo, d.HostName); err != nil {
				return fmt.Errorf("failed to create role for device: %w", err)
			}
		}
		uc.l.InfoContext(ctx, "created new device", slog.String("device name", d.HostName))

	}
	return nil
}
