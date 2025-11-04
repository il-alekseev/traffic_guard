package models

import (
	"time"

	"github.com/google/uuid"
)

// AnalysisRequest представляет структуру запроса на анализ
type AnalysisRequest struct {
	RequestID uuid.UUID   `json:"request_id"`
	Src       Src         `json:"src"`
	Dst       Destination `json:"dst"`
	Timestamp time.Time   `json:"timestamp"`
}

// Source содержит информацию об источнике запроса
type Src struct {
	IP string `json:"ip"`
}

// Destination содержит информацию о цели анализа
type Destination struct {
	Type      DestinationType `json:"type"`
	Resource  string          `json:"resource"`
	Port      *int            `json:"port,omitempty"`
	Proto     *Protocol       `json:"proto,omitempty"`
	ContentID *string         `json:"content_id,omitempty"`
}

// DestinationType определяет типы ресурсов
type DestinationType string

const (
	DestinationIP     DestinationType = "ip"
	DestinationDomain DestinationType = "domain"
	DestinationURL    DestinationType = "url"
)

// Protocol определяет поддерживаемые протоколы
type Protocol string

const (
	ProtocolHTTP  Protocol = "http"
	ProtocolHTTPS Protocol = "https"
)

// URLMetadataResult структура для результатов метаданных URL
type URLMetadataResult struct {
	RequestID   string     `json:"request_id"`
	Status      string     `json:"status"`
	Failure     string     `json:"failure"`
	Error       string     `json:"error"`
	Warnings    []string   `json:"warnings"`
	URL         string     `json:"url"`
	Domain      DomainInfo `json:"domain"`
	Strategy    string     `json:"strategy"`
	Metric      int        `json:"metric"`
	UserAgent   string     `json:"user_agent"`
	ProcessedAt string     `json:"processed_at"`
}

type DomainInfo struct {
	Name string `json:"name"`
	IP   string `json:"ip"`
	Geo  Geo    `json:"geo"`
}

type Geo struct {
	AsDomain      string `json:"as_domain"`
	AsName        string `json:"as_name"`
	ASN           string `json:"asn"`
	Continent     string `json:"continent"`
	ContinentCode string `json:"continent_code"`
	Country       string `json:"country"`
	CountryCode   string `json:"country_code"`
}

// MLAnalysisResult структура для результатов ML анализа
type MLAnalysisResult struct {
	RequestID       string `json:"REQUEST_ID"`
	ContentGUID     string `json:"CONTENT_GUID"`
	RecognisedClass string `json:"RECOGNISED_CLASS"`
}
