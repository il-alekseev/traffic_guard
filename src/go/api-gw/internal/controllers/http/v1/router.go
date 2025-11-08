package v1

import (
	_ "api-gateway/docs" // Импорт сгенерированных docs
	"api-gateway/internal/controllers/http/v1/middleware"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (s *Server) configureRouter() {
	// Swagger endpoint
	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	s.router.GET("/api/v1/healthcheck", s.healthcheck)
	//s.router.GET("/api/v1/version", s.version)

	// Global middleware
	s.router.Use(middleware.CorsMiddleware())
	// s.router.Use(middleware.RequestIDMiddleware()) // Функционал объединен с LoggingMiddleware()
	s.router.Use(middleware.LoggingMiddleware())

	// Auth routes
	auth := s.router.Group("/api/v1/auth")
	{
		auth.POST("/sign-in", s.signIn)
		auth.POST("/sign-out",
			middleware.CheckAuthHeader(),
			//middleware.VerifyToken(),
			s.signOut)
		auth.POST("/refresh",
			//middleware.CheckAuthHeader(),
			//middleware.VerifyToken(),
			s.refresh)
	}

	apiRouter := s.router.Group("")
	apiRouter.Use(middleware.CheckAuthHeader()) // Проверка наличия токена + формат
	// apiRouter.Use(middleware.VerifyToken(s.key)) // Проверка валидности токена
	apiRouter.Use(middleware.SetUserMetaData()) // Установка метаданных пользователя

	// Users routes
	users := apiRouter.Group("/api/v1/users")
	{
		users.GET("", s.getUsers)                          // SA, CA
		users.POST("", s.createUser)                       // SA, CA
		users.GET("/profile", s.getProfile)                // SA, CA, CO
		users.PUT("/profile/pass", s.updatePassword)       // SA, CA, CO
		users.GET("/count", s.getUsersCount)               // SA, CA 		// TODO: не очень то и нужно
		users.GET("/:userId", s.getUser)                   // SA, CA, CO
		users.PUT("/:userId", s.updateUser)                // SA, CA, CO
		users.DELETE("/:userId", s.deleteUser)             // SA, CA
		users.PUT("/:userId/pass/otp", s.resetPassword)    // SA, CA
		users.PUT("/:userId/roles", s.updateUserRoles)     // SA, CA
		users.DELETE("/:userId/roles", s.deleteUserRoles)  // SA, CA
		users.GET("/count-by-role", s.getCountUsersByRole) // SA
	}

	roles := apiRouter.Group("/api/v1/roles")
	{
		roles.GET("", s.getRoles)            // SA
		roles.GET("/count", s.getRolesCount) // SA
	}

	ctxcontrol := apiRouter.Group("/api/v1/contexts")
	{
		ctxcontrol.GET("", s.getContexts)                   // SA
		ctxcontrol.POST("", s.createContext)                // SA
		ctxcontrol.GET("/count", s.getContextsCount)        // SA
		ctxcontrol.GET("/:context_id", s.getContextByID)    // SA, CA
		ctxcontrol.PUT("/:context_id", s.updateContextByID) // SA, CA
		ctxcontrol.DELETE("/:context_id", s.deleteContext)  // SA
		ctxcontrol.GET("/free-ports", s.getFreePorts)
	}

	blog := apiRouter.Group("/api/v1")
	{
		blog.GET("/logs", s.logs) // SA, CA
	}

	analytics := apiRouter.Group("/api/v1/analytics")
	{
		// common
		analytics.GET("/categories", s.getCategories)
		analytics.GET("/devices", s.getDevices)

		//dashboards
		analytics.GET("/dashboards/anomalies", s.getV1DashboardsAnomalies)
		analytics.GET("/dashboards/devices", s.getV1DashboardsDevices)
		analytics.GET("/dashboards/requests", s.getDashboardsRequests)
		analytics.GET("/dashboards/top-categories", s.getDashboardsTopCategories)
		analytics.GET("/dashboards/top-unresolved_detections", s.getDashboardsTopUnresolvedDetections)
		analytics.GET("/dashboards/traffic", s.getDashboardsTraffic)

		//actions
		analytics.PATCH("/dashboards/act", s.patchV1DashboardsAct)

		//detections
		analytics.GET("/detections", s.getDetections)
		analytics.GET("/detections/stat", s.getDetectionsStat)

		//reports
		analytics.GET("/reports", s.getV1Reports)
		analytics.GET("/reports/:hostname", s.getV1ReportsHostname)

		//sessions
		analytics.GET("/sessions", s.getSessions)

		//utils
		// analytics.GET("/v1/healthcheck", s.getV1Healthcheck)
		// analytics.GET("/version", s.getVersion)

		//not implemented

	}
}
