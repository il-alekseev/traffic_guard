package models

import (
	"fmt"
	"strings"
)

type RequestStatus string

// Определяем перечисление с константами статусов
const (
	RequestStatusAllowed     RequestStatus = "allowed"      // Разрешен
	RequestStatusBlocked     RequestStatus = "blocked"      // Запрещен
	RequestStatusBeforeBlock RequestStatus = "before_block" // До блокировки
	RequestStatusPending     RequestStatus = "pending"      // Ожидает
)

// String возвращает строковое представление RequestStatus
func (s RequestStatus) String() string {
	switch s {
	case RequestStatusAllowed:
		return "Разрешен"
	case RequestStatusBlocked:
		return "Запрещен"
	case RequestStatusBeforeBlock:
		return "До блокировки"
	case RequestStatusPending:
		return "Ожидает"
	default:
		return string(s)
	}
}

// IsValid проверяет, является ли значение валидным RequestStatus
func (s RequestStatus) IsValid() bool {
	switch s {
	case RequestStatusAllowed, RequestStatusBlocked, RequestStatusBeforeBlock, RequestStatusPending:
		return true
	default:
		return false
	}
}

// Values возвращает все возможные значения RequestStatus
func (s RequestStatus) Values() []RequestStatus {
	return []RequestStatus{RequestStatusAllowed, RequestStatusBlocked, RequestStatusBeforeBlock, RequestStatusPending}
}

// ParseRequestStatus преобразует строку в RequestStatus
func ParseRequestStatus(str string) (RequestStatus, error) {
	switch strings.ToLower(str) {
	case "allowed":
		return RequestStatusAllowed, nil
	case "blocked":
		return RequestStatusBlocked, nil
	case "before_block":
		return RequestStatusBeforeBlock, nil
	case "pending":
		return RequestStatusPending, nil
	default:
		return "", fmt.Errorf("invalid Status: %s", str)
	}
}
