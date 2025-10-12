package v1

import (
	_ "tg-dbd/docs" // Импорт сгенерированных docs

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (s *Server) configureRouter() {
	// Сваггер
	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	v1 := s.router.Group("/v1")
	{
		// Вкладка dashboard
		dashboards := v1.Group("/dashboards")
		{
			dashboards.GET("/categories", s.getCategories)
			dashboards.GET("/resources", s.GetResources)
			dashboards.GET("/events", s.GetEvents)
			dashboards.GET("/anomalies", s.GetAnomalies)
			dashboards.GET("/proh_activity", s.GetProhActivity)
			dashboards.GET("/devices", s.GetDevicesStat)
			dashboards.GET("/traffic", s.GetTrafficStat)
			dashboards.GET("/requests", s.getCategories)
		}
		// Вкладка Выявления
		detections := v1.Group("/detections")
		{
			detections.GET("/", s.getDetections)
			detections.GET("/stat", s.getDetectionStat)
		}
		// Вкладка Сессии
		v1.GET("/sessions", s.getSessions)
		// Утилиты
		v1.GET("healthcheck", s.Healthcheck)
	}
}
