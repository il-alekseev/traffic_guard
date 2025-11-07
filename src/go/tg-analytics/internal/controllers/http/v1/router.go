package v1

import (
	"fmt"
	"os"
	_ "tg-an/docs" // Импорт сгенерированных docs
	"tg-an/internal/controllers/http/v1/middleware"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (s *Server) configureRouter() {
	s.router.Use(middleware.CorsMiddleware())
	// Сваггер
	// Динамический адрес для сваггера
	addr := os.Getenv("AN_SWAGGER")
	if addr == "" {
		addr = s.server.Addr
	}
	swaggerURL := ginSwagger.URL(fmt.Sprintf("http://%s/swagger/doc.json", addr))
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
			dashboards.GET("/top-unresolved_detections", s.GetTopUnresolvedDetections)
			dashboards.GET("/anomalies", s.GetAnomalies)
			dashboards.GET("/devices", s.GetDeviceStat)
			dashboards.GET("/traffic", s.GetTrafficStat)
			dashboards.GET("/requests", s.GetRequestStat)
		}
		// Вкладка Выявления
		detections := v1.Group("/detections")
		{
			detections.GET("/", s.GetDetections)
			detections.GET("/stat", s.GetDetectionStat)
			detections.PATCH("/act", s.Act)
		}
		// Вкладка Отыеты
		reports := v1.Group("/reports")
		{
			reports.GET("/", s.CreateReport)
			reports.GET("/:hostname", s.CreateReportForDevice)
		}
		// Утилиты
		v1.GET("/healthcheck", s.Healthcheck)
		v1.GET("/version", s.Version)
	}
}
