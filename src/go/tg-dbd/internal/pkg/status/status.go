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

// String реализует интерфейс fmt.Stringer
func (s Status) ToJSONString() string {
	switch s {
	case Allowed:
		return "allowed"
	case Blocked:
		return "blocked"
	case Prohibited:
		return "prohibited"
	case Waiting:
		return "waiting"
	default:
		return fmt.Sprintf("unknown(%d)", s)
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

// ParseStatus преобразует строку в значение Status
// Возвращает ошибку, если строка не соответствует допустимым значениям
func ParseStatus(input string) (Status, error) {
	switch input {
	case "allowed":
		return Allowed, nil
	case "blocked":
		return Blocked, nil
	case "prohibited":
		return Prohibited, nil
	case "waiting":
		return Waiting, nil
	default:
		return 0, fmt.Errorf("недопустимое значение статуса: %q (ожидается: allowed, blocked, prohibited, waiting)", input)
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
