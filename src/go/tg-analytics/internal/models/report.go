package models

import "time"

// Первая страница
type RequestReport struct {
	All         uint `json:"all"`
	BeforeBlock uint `json:"before_block"`
	Pending     uint `json:"pending"`
	AfterBlock  uint `json:"after_block"`
}

type CategoryStat struct {
	Category string        `json:"category"`
	Stat     RequestReport `json:"stat"`
}

type ResourceStat struct {
	Resource   string        `json:"resource"`
	Categories []string      `json:"categories"`
	Stat       RequestReport `json:"stat"`
}

type Traffic struct {
	Input  uint `json:"input"`
	Output uint `json:"output"`
}

type TrafficStatData struct {
	Data  TrafficStat `json:"data"`
	Count uint        `json:"count"`
}

type MainActivityPage struct {
	TopCategories []CategoryStat  `json:"top_categories"`
	TopResources  []ResourceStat  `json:"top_resources"`
	Traffic       TrafficStatData `json:"traffic"`
}

// Вторая страница
type DetectionReport struct {
	All        uint `json:"all"`
	Unresolved uint `json:"unresolved"`
	Blocked    uint `json:"blocked"`
	Allowed    uint `json:"allowed"`
}

type AnomalyBlockStat struct {
	Anomalies uint `json:"anomalies"`
	Blocks    uint `json:"blocks"`
	All       uint `json:"all"`
}

type DeviceReport struct {
	HostName         string           `json:"hostname"`
	Traffic          Traffic          `json:"traffic"`
	Requests         uint             `json:"requests"`
	AnomalyBlockStat AnomalyBlockStat `json:"anomaly_block_stat"`
	Detections       DetectionReport  `json:"detections"`
}

type DevicesAnalyticsPage struct {
	Analytics []DeviceReport `json:"analytics"`
}

// Третья страница
type AnomalyReport struct {
	URL       string        `json:"url"`
	LiveCount int64         `json:"live_count"`
	Status    string        `json:"status"`
	Traffic   Traffic       `json:"traffic"`
	Stat      RequestReport `json:"stat"`
}

type DeviceAnomaly struct {
	HostName    string          `json:"hostname"`
	AnomalyStat []AnomalyReport `json:"anomaly_stat"`
}

type DevicesAnomaliesListPage struct {
	Anomalies []DeviceAnomaly `json:"anomalies"`
}

// Четвертая страница
type DeviceAnomalyAnalytics struct {
	HostName         string           `json:"hostname"`
	Traffic          Traffic          `json:"traffic"`
	Requests         uint             `json:"requests"`
	AnomalyBlockStat AnomalyBlockStat `json:"anomaly_block_stat"`
	Detections       DetectionReport  `json:"detections"`
}

type TopAnomaliesPage struct {
	DeviceAnomaly []DeviceAnomalyAnalytics `json:"device_anomaly"`
}

// Пятая страница
type TopCategory struct {
	Category string        `json:"category"`
	Traffic  Traffic       `json:"traffic"`
	Stat     RequestReport `json:"stat"`
}

type TopCategoriesPage struct {
	Categories []TopCategory `json:"categories"`
}

type Report struct {
	From                time.Time                `json:"from"`
	To                  time.Time                `json:"to"`
	MainActivityPage    MainActivityPage         `json:"main_activity_page"`
	DeviceAnalyticsPage DevicesAnalyticsPage     `json:"device_analytics_page"`
	AnomaliesListPage   DevicesAnomaliesListPage `json:"anomalies_list_page"`
	TopAnomaliesPage    TopAnomaliesPage         `json:"top_anomalies_page"`
	TopCategoriesPage   TopCategoriesPage        `json:"top_categories_page"`
}
