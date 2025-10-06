package models

// Определяем вложенные структуры
type DayStat struct {
	Accepted int `json:"accepted"`
	Denied   int `json:"denied"`
}

type DateMap struct {
	// Используем map[string]Stat, так как даты могут быть разными
	// и их количество может варьироваться
	DateMap map[string]DayStat `json:"-"`
}

type RequestStat struct {
	RequestStat []DateMap `json:"request_stat"`
}

type Notification struct {
	Description string `json:"description"`
	Username    string `json:"username"`
	DateTime    string `json:"dateTime"`
}

type Detection struct {
	URL         string `json:"url"`
	DateTime    string `json:"dateTime"`
	Category    string `json:"category"`
	Description string `json:"description"`
	Host        string `json:"host"`
	NGFW        string `json:"ngfw"`
	AccessCount int64  `json:"accessCount"`
	Action      string `json:"action"`
	// TODO: добавить поле для контроля, что пользователь данное действие уже выполнил или нет и дату
}

type Session struct {
	Status      string `json:"status"`
	URL         string `json:"url"`
	IP          string `json:"ip"`
	NGFW        string `json:"ngfw"`
	Username    string `json:"username"`
	SessionType string `json:"sessionType"`
	Category    string `json:"category"`
	DateTime    string `json:"dateTime"`
}
