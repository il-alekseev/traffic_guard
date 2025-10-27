package dto

import "time"

type Session struct {
	ID          uint      `json:"id"`
	DatetimeUTC time.Time `json:"datetime_utc"`
	Type        string    `json:"type"`
	Status      string    `json:"status"`
	URL         string    `json:"url"`
	Proto       string    `json:"proto"`
	HostName    string    `json:"host_name"`
	SrcIP       string    `json:"src_ip"`
	SrcPort     int       `json:"src_port"`
	SrcCountry  string    `json:"src_country"`
	Username    string    `json:"username"`
	DstIP       string    `json:"dst_ip"`
	DstPort     int       `json:"dst_port"`
	DstCountry  string    `json:"dst_country"`
	Category    string    `json:"category"`
}

type Category struct {
	Name        string `json:"name"`
	AccessCount int    `json:"access_count"`
}

type Detection struct {
	IP                 string    `json:"ip"`
	Port               int       `json:"port"`
	Country            string    `json:"country"`
	Path               string    `json:"path"`
	AccessCount        int       `json:"access_count"`
	HostName           string    `json:"host_name"`
	Category           string    `json:"category"`
	Description        string    `json:"description"`
	Decision           string    `json:"decision"`
	LastAccessDatetime time.Time `json:"last_access_datetime"`
}

type DetectionStat struct {
	Detected   int `json:"detected"`
	Accepted   int `json:"accepted"`
	Denied     int `json:"denied"`
	Unresolved int `json:"unresolved"`
}

type RequestPoint struct {
	Date  time.Time `json:"date"`
	Value int       `json:"value"`
}
