package models

import "time"

type DeviceAnalyticsPage struct {
	Traffic   TrafficStat `json:"traffic"`
	Allowed   RequestStat `json:"allowed"`
	Blocked   RequestStat `json:"blocked"`
	Pending   RequestStat `json:"pending"`
	Anomalies uint        `json:"anomalies"`
	Blocks    uint        `json:"blocks"`
	All       uint        `json:"all"`
}

type ReportForDevice struct {
	From                time.Time              `json:"from"`
	To                  time.Time              `json:"to"`
	HostName            string                 `json:"host_name"`
	DeviceAnalyticsPage DeviceAnalyticsPage    `json:"device_analytics_page"`
	AnomaliesListPage   map[string]AnomalyStat `json:"anomalies_list_page"`
	CategoriesPage      map[string]TopCategory `json:"categories_page"`
}
