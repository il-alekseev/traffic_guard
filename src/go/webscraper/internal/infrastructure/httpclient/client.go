package httpclient

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	nethttp "net/http"
	"time"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/transform"
)

type Client struct {
	client *nethttp.Client
}

func NewClient(timeout time.Duration) *Client {
	if timeout <= 0 {
		timeout = 30 * time.Second
	}

	transport := nethttp.DefaultTransport.(*nethttp.Transport).Clone()
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	return &Client{
		client: &nethttp.Client{
			Timeout:   timeout,
			Transport: transport,
		},
	}
}

func (c *Client) Fetch(ctx context.Context, url string) (string, error) {
	req, err := nethttp.NewRequestWithContext(ctx, nethttp.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "Scrapper/0.1")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	encoding, _, _ := charset.DetermineEncoding(body, resp.Header.Get("Content-Type"))
	reader := transform.NewReader(bytes.NewReader(body), encoding.NewDecoder())

	decoded, err := io.ReadAll(reader)
	if err != nil {
		// Fallback to original bytes if decoding fails.
		return string(body), nil
	}

	return string(decoded), nil
}

// Client exposes the underlying *http.Client for components that need
// lower-level access (e.g., third-party libraries). Use sparingly.
func (c *Client) Client() *nethttp.Client { return c.client }
