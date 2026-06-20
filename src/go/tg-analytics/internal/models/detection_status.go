package models

import (
	"fmt"
	"strings"
)

type DetectionStatus string

// Определяем перечисление с константами типов сессий
const (
	DetectionStatusAllowed  DetectionStatus = "allowed"  // Разрешено
	DetectionStatusBlocked  DetectionStatus = "blocked"  // Заблокировано
	DetectionStatusUnsolved DetectionStatus = "unsolved" // Не решено
)

func (s DetectionStatus) String() string {
	switch s {
	case DetectionStatusAllowed:
		return "Разрешено"
	case DetectionStatusBlocked:
		return "Заблокировано"
	case DetectionStatusUnsolved:
		return "Не решено"
	default:
		return string(s)
	}
}

// IsValid проверяет, является ли значение валидным DetectionStatus
func (s DetectionStatus) IsValid() bool {
	switch s {
	case DetectionStatusAllowed, DetectionStatusBlocked, DetectionStatusUnsolved:
		return true
	default:
		return false
	}
}

// Values возвращает все возможные значения DetectionStatus
func (s DetectionStatus) Values() []DetectionStatus {
	return []DetectionStatus{DetectionStatusAllowed, DetectionStatusBlocked, DetectionStatusUnsolved}
}

// ParseDetectionStatus преобразует строку в DetectionStatus
func ParseDetectionStatus(str string) (DetectionStatus, error) {
	switch strings.ToLower(str) {
	case "Разрешено", "allowed", "allow":
		return DetectionStatusAllowed, nil
	case "Заблокировано", "blocked", "deny":
		return DetectionStatusBlocked, nil
	case "Не решено", "unsolved":
		return DetectionStatusUnsolved, nil
	default:
		return "", fmt.Errorf("invalid DetectionStatus: %s", str)
	}
}
