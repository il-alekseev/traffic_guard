package bootstrap

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"

	"scrapper/config"
	"scrapper/internal/app"
	"scrapper/internal/infrastructure/content"
	dnsresolver "scrapper/internal/infrastructure/dns"
	"scrapper/internal/infrastructure/geo"
	"scrapper/internal/infrastructure/httpclient"
	kafkaio "scrapper/internal/infrastructure/kafka"
	"scrapper/internal/infrastructure/logging"
	"scrapper/internal/infrastructure/parser"
	"scrapper/internal/infrastructure/quality"
	"scrapper/internal/infrastructure/storage"
	"scrapper/internal/metrics"
	"scrapper/internal/models"
	"scrapper/internal/processor"
	httpserver "scrapper/internal/server/http"
	"scrapper/internal/usecase"
)

const (
	defaultDBPath             = "data/ipinfo_lite.mmdb"
	defaultWorkersScale       = 2
	defaultRequestTO          = 15 * time.Second
	defaultStrategyName       = "http"
	defaultMinTextLength      = 200
	defaultMaxContentLength   = 50000
	defaultLogLevel           = "info"
	defaultUserAgentsPath     = "data/user_agents.txt"
	defaultAntiBotPhrasesPath = "data/anti_bot_phrases.txt"
	defaultKafkaClientID      = "scraper-worker"
	defaultKafkaPollTimeout   = 2 * time.Second
	defaultKafkaCommitPeriod  = 5 * time.Second
)

var defaultArchiveStrategies = []string{"wayback", "archive.today"}

const kafkaMessageKey = "__kafka_message"

type Params struct {
	Config *config.Config
}

func Run(ctx context.Context, p Params) error {
	if p.Config == nil {
		return errors.New("config is nil")
	}

	cfg := applyDefaults(*p.Config)

	logger := logging.New(cfg.Logging.Level)
	if logger == nil {
		logger = logging.New(defaultLogLevel)
	}
	bootLog := logger.With("component", "bootstrap")
	bootLog.Info("scraper runtime starting",
		"workers", cfg.Workers,
		"request_timeout", cfg.RequestTimeout.Duration.String(),
		"kafka_topic", cfg.Kafka.InputTopic,
		"log_level", cfg.Logging.Level,
	)

	userAgents, err := loadUserAgents(cfg.Content.UserAgentsPath)
	if err != nil {
		return fmt.Errorf("load user agents: %w", err)
	}
	bootLog.Info("user agents loaded", "count", len(userAgents), "path", cfg.Content.UserAgentsPath)

	antiBotPhrases, err := loadAntiBotPhrases(cfg.Content.AntiBotPhrasesPath)
	if err != nil {
		bootLog.Warn("anti-bot phrases not loaded", "error", err.Error(), "path", cfg.Content.AntiBotPhrasesPath)
	} else {
		bootLog.Info("anti-bot phrases loaded", "count", len(antiBotPhrases), "path", cfg.Content.AntiBotPhrasesPath)
	}

	resolver := dnsresolver.NewResolver()

	geoProvider, err := geo.NewMaxMindProvider(cfg.DBPath)
	if err != nil {
		return fmt.Errorf("geo provider: %w", err)
	}
	defer geoProvider.Close()

	httpClient := httpclient.NewClient(cfg.RequestTimeout.Duration, userAgents)
	cleaner := parser.NewHTMLCleaner()

	primaryStrategies, archiveStrategies := buildStrategies(cfg.Content.Strategies, cfg.Content.ArchiveStrategies, httpClient, userAgents)
	evaluator := quality.NewSimpleEvaluator(cfg.Content.MinTextLength)

	var repository usecase.AttemptRepository
	if repoPath := strings.TrimSpace(cfg.Content.RepositoryPath); repoPath != "" {
		fileRepo, repoErr := storage.NewFileRepository(repoPath)
		if repoErr != nil {
			return fmt.Errorf("repository: %w", repoErr)
		}
		repository = fileRepo
		bootLog.Info("attempt repository enabled", "path", repoPath)
	} else {
		bootLog.Info("attempt repository disabled")
	}

	metricsCollector := metrics.New()

	antiBotDetector := usecase.NewAntiBotDetector(antiBotPhrases, 500)

	service := usecase.NewScrapeService(
		primaryStrategies,
		archiveStrategies,
		resolver,
		geoProvider,
		cleaner,
		evaluator,
		cfg.Content.MaxContentLength,
		repository,
		logger,
		cfg.Workers,
		cfg.RequestTimeout.Duration,
		metricsCollector,
		antiBotDetector,
		cfg.Content.StripJSONFragments,
		userAgents,
	)
	scraperApp := app.NewScraper(service)

	kafkaConsumer, err := kafkaio.NewConsumer(cfg.Kafka)
	if err != nil {
		return fmt.Errorf("kafka consumer: %w", err)
	}
	defer kafkaConsumer.Close()

	metadataProducer, err := kafkaio.NewProducer(cfg.Kafka.MetadataTopic, cfg.Kafka)
	if err != nil {
		return fmt.Errorf("kafka metadata producer: %w", err)
	}
	defer metadataProducer.Close()

	contentProducer, err := kafkaio.NewProducer(cfg.Kafka.ContentTopic, cfg.Kafka)
	if err != nil {
		return fmt.Errorf("kafka content producer: %w", err)
	}
	defer contentProducer.Close()

	requestBuffer := cfg.Workers * 4
	if requestBuffer < cfg.Workers {
		requestBuffer = cfg.Workers
	}

	requests := make(chan models.ProcessingRequest, requestBuffer)
	outcomes := scraperApp.Stream(ctx, cfg.Workers, requests)

	errCh := make(chan error, 1)

	consumerLogger := logger.With("component", "consumer")
	go consumeLoop(ctx, kafkaConsumer, requests, errCh, consumerLogger)
	pipelineLogger := logger.With("component", "pipeline")
	httpSrv := httpserver.New(cfg.HTTP, cfg.Kafka, service, metricsCollector, logger)
	httpErrCh := make(chan error, 1)
	go func() {
		httpErrCh <- httpSrv.Start(ctx)
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errCh:
			if err != nil {
				return err
			}
			return nil
		case err := <-httpErrCh:
			if err != nil {
				return err
			}
		case outcome, ok := <-outcomes:
			if !ok {
				bootLog.Info("pipeline stopped: outcomes closed")
				return nil
			}
			reqLogger := pipelineLogger.With("request_id", outcome.Request.RequestID)
			reqLogger.Debug("processing outcome", "status", outcome.Best.Result.Status)
			contentID, err := handleOutcome(ctx, outcome, kafkaConsumer, metadataProducer, contentProducer, reqLogger)
			if err != nil {
				reqLogger.Error("outcome processing failed", "error", err.Error())
				return err
			}
			reqLogger.Debug("outcome processed", "content_id", contentID, "status", outcome.Best.Result.Status)
		}
	}
}

func applyDefaults(cfg config.Config) config.Config {
	if cfg.DBPath == "" {
		cfg.DBPath = defaultDBPath
	}

	if cfg.Workers <= 0 {
		cfg.Workers = defaultWorkers()
	}

	if cfg.RequestTimeout.Duration <= 0 {
		cfg.RequestTimeout.Duration = defaultRequestTO
	}

	if len(cfg.Content.Strategies) == 0 {
		cfg.Content.Strategies = []string{defaultStrategyName}
	}

	if cfg.Content.MinTextLength <= 0 {
		cfg.Content.MinTextLength = defaultMinTextLength
	}

	if cfg.Content.MaxContentLength <= 0 {
		cfg.Content.MaxContentLength = defaultMaxContentLength
	}

	if strings.TrimSpace(cfg.Content.UserAgentsPath) == "" {
		cfg.Content.UserAgentsPath = defaultUserAgentsPath
	}

	if len(cfg.Content.ArchiveStrategies) == 0 {
		cfg.Content.ArchiveStrategies = append([]string(nil), defaultArchiveStrategies...)
	}

	if strings.TrimSpace(cfg.Content.AntiBotPhrasesPath) == "" {
		cfg.Content.AntiBotPhrasesPath = defaultAntiBotPhrasesPath
	}

	if cfg.Logging.Level == "" {
		cfg.Logging.Level = defaultLogLevel
	}

	if cfg.Kafka.ClientID == "" {
		cfg.Kafka.ClientID = defaultKafkaClientID
	}

	if cfg.Kafka.PollTimeout.Duration <= 0 {
		cfg.Kafka.PollTimeout.Duration = defaultKafkaPollTimeout
	}

	if cfg.Kafka.CommitInterval.Duration <= 0 {
		cfg.Kafka.CommitInterval.Duration = defaultKafkaCommitPeriod
	}

	if cfg.HTTP.Address == "" {
		cfg.HTTP.Address = ":8010"
	}

	if cfg.HTTP.ReadTimeout.Duration <= 0 {
		cfg.HTTP.ReadTimeout.Duration = 5 * time.Second
	}

	if cfg.HTTP.WriteTimeout.Duration <= 0 {
		cfg.HTTP.WriteTimeout.Duration = 10 * time.Second
	}

	if cfg.HTTP.ShutdownTimeout.Duration <= 0 {
		cfg.HTTP.ShutdownTimeout.Duration = 10 * time.Second
	}

	return cfg
}

func defaultWorkers() int {
	workers := defaultWorkersScale * runtimeCPUCount()
	if workers < 1 {
		return 1
	}
	return workers
}

func runtimeCPUCount() int {
	return runtime.NumCPU()
}

func consumeLoop(
	ctx context.Context,
	consumer *kafkaio.Consumer,
	requests chan<- models.ProcessingRequest,
	errCh chan<- error,
	logger usecase.Logger,
) {
	defer close(requests)

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		msg, err := consumer.Fetch(ctx)
		if err != nil {
			if errors.Is(err, kafkaio.ErrEmptyPoll) {
				continue
			}
			if logger != nil {
				logger.Error("kafka fetch failed", "error", err.Error())
			}
			select {
			case errCh <- fmt.Errorf("fetch kafka message: %w", err):
			default:
			}
			return
		}

		req, err := processor.DecodeKafkaMessage(msg.Value)
		if err != nil {
			if logger != nil {
				logger.Warn("discarding invalid message", "error", err.Error())
			}
			if commitErr := consumer.Commit(ctx, msg); commitErr != nil && logger != nil {
				logger.Error("failed to commit invalid message", "error", commitErr.Error())
			}
			continue
		}

		if !msg.Time.IsZero() {
			req.Received = msg.Time
		}
		if req.Raw == nil {
			req.Raw = make(map[string]any, 1)
		}
		req.Raw[kafkaMessageKey] = msg
		if logger != nil {
			logger.Debug("message enqueued", "request_id", req.RequestID)
		}

		select {
		case <-ctx.Done():
			return
		case requests <- req:
		}
	}
}

func handleOutcome(
	ctx context.Context,
	outcome models.RequestOutcome,
	consumer *kafkaio.Consumer,
	metadataProducer *kafkaio.Producer,
	contentProducer *kafkaio.Producer,
	logger usecase.Logger,
) (string, error) {
	contentID := processor.EnsureContentID(outcome)
	best := outcome.Best.Result

	metadata := processor.BuildMetadata(outcome, contentID)
	if err := publishMetadata(ctx, metadataProducer, metadata); err != nil {
		return "", fmt.Errorf("publish metadata: %w", err)
	}

	if contentMsg, ok := processor.BuildContent(outcome, contentID); ok {
		if err := publishContent(ctx, contentProducer, contentMsg); err != nil {
			return "", fmt.Errorf("publish content: %w", err)
		}
	}

	msg, err := extractKafkaMessage(outcome.Request.Raw)
	if err != nil {
		return "", fmt.Errorf("extract kafka message: %w", err)
	}

	if err := consumer.Commit(ctx, msg); err != nil {
		return "", fmt.Errorf("commit kafka message: %w", err)
	}

	if logger != nil {
		logger.Info(
			"outcome committed",
			"request_id", outcome.Request.RequestID,
			"content_id", contentID,
			"status", best.Status,
			"strategy", best.Strategy,
		)
	}

	return contentID, nil
}

func publishMetadata(ctx context.Context, producer *kafkaio.Producer, metadata models.MetadataMessage) error {
	if producer == nil {
		return errors.New("metadata producer not initialized")
	}

	requestID := metadata.RequestID
	if requestID == "" {
		return errors.New("metadata missing request_id")
	}

	payload, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("marshal metadata: %w", err)
	}

	return producer.Publish(ctx, []byte(requestID), payload)
}

func publishContent(ctx context.Context, producer *kafkaio.Producer, message models.ContentMessage) error {
	if producer == nil {
		return nil
	}

	requestID := message.RequestID
	if requestID == "" || strings.TrimSpace(message.Content) == "" {
		return nil
	}

	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("marshal content: %w", err)
	}

	return producer.Publish(ctx, []byte(requestID), payload)
}

func extractKafkaMessage(raw map[string]any) (kafka.Message, error) {
	if raw == nil {
		return kafka.Message{}, errors.New("request metadata missing")
	}

	value, ok := raw[kafkaMessageKey]
	if !ok {
		return kafka.Message{}, errors.New("kafka message metadata missing")
	}

	switch v := value.(type) {
	case kafka.Message:
		return v, nil
	case *kafka.Message:
		if v == nil {
			return kafka.Message{}, errors.New("kafka message pointer nil")
		}
		return *v, nil
	default:
		return kafka.Message{}, errors.New("unexpected kafka message metadata type")
	}
}

func loadUserAgents(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 4096), 1<<20)

	seen := make(map[string]struct{})
	agents := make([]string, 0, 8)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if _, exists := seen[line]; exists {
			continue
		}
		seen[line] = struct{}{}
		agents = append(agents, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(agents) == 0 {
		return nil, fmt.Errorf("no user agents found in %s", path)
	}

	return agents, nil
}

func loadAntiBotPhrases(path string) ([]string, error) {
	if strings.TrimSpace(path) == "" {
		return nil, errors.New("anti-bot phrases path empty")
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	phrases := make([]string, 0, 16)
	seen := make(map[string]struct{})
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		lowered := strings.ToLower(line)
		if _, ok := seen[lowered]; ok {
			continue
		}
		seen[lowered] = struct{}{}
		phrases = append(phrases, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return phrases, nil
}

func buildStrategies(primaryNames, archiveNames []string, client *httpclient.Client, userAgents []string) ([]usecase.ContentStrategy, []usecase.ContentStrategy) {
	primary := buildStrategySet(primaryNames, client, userAgents, false)
	archive := buildStrategySet(archiveNames, client, userAgents, true)

	if len(primary) == 0 {
		primary = append(primary, content.NewHTTPStrategy(client))
	}

	return primary, archive
}

func buildStrategySet(names []string, client *httpclient.Client, userAgents []string, archive bool) []usecase.ContentStrategy {
	if len(names) == 0 {
		return nil
	}

	strategies := make([]usecase.ContentStrategy, 0, len(names))

	for _, name := range names {
		switch strings.ToLower(strings.TrimSpace(name)) {
		case "trafilatura":
			if archive {
				continue
			}
			if s, err := content.NewTrafilaturaStrategy(client.Client(), userAgents); err == nil {
				strategies = append(strategies, s)
			}
		case "", defaultStrategyName:
			if archive {
				continue
			}
			strategies = append(strategies, content.NewHTTPStrategy(client))
		case "wayback", "webarchive", "archive.org":
			if s := content.NewWaybackStrategy(client, userAgents); s != nil {
				strategies = append(strategies, s)
			}
		case "archive.today", "archivetoday", "archive_today":
			if s := content.NewArchiveTodayStrategy(client, userAgents); s != nil {
				strategies = append(strategies, s)
			}
		}
	}

	return strategies
}
