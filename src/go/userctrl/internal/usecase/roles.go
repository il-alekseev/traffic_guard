package usecase

import (
	"context"
	"fmt"
	"userctrl/internal/models"
)

// СreateRolesAndGroupsForContext — создаёт группу в Keycloak для указанного контекста (ctxID) и подгруппы CA/CO,
// назначает соответствующие роли, добавляет контекст в конфигурацию NASM через CLI и сохраняет изменения
// В случае ошибки возвращается обёрнутое исключение
func (uc *UseCase) СreateRolesAndGroupsForContext(ctx context.Context, meta *models.UserMeta, ctxID string) error {
	//log := uc.l.With(wsl.Label("method", "СreateRolesAndGroupsForContext"))
	// TODO: проверка, что такой группы нет
	// TODO: Ретрай для запросов - можно отдельную функцию напистать, чтобы при запросе
	_, err := uc.kc.CreateGroup(ctxID, []string{"CA", "CO"})
	if err != nil {
		return fmt.Errorf("CreateGroup [ctx: %s] - %w", ctxID, err)
	}
	ca := fmt.Sprintf("CA-%s", ctxID)
	co := fmt.Sprintf("CO-%s", ctxID)
	uc.kc.CreateRole(ca)
	uc.kc.CreateRole(co)

	uc.kc.AddRoleToGroupByPath(fmt.Sprintf("/%s/CA", ctxID), ca)
	uc.kc.AddRoleToGroupByPath(fmt.Sprintf("/%s/CO", ctxID), co)

	return nil
}

// DeleteRolesAndGroupsForContext — удаляет группу и связанные роли (CA-, CO-) по пути в Keycloak,
// удаляет контекст из конфигурации NASM и сохраняет изменения
func (uc *UseCase) DeleteRolesAndGroupsForContext(ctx context.Context, meta *models.UserMeta, ctxID string) error {
	//log := uc.l.With(wsl.Label("method", "DeleteRolesAndGroupsForContext"))

	if err := uc.kc.DeleteGroupByPath("/" + ctxID); err != nil {
		return err
	}
	ca := fmt.Sprintf("CA-%s", ctxID)
	co := fmt.Sprintf("CO-%s", ctxID)
	if err := uc.kc.DeleteRoleByName(ca); err != nil {
		return err
	}
	if err := uc.kc.DeleteRoleByName(co); err != nil {
		return err
	}

	return nil
}

// GetRoles — возвращает список всех ролей из Keycloak, преобразованных в формат models.Role
func (uc *UseCase) GetRoles(ctx context.Context, meta *models.UserMeta, page, limit int, search string) (*[]models.Role, int, error) {
	roles, err := uc.kc.GetRoles()
	if err != nil {
		return nil, 1, err
	}
	res := make([]models.Role, 0, len(roles))
	for _, v := range roles {
		res = append(res, models.Role{
			Role: v,
		})
	}
	return &res, len(res), nil
}

// GetRolesCount — возвращает общее количество ролей в системе
func (uc *UseCase) GetRolesCount(ctx context.Context, meta *models.UserMeta) (int, error) {
	roles, err := uc.kc.GetRoles()
	if err != nil {
		return 1, err
	}
	return len(roles), nil
}
