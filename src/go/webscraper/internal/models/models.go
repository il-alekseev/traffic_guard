package models

import (
	"errors"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"
)

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
	FailureDNS          FailureReason = "dns_error"
	FailureAntiBot      FailureReason = "anti_bot_detected"
	FailureShortContent FailureReason = "short_content"
	FailureValidation   FailureReason = "invalid_content"
	FailureQuality      FailureReason = "quality_rejected"
	FailureFiltered     FailureReason = "filtered_url"
	FailureUnavailable  FailureReason = "no_successful_strategy"
)

type GeoInfo map[string]any

type DomainInfo struct {
	Name string  `json:"name"`
	IP   string  `json:"ip,omitempty"`
	Geo  GeoInfo `json:"geo,omitempty"`
}

type ContentData struct {
	RawHTML   string `json:"raw_html,omitempty"`
	Text      string `json:"text,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
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
	URL       string        `json:"url"`
	Domain    DomainInfo    `json:"domain"`
	Status    ScrapeStatus  `json:"status"`
	Error     string        `json:"error,omitempty"`
	Content   string        `json:"content,omitempty"`
	Strategy  string        `json:"strategy,omitempty"`
	Metric    float64       `json:"metric,omitempty"`
	UserAgent string        `json:"user_agent,omitempty"`
	Warnings  []string      `json:"warnings,omitempty"`
	Failure   FailureReason `json:"failure,omitempty"`
}

type ProcessingRequest struct {
	RequestID string             `json:"request_id"`
	Src       RequestSource      `json:"src"`
	Dst       RequestDestination `json:"dst"`
	Received  time.Time          `json:"-"`
	Raw       map[string]any     `json:"-"`
}

type RequestSource struct {
	IP string `json:"ip"`
}

type RequestDestination struct {
	Type      DestinationType `json:"type"`
	Resource  string          `json:"resource"`
	Port      *int            `json:"port,omitempty"`
	Proto     *Protocol       `json:"proto,omitempty"`
	ContentID *string         `json:"content_id,omitempty"`
}

type DestinationType string

const (
	DestinationIP     DestinationType = "ip"
	DestinationDomain DestinationType = "domain"
	DestinationURL    DestinationType = "url"
)

type Protocol string

const (
	ProtocolHTTP  Protocol = "http"
	ProtocolHTTPS Protocol = "https"
)

type AttemptOutcome struct {
	URL    string       `json:"url"`
	Result ScrapeResult `json:"result"`
}

type RequestOutcome struct {
	Request     ProcessingRequest `json:"request"`
	Best        AttemptOutcome    `json:"best"`
	Attempts    []AttemptOutcome  `json:"attempts"`
	ProcessedAt time.Time         `json:"processed_at"`
	ContentID   string            `json:"content_id,omitempty"`
}

type MetadataMessage struct {
	RequestID   string        `json:"request_id"`
	ContentID   string        `json:"content_id,omitempty"`
	Status      ScrapeStatus  `json:"status"`
	Failure     FailureReason `json:"failure,omitempty"`
	Error       string        `json:"error,omitempty"`
	Warnings    []string      `json:"warnings,omitempty"`
	URL         string        `json:"url"`
	Domain      DomainInfo    `json:"domain"`
	Strategy    string        `json:"strategy,omitempty"`
	Metric      float64       `json:"metric,omitempty"`
	UserAgent   string        `json:"user_agent,omitempty"`
	ProcessedAt time.Time     `json:"processed_at"`
}

type ContentMessage struct {
	RequestID   string        `json:"request_id"`
	ContentID   string        `json:"content_id,omitempty"`
	URL         string        `json:"url"`
	Strategy    string        `json:"strategy,omitempty"`
	Metric      float64       `json:"metric,omitempty"`
	Content     string        `json:"content,omitempty"`
	Status      ScrapeStatus  `json:"status"`
	Failure     FailureReason `json:"failure,omitempty"`
	Error       string        `json:"error,omitempty"`
	UserAgent   string        `json:"user_agent,omitempty"`
	ProcessedAt time.Time     `json:"processed_at"`
}

func (d RequestDestination) CandidateURLs() []string {
	resource := strings.TrimSpace(d.Resource)
	if resource == "" {
		return nil
	}

	seen := make(map[string]struct{})
	candidates := make([]string, 0, 8)

	parsed, err := url.Parse(resource)
	if err == nil && parsed.Scheme != "" && parsed.Host != "" {
		normalized := parsed.String()
		seen[normalized] = struct{}{}
		candidates = append(candidates, normalized)
		return candidates
	}

	protos := d.candidateProtocols()
	ports := d.candidatePorts()
	portExplicit := d.Port != nil

	for _, proto := range protos {
		for _, port := range ports {
			value, err := d.buildURL(string(proto), port, portExplicit)
			if err != nil {
				continue
			}
			if _, ok := seen[value]; ok {
				continue
			}
			seen[value] = struct{}{}
			candidates = append(candidates, value)
		}
	}

	return candidates
}

func (d RequestDestination) candidateProtocols() []Protocol {
	if d.Proto != nil {
		p := Protocol(strings.ToLower(string(*d.Proto)))
		if p == "" {
			return []Protocol{ProtocolHTTPS, ProtocolHTTP}
		}
		return []Protocol{p}
	}
	return []Protocol{ProtocolHTTPS, ProtocolHTTP}
}

func (d RequestDestination) candidatePorts() []int {
	if d.Port != nil {
		return []int{*d.Port}
	}
	return []int{443, 8443, 80, 8080}
}

func (d RequestDestination) buildURL(proto string, port int, portExplicit bool) (string, error) {
	resource := strings.TrimSpace(d.Resource)
	if resource == "" {
		return "", errors.New("empty resource")
	}

	parsed, err := parseResource(resource, proto)
	if err != nil {
		return "", err
	}

	host := parsed.Host
	path := ensureLeadingSlash(parsed.Path)
	if parsed.RawPath != "" {
		path = parsed.RawPath
	}

	if host == "" {
		host = extractHostFallback(resource)
		if host == "" {
			return "", errors.New("destination host unresolved")
		}

		if idx := strings.Index(resource, "/"); idx >= 0 {
			path = resource[idx:]
		} else {
			path = ""
		}
	}

	baseHost, err := stripPort(host)
	if err != nil {
		return "", err
	}

	normalizedHost := trimIPv6Brackets(baseHost)

	hostPort := normalizedHost
	if portExplicit || !isStandardPort(proto, port) {
		hostPort = net.JoinHostPort(normalizedHost, strconv.Itoa(port))
	}

	u := url.URL{
		Scheme:   proto,
		Host:     hostPort,
		Path:     path,
		RawQuery: parsed.RawQuery,
		Fragment: parsed.Fragment,
	}

	// Preserve trailing slash for empty path if present in original.
	if path == "" && strings.HasSuffix(parsed.Path, "/") {
		u.Path = "/"
	}

	return u.String(), nil
}

func parseResource(resource, proto string) (*url.URL, error) {
	if strings.Contains(resource, "://") {
		return url.Parse(resource)
	}
	return url.Parse(proto + "://" + resource)
}

func ensureLeadingSlash(path string) string {
	if path == "" || strings.HasPrefix(path, "/") {
		return path
	}
	return "/" + path
}

func extractHostFallback(resource string) string {
	res := strings.TrimPrefix(resource, "//")
	if idx := strings.Index(res, "/"); idx >= 0 {
		return res[:idx]
	}
	return res
}

func stripPort(value string) (string, error) {
	if value == "" {
		return value, nil
	}

	host := value
	// Already bracketed IPv6 with optional port.
	if strings.HasPrefix(host, "[") {
		if strings.Contains(host, "]:") {
			h, _, err := net.SplitHostPort(host)
			if err != nil {
				return "", err
			}
			return h, nil
		}
		return host, nil
	}

	h, _, err := net.SplitHostPort(host)
	if err != nil {
		var addrErr *net.AddrError
		if errors.As(err, &addrErr) && strings.Contains(addrErr.Err, "missing port") {
			return host, nil
		}
		return host, nil
	}
	return h, nil
}

func trimIPv6Brackets(host string) string {
	if strings.HasPrefix(host, "[") && strings.HasSuffix(host, "]") {
		return host[1 : len(host)-1]
	}
	return host
}

func isStandardPort(proto string, port int) bool {
	switch strings.ToLower(proto) {
	case "https":
		return port == 443
	case "http":
		return port == 80
	default:
		return false
	}
}
