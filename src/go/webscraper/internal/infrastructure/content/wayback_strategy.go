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
	"strings"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/transform"

	"scrapper/internal/infrastructure/httpclient"
	"scrapper/internal/models"
)

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

	snapshotURL, err := s.lookupSnapshot(ctx, rawURL)
	if err != nil {
		return models.ContentData{}, err
	}

	agents := s.userAgents
	if len(agents) == 0 {
		agents = []string{httpclient.DefaultUserAgent}
	}

	var lastErr error
	for _, agent := range agents {
		data, err := s.downloadSnapshot(ctx, snapshotURL, agent)
		if err == nil {
			return data, nil
		}
		lastErr = err
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

	apiURL := fmt.Sprintf("https://archive.org/wayback/available?url=%s", url.QueryEscape(rawURL))
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
	snapshotURL, err := s.lookupSnapshot(ctx, rawURL)
	if err != nil {
		return models.ContentData{}, err
	}
	return s.downloadSnapshot(ctx, snapshotURL, userAgent)
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
	return value[:idx] + "/web/0id_" + value[start+end:]
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
