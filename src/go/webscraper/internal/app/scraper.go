package app

import (
	"context"

	"scrapper/internal/domain"
	"scrapper/internal/usecase"
)

type Scraper struct {
	service *usecase.ScrapeService
}

func NewScraper(service *usecase.ScrapeService) *Scraper {
	return &Scraper{service: service}
}

func (s *Scraper) Run(ctx context.Context, urls []string) []domain.ScrapeResult {
	if s == nil || s.service == nil {
		return nil
	}

	return s.service.Scrape(ctx, urls)
}
