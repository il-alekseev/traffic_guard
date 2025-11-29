package v1

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"tg-an/config"
	"tg-an/internal/controllers/http/v1/dto"

	"github.com/gin-gonic/gin"
)

// @title           Analytics API
// @version         1.0
// @description     API для аналитики
// @BasePath        /
// @schemes         http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token
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
	s.l.Info("tg-analytics http server started")
	s.server.ListenAndServe()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// ErrorResponse формирует стандартизированный JSON ответ с ошибкой.
// Также логирует ошибку с контекстом запроса.
//
// Параметры:
//   - c: контекст Gin
//   - code: HTTP статус код
//   - msg: сообщение об ошибке
//   - err: оригинальная ошибка (может быть nil)
func (s *Server) ErrorResponse(c *gin.Context, code int, msg string, err error) {
	m := ""
	if msg == "" {
		if err != nil {
			m = "error=" + err.Error()
		}
	} else {
		m = msg
		if err != nil {
			m = ", error=" + err.Error()
		}
	}
	s.l.ErrorContext(c.Request.Context(), "ErrorResponse", "code", code, "err", m)
	c.JSON(code,
		dto.ErrorResponse{
			Error: fmt.Sprintf("tg-analytics: %s", m),
		},
	)
}
