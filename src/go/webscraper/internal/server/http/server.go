package httpserver

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"scrapper/config"
	"scrapper/internal/infrastructure/kafka"
	"scrapper/internal/infrastructure/logging"
	"scrapper/internal/metrics"
	"scrapper/internal/models"
	"scrapper/internal/processor"
	"scrapper/internal/usecase"
)

type Server struct {
	config    config.HTTPConfig
	kafka     config.KafkaConfig
	service   *usecase.ScrapeService
	metrics   *metrics.Metrics
	logger    *logging.Logger
	startedAt time.Time
}

func New(cfg config.HTTPConfig, kafkaCfg config.KafkaConfig, service *usecase.ScrapeService, metrics *metrics.Metrics, logger *logging.Logger) *Server {
	return &Server{
		config:    cfg,
		kafka:     kafkaCfg,
		service:   service,
		metrics:   metrics,
		logger:    logger,
		startedAt: time.Now(),
	}
}

func (s *Server) Start(ctx context.Context) error {
	addr := s.config.Address
	if addr == "" {
		addr = ":8080"
	}

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(30 * time.Second))

	router.Get("/helth_check", s.handleHealth)
	router.Post("/get_content", s.handleContent)
	router.Get("/swagger.yaml", s.handleSwagger)
	router.Get("/swagger", s.handleSwaggerUI)
	router.Get("/swagger/", s.handleSwaggerUI)

	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  s.config.ReadTimeout.Duration,
		WriteTimeout: s.config.WriteTimeout.Duration,
	}

	if s.config.ReadTimeout.Duration <= 0 {
		server.ReadTimeout = 5 * time.Second
	}
	if s.config.WriteTimeout.Duration <= 0 {
		server.WriteTimeout = 10 * time.Second
	}

	errCh := make(chan error, 1)
	go func() {
		err := server.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		} else {
			errCh <- nil
		}
	}()

	if s.logger != nil {
		s.logger.Info("http server started", "addr", addr)
	}

	select {
	case <-ctx.Done():
		shutdownTimeout := s.config.ShutdownTimeout.Duration
		if shutdownTimeout <= 0 {
			shutdownTimeout = 10 * time.Second
		}
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		_ = server.Shutdown(shutdownCtx)
		return <-errCh
	case err := <-errCh:
		return err
	}
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	kafkaStatus := "ok"
	if err := checkKafka(ctx, s.kafka); err != nil {
		kafkaStatus = err.Error()
	}

	snapshot := s.metrics.Snapshot()
	overall := "ok"
	if kafkaStatus != "ok" {
		overall = "degraded"
	}
	workers := 0
	if s.service != nil {
		workers = s.service.Workers()
	}
	resp := map[string]any{
		"status": overall,
		"kafka": map[string]any{
			"status": kafkaStatus,
		},
		"workers": map[string]any{
			"configured": workers,
			"in_flight":  snapshot.InFlight,
			"processed":  snapshot.Processed,
			"failed":     snapshot.Failed,
		},
		"uptime":     time.Since(s.startedAt).String(),
		"started_at": snapshot.StartedAt,
		"processing": snapshot.Processing,
	}

	s.writeJSON(w, resp, http.StatusOK)
}

func (s *Server) handleContent(w http.ResponseWriter, r *http.Request) {
	if s.service == nil {
		s.writeError(w, http.StatusServiceUnavailable, errors.New("service unavailable"))
		return
	}

	var req models.ProcessingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeError(w, http.StatusBadRequest, err)
		return
	}

	normalized, err := processor.NormalizeRequest(req)
	if err != nil {
		s.writeError(w, http.StatusBadRequest, err)
		return
	}

	outcome := s.service.ProcessRequest(r.Context(), -1, normalized)
	contentID := processor.EnsureContentID(outcome)
	metadata := processor.BuildMetadata(outcome, contentID)
	contentMsg, hasContent := processor.BuildContent(outcome, contentID)

	resp := map[string]any{
		"metadata": metadata,
	}
	if hasContent {
		resp["content"] = contentMsg
	}

	s.writeJSON(w, resp, http.StatusOK)
}

func (s *Server) handleSwagger(w http.ResponseWriter, r *http.Request) {
	const swaggerPath = "docs/swagger.yaml"
	data, err := os.ReadFile(swaggerPath)
	if err != nil {
		s.writeError(w, http.StatusNotFound, err)
		return
	}
	w.Header().Set("Content-Type", "application/yaml")
	_, _ = w.Write(data)
}

func (s *Server) handleSwaggerUI(w http.ResponseWriter, r *http.Request) {
	const page = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Traffic Guard Webscraper API</title>
  <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/5.11.0/swagger-ui.min.css" integrity="sha512-4imNURhjMZvXZAN5aryPU4LFShYTVPim+5m5i4eFteqSZFnwGQxgrFWSdcZ3EW26Y2jSGB3b4Mk5sJf8h2ivRA==" crossorigin="anonymous" referrerpolicy="no-referrer" />
  <style>
    html, body { margin: 0; padding: 0; height: 100%; }
    #swagger-ui { height: 100%; }
  </style>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/5.11.0/swagger-ui-bundle.min.js" integrity="sha512-OYKKMqAQ1cEEzjkWIA7ji9EqW1Fo4KcoetdksFLYe0PwFWboSwE6mIZaHsiDfmVRhzZkL5egYFNmKcs5MpQJyw==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
  <script>
    window.onload = function() {
      SwaggerUIBundle({
        url: '/swagger.yaml',
        dom_id: '#swagger-ui'
      });
    };
  </script>
</body>
</html>`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(page))
}

func (s *Server) writeJSON(w http.ResponseWriter, payload any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func (s *Server) writeError(w http.ResponseWriter, status int, err error) {
	resp := map[string]any{
		"error": err.Error(),
	}
	s.writeJSON(w, resp, status)
}

func checkKafka(ctx context.Context, cfg config.KafkaConfig) error {
	if len(cfg.Brokers) == 0 {
		return errors.New("no kafka brokers configured")
	}

	dialer, err := kafkaio.NewDialer(cfg)
	if err != nil {
		return err
	}

	broker := cfg.Brokers[0]
	conn, err := dialer.DialContext(ctx, "tcp", broker)
	if err != nil {
		return err
	}
	if err := conn.Close(); err != nil {
		return err
	}
	return nil
}
