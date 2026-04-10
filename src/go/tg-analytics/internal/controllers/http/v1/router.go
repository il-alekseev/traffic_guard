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
	s.router.HandleMethodNotAllowed = true
	// Сваггер
	// Динамический адрес для сваггера
	addr := os.Getenv("AN_SWAGGER")
	if addr == "" {
		addr = s.server.Addr
	}
	swaggerURL := ginSwagger.URL(fmt.Sprintf("http://%s/swagger/doc.json", addr))
	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, swaggerURL))

	// Утилиты
	s.router.GET("/api/v1/healthcheck", s.Healthcheck)
	s.router.GET("/api/v1/version", s.Version)
	//Common
	s.router.GET("/api/v1/categories", s.GetContentCategories)

	// Middleware
	s.router.Use(middleware.CorsMiddleware())
	s.router.Use(middleware.RequestIDMiddleware())
	s.router.Use(middleware.CheckAuthHeader())
	s.router.Use(middleware.SetUserMetaData())

	v1 := s.router.Group("/api/v1")
	{
		// Common
		v1.GET("/devices", s.GetDevices)
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
			dashboards.GET("/proh-activity", s.GetProhActivity)
		}
		// Вкладка Выявления
		detections := v1.Group("/detections")
		{
			detections.GET("/", s.GetDetections)
			detections.GET("/stat", s.GetDetectionStat)
			detections.PATCH("/act", s.Act)
			// Версия запроса с получением процентов по всем негативным категориям домена
			detections.GET("/detections_v2", s.GetDetections_v2)
		}
		// Вкладка Отыеты
		reports := v1.Group("/reports")
		{
			reports.GET("/", s.CreateReport)
			reports.GET("/:hostname", s.CreateReportForDevice)
		}
	}
}
