package v1

import (
	"context"
	"userctrl/config"
	"userctrl/internal/controllers/http/v1/dto"
	"userctrl/pkg/slogger/wsl"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// @title User Control API
// @version 1.0
// @description API для управления пользователями, ролями и аутентификации

// @host localhost:8000
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token

type Server struct {
	url        string
	u          UseCaseInterface
	router     *gin.Engine
	httpServer *http.Server
	logger     *slog.Logger
}

// New - функция инициализации слоя http server
func New(cfg *config.Config, u UseCaseInterface, l slog.Logger) *Server {
	log := l.With(wsl.Label("layer", "http"))
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
	log := s.logger.With(wsl.Label("method", "Run"))
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
