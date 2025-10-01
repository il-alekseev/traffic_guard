package bootstrap

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	stdruntime "runtime"
	"strings"
	"time"

	"scrapper/config"
	"scrapper/internal/app"
	"scrapper/internal/infrastructure/content"
	dnsresolver "scrapper/internal/infrastructure/dns"
	"scrapper/internal/infrastructure/geo"
	"scrapper/internal/infrastructure/httpclient"
	"scrapper/internal/infrastructure/logging"
	"scrapper/internal/infrastructure/parser"
	"scrapper/internal/infrastructure/quality"
	"scrapper/internal/infrastructure/storage"
	"scrapper/internal/usecase"
)

const (
	defaultDBPath        = "data/ipinfo_lite.mmdb"
	defaultWorkersScale  = 2
	defaultRequestTO     = 15 * time.Second
	defaultHTTPTO        = 10 * time.Second
	defaultStrategyName  = "http"
	defaultRepository    = "data/attempts.jsonl"
	defaultMinTextLength = 200
	defaultLogPath       = "logs/scraper.log"
)

var ErrNoURLs = errors.New("no urls provided")

// Params encapsulates runtime dependencies needed to launch the scraper.
type Params struct {
	Config *config.Config
	URLs   []string
	Stdout io.Writer
	Stdin  io.Reader
}

// Run prepares infrastructure and executes the scraping workflow. It returns an
// error if configuration is invalid or the pipeline fails to start.
func Run(ctx context.Context, p Params) error {
	if p.Config == nil {
		return errors.New("config is nil")
	}

	stdout := p.Stdout
	if stdout == nil {
		stdout = os.Stdout
	}

	stdin := p.Stdin
	if stdin == nil {
		stdin = os.Stdin
	}

	cfg := applyDefaults(*p.Config)

	urls, err := prepareURLs(stdin, cfg.InputPath, p.URLs)
	if err != nil {
		return fmt.Errorf("prepare urls: %w", err)
	}
	if len(urls) == 0 {
		return ErrNoURLs
	}

	resolver := dnsresolver.NewResolver()

	geoProvider, err := geo.NewMaxMindProvider(cfg.DBPath)
	if err != nil {
		return fmt.Errorf("geo provider: %w", err)
	}
	defer geoProvider.Close()

	httpClient := httpclient.NewClient(cfg.HTTPTimeout.Duration)
	cleaner := parser.NewHTMLCleaner()

	logger, err := logging.NewFileLogger(cfg.Logging.File)
	if err != nil {
		return fmt.Errorf("logger: %w", err)
	}
	defer logger.Close()

	strategies := buildStrategies(cfg.Content.Strategies, httpClient)
	evaluator := quality.NewSimpleEvaluator(cfg.Content.MinTextLength)

	repository, err := storage.NewFileRepository(cfg.Content.RepositoryPath)
	if err != nil {
		return fmt.Errorf("repository: %w", err)
	}

	service := usecase.NewScrapeService(strategies, resolver, geoProvider, cleaner, evaluator, repository, logger, cfg.Workers, cfg.RequestTimeout.Duration)
	scraperApp := app.NewScraper(service)

	results := scraperApp.Run(ctx, urls)

	encoder := json.NewEncoder(stdout)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(results); err != nil {
		return fmt.Errorf("encode results: %w", err)
	}

	return nil
}

func applyDefaults(cfg config.Config) config.Config {
	if cfg.DBPath == "" {
		cfg.DBPath = defaultDBPath
	}

	if cfg.Workers <= 0 {
		cfg.Workers = stdruntime.NumCPU() * defaultWorkersScale
		if cfg.Workers < 1 {
			cfg.Workers = 1
		}
	}

	if cfg.RequestTimeout.Duration <= 0 {
		cfg.RequestTimeout.Duration = defaultRequestTO
	}

	if cfg.HTTPTimeout.Duration <= 0 {
		cfg.HTTPTimeout.Duration = defaultHTTPTO
	}

	if len(cfg.Content.Strategies) == 0 {
		cfg.Content.Strategies = []string{defaultStrategyName}
	}

	if cfg.Content.RepositoryPath == "" {
		cfg.Content.RepositoryPath = defaultRepository
	}

	if cfg.Content.MinTextLength <= 0 {
		cfg.Content.MinTextLength = defaultMinTextLength
	}

	if cfg.Logging.File == "" {
		cfg.Logging.File = defaultLogPath
	}

	return cfg
}

func prepareURLs(stdin io.Reader, inputPath string, cliURLs []string) ([]string, error) {
	urls := make([]string, 0)

	if inputPath != "" {
		loaded, err := loadURLs(stdin, inputPath)
		if err != nil {
			return nil, err
		}
		urls = append(urls, loaded...)
	}

	if len(cliURLs) > 0 {
		urls = append(urls, cliURLs...)
	}

	return filterURLs(urls), nil
}

func loadURLs(stdin io.Reader, path string) ([]string, error) {
	var data []byte
	var err error

	if path == "-" {
		data, err = io.ReadAll(stdin)
	} else {
		data, err = os.ReadFile(path)
	}

	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	return filterURLs(lines), nil
}

func filterURLs(values []string) []string {
	urls := make([]string, 0, len(values))
	seen := make(map[string]struct{})

	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}
		if _, exists := seen[trimmed]; exists {
			continue
		}
		seen[trimmed] = struct{}{}
		urls = append(urls, trimmed)
	}

	return urls
}

func buildStrategies(names []string, client *httpclient.Client) []usecase.ContentStrategy {
	strategies := make([]usecase.ContentStrategy, 0, len(names))

	for _, name := range names {
		switch strings.ToLower(strings.TrimSpace(name)) {
		case "trafilatura":
			if s, err := content.NewTrafilaturaStrategy(client.Client()); err == nil {
				strategies = append(strategies, s)
				break
			}
			// If not available (no build tag), fallthrough to next handlers.
		case "", defaultStrategyName:
			strategies = append(strategies, content.NewHTTPStrategy(client))
		}
	}

	if len(strategies) == 0 {
		strategies = append(strategies, content.NewHTTPStrategy(client))
	}

	return strategies
}
