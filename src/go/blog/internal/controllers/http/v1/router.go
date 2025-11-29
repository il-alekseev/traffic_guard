package v1

import (
	middleware2 "fiermon-blog/internal/controllers/http/v1/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// initRouter - инициализация роутера
func (s *Server) initRouter() {
	s.router.GET("/api/v1/healthcheck", s.healthcheck)
	s.router.GET("/api/v1/version", s.version)

	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := s.router.Group("/api/v1")
	api.Use(middleware2.RecoveryMiddleware())
	api.Use(middleware2.CorsMiddleware())
	api.Use(middleware2.RequestIDMiddleware())
	api.Use(middleware2.LoggingMiddleware(s.logger))
	//Use(s.authMiddleware())
	//api.Use(middleware2.SetUserMetaData(s.logger))

	{
		api.POST("/add", middleware2.SetServiceMeta(s.logger), s.postAddRecord)
		api.GET("/logs", middleware2.SetUserMetaData(s.logger), middleware2.CheckAuthHeader(), s.validateParams(), s.getLogs)
	}
}
