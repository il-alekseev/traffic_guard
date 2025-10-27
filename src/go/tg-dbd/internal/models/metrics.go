package models

import (
	"time"

	"gorm.io/datatypes"
)

type StatsJSON struct {
	ID     int            `json:"id" db:"id"`
	Date   time.Time      `json:"date" db:"date"`
	Stat   datatypes.JSON `json:"stat" db:"stat"`
	NodeID int            `json:"node_id" db:"node_id"`
}

// Указываем название таблицы
func (StatsJSON) TableName() string {
	return "stats_json"
}

// Структуры для парсинга JSON
type Statistics struct {
	Statistics struct {
		Network struct {
			Interfaces map[string]Interface `json:"interfaces"`
		} `json:"network"`
	} `json:"statistics"`
}

type Interface struct {
	Input struct {
		Bytes struct {
			Count int64 `json:"count"`
		} `json:"bytes"`
	} `json:"input"`
	Output struct {
		Bytes struct {
			Count int64 `json:"count"`
		} `json:"bytes"`
	} `json:"output"`
}
