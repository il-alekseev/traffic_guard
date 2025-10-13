package models

import (
	"time"
)

// Session представляет таблицу session
type Session struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	DatetimeUTC    time.Time `gorm:"column:datetime_utc" json:"datetime_utc"`
	DatetimeDevice time.Time `gorm:"column:datetime_device" json:"datetime_device"`
	DeviceID       uint      `gorm:"column:device_id" json:"device_id"`
	Type           string    `gorm:"type:varchar" json:"type"`
	StatusID       uint      `gorm:"column:status_id" json:"status_id"`
	URL            string    `gorm:"type:varchar" json:"url"`
	DstID          uint      `gorm:"column:dst_id" json:"dst_id"`
	SrcID          uint      `gorm:"column:src_id" json:"src_id"`
	Protocol       string    `gorm:"type:varchar" json:"protocol"`
	AttackHash     float64   `gorm:"column:attack_hash" json:"attack_hash"`
	TriggerCount   int       `gorm:"column:trigger_count" json:"trigger_count"`
}

// Device представляет таблицу device
type Device struct {
	ID   string `gorm:"primaryKey;type:varchar" json:"id"`
	Name string `gorm:"type:varchar" json:"name"`
	Host string `gorm:"type:varchar" json:"host"`
}

// Source представляет таблицу source
type Source struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	IP       string `gorm:"type:varchar" json:"ip"`
	Port     int    `gorm:"type:integer" json:"port"`
	Country  string `gorm:"type:varchar" json:"country"`
	Username string `gorm:"type:varchar" json:"username"`
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
	DecisionID             uint      `gorm:"column:decision_id" json:"decision_id"`
	DecisionDatetime       time.Time `gorm:"column:decision_datetime" json:"decision_datetime"`
	LastAccessDatetime     time.Time `gorm:"column:last_access_datetime" json:"last_access_datetime"`
}

// Decision представляет таблицу decision
type Decision struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Decision string `gorm:"type:varchar" json:"decision"`
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
	ID       uint    `gorm:"primaryKey" json:"id"`
	Category string  `gorm:"type:varchar" json:"category"`
	Percent  float64 `gorm:"type:float" json:"percent"`
}

// CategoryDomain представляет таблицу category_domain (связующая таблица)
type CategoryDomain struct {
	CategoryID uint `gorm:"primaryKey;column:category_id" json:"category_id"`
	DomainID   uint `gorm:"primaryKey;column:domain_id" json:"domain_id"`
}
