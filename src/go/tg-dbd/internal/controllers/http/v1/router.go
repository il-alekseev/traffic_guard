package v1

import (
	"fmt"
	_ "tg-dbd/docs" // Импорт сгенерированных docs
	"tg-dbd/internal/controllers/http/v1/middleware"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (s *Server) configureRouter() {
	s.router.Use(middleware.CorsMiddleware())
	// Сваггер
	// Динамический адрес для сваггера
	swaggerURL := ginSwagger.URL(fmt.Sprintf("http://%s/swagger/doc.json", s.server.Addr))
	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, swaggerURL))
	v1 := s.router.Group("/api/v1")
	{
		// Common
		v1.GET("/devices", s.GetDevices)
		v1.GET("/categories", s.GetContentCategories)
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
		v1.GET("/healthcheck", s.Healthcheck)
		v1.GET("/version", s.Version)
	}
}
