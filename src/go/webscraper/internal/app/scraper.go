package app

import (
	"context"
	"sync"

	"scrapper/internal/metrics"
	"scrapper/internal/models"
	"scrapper/internal/usecase"
)

type Scraper struct {
	service *usecase.ScrapeService
	metrics *metrics.Metrics
}

func NewScraper(service *usecase.ScrapeService, collector *metrics.Metrics) *Scraper {
	return &Scraper{
		service: service,
		metrics: collector,
	}
}

func (s *Scraper) Stream(ctx context.Context, workers int, requests <-chan models.ProcessingRequest) <-chan models.RequestOutcome {
	if s == nil || s.service == nil {
		ch := make(chan models.RequestOutcome)
		close(ch)
		return ch
	}

	if workers < 1 {
		workers = 1
	}

	output := make(chan models.RequestOutcome)
	wg := sync.WaitGroup{}
	wg.Add(workers)

	for workerID := 0; workerID < workers; workerID++ {
		id := workerID
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case req, ok := <-requests:
					if !ok {
						return
					}
					if s.metrics != nil {
						s.metrics.DecQueued()
					}
					outcome := s.service.ProcessRequest(ctx, id, req)
					select {
					case <-ctx.Done():
						return
					case output <- outcome:
					}
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(output)
	}()

	return output
}
