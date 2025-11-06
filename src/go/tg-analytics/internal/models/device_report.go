package models

import "time"

type RequestStatData struct {
	Time  []time.Time
	Data  []uint
	Count uint
}

type RequestsAnalytics struct {
	Allowed RequestStatData `json:"allowed"`
	Blocked RequestStatData `json:"blocked"`
	Pending RequestStatData `json:"pending"`
}

type DeviceAnalyticsPage struct {
	Traffic           TrafficStatData   `json:"traffic"`
	RequestsAnalytics RequestsAnalytics `json:"requests_analytics"`
	AnomalyBlockStat  AnomalyBlockStat  `json:"anomaly_block_stat"`
}

type ReportForDevice struct {
	From                time.Time           `json:"from"`
	To                  time.Time           `json:"to"`
	HostName            string              `json:"hostname"`
	DeviceAnalyticsPage DeviceAnalyticsPage `json:"device_analytics_page"`
	AnomaliesListPage   TopAnomaliesPage    `json:"anomalies_list_page"`
	CategoriesPage      TopCategoriesPage   `json:"categories_page"`
}
