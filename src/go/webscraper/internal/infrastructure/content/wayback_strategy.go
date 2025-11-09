package content

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/transform"

	"scrapper/internal/infrastructure/httpclient"
	"scrapper/internal/infrastructure/parser"
	"scrapper/internal/models"
)

const (
	waybackSnapshotLookback   = 7 * 24 * time.Hour
	waybackMaxSnapshots       = 5
	waybackMinCleanTextLength = 300
)

var waybackCleaner = parser.NewHTMLCleaner()

type WaybackStrategy struct {
	client     *httpclient.Client
	http       *http.Client
	userAgents []string
}

func NewWaybackStrategy(client *httpclient.Client, userAgents []string) *WaybackStrategy {
	if client == nil {
		return nil
	}
	httpClient := client.Client()
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	return &WaybackStrategy{
		client:     client,
		http:       httpClient,
		userAgents: append([]string(nil), userAgents...),
	}
}

func (s *WaybackStrategy) Name() string { return "wayback" }

func (s *WaybackStrategy) Fetch(ctx context.Context, rawURL string) (models.ContentData, error) {
	if s == nil || s.client == nil || s.http == nil {
		return models.ContentData{}, errors.New("wayback strategy not configured")
	}

	snapshots, err := s.lookupSnapshots(ctx, rawURL)
	if err != nil {
		return models.ContentData{}, err
	}

	agents := s.userAgents
	if len(agents) == 0 {
		agents = []string{httpclient.DefaultUserAgent}
	}

	var lastErr error
snapshotLoop:
	for _, snap := range snapshots {
		for _, agent := range agents {
			data, err := s.downloadSnapshot(ctx, snap.URL, agent)
			if err != nil {
				lastErr = err
				continue
			}
			if !waybackPayloadUsable(data.RawHTML) {
				lastErr = fmt.Errorf("wayback snapshot %s rejected: insufficient content", snap.Timestamp)
				continue snapshotLoop
			}
			return data, nil
		}
	}

	if lastErr == nil {
		lastErr = errors.New("no user agents available for wayback")
	}

	return models.ContentData{}, lastErr
}

func (s *WaybackStrategy) lookupSnapshot(ctx context.Context, rawURL string) (string, error) {
	if rawURL == "" {
		return "", errors.New("empty url")
	}

	now := time.Now().UTC().Format("20060102150405")
	apiURL := fmt.Sprintf("https://archive.org/wayback/available?url=%s&timestamp=%s", url.QueryEscape(rawURL), now)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return "", err
	}

	resp, err := s.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("wayback api status %d", resp.StatusCode)
	}

	var data struct {
		ArchivedSnapshots struct {
			Closest struct {
				URL       string `json:"url"`
				Available bool   `json:"available"`
				Status    string `json:"status"`
				Timestamp string `json:"timestamp"`
			} `json:"closest"`
		} `json:"archived_snapshots"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", fmt.Errorf("parse wayback response: %w", err)
	}

	closest := data.ArchivedSnapshots.Closest
	if !closest.Available || closest.URL == "" {
		return "", errors.New("no wayback snapshot available")
	}

	return normalizeWaybackSnapshotURL(closest.URL), nil
}

func (s *WaybackStrategy) FetchWithUserAgent(ctx context.Context, rawURL, userAgent string) (models.ContentData, error) {
	if s == nil || s.client == nil {
		return models.ContentData{}, errors.New("wayback strategy not configured")
	}
	snapshots, err := s.lookupSnapshots(ctx, rawURL)
	if err != nil {
		return models.ContentData{}, err
	}
	for _, snap := range snapshots {
		data, err := s.downloadSnapshot(ctx, snap.URL, userAgent)
		if err != nil {
			continue
		}
		if waybackPayloadUsable(data.RawHTML) {
			return data, nil
		}
	}
	return models.ContentData{}, errors.New("no usable wayback snapshot")
}

func normalizeWaybackSnapshotURL(value string) string {
	if value == "" {
		return value
	}
	const marker = "/web/"
	idx := strings.Index(value, marker)
	if idx == -1 {
		return value
	}
	start := idx + len(marker)
	end := strings.Index(value[start:], "/")
	if end == -1 {
		return value
	}
	timestamp := value[start : start+end]
	if timestamp == "" || strings.HasSuffix(timestamp, "id_") || strings.HasSuffix(timestamp, "if_") || strings.HasSuffix(timestamp, "0id_") {
		return value
	}
	return value[:idx] + "/web/" + timestamp + "id_" + value[start+end:]
}

func (s *WaybackStrategy) downloadSnapshot(ctx context.Context, snapshotURL, agent string) (models.ContentData, error) {
	if snapshotURL == "" {
		return models.ContentData{}, errors.New("empty snapshot url")
	}
	ua := strings.TrimSpace(agent)
	if ua == "" {
		ua = httpclient.DefaultUserAgent
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, snapshotURL, nil)
	if err != nil {
		return models.ContentData{}, err
	}
	req.Header.Set("User-Agent", ua)
	resp, err := s.http.Do(req)
	if err != nil {
		return models.ContentData{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		io.Copy(io.Discard, resp.Body)
		return models.ContentData{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.ContentData{}, err
	}
	contentType := resp.Header.Get("X-Archive-Orig-Content-Type")
	if contentType == "" {
		contentType = resp.Header.Get("X-Archive-Orig-Content-Encoding")
	}
	if contentType == "" {
		contentType = resp.Header.Get("Content-Type")
	}
	encoding, _, _ := charset.DetermineEncoding(body, contentType)
	reader := transform.NewReader(bytes.NewReader(body), encoding.NewDecoder())
	decoded, err := io.ReadAll(reader)
	if err != nil {
		decoded = body
	}
	return models.ContentData{RawHTML: string(decoded), UserAgent: ua}, nil
}

type waybackSnapshot struct {
	URL       string
	Timestamp string
}

func (s *WaybackStrategy) lookupSnapshots(ctx context.Context, rawURL string) ([]waybackSnapshot, error) {
	seen := make(map[string]struct{})
	result := make([]waybackSnapshot, 0, waybackMaxSnapshots)

	if snap, err := s.lookupClosestSnapshot(ctx, rawURL); err == nil {
		result = append(result, snap)
		seen[snap.URL] = struct{}{}
	}

	if len(result) < waybackMaxSnapshots {
		if extra, err := s.lookupCDXSnapshots(ctx, rawURL, seen, waybackMaxSnapshots-len(result)); err == nil {
			result = append(result, extra...)
		}
	}

	if len(result) == 0 {
		return nil, errors.New("no wayback snapshot available")
	}

	return result, nil
}

func (s *WaybackStrategy) lookupClosestSnapshot(ctx context.Context, rawURL string) (waybackSnapshot, error) {
	if rawURL == "" {
		return waybackSnapshot{}, errors.New("empty url")
	}

	now := time.Now().UTC().Format("20060102150405")
	apiURL := fmt.Sprintf("https://archive.org/wayback/available?url=%s&timestamp=%s", url.QueryEscape(rawURL), now)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return waybackSnapshot{}, err
	}

	resp, err := s.http.Do(req)
	if err != nil {
		return waybackSnapshot{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return waybackSnapshot{}, fmt.Errorf("wayback api status %d", resp.StatusCode)
	}

	var data struct {
		ArchivedSnapshots struct {
			Closest struct {
				URL       string `json:"url"`
				Available bool   `json:"available"`
				Status    string `json:"status"`
				Timestamp string `json:"timestamp"`
			} `json:"closest"`
		} `json:"archived_snapshots"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return waybackSnapshot{}, fmt.Errorf("parse wayback response: %w", err)
	}

	closest := data.ArchivedSnapshots.Closest
	if !closest.Available || closest.URL == "" {
		return waybackSnapshot{}, errors.New("no wayback snapshot available")
	}

	return waybackSnapshot{
		URL:       normalizeWaybackSnapshotURL(closest.URL),
		Timestamp: closest.Timestamp,
	}, nil
}

func (s *WaybackStrategy) lookupCDXSnapshots(ctx context.Context, rawURL string, seen map[string]struct{}, maxNeeded int) ([]waybackSnapshot, error) {
	if maxNeeded <= 0 {
		return nil, nil
	}

	now := time.Now().UTC()
	from := now.Add(-waybackSnapshotLookback)

	cdxURL := fmt.Sprintf(
		"https://web.archive.org/cdx/search/cdx?url=%s&output=json&filter=statuscode:200&collapse=digest&fl=timestamp,original&from=%s&to=%s&limit=%d",
		url.QueryEscape(rawURL),
		from.Format("20060102150405"),
		now.Format("20060102150405"),
		waybackMaxSnapshots*4,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, cdxURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("wayback cdx status %d", resp.StatusCode)
	}

	var rows [][]string
	if err := json.NewDecoder(resp.Body).Decode(&rows); err != nil {
		return nil, fmt.Errorf("parse wayback cdx response: %w", err)
	}
	if len(rows) <= 1 {
		return nil, errors.New("no wayback cdx snapshots")
	}

	candidates := make([]waybackSnapshot, 0, maxNeeded)
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) < 2 {
			continue
		}
		ts := strings.TrimSpace(row[0])
		orig := strings.TrimSpace(row[1])
		if ts == "" || orig == "" {
			continue
		}
		candidate := normalizeWaybackSnapshotURL(fmt.Sprintf("https://web.archive.org/web/%sid_/%s", ts, orig))
		if candidate == "" {
			continue
		}
		if _, exists := seen[candidate]; exists {
			continue
		}
		seen[candidate] = struct{}{}
		candidates = append(candidates, waybackSnapshot{URL: candidate, Timestamp: ts})
		if len(candidates) >= maxNeeded {
			break
		}
	}

	sort.SliceStable(candidates, func(i, j int) bool {
		return candidates[i].Timestamp > candidates[j].Timestamp
	})

	return candidates, nil
}

func waybackPayloadUsable(raw string) bool {
	if strings.TrimSpace(raw) == "" || waybackCleaner == nil {
		return false
	}
	text, err := waybackCleaner.Clean(raw)
	if err != nil {
		return len(strings.TrimSpace(text)) > waybackMinCleanTextLength
	}
	return len(strings.TrimSpace(text)) >= waybackMinCleanTextLength
}
