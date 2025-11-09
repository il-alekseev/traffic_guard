package models

import (
	"fmt"
	"strings"
)

type ActionType string

// Определяем перечисление с константами типов действий
const (
	ActionTypeAllow ActionType = "allow" // Разрешить
	ActionTypeDeny  ActionType = "deny"  // Запретить
)

// String возвращает строковое представление ActionType
func (a ActionType) String() string {
	switch a {
	case ActionTypeAllow:
		return "allow"
	case ActionTypeDeny:
		return "deny"
	default:
		return string(a)
	}
}

// IsValid проверяет, является ли значение валидным ActionType
func (a ActionType) IsValid() bool {
	switch a {
	case ActionTypeAllow, ActionTypeDeny:
		return true
	default:
		return false
	}
}

// Values возвращает все возможные значения ActionType
func (ActionType) Values() []ActionType {
	return []ActionType{
		ActionTypeAllow,
		ActionTypeDeny,
	}
}

// Values возвращает все возможные значения ActionType
func ActionTypeStringValues() []string {
	return []string{
		string(ActionTypeAllow),
		string(ActionTypeDeny),
	}
}

// ParseActionType преобразует строку в ActionType
func ParseActionType(str string) (ActionType, error) {
	switch strings.ToLower(strings.TrimSpace(str)) {
	case "allow":
		return ActionTypeAllow, nil
	case "deny":
		return ActionTypeDeny, nil
	default:
		return "", fmt.Errorf("invalid ActionType: %s", str)
	}
}
