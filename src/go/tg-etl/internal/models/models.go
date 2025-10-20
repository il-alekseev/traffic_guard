package models

import (
	"time"
)

// SecurityEvent представляет запись о событии безопасности
type SecurityEvent struct {
	ID               int       `json:"id" gorm:"primaryKey"`
	EventDateUTC     time.Time `json:"event_date_utc"`
	NodeSecurityDate time.Time `json:"node_security_date"`
	DeviceID         int       `json:"device_id"`
	EventID          int       `json:"event_id"`
	ImportanceID     int       `json:"importance_id"`
	SenderID         int       `json:"sender_id"`
	ReceiverID       int       `json:"receiver_id"`
	Protocol         string    `json:"protocol"`
	AttackHash       float64   `json:"attack_hash"`
	TriggerCount     int       `json:"trigger_count"`
	CreatedAt        time.Time `json:"created_at"`
}

// SecurityDevice представляет информацию об устройстве безопасности
type SecurityDevice struct {
	ID           int    `json:"id" gorm:"primaryKey"`
	SecurityNode string `json:"security_node"`
	Host         string `json:"host"`
}

// SendersInfo представляет информацию об отправителе
type SendersInfo struct {
	ID              int    `json:"id" gorm:"primaryKey"`
	SenderIPAddress string `json:"sender_ip_address"`
	SenderCountry   string `json:"sender_country"`
	SenderName      string `json:"sender_name"`
	SenderPort      int    `json:"sender_port"`
}

// ReceiversInfo представляет информацию о получателе
type ReceiversInfo struct {
	ID                int    `json:"id" gorm:"primaryKey"`
	ReceiverIPAddress string `json:"receiver_ip_address"`
	ReceiverCountry   string `json:"receiver_country"`
	ReceiverDomain    string `json:"receiver_domain"`
	ReceiverPort      int    `json:"receiver_port"`
}

// EventType представляет тип события
type EventType struct {
	ID         int    `json:"id" gorm:"primaryKey"`
	EventType  string `json:"event_type"`
	CategoryID int    `json:"category_id"`
	Action     string `json:"action"`
}

// Category представляет категорию события
type Category struct {
	ID       int     `json:"id" gorm:"primaryKey"`
	Category string  `json:"category"`
	Percent  float64 `json:"percent"`
}

// Importance представляет важность события
type Importance struct {
	ID         int    `json:"id" gorm:"primaryKey"`
	Importance string `json:"importance"`
}
