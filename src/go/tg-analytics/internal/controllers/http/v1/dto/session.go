package dto

import "time"

type Session struct {
	ID          uint      `json:"id"`
	DatetimeUTC time.Time `json:"datetime_utc"`
	Type        string    `json:"type"`
	Status      string    `json:"status"`
	URL         string    `json:"url"`
	Proto       string    `json:"proto"`
	HostName    string    `json:"hostname"`
	SrcIP       string    `json:"src_ip"`
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
	IP            string    `json:"ip"`
	Port          int       `json:"port"`
	Location      string    `json:"location"`
	Domain        string    `json:"domain"`
	RequestCount  int       `json:"request_count"`
	HostName      string    `json:"hostname"`
	Category      string    `json:"category"`
	Description   string    `json:"description"`
	Action        string    `json:"action"`
	CategorizedAt time.Time `json:"categorized_at"`
	NegRate       float32   `json:"neg_rate"`
}

type Detection_v2 struct {
	IP            string         `json:"ip"`
	Port          int            `json:"port"`
	Location      string         `json:"location"`
	Domain        string         `json:"domain"`
	RequestCount  int            `json:"request_count"`
	HostName      string         `json:"hostname"`
	Category      string         `json:"category"`
	Description   string         `json:"description"`
	Action        string         `json:"action"`
	CategorizedAt time.Time      `json:"categorized_at"`
	NegRate       float32        `json:"neg_rate"`
	Categories    []CategoryStat `json:"categories"`
}

type CategoryStat struct {
	Name string  `json:"name"`
	Rate float32 `json:"rate"`
}

type DetectionStat struct {
	Detected   int64 `json:"detected"`
	Allowed    int64 `json:"accepted"`
	Denied     int64 `json:"denied"`
	Unresolved int64 `json:"unresolved"`
}

type RequestPoint struct {
	Date  time.Time `json:"date"`
	Value int       `json:"value"`
}

type GetSessionsResponse struct {
	Data []Session `json:"data"` // Список сессий
	//Meta PaginationMeta `json:"meta"` // Метаданные пагинации // TODO: Костыль ниже! потом эту строчку расскоментировать, а нижние убрать
	Count uint `json:"count"` // Число сессий в ответе
	Total uint `json:"total"` // Общее количество сессий
	Pages uint `json:"pages"` // Общее количество страниц
}

type GetDetectionsResponse struct {
	Data []Detection    `json:"data"` // Список выявлений
	Meta PaginationMeta `json:"meta"` // Метаданные пагинации
}

type UnresolvedDetection struct {
	Domain         string `json:"domain"`
	RequestsAll    uint   `json:"requests_all"`
	RequestsBefore uint   `json:"requests_before"`
	RequestsAfter  uint   `json:"requests_after"`
}
