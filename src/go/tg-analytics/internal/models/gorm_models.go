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
	ID       uint   `gorm:"primaryKey" json:"id"`
	HostName string `gorm:"type:varchar" json:"host"`
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
	RequestID     uuid.UUID `gorm:"type:uuid;uniqueIndex" json:"request_id"`
}

// Category представляет категорию в базе данных
type Category struct {
	ID   uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"type:varchar(255);not null;uniqueIndex" json:"name"`
	Type string `gorm:"type:varchar(20);not null;default:'neutral'" json:"type"`
}

type Notification struct {
	Description string `json:"description"`
	Username    string `json:"username"`
	DateTime    string `json:"dateTime"`
}

type Resource struct {
	Name            string `json:"name"`
	ReqsBeforeBlock int    `json:"reqs_before_block"`
	ReqsAfterBlock  int    `json:"reqs_after_block"`
}

type ResourcePoint struct {
	Date    int64 `json:"date"`
	Locked  int   `json:"locked"`
	Delayed int   `json:"delayed"`
}

type DeviceStat struct {
	Name       string          `json:"name"`
	State      string          `json:"state"`
	Statistics []ResourcePoint `json:"statistics"`
}

type RequestPoint struct {
	Date  int64 `json:"date"`
	Count int   `json:"count"`
}

type RequestsStat struct {
	Accepted      []RequestPoint `json:"accepted"`
	Blocked       []RequestPoint `json:"blocked"`
	BeforeBlocked []RequestPoint `json:"before_blocked"`
	Delayed       []RequestPoint `json:"delayed"`
}

type AnomalyResourse struct {
	ResourseName      string `json:"resourse_name"`
	UnblockedRequests int    `json:"unblocked_requests"`
}

type Anomaly struct {
	BlockedResourcesCount int               `json:"blocked_resources_count"`
	FirewallsCount        int               `json:"firewalls_count"`
	Stat                  []AnomalyResourse `json:"stat"`
}
