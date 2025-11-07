package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"

	"scrapper/config"
	kafkaio "scrapper/internal/infrastructure/kafka"
	"scrapper/internal/models"
)

type options struct {
	envPath    string
	configPath string
	inputPath  string
	srcIP      string
	dstType    string
	dryRun     bool
}

func main() {
	opts := parseFlags()

	if err := loadDotEnv(opts.envPath); err != nil {
		log.Fatalf("load env file: %v", err)
	}

	cfg, err := config.Load(opts.configPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}
	config.ApplyEnvOverrides(cfg)

	if cfg.Kafka.InputTopic == "" {
		log.Fatalf("kafka input topic is not configured")
	}

	dstType := parseDestinationType(opts.dstType)

	ctx := context.Background()
	var producer *kafkaio.Producer
	if !opts.dryRun {
		producer, err = kafkaio.NewProducer(cfg.Kafka.InputTopic, cfg.Kafka)
		if err != nil {
			log.Fatalf("kafka producer: %v", err)
		}
		defer producer.Close()
	}

	file, err := os.Open(opts.inputPath)
	if err != nil {
		log.Fatalf("open input file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNo := 0
	sent := 0
	skipped := 0

	for scanner.Scan() {
		lineNo++
		raw := strings.TrimSpace(scanner.Text())
		if raw == "" || strings.HasPrefix(raw, "#") {
			continue
		}

		req := buildRequest(raw, opts.srcIP, dstType)
		payload, err := json.Marshal(req)
		if err != nil {
			log.Printf("line %d: marshal error: %v", lineNo, err)
			skipped++
			continue
		}

		if opts.dryRun {
			fmt.Printf("[dry-run] would publish request_id=%s url=%s\n", req.RequestID, req.Dst.Resource)
			sent++
			continue
		}

		if err := producer.Publish(ctx, []byte(req.RequestID), payload); err != nil {
			log.Printf("line %d: publish error: %v", lineNo, err)
			skipped++
			continue
		}
		sent++
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("read input: %v", err)
	}

	log.Printf("completed: sent=%d skipped=%d", sent, skipped)
}

func parseFlags() options {
	opt := options{}
	flag.StringVar(&opt.envPath, "env", ".env", "Path to .env file with Kafka settings")
	flag.StringVar(&opt.configPath, "config", "config/config.yaml", "Path to YAML config")
	flag.StringVar(&opt.inputPath, "input", "", "Path to file with URLs (one per line)")
	flag.StringVar(&opt.srcIP, "src-ip", "127.0.0.1", "Value for request.src.ip")
	flag.StringVar(&opt.dstType, "dst-type", "url", "Destination type (url|domain|ip)")
	flag.BoolVar(&opt.dryRun, "dry-run", false, "Do not publish to Kafka, only print")
	flag.Parse()

	if strings.TrimSpace(opt.inputPath) == "" {
		flag.Usage()
		log.Fatal("input file is required")
	}

	return opt
}

func loadDotEnv(path string) error {
	if strings.TrimSpace(path) == "" {
		return nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if key == "" {
			continue
		}
		_ = os.Setenv(key, value)
	}
	return scanner.Err()
}

func parseDestinationType(value string) models.DestinationType {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "ip":
		return models.DestinationIP
	case "domain":
		return models.DestinationDomain
	default:
		return models.DestinationURL
	}
}

func buildRequest(url string, srcIP string, dstType models.DestinationType) models.ProcessingRequest {
	return models.ProcessingRequest{
		RequestID: uuid.NewString(),
		Src: models.RequestSource{
			IP: strings.TrimSpace(srcIP),
		},
		Dst: models.RequestDestination{
			Type:     dstType,
			Resource: strings.TrimSpace(url),
		},
		Received: time.Now().UTC(),
		Raw:      map[string]any{"source": "url_publisher_tool"},
	}
}
