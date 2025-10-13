package content

import (
	"context"

	"scrapper/internal/domain"
	"scrapper/internal/infrastructure/httpclient"
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

func (s *HTTPStrategy) Fetch(ctx context.Context, url string) (domain.ContentData, error) {
	body, err := s.client.Fetch(ctx, url)
	if err != nil {
		return domain.ContentData{}, err
	}

	return domain.ContentData{RawHTML: body}, nil
}
