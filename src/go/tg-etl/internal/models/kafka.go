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
	RequestID   string            `json:"request_id"`
	URL         string            `json:"url"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	ImageURL    string            `json:"image_url"`
	Metadata    map[string]string `json:"metadata"`
	Timestamp   time.Time         `json:"timestamp"`
}

// MLAnalysisResult структура для результатов ML анализа
type MLAnalysisResult struct {
	RequestID    string                 `json:"request_id"`
	URL          string                 `json:"url"`
	Category     string                 `json:"category"`
	Sentiment    string                 `json:"sentiment"`
	Keywords     []string               `json:"keywords"`
	Confidence   float64                `json:"confidence"`
	AnalysisData map[string]interface{} `json:"analysis_data"`
	Timestamp    time.Time              `json:"timestamp"`
}
