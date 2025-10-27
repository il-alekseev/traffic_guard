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

	"scrapper/internal/metrics"
	"scrapper/internal/models"
)

// ContentStrategy describes a strategy capable of fetching content for a URL.
type ContentStrategy interface {
	Name() string
	Fetch(ctx context.Context, url string) (models.ContentData, error)
}

// DNSResolver resolves hostnames to IP addresses.
type DNSResolver interface {
	LookupIP(ctx context.Context, host string) ([]net.IP, error)
}

// GeoIPProvider enriches IP addresses with Geo information.
type GeoIPProvider interface {
	Lookup(ctx context.Context, ip net.IP) (models.GeoInfo, error)
}

// ContentCleaner produces plain text from raw HTML content.
type ContentCleaner interface {
	Clean(raw string) (string, error)
	CleanFragment(fragment string) (string, error)
}

// QualityEvaluator scores the extracted text and decides whether it is acceptable.
type QualityEvaluator interface {
	Evaluate(ctx context.Context, data models.ContentData) (models.QualityEvaluation, error)
}

// AttemptRepository persists every content acquisition attempt for later post-processing.
type AttemptRepository interface {
	SaveAttempt(ctx context.Context, attempt models.ContentAttempt) error
}

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

// ScrapeService orchestrates the scraping workflow across strategies and persistence.
type ScrapeService struct {
	strategies     []ContentStrategy
	dnsResolver    DNSResolver
	geoProvider    GeoIPProvider
	cleaner        ContentCleaner
	evaluator      QualityEvaluator
	maxContentLen  int
	repository     AttemptRepository
	logger         Logger
	workers        int
	requestTimeout time.Duration
	metrics        *metrics.Metrics
}

// NewScrapeService wires the core dependencies and returns a configured ScrapeService.
func NewScrapeService(strategies []ContentStrategy, dnsResolver DNSResolver, geoProvider GeoIPProvider, cleaner ContentCleaner, evaluator QualityEvaluator, maxContentLen int, repository AttemptRepository, logger Logger, workers int, requestTimeout time.Duration, metrics *metrics.Metrics) *ScrapeService {
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
		maxContentLen:  maxContentLen,
		repository:     repository,
		logger:         logger,
		workers:        workers,
		requestTimeout: requestTimeout,
		metrics:        metrics,
	}
}

func (s *ScrapeService) Scrape(ctx context.Context, urls []string) []models.ScrapeResult {
	results := make([]models.ScrapeResult, len(urls))
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
		result models.ScrapeResult
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
						result: models.ScrapeResult{
							URL:    job.url,
							Status: models.StatusError,
							Error:  err.Error(),
						},
					}
					continue
				}

				s.debug("job assigned", "worker_id", id, "url", job.url)
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

func (s *ScrapeService) ProcessRequest(ctx context.Context, workerID int, req models.ProcessingRequest) models.RequestOutcome {
	if s.metrics != nil {
		s.metrics.IncInFlight()
		defer s.metrics.DecInFlight()
		reqID := strings.TrimSpace(req.RequestID)
		s.metrics.TrackStart(reqID)
		defer s.metrics.TrackDone(reqID)
	}

	outcome := models.RequestOutcome{
		Request:     req,
		Attempts:    make([]models.AttemptOutcome, 0, 4),
		ProcessedAt: time.Now().UTC(),
	}
	if req.Dst.ContentID != nil {
		outcome.ContentID = strings.TrimSpace(*req.Dst.ContentID)
	}

	candidates := req.Dst.CandidateURLs()
	if len(candidates) == 0 {
		outcome.Best = models.AttemptOutcome{
			URL: "",
			Result: models.ScrapeResult{
				URL:    "",
				Status: models.StatusError,
				Error:  "no candidate urls resolved",
			},
		}
		return outcome
	}

	var best models.AttemptOutcome
	bestInitialized := false

	for _, candidate := range candidates {
		result := s.handle(ctx, workerID, candidate)
		attempt := models.AttemptOutcome{
			URL:    candidate,
			Result: result,
		}
		outcome.Attempts = append(outcome.Attempts, attempt)

		if !bestInitialized || isBetterResult(result, best.Result) {
			best = attempt
			bestInitialized = true
		}

		if result.Status == models.StatusOK {
			break
		}
	}

	if bestInitialized {
		outcome.Best = best
	}

	if s.metrics != nil {
		switch outcome.Best.Result.Status {
		case models.StatusOK, models.StatusPartial:
			s.metrics.IncProcessed()
		default:
			s.metrics.IncFailed()
		}
	}

	return outcome
}

func (s *ScrapeService) handle(parent context.Context, workerID int, rawURL string) models.ScrapeResult {
	result := models.ScrapeResult{
		URL:    rawURL,
		Status: models.StatusError,
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
		bestAttempt    *models.ContentAttempt
		bestEvaluation models.QualityEvaluation
		accepted       bool
	)
	reason := models.FailureNone

	for _, strategy := range s.strategies {
		s.debug("strategy start", "worker_id", workerID, "host", host, "strategy", strategy.Name())
		attempt := models.ContentAttempt{
			URL:       rawURL,
			Strategy:  strategy.Name(),
			Timestamp: time.Now().UTC(),
		}

		payload, fetchErr := strategy.Fetch(reqCtx, rawURL)
		if fetchErr != nil {
			attempt.Error = fetchErr.Error()
			s.warn("strategy fetch failed", "worker_id", workerID, "host", host, "strategy", strategy.Name(), "error", fetchErr.Error())
			if err := s.persistAttempt(reqCtx, attempt); err != nil {
				warnings = append(warnings, fmt.Sprintf("persist attempt (%s): %v", attempt.Strategy, err))
			}
			warnings = append(warnings, fmt.Sprintf("strategy %s fetch error: %v", strategy.Name(), fetchErr))
			reason = models.FailureFetch
			continue
		}

		attempt.Success = true
		attempt.Content.RawHTML = payload.RawHTML
		attempt.Content.Text = strings.TrimSpace(payload.Text)
		attempt.Content.UserAgent = payload.UserAgent

		if s.cleaner != nil {
			var (
				cleaned  string
				cleanErr error
				source   string
				original = attempt.Content.Text
			)

			if attempt.Content.RawHTML != "" {
				source = "html"
				cleaned, cleanErr = s.cleaner.Clean(attempt.Content.RawHTML)
			} else {
				source = "text"
				cleaned, cleanErr = s.cleaner.CleanFragment(attempt.Content.Text)
			}

			if cleanErr != nil {
				s.warn("content cleaning warning", "worker_id", workerID, "host", host, "strategy", strategy.Name(), "user_agent", attempt.Content.UserAgent, "source", source, "error", cleanErr.Error())
				warnings = append(warnings, fmt.Sprintf("strategy %s cleaning warning: %v", strategy.Name(), cleanErr))
			}

			if strings.TrimSpace(cleaned) == "" {
				attempt.Success = false
				attempt.Error = "clean content: empty"
				s.warn("content cleaning produced empty result", "worker_id", workerID, "host", host, "strategy", strategy.Name(), "user_agent", attempt.Content.UserAgent, "source", source)
				if err := s.persistAttempt(reqCtx, attempt); err != nil {
					warnings = append(warnings, fmt.Sprintf("persist attempt (%s): %v", attempt.Strategy, err))
				}
				warnings = append(warnings, fmt.Sprintf("strategy %s cleaning empty output", strategy.Name()))
				reason = models.FailureClean
				continue
			}

			attempt.Content.Text = cleaned
			if attempt.Content.RawHTML == "" {
				attempt.Content.RawHTML = original
			}
		}

		eval, evalErr := s.evaluator.Evaluate(reqCtx, attempt.Content)
		if evalErr != nil {
			attempt.Success = false
			attempt.Error = fmt.Sprintf("evaluate content: %v", evalErr)
			s.warn("content evaluation failed", "worker_id", workerID, "host", host, "strategy", strategy.Name(), "user_agent", attempt.Content.UserAgent, "error", evalErr.Error())
			if err := s.persistAttempt(reqCtx, attempt); err != nil {
				warnings = append(warnings, fmt.Sprintf("persist attempt (%s): %v", attempt.Strategy, err))
			}
			warnings = append(warnings, fmt.Sprintf("strategy %s evaluation error: %v", strategy.Name(), evalErr))
			reason = models.FailureEvaluate
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
			s.info("strategy accepted", "worker_id", workerID, "host", host, "strategy", strategy.Name(), "user_agent", attempt.Content.UserAgent, "metric", eval.Score, "text_len", len(attempt.Content.Text))
			reason = models.FailureNone
			break
		}

		if attempt.Content.Text == "" {
			continue
		}

		if bestAttempt == nil || eval.Score > bestEvaluation.Score {
			copy := attempt
			bestAttempt = &copy
			bestEvaluation = eval
			s.debug("strategy candidate", "worker_id", workerID, "host", host, "strategy", strategy.Name(), "user_agent", attempt.Content.UserAgent, "metric", eval.Score, "text_len", len(attempt.Content.Text))
		}
	}

	if bestAttempt != nil {
		contentText, truncated := s.truncate(bestAttempt.Content.Text)
		result.Content = contentText
		result.Strategy = bestAttempt.Strategy
		result.Metric = bestEvaluation.Score
		result.UserAgent = bestAttempt.Content.UserAgent

		if result.Content == "" {
			warnings = append(warnings, "no textual content extracted")
			result.Status = models.StatusError
			reason = models.FailureEmpty
		} else if accepted && len(warnings) == 0 {
			result.Status = models.StatusOK
		} else {
			if !accepted {
				warnings = append(warnings, "content requirements not satisfied")
				reason = models.FailureRequirements
			}
			result.Status = models.StatusPartial
		}
		if truncated {
			warnings = append(warnings, fmt.Sprintf("content truncated to %d characters", s.maxContentLen))
		}
	} else {
		warnings = append(warnings, "no strategy succeeded")
		result.Status = models.StatusError
		reason = models.FailureUnavailable
		s.error("scrape failed", "worker_id", workerID, "host", host)
	}

	if len(warnings) > 0 {
		result.Warnings = append(result.Warnings, warnings...)
	}

	if reason == models.FailureNone && result.Status != models.StatusOK {
		if len(result.Warnings) > 0 {
			reason = models.FailureRequirements
		} else {
			reason = models.FailureUnavailable
		}
	}

	if result.Strategy != "" {
		args := []any{"worker_id", workerID, "host", host, "status", result.Status, "strategy", result.Strategy, "metric", result.Metric, "reason", reason}
		if result.UserAgent != "" {
			args = append(args, "user_agent", result.UserAgent)
		}
		s.info("scrape completed", args...)
	} else {
		args := []any{"worker_id", workerID, "host", host, "status", result.Status, "reason", reason}
		if result.UserAgent != "" {
			args = append(args, "user_agent", result.UserAgent)
		}
		s.info("scrape completed", args...)
	}

	if len(result.Warnings) == 0 {
		result.Warnings = nil
	}

	if reason != models.FailureNone {
		result.Failure = reason
		result.Error = failureMessage(reason)
	} else {
		result.Error = ""
	}

	return result
}

func (s *ScrapeService) persistAttempt(ctx context.Context, attempt models.ContentAttempt) error {
	if s.repository == nil {
		return nil
	}

	return s.repository.SaveAttempt(ctx, attempt)
}

func (s *ScrapeService) debug(msg string, args ...any) {
	if s.logger == nil {
		return
	}
	s.logger.Debug(msg, args...)
}

func (s *ScrapeService) info(msg string, args ...any) {
	if s.logger == nil {
		return
	}
	s.logger.Info(msg, args...)
}

func (s *ScrapeService) warn(msg string, args ...any) {
	if s.logger == nil {
		return
	}
	s.logger.Warn(msg, args...)
}

func (s *ScrapeService) error(msg string, args ...any) {
	if s.logger == nil {
		return
	}
	s.logger.Error(msg, args...)
}

func (s *ScrapeService) truncate(text string) (string, bool) {
	if s.maxContentLen <= 0 || len(text) == 0 {
		return text, false
	}
	count := 0
	for idx := range text {
		if count == s.maxContentLen {
			return text[:idx], true
		}
		count++
	}
	return text, false
}

func (s *ScrapeService) Workers() int {
	if s == nil {
		return 0
	}
	return s.workers
}

func isBetterResult(current, previous models.ScrapeResult) bool {
	currentRank := rankStatus(current.Status)
	previousRank := rankStatus(previous.Status)
	if currentRank != previousRank {
		return currentRank > previousRank
	}

	if current.Metric != previous.Metric {
		return current.Metric > previous.Metric
	}

	return len(current.Content) > len(previous.Content)
}

func rankStatus(status models.ScrapeStatus) int {
	switch status {
	case models.StatusOK:
		return 3
	case models.StatusPartial:
		return 2
	case models.StatusError:
		return 1
	default:
		return 0
	}
}

func failureMessage(reason models.FailureReason) string {
	switch reason {
	case models.FailureFetch:
		return "failed to fetch content"
	case models.FailureClean:
		return "failed to clean content"
	case models.FailureEvaluate:
		return "content did not pass evaluation"
	case models.FailureEmpty:
		return "extracted content empty"
	case models.FailureRequirements:
		return "content requirements not satisfied"
	case models.FailureUnavailable:
		return "no strategy succeeded"
	default:
		return ""
	}
}
