package models

import "time"

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
	ID       uint   `gorm:"primaryKey" json:"id"`
	IP       string `gorm:"type:varchar" json:"ip"`
	Port     int    `gorm:"type:integer" json:"port"`
	Country  string `gorm:"type:varchar" json:"country"`
	Username string `gorm:"type:varchar" json:"username"`
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

type TrafficPoint struct {
	Date   int64 `json:"date"`
	Input  int   `json:"input"`
	Output int   `json:"output"`
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
