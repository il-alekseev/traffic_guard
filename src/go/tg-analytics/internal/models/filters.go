package models

// Sessions
type SessionFilter struct {
	HostName string `form:"hostname"`
	Category string `form:"category"`
	Type     string `form:"type"`
	Status   string `form:"status"`
}

// Dashboards
type CategoryFilter struct {
	HostName string `form:"hostname"`
	Type     string `form:"type"`
}

type DashboardFilter struct {
	HostName    string `form:"hostname"`
	TopCategory string `form:"top_category"`
}

// Detections
type DetectionFilter struct {
	HostName    string `form:"hostname"`
	TopCategory string `form:"top_category"`
	Status      string `form:"_status"` // TODO: убрать _ после обновления фронта
	// TODO: добавить фильтры по комментарию:
	//В фильтрах можно скрыть решенные выявления (заблокированные/разрешенные), указать NGFW, выбор локации, выбор категории (наркотики/Экстремизм...)
}
