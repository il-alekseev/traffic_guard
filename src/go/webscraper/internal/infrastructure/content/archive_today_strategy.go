package content

import (
	"context"
	"errors"
	"strings"

	"scrapper/internal/infrastructure/httpclient"
	"scrapper/internal/models"
)

type ArchiveTodayStrategy struct {
	client     *httpclient.Client
	userAgents []string
}

func NewArchiveTodayStrategy(client *httpclient.Client, userAgents []string) *ArchiveTodayStrategy {
	if client == nil {
		return nil
	}
	return &ArchiveTodayStrategy{
		client:     client,
		userAgents: append([]string(nil), userAgents...),
	}
}

func (s *ArchiveTodayStrategy) Name() string { return "archive.today" }

func (s *ArchiveTodayStrategy) Fetch(ctx context.Context, rawURL string) (models.ContentData, error) {
	if s == nil || s.client == nil {
		return models.ContentData{}, errors.New("archive.today strategy not configured")
	}

	target := strings.TrimSpace(rawURL)
	if target == "" {
		return models.ContentData{}, errors.New("empty url")
	}

	archiveURL := "https://archive.today/" + target

	agents := s.userAgents
	if len(agents) == 0 {
		agents = []string{httpclient.DefaultUserAgent}
	}

	var lastErr error
	for _, agent := range agents {
		body, usedAgent, err := s.client.FetchWithCustomAgent(ctx, archiveURL, agent)
		if err == nil {
			return models.ContentData{RawHTML: body, UserAgent: usedAgent}, nil
		}
		lastErr = err
	}

	if lastErr == nil {
		lastErr = errors.New("no user agents available for archive.today")
	}

	return models.ContentData{}, lastErr
}

func (s *ArchiveTodayStrategy) FetchWithUserAgent(ctx context.Context, rawURL, userAgent string) (models.ContentData, error) {
	if s == nil || s.client == nil {
		return models.ContentData{}, errors.New("archive.today strategy not configured")
	}

	target := strings.TrimSpace(rawURL)
	if target == "" {
		return models.ContentData{}, errors.New("empty url")
	}

	archiveURL := "https://archive.today/" + target
	body, usedAgent, err := s.client.FetchWithCustomAgent(ctx, archiveURL, userAgent)
	if err != nil {
		return models.ContentData{}, err
	}

	return models.ContentData{RawHTML: body, UserAgent: usedAgent}, nil
}
