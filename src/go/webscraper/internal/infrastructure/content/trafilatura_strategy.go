package content

// TrafilaturaStrategy uses github.com/markusmobius/go-trafilatura to fetch and
// extract the readable text of a page. It performs its own HTTP access and
// returns cleaned text via ContentData.Text (leaving RawHTML empty).

import (
	"bytes"
	"context"
	"io"
	"net/http"
	nurl "net/url"
	"time"

	trafilatura "github.com/markusmobius/go-trafilatura"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/transform"

	"scrapper/internal/domain"
)

type TrafilaturaStrategy struct{ client *http.Client }

func NewTrafilaturaStrategy(client *http.Client) (*TrafilaturaStrategy, error) { // nolint:revive
	if client == nil {
		client = &http.Client{Timeout: 30 * time.Second}
	}
	return &TrafilaturaStrategy{client: client}, nil
}

func (s *TrafilaturaStrategy) Name() string { return "trafilatura" }

func (s *TrafilaturaStrategy) Fetch(ctx context.Context, url string) (domain.ContentData, error) {
	// Build HTTP request with context and a sane timeout.
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return domain.ContentData{}, err
	}
	req.Header.Set("User-Agent", "Scrapper/0.1 (trafilatura)")

	resp, err := s.client.Do(req)
	if err != nil {
		return domain.ContentData{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return domain.ContentData{}, err
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
		return domain.ContentData{}, err
	}

	return domain.ContentData{Text: result.ContentText}, nil
}
