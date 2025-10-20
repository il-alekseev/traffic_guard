package status

import (
	"encoding/json"
	"fmt"
)

type Status int

const (
	Allowed Status = iota + 1
	Blocked
	Prohibited
	Waiting
)

// String реализует интерфейс fmt.Stringer
func (s Status) String() string {
	switch s {
	case Allowed:
		return "Разрешено"
	case Blocked:
		return "Заблокировано"
	case Prohibited:
		return "Запрещено"
	case Waiting:
		return "Ожидает"
	default:
		return fmt.Sprintf("Неизвестно(%d)", s)
	}
}

// Валидация значения
func (s Status) IsValid() bool {
	switch s {
	case Allowed, Blocked, Prohibited, Waiting:
		return true
	default:
		return false
	}
}

// MarshalJSON реализует json.Marshaler
func (s Status) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

// UnmarshalJSON реализует json.Unmarshaler
func (s *Status) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	switch str {
	case "allowed":
		*s = Allowed
	case "blocked":
		*s = Blocked
	case "prohibited":
		*s = Prohibited
	case "waiting":
		*s = Waiting
	default:
		return fmt.Errorf("invalid status value: %s", str)
	}
	return nil
}
