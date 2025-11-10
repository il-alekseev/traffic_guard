package server

import (
	"fiermon-blog/internal/server/middleware"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// initRouter - инициализация роутера
func (s *Server) initRouter() {
	s.router.GET("/api/v1/healthcheck", s.healthcheck)

	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := s.router.Group("/api/v1")
	api.Use(middleware.CorsMiddleware())
	api.Use(middleware.RequestIDMiddleware())
	api.Use(middleware.LoggingMiddleware(s.logger))
	api.Use(middleware.RecoveryMiddleware())
	//Use(s.authMiddleware())
	api.Use(middleware.CheckAuthHeader())
	api.Use(middleware.SetUserMetaData(s.logger))

	{
		//api.GET("/healthcheck", s.healthcheck)
		api.GET("/logs", s.validateParams(), s.logs)
	}
}
