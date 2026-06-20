package v1

import (
	"fmt"
	"os"
	_ "tg-etl/docs" // Импорт сгенерированных docs
	"tg-etl/internal/controllers/http/v1/middleware"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (s *Server) configureRouter() {
	s.router.Use(middleware.CorsMiddleware())
	// Сваггер
	// Динамический адрес для сваггера
	addr := os.Getenv("ETL_SWAGGER")
	if addr == "" {
		addr = s.server.Addr
	}
	swaggerURL := ginSwagger.URL(fmt.Sprintf("http://%s/swagger/doc.json", addr))
	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, swaggerURL))
	v1 := s.router.Group("/api/v1")
	{
		// Утилиты
		v1.GET("/healthcheck", s.Healthcheck)
		v1.GET("/version", s.Version)
	}
}
