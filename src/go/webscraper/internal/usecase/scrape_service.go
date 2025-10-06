package usecase

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"runtime"
	"strings"
	"sync"
	"time"

	"scrapper/internal/domain"
)

// ContentStrategy describes a strategy capable of fetching content for a URL.
type ContentStrategy interface {
	Name() string
	Fetch(ctx context.Context, url string) (domain.ContentData, error)
}

// DNSResolver resolves hostnames to IP addresses.
type DNSResolver interface {
	LookupIP(ctx context.Context, host string) ([]net.IP, error)
}

// GeoIPProvider enriches IP addresses with Geo information.
type GeoIPProvider interface {
	Lookup(ctx context.Context, ip net.IP) (domain.GeoInfo, error)
}

// ContentCleaner produces plain text from raw HTML content.
type ContentCleaner interface {
	Clean(raw string) (string, error)
}

// QualityEvaluator scores the extracted text and decides whether it is acceptable.
type QualityEvaluator interface {
	Evaluate(ctx context.Context, data domain.ContentData) (domain.QualityEvaluation, error)
}

// AttemptRepository persists every content acquisition attempt for later post-processing.
type AttemptRepository interface {
	SaveAttempt(ctx context.Context, attempt domain.ContentAttempt) error
}

type Logger interface {
	Logf(format string, args ...interface{})
}

// ScrapeService orchestrates the scraping workflow across strategies and persistence.
type ScrapeService struct {
	strategies     []ContentStrategy
	dnsResolver    DNSResolver
	geoProvider    GeoIPProvider
	cleaner        ContentCleaner
	evaluator      QualityEvaluator
	repository     AttemptRepository
	logger         Logger
	workers        int
	requestTimeout time.Duration
}

// NewScrapeService wires the core dependencies and returns a configured ScrapeService.
func NewScrapeService(strategies []ContentStrategy, dnsResolver DNSResolver, geoProvider GeoIPProvider, cleaner ContentCleaner, evaluator QualityEvaluator, repository AttemptRepository, logger Logger, workers int, requestTimeout time.Duration) *ScrapeService {
	if len(strategies) == 0 {
		panic("scrape service requires at least one content strategy")
	}

	if workers <= 0 {
		workers = runtime.NumCPU()
		if workers < 1 {
			workers = 1
		}
	}

	return &ScrapeService{
		strategies:     strategies,
		dnsResolver:    dnsResolver,
		geoProvider:    geoProvider,
		cleaner:        cleaner,
		evaluator:      evaluator,
		repository:     repository,
		logger:         logger,
		workers:        workers,
		requestTimeout: requestTimeout,
	}
}

func (s *ScrapeService) Scrape(ctx context.Context, urls []string) []domain.ScrapeResult {
	results := make([]domain.ScrapeResult, len(urls))
	if len(urls) == 0 {
		return results
	}

	workerCount := s.workers
	if workerCount > len(urls) {
		workerCount = len(urls)
	}
	if workerCount < 1 {
		workerCount = 1
	}

	type job struct {
		index int
		url   string
	}

	type item struct {
		index  int
		result domain.ScrapeResult
	}

	jobs := make(chan job)
	output := make(chan item)

	var wg sync.WaitGroup
	wg.Add(workerCount)

	for i := 0; i < workerCount; i++ {
		workerID := i
		go func(id int) {
			defer wg.Done()
			for job := range jobs {
				if err := ctx.Err(); err != nil {
					output <- item{
						index: job.index,
						result: domain.ScrapeResult{
							URL:    job.url,
							Status: domain.StatusError,
							Error:  err.Error(),
						},
					}
					continue
				}

				s.logf("worker=%d url=%s assigned", id, job.url)
				res := s.handle(ctx, id, job.url)
				output <- item{index: job.index, result: res}
			}
		}(workerID)
	}

	go func() {
		wg.Wait()
		close(output)
	}()

	go func() {
		defer close(jobs)
		for idx, rawURL := range urls {
			select {
			case <-ctx.Done():
				return
			case jobs <- job{index: idx, url: rawURL}:
			}
		}
	}()

	for item := range output {
		results[item.index] = item.result
	}

	return results
}

func (s *ScrapeService) handle(parent context.Context, workerID int, rawURL string) domain.ScrapeResult {
	result := domain.ScrapeResult{
		URL:    rawURL,
		Status: domain.StatusError,
	}

	if err := parent.Err(); err != nil {
		result.Error = err.Error()
		return result
	}

	reqCtx := parent
	var cancel context.CancelFunc
	if s.requestTimeout > 0 {
		reqCtx, cancel = context.WithTimeout(parent, s.requestTimeout)
		defer cancel()
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		result.Error = fmt.Sprintf("invalid url: %v", err)
		return result
	}

	host := parsed.Hostname()
	if host == "" {
		result.Error = "missing host in url"
		return result
	}

	result.Domain.Name = host

	ips, err := s.dnsResolver.LookupIP(reqCtx, host)
	if err != nil {
		result.Error = fmt.Sprintf("dns lookup: %v", err)
		return result
	}

	var primaryIP net.IP
	for _, ip := range ips {
		if ip != nil {
			primaryIP = ip
			break
		}
	}

	warnings := make([]string, 0, 1)

	if primaryIP != nil {
		result.Domain.IP = primaryIP.String()
		if s.geoProvider != nil {
			geoInfo, geoErr := s.geoProvider.Lookup(reqCtx, primaryIP)
			if geoErr != nil {
				warnings = append(warnings, fmt.Sprintf("geo lookup: %v", geoErr))
			} else if len(geoInfo) > 0 {
				result.Domain.Geo = geoInfo
			}
		}
	} else {
		warnings = append(warnings, "no ip addresses found for host")
	}

	var (
		bestAttempt    *domain.ContentAttempt
		bestEvaluation domain.QualityEvaluation
		accepted       bool
	)
	reason := domain.FailureNone

	for _, strategy := range s.strategies {
		s.logf("worker=%d host=%s strategy=%s status=start", workerID, host, strategy.Name())
		attempt := domain.ContentAttempt{
			URL:       rawURL,
			Strategy:  strategy.Name(),
			Timestamp: time.Now().UTC(),
		}

		payload, fetchErr := strategy.Fetch(reqCtx, rawURL)
		if fetchErr != nil {
			attempt.Error = fetchErr.Error()
			s.logf("worker=%d host=%s strategy=%s status=fetch_error error=%q", workerID, host, strategy.Name(), fetchErr)
			if err := s.persistAttempt(reqCtx, attempt); err != nil {
				warnings = append(warnings, fmt.Sprintf("persist attempt (%s): %v", attempt.Strategy, err))
			}
			warnings = append(warnings, fmt.Sprintf("strategy %s fetch error: %v", strategy.Name(), fetchErr))
			reason = domain.FailureFetch
			continue
		}

		attempt.Success = true
		attempt.Content.RawHTML = payload.RawHTML
		attempt.Content.Text = strings.TrimSpace(payload.Text)

		if s.cleaner != nil && attempt.Content.RawHTML != "" {
			cleaned, cleanErr := s.cleaner.Clean(attempt.Content.RawHTML)
			if cleanErr != nil {
				attempt.Success = false
				attempt.Error = fmt.Sprintf("clean content: %v", cleanErr)
				s.logf("worker=%d host=%s strategy=%s status=clean_error error=%q", workerID, host, strategy.Name(), cleanErr)
				if err := s.persistAttempt(reqCtx, attempt); err != nil {
					warnings = append(warnings, fmt.Sprintf("persist attempt (%s): %v", attempt.Strategy, err))
				}
				warnings = append(warnings, fmt.Sprintf("strategy %s clean error: %v", strategy.Name(), cleanErr))
				reason = domain.FailureClean
				continue
			}
			attempt.Content.Text = cleaned
		}

		eval, evalErr := s.evaluator.Evaluate(reqCtx, attempt.Content)
		if evalErr != nil {
			attempt.Success = false
			attempt.Error = fmt.Sprintf("evaluate content: %v", evalErr)
			s.logf("worker=%d host=%s strategy=%s status=evaluate_error error=%q", workerID, host, strategy.Name(), evalErr)
			if err := s.persistAttempt(reqCtx, attempt); err != nil {
				warnings = append(warnings, fmt.Sprintf("persist attempt (%s): %v", attempt.Strategy, err))
			}
			warnings = append(warnings, fmt.Sprintf("strategy %s evaluation error: %v", strategy.Name(), evalErr))
			reason = domain.FailureEvaluate
			continue
		}

		attempt.Metric = eval.Score
		attempt.Accepted = eval.Accepted

		if err := s.persistAttempt(reqCtx, attempt); err != nil {
			warnings = append(warnings, fmt.Sprintf("persist attempt (%s): %v", attempt.Strategy, err))
		}

		if eval.Accepted {
			copy := attempt
			bestAttempt = &copy
			bestEvaluation = eval
			accepted = true
			s.logf("worker=%d host=%s strategy=%s status=accepted metric=%.2f text_len=%d", workerID, host, strategy.Name(), eval.Score, len(attempt.Content.Text))
			reason = domain.FailureNone
			break
		}

		if attempt.Content.Text == "" {
			continue
		}

		if bestAttempt == nil || eval.Score > bestEvaluation.Score {
			copy := attempt
			bestAttempt = &copy
			bestEvaluation = eval
			s.logf("worker=%d host=%s strategy=%s status=candidate metric=%.2f text_len=%d", workerID, host, strategy.Name(), eval.Score, len(attempt.Content.Text))
		}
	}

	if bestAttempt != nil {
		result.Content = bestAttempt.Content.Text
		result.Strategy = bestAttempt.Strategy
		result.Metric = bestEvaluation.Score

		if result.Content == "" {
			warnings = append(warnings, "no textual content extracted")
			result.Status = domain.StatusError
			reason = domain.FailureEmpty
		} else if accepted && len(warnings) == 0 {
			result.Status = domain.StatusOK
		} else {
			if !accepted {
				warnings = append(warnings, "content requirements not satisfied")
				reason = domain.FailureRequirements
			}
			result.Status = domain.StatusPartial
		}
	} else {
		warnings = append(warnings, "no strategy succeeded")
		result.Status = domain.StatusError
		reason = domain.FailureUnavailable
		s.logf("worker=%d host=%s status=failed", workerID, host)
	}

	if len(warnings) > 0 {
		result.Warnings = append(result.Warnings, warnings...)
	}

	if reason == domain.FailureNone && result.Status != domain.StatusOK {
		if len(result.Warnings) > 0 {
			reason = domain.FailureRequirements
		} else {
			reason = domain.FailureUnavailable
		}
	}

	if result.Strategy != "" {
		s.logf("worker=%d host=%s status=%s strategy=%s metric=%.2f reason=%s", workerID, host, result.Status, result.Strategy, result.Metric, reason)
	} else {
		s.logf("worker=%d host=%s status=%s reason=%s", workerID, host, result.Status, reason)
	}

	if len(result.Warnings) == 0 {
		result.Warnings = nil
	}

	if reason != domain.FailureNone {
		result.Failure = reason
		result.Error = failureMessage(reason)
	} else {
		result.Error = ""
	}

	return result
}

func (s *ScrapeService) persistAttempt(ctx context.Context, attempt domain.ContentAttempt) error {
	if s.repository == nil {
		return nil
	}

	return s.repository.SaveAttempt(ctx, attempt)
}

func (s *ScrapeService) logf(format string, args ...interface{}) {
	if s.logger == nil {
		return
	}
	s.logger.Logf(format, args...)
}

func failureMessage(reason domain.FailureReason) string {
	switch reason {
	case domain.FailureFetch:
		return "failed to fetch content"
	case domain.FailureClean:
		return "failed to clean content"
	case domain.FailureEvaluate:
		return "content did not pass evaluation"
	case domain.FailureEmpty:
		return "extracted content empty"
	case domain.FailureRequirements:
		return "content requirements not satisfied"
	case domain.FailureUnavailable:
		return "no strategy succeeded"
	default:
		return ""
	}
}
