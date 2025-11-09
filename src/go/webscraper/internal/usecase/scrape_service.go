package usecase

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"net"
	"net/url"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"scrapper/internal/metrics"
	"scrapper/internal/models"
)

// ContentStrategy describes a strategy capable of fetching content for a URL.
type ContentStrategy interface {
	Name() string
	Fetch(ctx context.Context, url string) (models.ContentData, error)
}

type CustomAgentStrategy interface {
	FetchWithUserAgent(ctx context.Context, url, userAgent string) (models.ContentData, error)
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

const (
	binaryProbeLimit          = 2048
	scriptProbeLimit          = 2048
	binaryEntropyThreshold    = 4.5
	printableRatioThreshold   = 0.85
	shortContentFallbackLimit = 150
	maxUserAgentRetries       = 3
)

var (
	disallowedURLExtensions = []struct {
		suffix   string
		category string
	}{
		{suffix: ".jpg", category: "image"},
		{suffix: ".jpeg", category: "image"},
		{suffix: ".png", category: "image"},
		{suffix: ".gif", category: "image"},
		{suffix: ".bmp", category: "image"},
		{suffix: ".svg", category: "image"},
		{suffix: ".webp", category: "image"},
		{suffix: ".ico", category: "image"},
		{suffix: ".mp4", category: "video"},
		{suffix: ".mkv", category: "video"},
		{suffix: ".avi", category: "video"},
		{suffix: ".mov", category: "video"},
		{suffix: ".wmv", category: "video"},
		{suffix: ".flv", category: "video"},
		{suffix: ".webm", category: "video"},
		{suffix: ".mp3", category: "audio"},
		{suffix: ".wav", category: "audio"},
		{suffix: ".ogg", category: "audio"},
		{suffix: ".m4a", category: "audio"},
		{suffix: ".m4s", category: "audio"},
		{suffix: ".pdf", category: "document"},
		{suffix: ".doc", category: "document"},
		{suffix: ".docx", category: "document"},
		{suffix: ".ppt", category: "document"},
		{suffix: ".pptx", category: "document"},
		{suffix: ".xls", category: "document"},
		{suffix: ".xlsx", category: "document"},
		{suffix: ".ods", category: "document"},
		{suffix: ".odt", category: "document"},
		{suffix: ".txt", category: "text"},
		{suffix: ".zip", category: "archive"},
		{suffix: ".rar", category: "archive"},
		{suffix: ".7z", category: "archive"},
		{suffix: ".tar", category: "archive"},
		{suffix: ".gz", category: "archive"},
		{suffix: ".bz2", category: "archive"},
		{suffix: ".iso", category: "archive"},
		{suffix: ".js", category: "script"},
		{suffix: ".css", category: "stylesheet"},
		{suffix: ".map", category: "sourcemap"},
		{suffix: ".exe", category: "binary"},
		{suffix: ".dll", category: "binary"},
		{suffix: ".bin", category: "binary"},
		{suffix: ".apk", category: "binary"},
		{suffix: ".deb", category: "binary"},
		{suffix: ".rpm", category: "binary"},
		{suffix: ".ttf", category: "font"},
		{suffix: ".woff", category: "font"},
		{suffix: ".woff2", category: "font"},
		{suffix: ".eot", category: "font"},
	}

	disallowedHostPrefixes = []struct {
		prefix  string
		message string
	}{
		{prefix: "cdn.", message: "domain %s looks like CDN/static host; skipping"},
		{prefix: "assets.", message: "domain %s looks like CDN/static host; skipping"},
		{prefix: "static.", message: "domain %s looks like CDN/static host; skipping"},
		{prefix: "media.", message: "domain %s serves media assets; skipping"},
		{prefix: "files.", message: "domain %s serves files; skipping"},
		{prefix: "uploads.", message: "domain %s looks like file storage; skipping"},
		{prefix: "storage.", message: "domain %s looks like file storage; skipping"},
		{prefix: "s3.", message: "domain %s looks like S3/CloudFront resource; skipping"},
		{prefix: "download.", message: "domain %s is for downloads; skipping"},
	}

	disallowedHostFragments = []struct {
		fragment string
		message  string
	}{
		{fragment: ".s3.", message: "domain %s looks like S3/CloudFront resource; skipping"},
		{fragment: ".cloudfront.", message: "domain %s looks like CDN/CloudFront resource; skipping"},
		{fragment: ".storage.googleapis.com", message: "domain %s looks like GCS/S3 resource; skipping"},
		{fragment: ".cdn.", message: "domain %s looks like CDN; skipping"},
	}

	disallowedPathFragments = []struct {
		fragment string
		message  string
	}{
		{fragment: "/download", message: "path contains \"%s\" and looks like file download"},
		{fragment: "/attachment", message: "path contains \"%s\" and looks like attachment"},
		{fragment: "/files/", message: "path contains \"%s\" and looks like file storage"},
		{fragment: "/static/", message: "path contains \"%s\" and serves static assets"},
		{fragment: "/assets/", message: "path contains \"%s\" and serves static assets"},
		{fragment: "/media/", message: "path contains \"%s\" and serves media catalog"},
		{fragment: "/uploads/", message: "path contains \"%s\" and serves uploaded files"},
		{fragment: "/upload/", message: "path contains \"%s\" and serves uploaded files"},
		{fragment: "/videos/", message: "path contains \"%s\" and serves video catalog"},
		{fragment: "/images/", message: "path contains \"%s\" and serves image catalog"},
		{fragment: "/thumbnails/", message: "path contains \"%s\" and serves thumbnails"},
		{fragment: "/sitemap.xml", message: "path contains \"%s\" (sitemap) and is skipped"},
	}

	disallowedQueryFragments = []struct {
		fragment string
		message  string
	}{
		{fragment: "download=", message: "query contains \"%s\" and requests file download"},
		{fragment: "attachment=", message: "query contains \"%s\" and returns attachment"},
		{fragment: "file=", message: "query contains \"%s\" and returns file"},
		{fragment: "action=download", message: "query contains \"%s\" and forces download"},
		{fragment: "mode=download", message: "query contains \"%s\" and forces download"},
		{fragment: "type=download", message: "query contains \"%s\" and forces download"},
		{fragment: "type=file", message: "query contains \"%s\" and returns file"},
		{fragment: "format=pdf", message: "query contains \"%s\" and requests PDF"},
		{fragment: "format=doc", message: "query contains \"%s\" and requests document"},
		{fragment: "format=docx", message: "query contains \"%s\" and requests document"},
		{fragment: "format=ppt", message: "query contains \"%s\" and requests presentation"},
		{fragment: "format=pptx", message: "query contains \"%s\" and requests presentation"},
		{fragment: "format=xls", message: "query contains \"%s\" and requests spreadsheet"},
		{fragment: "format=xlsx", message: "query contains \"%s\" and requests spreadsheet"},
		{fragment: "format=ods", message: "query contains \"%s\" and requests spreadsheet"},
		{fragment: "format=odt", message: "query contains \"%s\" and requests document"},
		{fragment: "format=zip", message: "query contains \"%s\" and requests archive"},
		{fragment: "format=rar", message: "query contains \"%s\" and requests archive"},
		{fragment: "format=7z", message: "query contains \"%s\" and requests archive"},
		{fragment: "format=tar", message: "query contains \"%s\" and requests archive"},
		{fragment: "format=gz", message: "query contains \"%s\" and requests archive"},
		{fragment: "format=bz2", message: "query contains \"%s\" and requests archive"},
		{fragment: "format=iso", message: "query contains \"%s\" and requests disk image"},
	}

	binarySignatures = []struct {
		prefix  string
		message string
	}{
		{prefix: "\x89PNG\r\n\x1a\n", message: "payload matches PNG signature"},
		{prefix: "GIF87a", message: "payload matches GIF signature"},
		{prefix: "GIF89a", message: "payload matches GIF signature"},
		{prefix: "RIFF", message: "payload matches RIFF/WEBP or media container"},
		{prefix: "ftyp", message: "payload matches MP4/ISO BMFF signature"},
		{prefix: "moof", message: "payload matches fragmented MP4 signature"},
		{prefix: "%PDF-", message: "payload matches PDF signature"},
		{prefix: "PK\x03\x04", message: "payload matches ZIP/Office archive"},
		{prefix: "\x1f\x8b\x08", message: "payload matches gzip stream"},
		{prefix: "ID3", message: "payload matches MP3 signature"},
		{prefix: "\x00\x00\x01\xba", message: "payload matches MPEG transport stream"},
	}

	scriptKeywords = []string{
		"function ",
		"function(",
		"const ",
		"let ",
		"var ",
		"(()=>",
		"=>",
		"window.",
		"document.",
		"export ",
		"import ",
		"return ",
		"this.",
		"prototype",
		"object.",
		"require(",
		"class ",
	}

	scriptOperatorRunes = map[rune]struct{}{
		'{': {},
		'}': {},
		'(': {},
		')': {},
		'[': {},
		']': {},
		';': {},
		':': {},
		',': {},
		'.': {},
		'+': {},
		'-': {},
		'*': {},
		'/': {},
		'=': {},
		'!': {},
		'>': {},
		'<': {},
		'&': {},
		'|': {},
	}
	retryableStatusCodes = map[int]struct{}{
		401: {},
		403: {},
		429: {},
		451: {},
	}
)

// ScrapeService orchestrates the scraping workflow across strategies and persistence.
type ScrapeService struct {
	strategies          []ContentStrategy
	archiveStrategies   []ContentStrategy
	dnsResolver         DNSResolver
	geoProvider         GeoIPProvider
	cleaner             ContentCleaner
	evaluator           QualityEvaluator
	maxContentLen       int
	repository          AttemptRepository
	logger              Logger
	workers             int
	requestTimeout      time.Duration
	metrics             *metrics.Metrics
	antiBotDetector     *AntiBotDetector
	shortContentLimit   int
	stripJSON           bool
	userAgents          []string
	uaRand              *rand.Rand
	randMu              sync.Mutex
	maxUserAgentRetries int
}

// NewScrapeService wires the core dependencies and returns a configured ScrapeService.
func NewScrapeService(strategies []ContentStrategy, archiveStrategies []ContentStrategy, dnsResolver DNSResolver, geoProvider GeoIPProvider, cleaner ContentCleaner, evaluator QualityEvaluator, maxContentLen int, repository AttemptRepository, logger Logger, workers int, requestTimeout time.Duration, metrics *metrics.Metrics, detector *AntiBotDetector, stripJSON bool, userAgents []string) *ScrapeService {
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
		strategies:          strategies,
		archiveStrategies:   archiveStrategies,
		dnsResolver:         dnsResolver,
		geoProvider:         geoProvider,
		cleaner:             cleaner,
		evaluator:           evaluator,
		maxContentLen:       maxContentLen,
		repository:          repository,
		logger:              logger,
		workers:             workers,
		requestTimeout:      requestTimeout,
		metrics:             metrics,
		antiBotDetector:     detector,
		shortContentLimit:   shortContentFallbackLimit,
		stripJSON:           stripJSON,
		userAgents:          append([]string(nil), userAgents...),
		uaRand:              rand.New(rand.NewSource(time.Now().UnixNano())),
		maxUserAgentRetries: maxUserAgentRetries,
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
	return s.handleWithFallback(parent, workerID, rawURL, true)
}

func (s *ScrapeService) handleWithFallback(parent context.Context, workerID int, rawURL string, allowRootFallback bool) models.ScrapeResult {
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
		result.Failure = models.FailureDNS
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

	if reasons := shouldSkipURL(parsed); len(reasons) > 0 {
		result.Status = models.StatusError
		result.Failure = models.FailureFiltered
		result.Error = strings.Join(reasons, "; ")
		warnings = append(warnings, reasons...)
		result.Warnings = append(result.Warnings, warnings...)
		s.warn("url skipped due to disallowed resource", "worker_id", workerID, "host", host, "url", rawURL, "reasons", strings.Join(reasons, "; "))
		return result
	}

	var (
		bestAttempt    *models.ContentAttempt
		bestEvaluation models.QualityEvaluation
		accepted       bool
	)
	reason := models.FailureNone

	strategyQueue := append([]ContentStrategy(nil), s.strategies...)
	archiveAdded := false

	for idx := 0; idx < len(strategyQueue); idx++ {
		strategy := strategyQueue[idx]
		s.debug("strategy start", "worker_id", workerID, "host", host, "strategy", strategy.Name())
		attempt := models.ContentAttempt{
			URL:       rawURL,
			Strategy:  strategy.Name(),
			Timestamp: time.Now().UTC(),
		}

		payload, fetchErr := strategy.Fetch(reqCtx, rawURL)
		if fetchErr != nil {
			if statusCode, hasStatus := parseStatusCodeFromError(fetchErr); hasStatus {
				if _, retryable := retryableStatusCodes[statusCode]; retryable {
					if newPayload, ok := s.retryWithUserAgent(reqCtx, strategy, rawURL, ""); ok {
						s.info("strategy retry due to http status", "worker_id", workerID, "host", host, "strategy", strategy.Name(), "status_code", statusCode)
						payload = newPayload
						goto payloadLoop
					}
					if len(s.archiveStrategies) > 0 && !archiveAdded {
						strategyQueue = append(strategyQueue, s.archiveStrategies...)
						archiveAdded = true
					}
					warnings = append(warnings, fmt.Sprintf("strategy %s http status %d, archives scheduled", strategy.Name(), statusCode))
				} else if len(s.archiveStrategies) > 0 && !archiveAdded {
					if statusCode == 404 {
						warnings = append(warnings, fmt.Sprintf("strategy %s http status %d, archives skipped", strategy.Name(), statusCode))
					} else {
						strategyQueue = append(strategyQueue, s.archiveStrategies...)
						archiveAdded = true
					}
				}
			} else if len(s.archiveStrategies) > 0 && !archiveAdded {
				strategyQueue = append(strategyQueue, s.archiveStrategies...)
				archiveAdded = true
			}

			attempt.Error = fetchErr.Error()
			s.warn("strategy fetch failed", "worker_id", workerID, "host", host, "strategy", strategy.Name(), "error", fetchErr.Error())
			if err := s.persistAttempt(reqCtx, attempt); err != nil {
				warnings = append(warnings, fmt.Sprintf("persist attempt (%s): %v", attempt.Strategy, err))
			}
			warnings = append(warnings, fmt.Sprintf("strategy %s fetch error: %v", strategy.Name(), fetchErr))
			reason = models.FailureFetch
			continue
		}

	payloadLoop:
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
				if len(s.archiveStrategies) > 0 && !archiveAdded {
					strategyQueue = append(strategyQueue, s.archiveStrategies...)
					archiveAdded = true
				}
				continue
			}

			attempt.Content.Text = cleaned
			if attempt.Content.RawHTML == "" {
				attempt.Content.RawHTML = original
			}
		}

		if s.stripJSON {
			stripped, removed := stripJSONFragments(attempt.Content.Text)
			if removed > 0 {
				attempt.Content.Text = stripped
				warnings = append(warnings, fmt.Sprintf("removed %d json fragments", removed))
			}
		}

		if s.antiBotDetector != nil {
			if phrase := s.antiBotDetector.Check(attempt.Content.Text); phrase != "" {
				if newPayload, ok := s.retryWithUserAgent(reqCtx, strategy, rawURL, attempt.Content.UserAgent); ok {
					s.info("anti-bot detected, retrying with different user agent", "worker_id", workerID, "host", host, "strategy", strategy.Name(), "phrase", phrase)
					payload = newPayload
					goto payloadLoop
				}
				attempt.Success = false
				attempt.Error = "anti-bot content detected"
				warning := fmt.Sprintf("anti-bot phrase detected (%s)", phrase)
				s.warn("anti-bot content detected", "worker_id", workerID, "host", host, "strategy", strategy.Name(), "phrase", phrase)
				if err := s.persistAttempt(reqCtx, attempt); err != nil {
					warnings = append(warnings, fmt.Sprintf("persist attempt (%s): %v", attempt.Strategy, err))
				}
				warnings = append(warnings, warning)
				reason = models.FailureAntiBot
				if len(s.archiveStrategies) > 0 && !archiveAdded {
					strategyQueue = append(strategyQueue, s.archiveStrategies...)
					archiveAdded = true
				}
				continue
			}
		}

		if s.shortContentLimit > 0 {
			runeCount := utf8.RuneCountInString(attempt.Content.Text)
			if runeCount > 0 && runeCount < s.shortContentLimit {
				if newPayload, ok := s.retryWithUserAgent(reqCtx, strategy, rawURL, attempt.Content.UserAgent); ok {
					s.info("content too short, retrying with different user agent", "worker_id", workerID, "host", host, "strategy", strategy.Name(), "runes", runeCount)
					payload = newPayload
					goto payloadLoop
				}
				attempt.Success = false
				attempt.Error = fmt.Sprintf("content too short (%d < %d characters)", runeCount, s.shortContentLimit)
				s.warn("content too short, switching to archive strategies", "worker_id", workerID, "host", host, "strategy", strategy.Name(), "runes", runeCount)
				if err := s.persistAttempt(reqCtx, attempt); err != nil {
					warnings = append(warnings, fmt.Sprintf("persist attempt (%s): %v", attempt.Strategy, err))
				}
				warnings = append(warnings, fmt.Sprintf("strategy %s content too short (%d < %d); trying archive sources", strategy.Name(), runeCount, s.shortContentLimit))
				reason = models.FailureShortContent
				if len(s.archiveStrategies) > 0 && !archiveAdded {
					strategyQueue = append(strategyQueue, s.archiveStrategies...)
					archiveAdded = true
				}
				continue
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

		validationReasons := detectInvalidContent(attempt.Content.RawHTML, attempt.Content.Text)
		if len(validationReasons) > 0 {
			attempt.Success = false
			attempt.Accepted = false
			attempt.Error = strings.Join(validationReasons, "; ")
		}

		if err := s.persistAttempt(reqCtx, attempt); err != nil {
			warnings = append(warnings, fmt.Sprintf("persist attempt (%s): %v", attempt.Strategy, err))
		}

		if len(validationReasons) > 0 {
			warnings = append(warnings, fmt.Sprintf("strategy %s rejected after validation: %s", strategy.Name(), strings.Join(validationReasons, "; ")))
			s.warn("strategy payload rejected after validation", "worker_id", workerID, "host", host, "strategy", strategy.Name(), "user_agent", attempt.Content.UserAgent, "reasons", strings.Join(validationReasons, "; "))
			reason = models.FailureValidation
			continue
		}

		if attempt.Accepted {
			copy := attempt
			bestAttempt = &copy
			bestEvaluation = eval
			accepted = true
			s.info("strategy accepted", "worker_id", workerID, "host", host, "strategy", strategy.Name(), "user_agent", attempt.Content.UserAgent, "metric", eval.Score, "text_len", len(attempt.Content.Text))
			reason = models.FailureNone
			break
		}

		if attempt.Content.Text == "" {
			if len(s.archiveStrategies) > 0 && !archiveAdded {
				strategyQueue = append(strategyQueue, s.archiveStrategies...)
				archiveAdded = true
			}
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
		} else if accepted {
			result.Status = models.StatusOK
		} else {
			if !accepted {
				warnings = append(warnings, "content requirements not satisfied")
				reason = models.FailureQuality
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

	if allowRootFallback && reason == models.FailureUnavailable && shouldFallbackToRoot(parsed) {
		rootURL := buildRootURL(parsed)
		if rootURL != "" && rootURL != rawURL {
			s.warn("no strategy succeeded, retrying with root url", "worker_id", workerID, "host", host, "url", rawURL, "fallback_url", rootURL)
			fallback := s.handleWithFallback(parent, workerID, rootURL, false)
			note := fmt.Sprintf("fallback to root %s after failure for %s", rootURL, rawURL)
			fallback.Warnings = append([]string{note}, fallback.Warnings...)
			return fallback
		}
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
	case models.FailureDNS:
		return "dns lookup failed"
	case models.FailureAntiBot:
		return "anti-bot content detected"
	case models.FailureShortContent:
		return "content too short"
	case models.FailureValidation:
		return "invalid or binary content"
	case models.FailureQuality:
		return "content quality score rejected"
	case models.FailureFiltered:
		return "url filtered by policy"
	case models.FailureUnavailable:
		return "no strategy succeeded"
	default:
		return ""
	}
}

func shouldFallbackToRoot(u *url.URL) bool {
	if u == nil {
		return false
	}
	path := strings.Trim(u.EscapedPath(), "/")
	if path == "" {
		return false
	}
	segments := 0
	for _, part := range strings.Split(path, "/") {
		if part != "" {
			segments++
		}
	}
	if segments >= 2 {
		return true
	}
	if len(path) >= 48 {
		return true
	}
	return len(u.RawQuery) >= 48
}

func buildRootURL(u *url.URL) string {
	if u == nil {
		return ""
	}
	root := &url.URL{
		Scheme: u.Scheme,
		Host:   u.Host,
		Path:   "/",
	}
	if root.Scheme == "" {
		root.Scheme = "https"
	}
	return root.String()
}

func shouldSkipURL(u *url.URL) []string {
	if u == nil {
		return nil
	}

	reasons := make([]string, 0, 4)

	host := strings.ToLower(u.Hostname())
	if host != "" {
		for _, rule := range disallowedHostPrefixes {
			if strings.HasPrefix(host, rule.prefix) {
				reasons = append(reasons, fmt.Sprintf(rule.message, host))
			}
		}
		for _, rule := range disallowedHostFragments {
			if strings.Contains(host, rule.fragment) {
				reasons = append(reasons, fmt.Sprintf(rule.message, host))
			}
		}
	}

	path := strings.ToLower(u.Path)
	if path != "" {
		leaf := path
		if idx := strings.LastIndex(leaf, "/"); idx >= 0 {
			leaf = leaf[idx+1:]
		}
		leaf = strings.TrimSpace(leaf)
		if leaf != "" {
			for _, item := range disallowedURLExtensions {
				if strings.HasSuffix(leaf, item.suffix) {
					reasons = append(reasons, fmt.Sprintf("resource with extension %s (%s) is not processed", item.suffix, item.category))
				}
			}
		}

		for _, rule := range disallowedPathFragments {
			if strings.Contains(path, rule.fragment) {
				reasons = append(reasons, fmt.Sprintf(rule.message, rule.fragment))
			}
		}
	}

	rawQuery := strings.ToLower(u.RawQuery)
	if rawQuery != "" {
		for _, rule := range disallowedQueryFragments {
			if strings.Contains(rawQuery, rule.fragment) {
				reasons = append(reasons, fmt.Sprintf(rule.message, rule.fragment))
			}
		}
	}

	if len(reasons) == 0 {
		return nil
	}

	return reasons
}

func detectInvalidContent(rawHTML, text string) []string {
	seen := make(map[string]struct{})
	reasons := make([]string, 0, 2)

	for _, candidate := range []string{rawHTML, text} {
		if msg := detectBinaryPayload(candidate); msg != "" {
			reason := fmt.Sprintf("binary payload detected: %s", msg)
			if _, ok := seen[reason]; !ok {
				seen[reason] = struct{}{}
				reasons = append(reasons, reason)
			}
		}
	}

	if msg := detectScriptPayload(text); msg != "" {
		reason := fmt.Sprintf("javascript payload detected: %s", msg)
		if _, ok := seen[reason]; !ok {
			seen[reason] = struct{}{}
			reasons = append(reasons, reason)
		}
	}

	if len(reasons) == 0 {
		return nil
	}

	sort.Strings(reasons)
	return reasons
}

func detectBinaryPayload(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return ""
	}

	sample := trimmed
	if len(sample) > binaryProbeLimit {
		sample = sample[:binaryProbeLimit]
	}

	sampleBytes := []byte(sample)

	for _, sig := range binarySignatures {
		if strings.HasPrefix(sample, sig.prefix) {
			return sig.message
		}
	}

	nonPrintable := 0
	total := len(sampleBytes)
	for i := 0; i < len(sampleBytes); i++ {
		b := sampleBytes[i]
		if b == 0x00 {
			return "content contains null bytes"
		}
		if b < 32 && b != '\n' && b != '\r' && b != '\t' {
			nonPrintable++
		}
	}

	if total > 0 {
		ratio := float64(nonPrintable) / float64(total)
		if ratio >= 0.03 {
			return "content contains many non-printable bytes"
		}

		printableRatio := 1 - ratio
		entropy := shannonEntropy(sampleBytes)
		if entropy >= binaryEntropyThreshold && printableRatio <= printableRatioThreshold {
			return fmt.Sprintf("content entropy %.2f bits/byte indicates binary payload", entropy)
		}
	}

	return ""
}

func detectScriptPayload(text string) string {
	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return ""
	}

	runes := []rune(trimmed)
	if len(runes) > scriptProbeLimit {
		runes = runes[:scriptProbeLimit]
	}
	runeCount := len(runes)
	if runeCount == 0 {
		return ""
	}

	sample := strings.ToLower(string(runes))

	keywordHits := 0
	for _, kw := range scriptKeywords {
		if strings.Contains(sample, kw) {
			keywordHits++
		}
	}

	if keywordHits == 0 {
		return ""
	}

	operatorCount := 0
	for _, r := range runes {
		if _, ok := scriptOperatorRunes[r]; ok {
			operatorCount++
		}
	}

	operatorRatio := float64(operatorCount) / float64(runeCount)
	if keywordHits >= 3 && operatorRatio >= 0.18 {
		return fmt.Sprintf("%d script keywords and operator ratio %.2f", keywordHits, operatorRatio)
	}

	if keywordHits >= 6 {
		return fmt.Sprintf("%d script keywords detected", keywordHits)
	}

	return ""
}

func stripJSONFragments(text string) (string, int) {
	if text == "" {
		return text, 0
	}

	runes := []rune(text)
	var builder strings.Builder
	builder.Grow(len(text))

	removed := 0
	length := len(runes)

	for i := 0; i < length; {
		r := runes[i]
		if r == '{' || r == '[' {
			start := i
			depth := 1
			i++
			for i < length && depth > 0 {
				switch runes[i] {
				case '{', '[':
					depth++
				case '}':
					if runes[start] == '{' {
						depth--
					}
				case ']':
					if runes[start] == '[' {
						depth--
					}
				}
				i++
			}
			fragment := string(runes[start:i])
			if looksLikeJSONFragment(fragment) {
				removed++
				continue
			}
			builder.WriteString(fragment)
			continue
		}

		builder.WriteRune(r)
		i++
	}

	if removed == 0 {
		return text, 0
	}

	return strings.TrimSpace(builder.String()), removed
}

func looksLikeJSONFragment(fragment string) bool {
	trimmed := strings.TrimSpace(fragment)
	if len(trimmed) < 120 {
		return false
	}
	if trimmed[0] != '{' && trimmed[0] != '[' {
		return false
	}

	colons := strings.Count(trimmed, ":")
	quotes := strings.Count(trimmed, "\"")
	if colons < 3 || quotes < 6 {
		return false
	}

	return true
}

func (s *ScrapeService) retryWithUserAgent(ctx context.Context, strategy ContentStrategy, url string, currentUA string) (models.ContentData, bool) {
	custom, ok := strategy.(CustomAgentStrategy)
	if !ok || len(s.userAgents) == 0 || s.maxUserAgentRetries <= 0 {
		return models.ContentData{}, false
	}

	exclude := make(map[string]struct{})
	if trimmed := strings.TrimSpace(currentUA); trimmed != "" {
		exclude[trimmed] = struct{}{}
	}

	candidates := s.randomUserAgents(exclude)
	for _, ua := range candidates {
		payload, err := custom.FetchWithUserAgent(ctx, url, ua)
		if err != nil {
			continue
		}
		return payload, true
	}

	return models.ContentData{}, false
}

func (s *ScrapeService) randomUserAgents(exclude map[string]struct{}) []string {
	available := make([]string, 0, len(s.userAgents))
	for _, ua := range s.userAgents {
		trimmed := strings.TrimSpace(ua)
		if trimmed == "" {
			continue
		}
		if exclude != nil {
			if _, ok := exclude[trimmed]; ok {
				continue
			}
		}
		available = append(available, trimmed)
	}

	if len(available) == 0 {
		return nil
	}

	s.randMu.Lock()
	s.uaRand.Shuffle(len(available), func(i, j int) {
		available[i], available[j] = available[j], available[i]
	})
	s.randMu.Unlock()

	if len(available) > s.maxUserAgentRetries {
		available = available[:s.maxUserAgentRetries]
	}

	return available
}

func parseStatusCodeFromError(err error) (int, bool) {
	if err == nil {
		return 0, false
	}
	const marker = "status code: "
	msg := err.Error()
	idx := strings.LastIndex(msg, marker)
	if idx == -1 {
		return 0, false
	}
	idx += len(marker)
	if idx >= len(msg) {
		return 0, false
	}

	end := idx
	for end < len(msg) && msg[end] >= '0' && msg[end] <= '9' {
		end++
	}
	if end == idx {
		return 0, false
	}

	code, convErr := strconv.Atoi(msg[idx:end])
	if convErr != nil {
		return 0, false
	}
	return code, true
}

func shannonEntropy(data []byte) float64 {
	if len(data) == 0 {
		return 0
	}

	counts := make(map[byte]int, 32)
	for _, b := range data {
		counts[b]++
	}

	var entropy float64
	total := float64(len(data))
	for _, c := range counts {
		p := float64(c) / total
		entropy -= p * math.Log2(p)
	}
	return entropy
}
