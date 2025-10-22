package content

// TrafilaturaStrategy uses github.com/markusmobius/go-trafilatura to fetch and
// extract the readable text of a page. It performs its own HTTP access and
// returns cleaned text via ContentData.Text (leaving RawHTML empty).

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	nurl "net/url"
	"time"

	trafilatura "github.com/markusmobius/go-trafilatura"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/transform"

	"scrapper/internal/infrastructure/httpclient"
	"scrapper/internal/models"
)

type TrafilaturaStrategy struct {
	client     *http.Client
	userAgents []string
}

func NewTrafilaturaStrategy(client *http.Client, userAgents []string) (*TrafilaturaStrategy, error) { // nolint:revive
	if client == nil {
		client = &http.Client{Timeout: 30 * time.Second}
	}
	return &TrafilaturaStrategy{
		client:     client,
		userAgents: append([]string(nil), userAgents...),
	}, nil
}

func (s *TrafilaturaStrategy) Name() string { return "trafilatura" }

func (s *TrafilaturaStrategy) Fetch(ctx context.Context, url string) (models.ContentData, error) {
	agents := s.userAgents
	if len(agents) == 0 {
		agents = []string{httpclient.DefaultUserAgent}
	}

	var lastErr error
	for _, agent := range agents {
		data, err := s.fetchWithAgent(ctx, url, agent)
		if err == nil {
			return data, nil
		}
		lastErr = fmt.Errorf("user-agent %q: %w", agent, err)
	}

	if lastErr == nil {
		lastErr = errors.New("no user agents available")
	}

	return models.ContentData{}, lastErr
}

func (s *TrafilaturaStrategy) fetchWithAgent(ctx context.Context, url, agent string) (models.ContentData, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return models.ContentData{}, err
	}
	req.Header.Set("User-Agent", agent)

	resp, err := s.client.Do(req)
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

	encoding, _, _ := charset.DetermineEncoding(body, resp.Header.Get("Content-Type"))
	reader := transform.NewReader(bytes.NewReader(body), encoding.NewDecoder())
	decoded, err := io.ReadAll(reader)
	if err != nil {
		decoded = body
	}

	parsedURL, _ := nurl.ParseRequestURI(url)
	opts := trafilatura.Options{OriginalURL: parsedURL}
	result, err := trafilatura.Extract(bytes.NewReader(decoded), opts)
	if err != nil {
		return models.ContentData{}, err
	}

	return models.ContentData{Text: result.ContentText, UserAgent: agent}, nil
}
