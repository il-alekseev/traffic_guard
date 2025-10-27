package v1

import (
	"context"
	"dashboard_serv/config"
	"dashboard_serv/internal/controllers/http/v1/dto"
	"dashboard_serv/internal/usecase"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// @title NGFW Dashboard API
// @version 1.0
// @description API для дашборда статистики NGFW
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host 127.0.0.1:8081
// @BasePath /api/v1

type Server struct {
	url        string
	u          usecase.UseCase
	router     *gin.Engine
	httpServer *http.Server
	logger     *slog.Logger
}

// New - функция инициализации слоя http server
func New(cfg *config.Config, u usecase.UseCase, l slog.Logger) *Server {
	log := slog.New(&slog.JSONHandler{})
	router := gin.New()

	s := &Server{
		u:      u,
		router: router,
		url:    ":" + cfg.HTTP.Port,
		httpServer: &http.Server{
			Addr:    cfg.HTTP.Host + ":" + cfg.HTTP.Port,
			Handler: router,
		},
		logger: log,
	}

	// настройка маршрутов и middleware
	s.configureRouter()
	return s
}

// Run - запуск http server
func (s *Server) Run() error {
	log := slog.New(&slog.JSONHandler{})
	log.Info(fmt.Sprintf("Server started on %s", s.httpServer.Addr))
	return s.httpServer.ListenAndServe()
}

// Stop - остановка http server
func (s *Server) Stop(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	// TODO usercase.Close()
	defer cancel()

	return s.httpServer.Shutdown(ctx)
}

// ErrorResponse - обертка ответа ошибки
func (s *Server) ErrorResponse(c *gin.Context, code int, msg string, err error) {
	errMsg := ""
	if err != nil {
		errMsg = ": " + err.Error()
	}
	c.JSON(code,
		dto.ErrorResponse{
			Error: fmt.Sprintf("usercontrol: %s%s", msg, errMsg),
		},
	)
}
