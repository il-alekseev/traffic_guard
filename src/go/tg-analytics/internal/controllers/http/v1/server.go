package v1

import (
	"context"
	"log/slog"
	"net/http"
	"tg-an/config"

	"github.com/gin-gonic/gin"
)

// @title           Analytics API
// @version         1.0
// @description     API для аналитики
// @BasePath        /
// @schemes         http
type Server struct {
	devVersion string
	server     *http.Server
	router     *gin.Engine
	l          slog.Logger
	u          UseCaseInterface
}

func New(cfg *config.Config, l slog.Logger, u UseCaseInterface) *Server {
	gin.SetMode(gin.DebugMode)

	r := gin.New()

	s := http.Server{
		Addr:    cfg.HTTP.Host + ":" + cfg.HTTP.Port,
		Handler: r,
	}
	server := Server{
		devVersion: cfg.App.DevVersion,
		server:     &s,
		router:     r,
		l:          l,
		u:          u,
	}
	server.configureRouter()
	return &server
}

func (s *Server) Start() error {
	s.l.Info("dashboard http server started")
	s.server.ListenAndServe()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
