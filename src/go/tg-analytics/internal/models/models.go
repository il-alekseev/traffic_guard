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
