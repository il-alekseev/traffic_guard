package models

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

type DetectionStat struct {
	Detected   int `json:"detected"`
	Accepted   int `json:"accepted"`
	Denied     int `json:"denied"`
	Unresolved int `json:"unresolved"`
}

type Resource struct {
	Name            string `json:"name"`
	ReqsBeforeBlock int    `json:"reqs_before_block"`
	ReqsAfterBlock  int    `json:"reqs_after_block"`
}

type ResourcePoint struct {
	Date    int64 `json:"date"`
	Locked  int   `json:"locked"`
	Delayed int   `json:"delayed"`
}

type DeviceStat struct {
	Name       string          `json:"name"`
	State      string          `json:"state"`
	Statistics []ResourcePoint `json:"statistics"`
}

type TrafficPoint struct {
	Date   int64 `json:"date"`
	Input  int   `json:"input"`
	Output int   `json:"output"`
}

type RequestPoint struct {
	Date  int64 `json:"date"`
	Count int   `json:"count"`
}

type RequestsStat struct {
	Accepted      []RequestPoint `json:"accepted"`
	Blocked       []RequestPoint `json:"blocked"`
	BeforeBlocked []RequestPoint `json:"before_blocked"`
	Delayed       []RequestPoint `json:"delayed"`
}

type AnomalyResourse struct {
	ResourseName      string `json:"resourse_name"`
	UnblockedRequests int    `json:"unblocked_requests"`
}

type Anomaly struct {
	BlockedResourcesCount int               `json:"blocked_resources_count"`
	FirewallsCount        int               `json:"firewalls_count"`
	Stat                  []AnomalyResourse `json:"stat"`
}
