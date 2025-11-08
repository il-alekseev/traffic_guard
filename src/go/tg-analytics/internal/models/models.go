package models

import "time"

type Pagination struct {
	Page  int `form:"page"`
	Limit int `form:"limit"`
}

type Sorting struct {
	OrderBy  string `form:"order_by"`
	OrderDir string `form:"order_dir"`
}

// CategoryCount представляет результат подсчета категорий
type CategoryCount struct {
	Category string `json:"category"`
	Count    int64  `json:"count"`
}

type TrafficStat struct {
	Time   []time.Time
	Input  []uint
	Output []uint
}

type RequestStat struct {
	Time []time.Time
	Data []uint
}

type DeviceRequestStat struct {
	HostName string `json:"hostname"`
	Status   string `json:"status"`
	Blocked  []uint `json:"blocked"`
	Pending  []uint `json:"pending"`
}

type HostAnomalies struct {
	HostName     string   `json:"hostname"`
	AnomalyCount uint     `json:"anomaly_count"`
	Domains      []string `json:"domains"`
}
