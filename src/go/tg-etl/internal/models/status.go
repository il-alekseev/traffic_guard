package models

import (
	"fmt"
	"strings"
)

type Status string

// Определяем перечисление с константами статусов
const (
	Allowed   Status = "allowed"   // Разрешен
	Blocked   Status = "blocked"   // Заблокирован
	Forbidden Status = "forbidden" // Запрещен
	Pending   Status = "pending"   // Ожидает
)

// String возвращает строковое представление Status
func (s Status) String() string {
	switch s {
	case Allowed:
		return "Разрешен"
	case Blocked:
		return "Заблокирован"
	case Forbidden:
		return "Запрещен"
	case Pending:
		return "Ожидает"
	default:
		return string(s)
	}
}

// IsValid проверяет, является ли значение валидным Status
func (s Status) IsValid() bool {
	switch s {
	case Allowed, Blocked, Forbidden, Pending:
		return true
	default:
		return false
	}
}

// Values возвращает все возможные значения Status
func (s Status) Values() []Status {
	return []Status{Allowed, Blocked, Forbidden, Pending}
}

// ParseStatus преобразует строку в Status
func ParseStatus(str string) (Status, error) {
	switch strings.ToLower(str) {
	case "allowed":
		return Allowed, nil
	case "blocked":
		return Blocked, nil
	case "forbidden":
		return Forbidden, nil
	case "pending":
		return Pending, nil
	default:
		return "", fmt.Errorf("invalid Status: %s", str)
	}
}

// IsActive проверяет, активен ли статус (разрешен)
func (s Status) IsActive() bool {
	return s == Allowed
}

// IsInactive проверяет, неактивен ли статус (заблокирован или запрещен)
func (s Status) IsInactive() bool {
	return s == Blocked || s == Forbidden
}

// IsPending проверяет, находится ли статус в ожидании
func (s Status) IsPending() bool {
	return s == Pending
}
