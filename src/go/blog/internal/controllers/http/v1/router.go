package v1

import (
	middleware2 "fiermon-blog/internal/controllers/http/v1/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// initRouter - инициализация роутера
func (s *Server) initRouter() {
	s.router.GET("/api/v1/healthcheck", s.healthcheck)

	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := s.router.Group("/api/v1")
	api.Use(middleware2.CorsMiddleware())
	api.Use(middleware2.RequestIDMiddleware())
	api.Use(middleware2.LoggingMiddleware(s.logger))
	api.Use(middleware2.RecoveryMiddleware())
	//Use(s.authMiddleware())
	api.Use(middleware2.CheckAuthHeader())
	api.Use(middleware2.SetUserMetaData(s.logger))

	{
		api.POST("/add", s.postAddRecord)
		api.GET("/logs", s.validateParams(), s.getLogs)
	}
}
