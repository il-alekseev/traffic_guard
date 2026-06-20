package models

import (
	"fmt"
	"strings"
)

type SessionType string

// Определяем перечисление с константами типов сессий
const (
	SessionTypeAllowed SessionType = "allowed" // Разрешен
	SessionTypeBlocked SessionType = "blocked" // Запрещен
	SessionTypeVPN     SessionType = "vpn"     // VPN
)

func (s SessionType) String() string {
	switch s {
	case SessionTypeAllowed:
		return "Разрешен"
	case SessionTypeBlocked:
		return "Запрещен"
	case SessionTypeVPN:
		return "VPN"
	default:
		return string(s)
	}
}

// IsValid проверяет, является ли значение валидным SessionType
func (s SessionType) IsValid() bool {
	switch s {
	case SessionTypeAllowed, SessionTypeBlocked, SessionTypeVPN:
		return true
	default:
		return false
	}
}

// Values возвращает все возможные значения SessionType
func (s SessionType) Values() []SessionType {
	return []SessionType{SessionTypeAllowed, SessionTypeBlocked, SessionTypeVPN}
}

// ParseSessionType преобразует строку в SessionType
func ParseSessionType(str string) (SessionType, error) {
	switch strings.ToLower(str) {
	case "Разрешен", "allowed":
		return SessionTypeAllowed, nil
	case "Запрещен", "blocked":
		return SessionTypeBlocked, nil
	case "VPN", "vpn":
		return SessionTypeVPN, nil
	default:
		return "", fmt.Errorf("invalid SessionType: %s", str)
	}
}
