package v1

import (
	"context"
	"log/slog"
	"net/http"
	"tg-etl/config"

	"github.com/gin-gonic/gin"
)

// @title           ETL API
// @version         1.0
// @description     API для ETL
// @host            127.0.0.1:7001
// @BasePath        /api/v1
// @schemes         http
type Server struct {
	devVersion string
	server     *http.Server
	router     *gin.Engine
	l          slog.Logger
}

func New(cfg *config.Config, l slog.Logger) *Server {
	gin.SetMode(gin.DebugMode)

	r := gin.New()

	s := http.Server{
		Addr:    cfg.HTTP.Host + ":" + cfg.HTTP.Port,
		Handler: r,
	}
	server := Server{
		devVersion: cfg.DevVersion,
		server:     &s,
		router:     r,
		l:          l,
	}
	server.configureRouter()
	return &server
}

func (s *Server) Start() error {
	s.l.Info("etl http server started")
	s.server.ListenAndServe()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
