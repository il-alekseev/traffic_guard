package v1

import (
	"fiermon-blog/internal/controllers/http/v1/dto"
	"fiermon-blog/pkg/fslog/wsl"
	"fmt"
	"log/slog"
	"net/http"

	"fiermon-blog/config"
	"fiermon-blog/internal/service"

	"github.com/gin-gonic/gin"
)

// @title           app for get logs
// @version         1.0
// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token

type Server struct {
	httpServer *http.Server
	router     *gin.Engine
	logger     *slog.Logger
	serv       *service.Service
	devVersion string
}

func NewServer(cfg config.HTTP, log *slog.Logger, serv *service.Service, version string) *Server {
	log = log.With(wsl.Label("layer", "http_server"))

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	httpServer := &http.Server{
		Addr:              cfg.Host + ":" + cfg.Port,
		ReadHeaderTimeout: cfg.Timeout,
		IdleTimeout:       cfg.IdleTimeout,
		Handler:           router,
	}

	s := &Server{
		httpServer: httpServer,
		logger:     log,
		serv:       serv,
		router:     router,
		devVersion: version,
	}

	s.initRouter()

	return s
}

// Start - запуск http сервера
func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

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
