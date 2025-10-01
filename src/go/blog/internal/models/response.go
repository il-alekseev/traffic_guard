package models

type Logs struct {
	Data []BusinessLog `json:"data"`
	Meta *Meta         `json:"meta"`
}
type Meta struct {
	Limit int `json:"limit"`
	Page  int `json:"page"`
	Pages int `json:"pages"`
	Total int `json:"total"`
}
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
