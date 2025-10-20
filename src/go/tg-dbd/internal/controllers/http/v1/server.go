package v1

import (
	"context"
	"log/slog"
	"net/http"
	"tg-dbd/config"

	"github.com/gin-gonic/gin"
)

// @title           Dashboard API
// @version         1.0
// @description     API для получения дашбордов
// @host            127.0.0.1:7000
// @BasePath        /api/v1
// @schemes         http
type Server struct {
	server *http.Server
	router *gin.Engine
	l      slog.Logger
	u      UseCaseInterface
}

func New(cfg *config.HTTP, l slog.Logger, u UseCaseInterface) *Server {
	gin.SetMode(gin.DebugMode)

	r := gin.New()

	s := http.Server{
		Addr:    cfg.Host + ":" + cfg.Port,
		Handler: r,
	}
	server := Server{
		server: &s,
		router: r,
		l:      l,
		u:      u,
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
