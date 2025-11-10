package v1

import (
	"context"
	"userctrl/internal/controllers/http/v1/dto"
	"userctrl/internal/models"
)

// UseCaseAuth — методы для управления сессиями: вход, выход, обновление и валидация токенов
type UseCaseAuth interface {
	RefreshToken(ctx context.Context, refreshToken string) (models.Token, error)
	SignIn(ctx context.Context, username string, password string) (models.Token, error)
	SignOut(ctx context.Context, refreshToken string) error
	ValidateToken(ctx context.Context, accessToken string) (int, error)
}

// UseCaseUser — методы для управления пользователями: создание, обновление, удаление, поиск, сброс пароля и др.
type UseCaseUser interface {
	Count(ctx context.Context, meta *models.UserMeta) (int, error)
	Create(ctx context.Context, meta *models.UserMeta, u dto.UserCreateRequest) (models.User, error)
	Delete(ctx context.Context, meta *models.UserMeta, userID string) error
	GetByID(ctx context.Context, meta *models.UserMeta, userID string) (models.User, error)
	GetCountUsersByRole(ctx context.Context, meta *models.UserMeta) (*models.CountUserByRole, error)
	GetProfile(ctx context.Context, meta *models.UserMeta, userID string) (models.User, error)
	GetUsers(ctx context.Context, meta *models.UserMeta, page int, limit int, role string, contextID string, search string) ([]models.User, int, error)
	GetUsersCountForContext(ctx context.Context, meta *models.UserMeta, ctxID string) (*models.MetaDataForContext, error)
	ResetPassword(ctx context.Context, meta *models.UserMeta, userID string, newPassword string) error
	Update(ctx context.Context, meta *models.UserMeta, userID string, u dto.UserUpdateData) error
	UpdatePassword(ctx context.Context, meta *models.UserMeta, userID string, oldPassword string, newPassword string) error
}

// UseCaseRole — методы для управления ролями и группами: назначение, удаление, получение списка, а также создание/удаление ролей по контексту
type UseCaseRole interface {
	DeleteRolesAndGroupsForContext(ctx context.Context, meta *models.UserMeta, ctxID string) error
	GetRoles(ctx context.Context, meta *models.UserMeta, page int, limit int, search string) (*[]models.Role, int, error)
	GetRolesCount(ctx context.Context, meta *models.UserMeta) (int, error)
	RemoveRole(ctx context.Context, meta *models.UserMeta, userID string, role string) error
	UpdateRole(ctx context.Context, meta *models.UserMeta, userID string, role string) error
	СreateRolesAndGroupsForContext(ctx context.Context, meta *models.UserMeta, ctxID string) error
}

// UseCaseInterface - определяет контракт для бизнес-логики приложения
type UseCaseInterface interface {
	UseCaseAuth
	UseCaseUser
	UseCaseRole
}
