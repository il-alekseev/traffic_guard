package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"scrapper/internal/infrastructure/content"
	"scrapper/internal/infrastructure/httpclient"
	"scrapper/internal/models"
)

func main() {
	var (
		rawURL        string
		userAgents    string
		timeout       time.Duration
		showSnippet   bool
		strategyNames string
	)

	flag.StringVar(&rawURL, "url", "", "URL to fetch via archive strategies (required)")
	flag.StringVar(&userAgents, "user-agents", "data/user_agents.txt", "Path to user agents list")
	flag.DurationVar(&timeout, "timeout", 30*time.Second, "Per-strategy timeout")
	flag.BoolVar(&showSnippet, "snippet", true, "Print first 400 characters of fetched payload")
	flag.StringVar(&strategyNames, "strategies", "wayback,archive.today", "Comma-separated list of archive strategies to run")
	flag.Parse()

	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		flag.Usage()
		log.Fatal("missing --url")
	}

	agents, err := loadUserAgents(userAgents)
	if err != nil {
		log.Fatalf("load user agents: %v", err)
	}

	httpClient := httpclient.NewClient(timeout, agents)
	strategies := buildStrategies(strings.Split(strategyNames, ","), httpClient, agents)
	if len(strategies) == 0 {
		log.Fatal("no archive strategies enabled")
	}

	fmt.Printf("Checking archives for %s (timeout per strategy %s)\n", rawURL, timeout)
	for _, strategy := range strategies {
		if strategy == nil {
			continue
		}
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		data, err := strategy.Fetch(ctx, rawURL)
		cancel()

		fmt.Printf("\n[%s]\n", strings.ToUpper(strategy.Name()))
		if err != nil {
			fmt.Printf("  error: %v\n", err)
			continue
		}

		length := len(data.RawHTML)
		if length == 0 && data.Text != "" {
			length = len(data.Text)
		}
		fmt.Printf("  success, bytes=%d user_agent=%s\n", length, data.UserAgent)
		if showSnippet {
			preview := snippetText(data.RawHTML)
			if preview == "" {
				preview = snippetText(data.Text)
			}
			if preview != "" {
				fmt.Printf("  snippet:\n%s\n", preview)
			}
		}
	}
}

type archiveStrategy interface {
	Name() string
	Fetch(ctx context.Context, url string) (models.ContentData, error)
}

func buildStrategies(names []string, client *httpclient.Client, userAgents []string) []archiveStrategy {
	result := make([]archiveStrategy, 0, len(names))
	for _, raw := range names {
		name := strings.ToLower(strings.TrimSpace(raw))
		switch name {
		case "wayback", "webarchive", "archive.org":
			if s := content.NewWaybackStrategy(client, userAgents); s != nil {
				result = append(result, s)
			}
		case "archive.today", "archive_today", "archivetoday":
			if s := content.NewArchiveTodayStrategy(client, userAgents); s != nil {
				result = append(result, s)
			}
		default:
			// ignore unknown names to keep flag flexible
		}
	}
	return result
}

func loadUserAgents(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	agents := make([]string, 0, 32)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		trimmed := strings.TrimSpace(scanner.Text())
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		agents = append(agents, trimmed)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if len(agents) == 0 {
		agents = []string{httpclient.DefaultUserAgent}
	}
	return agents, nil
}

func snippetText(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	if len(raw) > 400 {
		return raw[:400] + "..."
	}
	return raw
}
