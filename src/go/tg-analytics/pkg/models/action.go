package models

import "time"

// ActionType представляет тип действия/решения
type ActionType string

const (
	ActionTypeAllowed ActionType = "allowed" // Разрешено
	ActionTypeDenied  ActionType = "denied"  // Запрещено
)

func (t ActionType) String() string {
	switch t {
	case ActionTypeAllowed:
		return "Разрешено"
	case ActionTypeDenied:
		return "Запрещено"
	default:
		return "Неизвестный тип"
	}
}

// Action представляет действие/решение в базе данных
type Action struct {
	ID        uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Action    ActionType `gorm:"type:varchar(20);not null;uniqueIndex" json:"action"`
	CreatedAt time.Time  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	CreatedBy string     `gorm:"type:varchar(100);not null;default:'system'" json:"created_by"`
}

// String возвращает строковое представление действия
func (a Action) String() string {
	return a.Action.String()
}

// IsValid проверяет, является ли действие допустимым
func (a Action) IsValid() bool {
	return a.Action == ActionTypeAllowed || a.Action == ActionTypeDenied
}

// IsAllowed проверяет, является ли действие разрешающим
func (a Action) IsAllowed() bool {
	return a.Action == ActionTypeAllowed
}

// IsDenied проверяет, является ли действие запрещающим
func (a Action) IsDenied() bool {
	return a.Action == ActionTypeDenied
}
