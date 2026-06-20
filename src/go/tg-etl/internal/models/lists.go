package models

import (
	"fmt"
	"strings"
)

type ListType string

// Определяем перечисление с константами
const (
	Whitelist ListType = "whitelist"
	Blacklist ListType = "blacklist"
)

// String возвращает строковое представление ListType
func (l ListType) String() string {
	return string(l)
}

// IsValid проверяет, является ли значение валидным ListType
func (l ListType) IsValid() bool {
	switch l {
	case Whitelist, Blacklist:
		return true
	default:
		return false
	}
}

// Values возвращает все возможные значения ListType
func (l ListType) Values() []ListType {
	return []ListType{Whitelist, Blacklist}
}

// ParseListType преобразует строку в ListType
func ParseListType(s string) (ListType, error) {
	switch strings.ToLower(s) {
	case "whitelist":
		return Whitelist, nil
	case "blacklist":
		return Blacklist, nil
	default:
		return "", fmt.Errorf("invalid ListType: %s", s)
	}
}
