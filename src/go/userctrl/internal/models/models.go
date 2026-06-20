package models

// User — модель пользователя с полями: ID, логин, email, ФИО, роль, метаданные и флаги
type User struct {
	ID                     string `json:"user_id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" example:"a71032a2-87eb-4d4e-a2a0-d86404672799"`
	Login                  string `json:"login" gorm:"unique;not null;size:50" example:"johndoe"`
	Email                  string `json:"email" gorm:"unique;not null;size:255" example:"john.doe@example.com"`
	FirstName              string `json:"first_name" gorm:"size:100;not null" example:"John"`
	LastName               string `json:"last_name" gorm:"size:100;not null" example:"Doe"`
	Patronymic             string `json:"patronymic" gorm:"size:100" example:"Ivanovich"`
	Role                   string `json:"role" gorm:"size:50;not null;default:'user'" example:"SA"`
	CreatedAt              string `json:"created_at" gorm:"type:timestamp;default:now()" example:"2023-04-01T12:00:00Z"`
	IsNeedToChangePassword bool   `json:"is_need_to_change_password" gorm:"default:true" example:"true"`
	IsSuperAdmin           bool   `json:"is_super_admin" gorm:"default:false" example:"false"`
}

// Role — модель роли с единственным полем role
type Role struct {
	Role string `json:"role" binding:"required"`
}

// Token — данные аутентификации: access и refresh токены, сессия Grafana и срок действия
type Token struct {
	AccessToken          string `json:"access_token"`
	RefreshToken         string `json:"refresh_token"`
	GrafanaSession       string `json:"grafana_session"`
	GrafanaSessionExpiry string `json:"grafana_session_expiry"`
	ExpiresIn            int    `json:"expires_in"`
}

// MetaDataForContext — информация о количестве пользователей (CA, CO) в определённом контексте
type MetaDataForContext struct {
	ContextID string `json:"context_id"`
	CountCA   int    `json:"count_ca"` // число CA в контексте
	CountCO   int    `json:"count_co"` // число CO в контексте
}

// CountUserByRole — количество пользователей по ролям (SA, CA, CO)
type CountUserByRole struct {
	CountSA int `json:"count_sa"`
	CountCA int `json:"count_ca"`
	CountCO int `json:"count_co"`
}

// UserMeta — метаинформация о текущем пользователе (его UUID, имя, роль, контекст) для передачи между слоями приложения
type UserMeta struct {
	UUID       string `json:"uuid"`
	Username   string `json:"username"`
	ClientRole string `json:"client_role"`
	ShortRole  string `json:"short_role"`
	ContextID  string `json:"context_id"`
}
