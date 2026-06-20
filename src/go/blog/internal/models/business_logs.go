package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type BusinessLog struct {
	ID          int       `json:"id" gorm:"column:id;primaryKey;autoIncrement;not null"`
	Timestamp   time.Time `json:"timestamp,omitzero" gorm:"column:timestamp;type:timestamp with time zone;not null"`
	EventType   string    `json:"event_type,omitempty" gorm:"column:event_type;type:varchar(40);not null"` // CREATE, UPDATE, DELETE
	Entity      string    `json:"entity,omitempty" gorm:"column:entity;type:varchar(40);not null"`         // user, context
	Username    string    `json:"user_name,omitempty" gorm:"column:user_name;type:varchar(36);not null"`
	UserRole    string    `json:"user_role,omitempty" gorm:"column:user_role;type:varchar(10);not null"` // SA, CA
	Context     string    `json:"context,omitempty" gorm:"column:context;type:varchar(50)"`
	EntityID    string    `json:"entity_id,omitempty" gorm:"column:entity_id;type:varchar(50);not null"` // userID or ContextID
	OldValue    JSONB     `json:"old_value,omitempty" gorm:"column:old_value;type:jsonb"`
	NewValue    JSONB     `json:"new_value,omitempty" gorm:"column:new_value;type:jsonb"`
	Description string    `json:"description,omitempty" gorm:"column:description;type:text"`
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
