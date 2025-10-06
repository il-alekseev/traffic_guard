package v1

func (s *Server) configureRouter() {
	v1 := s.router.Group("/v1")
	{
		//dashboards := v1.Group("/dashboards")
		//{
		//	dashboards.GET()
		//}
		v1.GET("healthcheck", s.Healthcheck)
	}
}
