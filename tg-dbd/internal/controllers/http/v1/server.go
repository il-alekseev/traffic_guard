package v1

import (
	"context"
	"net/http"
	"tg-dbd/config"
	"tg-dbd/pkg/logger"

	"github.com/gin-gonic/gin"
)

// TODO: переделать под slog
type Server struct {
	server *http.Server
	router *gin.Engine
	logger logger.Interface
	u      UseCaseInterface
}

func New(cfg *config.HTTP, l logger.Interface, u UseCaseInterface) *Server {
	gin.SetMode(gin.DebugMode)

	r := gin.New()

	s := http.Server{
		Addr:    cfg.Host + ":" + cfg.Port,
		Handler: r,
	}
	server := Server{
		server: &s,
		router: r,
		logger: l,
		u:      u,
	}
	server.configureRouter()
	return &server
}

func (s *Server) Start() error {
	s.logger.Info("dashboard http server started")
	s.server.ListenAndServe()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}