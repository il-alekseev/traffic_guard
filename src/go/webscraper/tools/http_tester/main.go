package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

type requestPayload struct {
	RequestID string             `json:"request_id"`
	Src       requestSource      `json:"src"`
	Dst       requestDestination `json:"dst"`
}

type requestSource struct {
	IP string `json:"ip"`
}

type requestDestination struct {
	Type     string `json:"type"`
	Resource string `json:"resource"`
}

type responseBody struct {
	Metadata struct {
		RequestID string   `json:"request_id"`
		Status    string   `json:"status"`
		Failure   string   `json:"failure"`
		Error     string   `json:"error"`
		Warnings  []string `json:"warnings"`
		Strategy  string   `json:"strategy"`
		ContentID string   `json:"content_id"`
	} `json:"metadata"`
	Content json.RawMessage `json:"content"`
}

func main() {
	var (
		endpoint   string
		inputPath  string
		srcIP      string
		dstType    string
		reqID      string
		delay      time.Duration
		maxRetries int
		retryDelay time.Duration
	)
	flag.StringVar(&endpoint, "endpoint", "http://localhost:8010/get_content", "HTTP endpoint to call")
	flag.StringVar(&inputPath, "input", "", "Path to file with URLs (one per line)")
	flag.StringVar(&srcIP, "src-ip", "127.0.0.1", "Value for src.ip")
	flag.StringVar(&dstType, "dst-type", "url", "Destination type (url|domain|ip)")
	flag.StringVar(&reqID, "request-id", "", "Optional fixed request_id; if empty a uuid is generated for each request")
	flag.DurationVar(&delay, "delay", 0, "Delay between requests (e.g., 500ms, 2s)")
	flag.IntVar(&maxRetries, "retries", 3, "Maximum HTTP retries on transport errors")
	flag.DurationVar(&retryDelay, "retry-delay", 250*time.Millisecond, "Delay between retries on transport errors")
	flag.Parse()

	if strings.TrimSpace(inputPath) == "" {
		flag.Usage()
		log.Fatal("input file is required")
	}

	file, err := os.Open(inputPath)
	if err != nil {
		log.Fatalf("open input file: %v", err)
	}
	defer file.Close()

	client := &http.Client{Timeout: 120 * time.Second}
	scanner := bufio.NewScanner(file)
	line := 0

	for scanner.Scan() {
		line++
		raw := strings.TrimSpace(scanner.Text())
		if raw == "" || strings.HasPrefix(raw, "#") {
			continue
		}

		currentID := strings.TrimSpace(reqID)
		if currentID == "" {
			currentID = uuid.NewString()
		}

		payload := requestPayload{
			RequestID: currentID,
			Src: requestSource{
				IP: srcIP,
			},
			Dst: requestDestination{
				Type:     dstType,
				Resource: raw,
			},
		}
		bodyBytes, err := json.Marshal(payload)
		if err != nil {
			log.Printf("[line %d] marshal error: %v", line, err)
			continue
		}

		var (
			resp    *http.Response
			rawBody []byte
		)
		success := false
		for attempt := 0; attempt < maxRetries; attempt++ {
			ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
			req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(bodyBytes))
			if err != nil {
				cancel()
				log.Printf("[line %d] build request: %v", line, err)
				break
			}
			req.Header.Set("Content-Type", "application/json")

			resp, err = client.Do(req)
			cancel()
			if err != nil {
				if attempt == maxRetries-1 {
					log.Printf("[line %d] request error: %v", line, err)
				} else {
					log.Printf("[line %d] request error (retry %d/%d): %v", line, attempt+1, maxRetries, err)
					time.Sleep(retryDelay)
				}
				continue
			}

			rawBody, err = io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				log.Printf("[line %d] read response: %v", line, err)
				continue
			}
			success = true
			break
		}

		if !success {
			continue
		}

		var respData responseBody
		if err := json.Unmarshal(rawBody, &respData); err != nil {
			log.Printf("[line %d] decode response: %v", line, err)
			continue
		}

		meta := respData.Metadata
		fmt.Printf("line=%d url=%s status=%s failure=%s error=%s strategy=%s request_id=%s\n", line, raw, meta.Status, meta.Failure, meta.Error, meta.Strategy, meta.RequestID)
		if len(meta.Warnings) > 0 {
			fmt.Printf("  warnings: %s\n", strings.Join(meta.Warnings, " | "))
		}
		if meta.ContentID != "" {
			fmt.Printf("  content_id: %s\n", meta.ContentID)
		}
		fmt.Printf("  response: %s\n", strings.TrimSpace(string(rawBody)))

		if delay > 0 {
			time.Sleep(delay)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("read input: %v", err)
	}
}
