package content

import (
	"context"

	"scrapper/internal/infrastructure/httpclient"
	"scrapper/internal/models"
)

type HTTPStrategy struct {
	client *httpclient.Client
}

func NewHTTPStrategy(client *httpclient.Client) *HTTPStrategy {
	return &HTTPStrategy{client: client}
}

func (s *HTTPStrategy) Name() string {
	return "http"
}

func (s *HTTPStrategy) Fetch(ctx context.Context, url string) (models.ContentData, error) {
	body, agent, err := s.client.Fetch(ctx, url)
	if err != nil {
		return models.ContentData{}, err
	}

	return models.ContentData{RawHTML: body, UserAgent: agent}, nil
}

func (s *HTTPStrategy) FetchWithUserAgent(ctx context.Context, url, agent string) (models.ContentData, error) {
	body, usedAgent, err := s.client.FetchWithCustomAgent(ctx, url, agent)
	if err != nil {
		return models.ContentData{}, err
	}
	return models.ContentData{RawHTML: body, UserAgent: usedAgent}, nil
}
