package models

import "time"

type IDsLog struct {
	ID             int64     `json:"id" gorm:"primaryKey;column:id"`
	Timestamp      time.Time `json:"timestamp" gorm:"column:timestamp;type:timestamptz;not null"`
	TimestampStr   string    `json:"timestamp_str" gorm:"column:timestamp_str;type:varchar"`
	SensorID       int       `json:"sensor_id" gorm:"column:sensor_id;type:integer"`
	EventType      string    `json:"event_type" gorm:"column:event_type;type:varchar"`
	SrcIP          string    `json:"src_ip" gorm:"column:src_ip;type:varchar"`
	SrcPort        int       `json:"src_port" gorm:"column:src_port;type:integer"`
	DestIP         string    `json:"dest_ip" gorm:"column:dest_ip;type:varchar"`
	DestPort       int       `json:"dest_port" gorm:"column:dest_port;type:integer"`
	Proto          string    `json:"proto" gorm:"column:proto;type:text"`
	InIface        string    `json:"in_iface" gorm:"column:in_iface;type:varchar"`
	Action         string    `json:"action" gorm:"column:action;type:varchar"`
	SignatureID    int       `json:"signature_id" gorm:"column:signature_id;type:integer"`
	Signature      string    `json:"signature" gorm:"column:signature;type:text"`
	Category       string    `json:"category" gorm:"column:category;type:text"`
	Severity       int16     `json:"severity" gorm:"column:severity;type:smallint"`
	Payload        string    `json:"payload" gorm:"column:payload;type:text"`
	Packet         string    `json:"packet" gorm:"column:packet;type:text"`
	Msgrepeatcount int       `json:"msgrepeatcount" gorm:"column:msgrepeatcount;type:integer"`
	DestDomain     string    `json:"dest_domain" gorm:"column:dest_domain;type:varchar"`
	Revision       int       `json:"revision" gorm:"column:revision;type:integer"`
	SignatureBody  string    `json:"signature_body" gorm:"column:signature_body;type:text"`
	Sourcename     string    `json:"sourcename" gorm:"column:sourcename;type:varchar"`
	Hostname       string    `json:"hostname" gorm:"column:hostname;type:varchar"`
	DestCountry    string    `json:"dest_country" gorm:"column:dest_country;type:varchar"`
	SrcCountry     string    `json:"src_country" gorm:"column:src_country;type:varchar"`
	Username       string    `json:"username" gorm:"column:username;type:varchar"`
	VRF            string    `json:"vrf" gorm:"column:vrf;type:varchar;not null"`
}
