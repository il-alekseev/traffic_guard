package models

import (
	"fmt"
	"strings"
)

type Status string

// Определяем перечисление с константами статусов
const (
	StatusAllowed Status = "allowed" // Разрешен
	StatusBlocked Status = "blocked" // Запрещен
	StatusAnomaly Status = "anomaly" // Аномалия
	StatusPending Status = "pending" // Ожидает
)

// String возвращает строковое представление Status
func (s Status) String() string {
	switch s {
	case StatusAllowed:
		return "Разрешен"
	case StatusBlocked:
		return "Запрещен"
	case StatusAnomaly:
		return "Аномалия"
	case StatusPending:
		return "Ожидает"
	default:
		return string(s)
	}
}

// IsValid проверяет, является ли значение валидным Status
func (s Status) IsValid() bool {
	switch s {
	case StatusAllowed, StatusBlocked, StatusAnomaly, StatusPending:
		return true
	default:
		return false
	}
}

// Values возвращает все возможные значения Status
func (s Status) Values() []Status {
	return []Status{StatusAllowed, StatusBlocked, StatusAnomaly, StatusPending}
}

// ParseStatus преобразует строку в Status
func ParseStatus(str string) (Status, error) {
	switch strings.ToLower(str) {
	case "Разрешен", "allowed":
		return StatusAllowed, nil
	case "Запрещен", "blocked":
		return StatusBlocked, nil
	case "Аномалия", "anomaly":
		return StatusAnomaly, nil
	case "Ожидает", "pending":
		return StatusPending, nil
	default:
		return "", fmt.Errorf("invalid Status: %s", str)
	}
}
