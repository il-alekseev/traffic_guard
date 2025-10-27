package v1

import (
	"api-gateway/config"
	"api-gateway/pkg/analytics"
	"api-gateway/pkg/blogserv"
	"api-gateway/pkg/ctxcontrol"

	"api-gateway/pkg/models"
	"api-gateway/pkg/slogger"
	"api-gateway/pkg/usercontrol"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// @title API Gw
// @version 1.0

// @host localhost:8001
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token

type Server struct {
	userCl      *usercontrol.Usercontrol
	ctxCl       *ctxcontrol.Ctxcontrol
	blogCL      *blogserv.Blogserv
	analyticsCL *analytics.Analytics
	router      *gin.Engine
	domain      string
	httpServer  *http.Server
	key         string // для проверки подписи токенов
}

func New(
	ctx context.Context,
	cfg *config.Config,
	userCl *usercontrol.Usercontrol,
	ctxCl *ctxcontrol.Ctxcontrol,
	blogCL *blogserv.Blogserv,
	analyticsCL *analytics.Analytics,
) (*Server, error) {
	key, err := os.ReadFile(cfg.KeyCloak.PemFile)
	if err != nil {
		return nil, slogger.WrapError(ctx, err)
	}

	router := gin.New()
	apigw := Server{
		router:      router,
		userCl:      userCl,
		ctxCl:       ctxCl,
		blogCL:      blogCL,
		analyticsCL: analyticsCL,
		domain:      cfg.Swagger.Host,
		httpServer: &http.Server{
			Addr:    cfg.HTTP.Host + ":" + cfg.HTTP.Port,
			Handler: router,
		},
		key: string(key),
	}
	apigw.configureRouter()
	return &apigw, nil
}

func (s *Server) Run() error {
	slog.Info("Server started", "addr", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return s.httpServer.Shutdown(ctx)
}

func (s *Server) ErrorResponse(c *gin.Context, code int, msg string, err error) {
	errMsg := ""
	if err != nil {
		errMsg = ": " + err.Error()
	}
	slog.ErrorContext(c.Request.Context(), "Error response", "msg", msg, "err", err)
	c.JSON(code,
		models.DtoErrorResponse{
			Error: fmt.Sprintf("api-gw: %s%s", msg, errMsg),
		},
	)
}
