package models

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	DatetimeUTC time.Time `gorm:"type:timestamp;not null;index" json:"datetime_utc"`
	DeviceID    uint      `gorm:"type:integer;not null;index" json:"device_id"`
	Type        string    `gorm:"type:varchar(255);not null" json:"type"`
	Status      string    `gorm:"type:text;not null" json:"status"`
	IP          string    `gorm:"type:varchar(45);not null" json:"ip"`
	URLID       uint      `gorm:"type:integer;not null;index" json:"url_id"`
	DomainID    uint      `gorm:"type:integer;not null;index" json:"domain_id"`
	SrcID       uint      `gorm:"type:integer;not null;index" json:"src_id"`
}

// Device представляет таблицу device
type Device struct {
	ID       uint   `gorm:"primaryKey;column:id" json:"id"`
	HostName string `gorm:"column:hostname;type:varchar" json:"hostname"`
}

// Source представляет таблицу source
type Source struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	IP       string `gorm:"type:varchar" json:"ip"`
	Country  string `gorm:"type:varchar" json:"country"`
	Username string `gorm:"type:varchar" json:"username"`
}

// Domain представляет таблицу domain
type Domain struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	IP            string    `gorm:"type:varchar" json:"ip"`
	Port          int       `gorm:"type:integer" json:"port"`
	Country       string    `gorm:"type:varchar" json:"country"`
	Path          string    `gorm:"type:varchar" json:"path"`
	CategoryID    int       `gorm:"type:integer;default:1" json:"category_id"`
	CategorizedAt time.Time `gorm:"type:timestamp;" json:"categorized_at"`
	ActionID      uint      `gorm:"column:action_id;default:0" json:"action_id"`
	AnalysisCount uint      `gorm:"column:analysis_count;default:0" json:"analysis_count"`
}

// Action представляет таблицу action
type Action struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Action    string    `gorm:"type:varchar" json:"action"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	CreatedBy string    `gorm:"column:created_by" json:"created_by"`
}

// Пока приходит только одна категория контента, поэтому эти таблицы не нужны
// Category представляет таблицу content_category
/*type ContentCategory struct {
	ID         uint    `gorm:"primaryKey" json:"id"`
	CategoryID uint    `gorm:"type:integer" json:"category_id"`
	Percent    float64 `gorm:"type:float" json:"percent"`
}

// CategoryDomain представляет таблицу category_domain (связующая таблица)
type CategoryDomain struct {
	ContentCategoryID uint `gorm:"primaryKey;column:content_category_id" json:"content_category_id"`
	DomainID          uint `gorm:"primaryKey;column:domain_id" json:"domain_id"`
}*/

// DomainControlLists представляет таблицу domain_control_lists белого и черного списка доменов
type DomainControlLists struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	DomainID  uint      `gorm:"type:integer;not null;index" json:"domain_id"`
	CreatedAt time.Time `gorm:"type:timestamp;not null;" json:"created_at"`
	Type      string    `gorm:"type:varchar;not null;" json:"type"`
}

type URL struct {
	ID            uint      `gorm:"primaryKey;column:id" json:"id"`
	Path          string    `gorm:"type:varchar(2048);not null" json:"path"`
	Proto         string    `gorm:"type:varchar(10);not null" json:"proto"`
	DomainID      uint      `gorm:"type:integer;not null;index" json:"domain_id"`
	IDSLogsAt     time.Time `gorm:"column:ids_logs_at;type:timestamp" json:"ids_logs_at"`
	PutKafkaAt    time.Time `gorm:"column:put_kafka_at;type:timestamp" json:"put_kafka_at"`
	GetMetaDataAt time.Time `gorm:"column:get_metadata_at;type:timestamp" json:"get_metadata_at"`
	GetCategoryAt time.Time `gorm:"column:get_category_at;type:timestamp" json:"get_category_at"`
	//RequestID     uuid.UUID `gorm:"type:uuid;uniqueIndex" json:"request_id"`
	// TODO: отладить, почему повторяются RequestID
	RequestID uuid.UUID `gorm:"type:uuid" json:"request_id"`
}

type LastLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Timestamp time.Time `gorm:"type:timestamptz;not null" json:"timestamp"`
}
