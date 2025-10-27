package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Session struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	DatetimeUTC time.Time `gorm:"type:timestamp;not null;index" json:"datetime_utc"`
	DeviceID    uint      `gorm:"type:integer;not null;index" json:"device_id"`
	Type        string    `gorm:"type:varchar(255);not null" json:"type"`
	// TODO: брать из домена
	//StatusID    uint      `gorm:"type:integer;not null;index" json:"status_id"`
	Status   string `gorm:"type:text;not null" json:"status"`
	URL      string `gorm:"type:text" json:"url"`
	IP       string `gorm:"type:varchar(45);not null" json:"ip"`
	DomainID uint   `gorm:"type:integer;not null;index" json:"domain_id"`
	SrcID    uint   `gorm:"type:integer;not null;index" json:"src_id"`
	Proto    string `gorm:"type:varchar(10);not null" json:"proto"`
}

// Device представляет таблицу device
type Device struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	HostName string `gorm:"type:varchar" json:"host"`
}

// Source представляет таблицу source
type Source struct {
	ID uint   `gorm:"primaryKey" json:"id"`
	IP string `gorm:"type:varchar" json:"ip"`
	//Port     int    `gorm:"type:integer" json:"port"`
	Country  string `gorm:"type:varchar" json:"country"`
	Username string `gorm:"type:varchar" json:"username"`
}

// BeforeCreate - GORM hook для автоматической генерации UUID перед созданием записи
func (domain *Domain) BeforeCreate(tx *gorm.DB) error {
	if domain.UUID == uuid.Nil {
		domain.UUID = uuid.New()
	}
	return nil
}

// Domain представляет таблицу domain
type Domain struct {
	ID                     uint      `gorm:"primaryKey" json:"id"`
	IP                     string    `gorm:"type:varchar" json:"ip"`
	Port                   int       `gorm:"type:integer" json:"port"`
	Country                string    `gorm:"type:varchar" json:"country"`
	Path                   string    `gorm:"type:varchar" json:"url"`
	AccessCount            int       `gorm:"column:access_count" json:"access_count"`
	AnalysisAttemptsCount  int       `gorm:"column:analysis_attemps_count" json:"analysis_attempts_count"`
	ContentAnalysisCounter int       `gorm:"column:content_analysis_counter" json:"content_analysis_counter"`
	CategotyID             int       `gorm:"type:integer" json:"categoty_id"`
	DecisionID             uint      `gorm:"column:decision_id" json:"decision_id"`
	LastAccessDatetime     time.Time `gorm:"column:last_access_datetime" json:"last_access_datetime"`
	PutKafkaDateTime       time.Time `gorm:"column:put_kafka_datetime" json:"put_kafka_datetime"`
	UUID                   uuid.UUID `gorm:"type:uuid;uniqueIndex" json:"uuid"`
}

// Decision представляет таблицу decision
type Decision struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Decision  string    `gorm:"type:varchar" json:"decision"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	CreatedBy string    `gorm:"column:created_by" json:"created_by"`
}

// Detection представляет таблицу detection
type Detection struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Description string `gorm:"type:varchar" json:"description"`
	SessionID   uint   `gorm:"column:session_id" json:"session_id"`
	ActionID    int    `gorm:"column:action_id" json:"action_id"`
}

// Status представляет таблицу status
type Status struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	Status string `gorm:"type:varchar" json:"status"`
}

// Category представляет таблицу category
type Category struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"type:varchar;unique" json:"name"`
}

// Category представляет таблицу content_category
type ContentCategory struct {
	ID         uint    `gorm:"primaryKey" json:"id"`
	CategoryID uint    `gorm:"type:integer" json:"category_id"`
	Percent    float64 `gorm:"type:float" json:"percent"`
}

// CategoryDomain представляет таблицу category_domain (связующая таблица)
type CategoryDomain struct {
	ContentCategoryID uint `gorm:"primaryKey;column:content_category_id" json:"content_category_id"`
	DomainID          uint `gorm:"primaryKey;column:domain_id" json:"domain_id"`
}
