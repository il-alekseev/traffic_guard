package models

import "time"

// IdsLog представляет таблицу ids_log
type IdsLog struct {
	ID             int64     `gorm:"primaryKey" json:"id"`
	Timestamp      time.Time `gorm:"type:timestamptz;not null" json:"timestamp"`
	TimestampStr   string    `gorm:"type:varchar" json:"timestamp_str"`
	SensorID       int       `gorm:"column:sensor_id" json:"sensor_id"`
	EventType      string    `gorm:"type:varchar" json:"event_type"`
	SrcIP          string    `gorm:"type:varchar" json:"src_ip"`
	SrcPort        int       `gorm:"column:src_port" json:"src_port"`
	DestIP         string    `gorm:"type:varchar" json:"dest_ip"`
	DestPort       int       `gorm:"column:dest_port" json:"dest_port"`
	Proto          string    `gorm:"type:text" json:"proto"`
	InIface        string    `gorm:"type:varchar" json:"in_iface"`
	Action         string    `gorm:"type:varchar" json:"action"`
	SignatureID    int       `gorm:"column:signature_id" json:"signature_id"`
	Signature      string    `gorm:"type:text" json:"signature"`
	Category       string    `gorm:"type:text" json:"category"`
	Severity       int16     `gorm:"type:smallint" json:"severity"`
	Payload        string    `gorm:"type:text" json:"payload"`
	Packet         string    `gorm:"type:text" json:"packet"`
	Msgrepeatcount int       `gorm:"column:msgrepeatcount" json:"msgrepeatcount"`
	DestDomain     string    `gorm:"type:varchar" json:"dest_domain"`
	Revision       int       `gorm:"type:integer" json:"revision"`
	SignatureBody  string    `gorm:"type:text" json:"signature_body"`
	Sourcename     string    `gorm:"type:varchar" json:"sourcename"`
	Hostname       string    `gorm:"type:varchar" json:"hostname"`
	DestCountry    string    `gorm:"type:varchar" json:"dest_country"`
	SrcCountry     string    `gorm:"type:varchar" json:"src_country"`
	Username       string    `gorm:"type:varchar" json:"username"`
	Vrf            string    `gorm:"type:varchar;not null" json:"vrf"`
}

// TableName задает имя таблицы в БД
func (IdsLog) TableName() string {
	return "ids_log"
}
