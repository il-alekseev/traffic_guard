package models

import "time"

type RequestReport struct {
	BeforeBlock uint `json:"before_block"`
	Pending     uint `json:"pending"`
	AfterBlock  uint `json:"after_block"`
}

type ResourceStat struct {
	Category string        `json:"category"`
	Stat     RequestReport `json:"stat"`
}

type MainActivityPage struct {
	TopCategories map[string]RequestReport `json:"top_categories"`
	TopResources  map[string]ResourceStat  `json:"top_resources"`
	Traffic       TrafficStat              `json:"traffic"`
}

type DetectionReport struct {
	All        uint `json:"all"`
	Unresolved uint `json:"unresolved"`
	Blocked    uint `json:"blocked"`
	Allowed    uint `json:"allowed"`
}

type DeviceReport struct {
	Input      uint            `json:"input"`
	Output     uint            `json:"output"`
	Requests   uint            `json:"requests"`
	Anomalies  uint            `json:"anomalies"`
	Blocks     uint            `json:"blocks"`
	All        uint            `json:"all"`
	Detections DetectionReport `json:"detections"`
}

type DevicesAnalyticsPage struct {
	Analytics map[string]DeviceReport `json:"analytics"`
}

type AnomalyStat struct {
	LiveCount   int64 `json:"live_count"`
	Input       uint  `json:"input"`
	Output      uint  `json:"output"`
	Requests    uint  `json:"requests"`
	BeforeBlock uint  `json:"before_block"`
	Pending     uint  `json:"pending"`
	AfterBlock  uint  `json:"after_block"`
}

type AnomaliesListPage struct {
	Anomalies map[string]AnomalyStat `json:"anomalies"`
}

type TopAnomaly struct {
	Input      uint            `json:"input"`
	Output     uint            `json:"output"`
	Requests   uint            `json:"requests"`
	Anomalies  uint            `json:"anomalies"`
	Blocks     uint            `json:"blocks"`
	All        uint            `json:"all"`
	Detections DetectionReport `json:"detections"`
}

type TopAnomaliesPage struct {
	Anomalies map[string]TopAnomaly `json:"anomalies"`
}

type TopCategory struct {
	Input       uint `json:"input"`
	Output      uint `json:"output"`
	Requests    uint `json:"requests"`
	BeforeBlock uint `json:"before_block"`
	Pending     uint `json:"pending"`
	AfterBlock  uint `json:"after_block"`
}

type TopCategoriesPage struct {
	Categories map[string]TopCategory `json:"categories"`
}

type Report struct {
	From                time.Time                    `json:"from"`
	To                  time.Time                    `json:"to"`
	MainActivityPage    MainActivityPage             `json:"main_activity_page"`
	DeviceAnalyticsPage DevicesAnalyticsPage         `json:"device_analytics_page"`
	AnomaliesListPage   map[string]AnomaliesListPage `json:"anomalies_list_page"`
	TopAnomaliesPage    TopAnomaliesPage             `json:"top_anomalies_page"`
	TopCategoriesPage   map[string]TopCategoriesPage `json:"top_categories_page"`
}
