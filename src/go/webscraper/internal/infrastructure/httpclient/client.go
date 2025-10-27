package httpclient

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	nethttp "net/http"
	"net/url"
	"sync"
	"time"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/transform"
)

const DefaultUserAgent = "Scrapper/0.1"

type Client struct {
	client     *nethttp.Client
	userAgents []string
	mu         sync.RWMutex
	hostAgent  map[string]string
}

func NewClient(timeout time.Duration, userAgents []string) *Client {
	if timeout <= 0 {
		timeout = 30 * time.Second
	}

	transport := nethttp.DefaultTransport.(*nethttp.Transport).Clone()
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	copiedAgents := append([]string(nil), userAgents...)

	return &Client{
		client: &nethttp.Client{
			Timeout:   timeout,
			Transport: transport,
		},
		userAgents: copiedAgents,
		hostAgent:  make(map[string]string),
	}
}

func (c *Client) Fetch(ctx context.Context, url string) (string, string, error) {
	agents, host := c.orderedAgents(url)

	var lastErr error
	for _, agent := range agents {
		body, err := c.fetchWithAgent(ctx, url, agent)
		if err == nil {
			if host != "" {
				c.setPreferredAgent(host, agent)
			}
			return body, agent, nil
		}
		lastErr = fmt.Errorf("user-agent %q: %w", agent, err)
		if host != "" {
			c.clearPreferredAgent(host, agent)
		}
	}

	if lastErr == nil {
		lastErr = fmt.Errorf("no user agents available")
	}

	return "", "", lastErr
}

func (c *Client) fetchWithAgent(ctx context.Context, url, agent string) (string, error) {
	req, err := nethttp.NewRequestWithContext(ctx, nethttp.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", agent)

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		io.Copy(io.Discard, resp.Body)
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
		return string(body), nil
	}

	return string(decoded), nil
}

// Client exposes the underlying *http.Client for components that need
// lower-level access (e.g., third-party libraries). Use sparingly.
func (c *Client) Client() *nethttp.Client { return c.client }

func (c *Client) orderedAgents(rawURL string) ([]string, string) {
	base := c.userAgents
	if len(base) == 0 {
		base = []string{DefaultUserAgent}
	}

	u, err := url.Parse(rawURL)
	if err != nil || u.Hostname() == "" {
		return append([]string(nil), base...), ""
	}

	host := u.Hostname()

	c.mu.RLock()
	preferred, ok := c.hostAgent[host]
	c.mu.RUnlock()
	if !ok || preferred == "" {
		return append([]string(nil), base...), host
	}

	ordered := make([]string, 0, len(base)+1)
	ordered = append(ordered, preferred)
	for _, agent := range base {
		if agent != preferred {
			ordered = append(ordered, agent)
		}
	}

	return ordered, host
}

func (c *Client) setPreferredAgent(host, agent string) {
	c.mu.Lock()
	c.hostAgent[host] = agent
	c.mu.Unlock()
}

func (c *Client) clearPreferredAgent(host, failed string) {
	c.mu.Lock()
	if current, ok := c.hostAgent[host]; ok && current == failed {
		delete(c.hostAgent, host)
	}
	c.mu.Unlock()
}
