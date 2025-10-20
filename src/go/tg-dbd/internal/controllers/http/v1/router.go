package v1

import (
	_ "tg-dbd/docs" // Импорт сгенерированных docs
	"tg-dbd/internal/controllers/http/v1/middleware"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (s *Server) configureRouter() {
	s.router.Use(middleware.CorsMiddleware())
	// Сваггер
	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	v1 := s.router.Group("/api/v1")
	{
		// Вкладка Сессии
		v1.GET("/sessions", s.GetSessions)
		// Вкладка dashboard
		dashboards := v1.Group("/dashboards")
		{
			dashboards.GET("/top-categories", s.GetTopCategories)
			dashboards.GET("/resources", s.GetResources)
			dashboards.GET("/anomalies", s.GetAnomalies)
			dashboards.GET("/proh_activity", s.GetProhActivity)
			dashboards.GET("/devices", s.GetDevicesStat)
			dashboards.GET("/traffic", s.GetTrafficStat)
			dashboards.GET("/requests", s.GetRequestsStat)
		}
		// Вкладка Выявления
		detections := v1.Group("/detections")
		{
			detections.GET("/", s.GetTopDetections)
			detections.GET("/stat", s.GetDetectionStat)
		}
		// Утилиты
		v1.GET("healthcheck", s.Healthcheck)
	}
}
