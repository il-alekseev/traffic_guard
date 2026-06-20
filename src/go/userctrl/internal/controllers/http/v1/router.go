package v1

import (
	_ "userctrl/docs" // Импорт сгенерированных docs
	"userctrl/internal/controllers/http/v1/middleware"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (s *Server) configureRouter() {
	// Swagger endpoint
	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	s.router.GET("/v1/healthcheck", s.healthcheck)

	// Global middleware
	s.router.Use(middleware.CorsMiddleware())
	s.router.Use(middleware.RequestIDMiddleware())
	s.router.Use(middleware.LoggingMiddleware())
	s.router.Use(middleware.RecoveryMiddleware())

	// Auth routes
	auth := s.router.Group("/v1/auth") // SA, CA, CO
	{
		auth.POST("/sign-in", s.signIn)
		auth.POST("/sign-out",
			middleware.CheckAuthHeader(),
			s.signOut)
		auth.POST("/refresh",
			//middleware.CheckAuthHeader(),
			s.refresh)
	}

	apiRouter := s.router.Group("/v1")

	//Use(s.authMiddleware())
	apiRouter.Use(middleware.CheckAuthHeader())
	apiRouter.Use(middleware.SetUserMetaData(*s.logger))

	users := apiRouter.Group("/users")
	{
		users.GET("", s.getUsers)                                         // SA, CA
		users.POST("", s.createUser)                                      // SA, CA
		users.GET("/profile", s.getProfile)                               // SA, CA, CO
		users.PUT("/profile/pass", s.updatePassword)                      // SA, CA, CO
		users.GET("/count", s.getUsersCount)                              // SA
		users.GET("/count/context/:contextID", s.getUsersCountForContext) // SA, CA
		users.GET("/:userId", s.getUser)                                  // SA, CA, CO
		users.PUT("/:userId", s.updateUser)                               // SA, CA, CO
		users.DELETE("/:userId", s.deleteUser)                            // SA, CA
		users.PUT("/:userId/pass/otp", s.resetPassword)                   // SA, CA
		users.PUT("/:userId/roles", s.updateUserRoles)                    // SA, CA
		users.DELETE("/:userId/roles", s.deleteUserRoles)                 // SA, CA
		users.GET("/count-by-role", s.getCountUsersByRole)                // SA, CA
	}

	roles := apiRouter.Group("/roles")
	{
		roles.GET("", s.getRoles)                                  // SA
		roles.GET("/count", s.getRolesCount)                       // SA
		roles.PUT("/context", s.createRolesAndGroupsForContext)    // SA
		roles.DELETE("/context", s.deleteRolesAndGroupsForContext) // SA
	}
}
