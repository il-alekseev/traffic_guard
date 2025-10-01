package domain

import "time"

type ScrapeStatus string

const (
	StatusOK      ScrapeStatus = "ok"
	StatusPartial ScrapeStatus = "partial"
	StatusError   ScrapeStatus = "error"
)

type FailureReason string

const (
	FailureNone         FailureReason = "none"
	FailureFetch        FailureReason = "fetch_error"
	FailureClean        FailureReason = "clean_error"
	FailureEvaluate     FailureReason = "evaluate_error"
	FailureEmpty        FailureReason = "empty_content"
	FailureRequirements FailureReason = "requirements_not_met"
	FailureUnavailable  FailureReason = "no_successful_strategy"
)

type GeoInfo map[string]any

type DomainInfo struct {
	Name string  `json:"name"`
	IP   string  `json:"ip,omitempty"`
	Geo  GeoInfo `json:"geo,omitempty"`
}

type ContentData struct {
	RawHTML string `json:"raw_html,omitempty"`
	Text    string `json:"text,omitempty"`
}

type QualityEvaluation struct {
	Score    float64 `json:"score"`
	Accepted bool    `json:"accepted"`
}

type ContentAttempt struct {
	URL       string      `json:"url"`
	Strategy  string      `json:"strategy"`
	Success   bool        `json:"success"`
	Error     string      `json:"error,omitempty"`
	Metric    float64     `json:"metric"`
	Accepted  bool        `json:"accepted"`
	Content   ContentData `json:"content"`
	Timestamp time.Time   `json:"timestamp"`
}

type ScrapeResult struct {
	URL      string        `json:"url"`
	Domain   DomainInfo    `json:"domain"`
	Status   ScrapeStatus  `json:"status"`
	Error    string        `json:"error,omitempty"`
	Content  string        `json:"content,omitempty"`
	Strategy string        `json:"strategy,omitempty"`
	Metric   float64       `json:"metric,omitempty"`
	Warnings []string      `json:"warnings,omitempty"`
	Failure  FailureReason `json:"failure,omitempty"`
}
