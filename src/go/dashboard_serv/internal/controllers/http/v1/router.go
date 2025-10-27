package v1

import (
	//_ "dashboard_serv/docs" // Импорт сгенерированных docs

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (s *Server) configureRouter() {
	// Swagger endpoint
	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// // Группа API для дашбордов
	// api := s.router.Group("/api/v1")
	// {
	// 	dashboard := api.Group("/dashboard")
	// 	{
	// 		// dashboard.POST("/requests", s.GetRequestStats)
	// 		// dashboard.POST("/traffic", s.GetTrafficStats)
	// 		// dashboard.POST("/top-resources", s.GetTopResources)
	// 		// dashboard.POST("/blocked-categories", s.GetBlockedCategories)
	// 		// dashboard.POST("/blocked-resources", s.GetBlockedResources)
	// 		// dashboard.GET("/ngfw", s.GetNGFWList)
	// 		// dashboard.GET("/categories", s.GetCategoriesList)
	// 	}
	// }

}
