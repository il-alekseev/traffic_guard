package bizlogger

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// Constants for EventType values
const (
	EventTypeCreate string = "CREATE"
	EventTypeUpdate string = "UPDATE"
	EventTypeDelete string = "DELETE"
)

// Constants for UserRole values
const (
	UserRoleSA string = "SA" // UserRoleSA represents a System Administrator role
	UserRoleCA string = "CA" // UserRoleCA represents a Client Administrator role
)

// Constants for Entity values
const (
	EntityUser    string = "user"    // EntityUser represents user account entities.
	EntityContext string = "context" // EntityContext represents context entities.
)

// BizLog represents a single business log with full information.
// JSON tags ensure proper serialization for file storage and output to terminal.
type BusinessLog struct {
	ID          int       `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Timestamp   time.Time `json:"timestamp" gorm:"column:timestamp;type:timestamp with time zone;not null"`
	EventType   string    `json:"event_type" gorm:"column:event_type;type:varchar(40);not null"` // CREATE, UPDATE, DELETE
	Entity      string    `json:"entity" gorm:"column:entity;type:varchar(40);not null"`         // user, context
	Username    string    `json:"username" gorm:"column:user_name;type:varchar(36);not null"`
	UserRole    string    `json:"user_role" gorm:"column:user_role;type:varchar(10);not null"` // SA, CA
	ContextID   string    `json:"context_id" gorm:"column:context;type:varchar(50)"`
	EntityID    string    `json:"entity_id" gorm:"column:entity_id;type:varchar(50);not null"` // userID or ContextID
	OldValue    JSONB     `json:"old_value" gorm:"column:old_value;type:jsonb"`
	NewValue    JSONB     `json:"new_value" gorm:"column:new_value;type:jsonb"`
	Description string    `json:"description" gorm:"column:description;type:text"`
}

// JSONB тип для работы с jsonb в GORM
type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &j)
}

func convertToJSONB(value any) (JSONB, error) {
	if value == nil {
		return nil, nil
	}

	// Если значение уже JSONB, просто возвращаем его
	if j, ok := value.(JSONB); ok {
		return j, nil
	}

	// Если значение уже map[string]interface{}, приводим к JSONB
	if m, ok := value.(map[string]interface{}); ok {
		return JSONB(m), nil
	}

	// Сериализуем в JSON
	jsonData, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	// Десериализуем в map[string]interface{}
	var result JSONB
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
