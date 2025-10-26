package models

import (
	"fmt"
	"strings"
)

type SessionType string

// Определяем перечисление с константами типов сессий
const (
	Firewall SessionType = "firewall" // Файервол
	Anomaly  SessionType = "anomaly"  // Аномалия
	Blocking SessionType = "blocking" // Блокировка
	Waiting  SessionType = "waiting"  // Ожидание
	VPN      SessionType = "vpn"      // VPN
	Allowing SessionType = "allowing" // Разрешен
)

func (s SessionType) String() string {
	switch s {
	case Firewall:
		return "Файервол"
	case Anomaly:
		return "Аномалия"
	case Blocking:
		return "Блокировка"
	case Waiting:
		return "Ожидание"
	case VPN:
		return "VPN"
	case Allowing:
		return "Разрешен"
	default:
		return string(s)
	}
}

// IsValid проверяет, является ли значение валидным SessionType
func (s SessionType) IsValid() bool {
	switch s {
	case Firewall, Anomaly, Blocking, Waiting, VPN, Allowing:
		return true
	default:
		return false
	}
}

// Values возвращает все возможные значения SessionType
func (s SessionType) Values() []SessionType {
	return []SessionType{Firewall, Anomaly, Blocking, Waiting, VPN, Allowing}
}

// ParseSessionType преобразует строку в SessionType
func ParseSessionType(str string) (SessionType, error) {
	switch strings.ToLower(str) {
	case "firewall":
		return Firewall, nil
	case "anomaly":
		return Anomaly, nil
	case "blocking":
		return Blocking, nil
	case "waiting":
		return Waiting, nil
	case "vpn":
		return VPN, nil
	case "allowed":
		return Allowing, nil
	default:
		return "", fmt.Errorf("invalid SessionType: %s", str)
	}
}

// IsSecurityRelated проверяет, связан ли тип сессии с безопасностью
func (s SessionType) IsSecurityRelated() bool {
	return s == Firewall || s == Anomaly || s == Blocking
}

// IsNetworkRelated проверяет, связан ли тип сессии с сетью
func (s SessionType) IsNetworkRelated() bool {
	return s == VPN
}

// RequiresMonitoring проверяет, требует ли тип сессии мониторинга
func (s SessionType) RequiresMonitoring() bool {
	return s == Anomaly || s == Waiting
}

// IsAllowed проверяет, является ли тип сессии разрешенным
func (s SessionType) IsAllowed() bool {
	return s == Allowing
}
