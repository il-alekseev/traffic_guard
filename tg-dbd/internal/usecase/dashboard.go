package usecase

import (
	"context"
	"tg-dbd/internal/models"
	"time"
)

// GetSessions возвращает список сессий с пагинацией
func (u *Usecase) GetSessions(ctx context.Context, start time.Time, end time.Time, page int, limit int, filter string, orderBy string, orderDir string) ([]models.Session, int64, error) {
	u.l.Info("GetSessions called",
		"start", start, "end", end,
		"page", page, "limit", limit, "filter", filter,
		"orderBy", orderBy, "orderDir", orderDir)

	if err := ctx.Err(); err != nil {
		return nil, 0, err
	}

	// Тестовые данные сессий
	sessions := []models.Session{
		{
			Status:      "active",
			URL:         "https://corporate-app.com/dashboard",
			IP:          "192.168.1.100",
			NGFW:        "fw-01",
			Username:    "ivanov",
			SessionType: "web",
			Category:    "Корпоративные приложения",
			DateTime:    "2024-01-20 14:35:20",
		},
		{
			Status:      "blocked",
			URL:         "http://malicious-site.com",
			IP:          "192.168.1.101",
			NGFW:        "fw-02",
			Username:    "petrov",
			SessionType: "web",
			Category:    "Вредоносные сайты",
			DateTime:    "2024-01-20 14:20:15",
		},
		{
			Status:      "completed",
			URL:         "https://email-service.com/inbox",
			IP:          "192.168.1.102",
			NGFW:        "fw-01",
			Username:    "sidorov",
			SessionType: "web",
			Category:    "Электронная почта",
			DateTime:    "2024-01-20 14:10:45",
		},
		{
			Status:      "active",
			URL:         "https://cloud-storage.com/files",
			IP:          "192.168.1.103",
			NGFW:        "fw-03",
			Username:    "smirnov",
			SessionType: "web",
			Category:    "Облачные хранилища",
			DateTime:    "2024-01-20 13:55:30",
		},
		{
			Status:      "warned",
			URL:         "https://social-network.com/feed",
			IP:          "192.168.1.104",
			NGFW:        "fw-02",
			Username:    "kuznetsov",
			SessionType: "web",
			Category:    "Социальные сети",
			DateTime:    "2024-01-20 13:40:10",
		},
	}

	// Применяем пагинацию
	startIdx := (page - 1) * limit
	endIdx := startIdx + limit

	if startIdx >= len(sessions) {
		return []models.Session{}, int64(len(sessions)), nil
	}

	if endIdx > len(sessions) {
		endIdx = len(sessions)
	}

	paginatedSessions := sessions[startIdx:endIdx]
	totalCount := int64(len(sessions))

	return paginatedSessions, totalCount, nil
}

// GetTopDetectionsInfo возвращает информацию о топ обнаружениях
func (u *Usecase) GetDetections(ctx context.Context, start time.Time, end time.Time, count int, filter string) ([]models.Detection, error) {
	u.l.Info("GetDetections called", "start", start, "end", end, "filter", filter, "topCount", count)

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	detections := []models.Detection{
		{
			URL:         "http://malicious-site.com/download.exe",
			DateTime:    "2024-01-20 14:25:30",
			Category:    "Вредоносное ПО",
			Description: "Попытка загрузки потенциально опасного файла",
			Host:        "ws-1245",
			NGFW:        "fw-01",
			AccessCount: 15,
			Action:      "blocked",
		},
		{
			URL:         "https://phishing-bank.com/login",
			DateTime:    "2024-01-20 13:40:22",
			Category:    "Фишинг",
			Description: "Доступ к фишинговому сайту",
			Host:        "ws-0678",
			NGFW:        "fw-02",
			AccessCount: 8,
			Action:      "blocked",
		},
		{
			URL:         "http://torrent-tracker.org/movie.torrent",
			DateTime:    "2024-01-20 12:15:45",
			Category:    "P2P",
			Description: "Попытка доступа к торрент-трекеру",
			Host:        "ws-3321",
			NGFW:        "fw-01",
			AccessCount: 12,
			Action:      "warned",
		},
		{
			URL:         "https://social-media.com/private",
			DateTime:    "2024-01-20 11:30:15",
			Category:    "Социальные сети",
			Description: "Доступ к социальной сети в рабочее время",
			Host:        "ws-4456",
			NGFW:        "fw-03",
			AccessCount: 45,
			Action:      "allowed",
		},
	}

	// Ограничиваем количество согласно topCount
	if count > 0 && int(count) < len(detections) {
		detections = detections[:count]
	}

	return detections, nil
}

// GetDetectionsStat возвращает статистику обнаружений за указанный период
func (u *Usecase) GetDetectionsStat(ctx context.Context, start time.Time, end time.Time, filter string) (models.DetectionStat, error) {
	u.l.Info("GetDetectionsStat called",
		"start", start,
		"end", end,
		"filter", filter)

	if err := ctx.Err(); err != nil {
		return models.DetectionStat{}, err
	}

	// Заглушка с тестовыми данными
	stat := models.DetectionStat{
		Detected:   156,
		Accepted:   1240,
		Denied:     89,
		Unresolved: 23,
	}

	return stat, nil
}

// GetCategories возвращает топ категорий
func (u *Usecase) GetCategories(ctx context.Context, start time.Time, end time.Time, count int, filter string) (map[string]int, error) {
	u.l.Info("GetTopCategories called", "start", start, "end", end, "filter", filter, "topCount", count)

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	categories := map[string]int{
		"Социальные сети":    156,
		"Электронная почта":  142,
		"Облачные хранилища": 128,
		"Финансы":            115,
		"Новости":            98,
		"Развлечения":        87,
		"Игры":               76,
		"Образование":        65,
	}

	// Ограничиваем количество согласно topCount
	if count > 0 {
		limited := make(map[string]int)
		count1 := 0
		for k, v := range categories {
			if count1 >= count {
				break
			}
			limited[k] = v
			count1++
		}
		return limited, nil
	}

	return categories, nil
}

// GetResources возвращает список ресурсов с статистикой запросов
func (u *Usecase) GetResources(ctx context.Context, start time.Time, end time.Time, count int, filter string) ([]models.Resource, error) {
	u.l.Info("GetResources called",
		"start", start,
		"end", end,
		"count", count,
		"filter", filter)

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	// Заглушка с тестовыми данными
	resources := []models.Resource{
		{
			Name:            "youtube.com",
			ReqsBeforeBlock: 245,
			ReqsAfterBlock:  12,
		},
		{
			Name:            "facebook.com",
			ReqsBeforeBlock: 198,
			ReqsAfterBlock:  8,
		},
		{
			Name:            "instagram.com",
			ReqsBeforeBlock: 176,
			ReqsAfterBlock:  15,
		},
		{
			Name:            "twitter.com",
			ReqsBeforeBlock: 154,
			ReqsAfterBlock:  5,
		},
		{
			Name:            "tiktok.com",
			ReqsBeforeBlock: 132,
			ReqsAfterBlock:  22,
		},
		{
			Name:            "netflix.com",
			ReqsBeforeBlock: 121,
			ReqsAfterBlock:  3,
		},
		{
			Name:            "spotify.com",
			ReqsBeforeBlock: 109,
			ReqsAfterBlock:  7,
		},
		{
			Name:            "reddit.com",
			ReqsBeforeBlock: 98,
			ReqsAfterBlock:  11,
		},
	}

	// Ограничиваем количество согласно count
	if count > 0 && count < len(resources) {
		resources = resources[:count]
	}

	return resources, nil
}

// GetEvents возвращает список событий за указанный период
func (u *Usecase) GetEvents(ctx context.Context, start time.Time, end time.Time, count int, filter string) ([]models.Notification, error) {
	u.l.Info("GetEvents called",
		"start", start,
		"end", end,
		"count", count,
		"filter", filter)

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	// Заглушка с тестовыми данными событий
	events := []models.Notification{
		{
			Description: "Обнаружена подозрительная активность пользователя",
			Username:    "ivanov",
			DateTime:    "2024-01-20 14:30:25",
		},
		{
			Description: "Превышено количество попыток входа в систему",
			Username:    "petrov",
			DateTime:    "2024-01-20 13:15:10",
		},
		{
			Description: "Попытка доступа к запрещенному ресурсу",
			Username:    "sidorov",
			DateTime:    "2024-01-20 12:45:33",
		},
		{
			Description: "Необычная активность в нерабочее время",
			Username:    "smirnov",
			DateTime:    "2024-01-20 11:20:47",
		},
		{
			Description: "Обнаружена DDoS атака на внешний интерфейс",
			Username:    "system",
			DateTime:    "2024-01-20 10:35:15",
		},
		{
			Description: "Подозрительный трафик с рабочей станции",
			Username:    "kuznetsov",
			DateTime:    "2024-01-20 09:50:22",
		},
		{
			Description: "Попытка обхода системы безопасности",
			Username:    "popov",
			DateTime:    "2024-01-20 08:15:18",
		},
	}

	// Ограничиваем количество согласно count
	if count > 0 && count < len(events) {
		events = events[:count]
	}

	return events, nil
}

// GetDevicesStat возвращает статистику по устройствам за указанный период
func (u *Usecase) GetDevicesStat(ctx context.Context, start time.Time, end time.Time) ([]models.DeviceStat, error) {
	u.l.Info("GetDevicesStat called",
		"start", start,
		"end", end)

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	// Заглушка с тестовыми данными устройств
	devices := []models.DeviceStat{
		{
			Name:  "NGFW-01",
			State: "active",
			Statistics: []models.ResourcePoint{
				{Date: 1705759200, Locked: 45, Delayed: 12},
				{Date: 1705845600, Locked: 38, Delayed: 8},
				{Date: 1705932000, Locked: 52, Delayed: 15},
				{Date: 1706018400, Locked: 41, Delayed: 10},
				{Date: 1706104800, Locked: 49, Delayed: 13},
			},
		},
		{
			Name:  "NGFW-02",
			State: "active",
			Statistics: []models.ResourcePoint{
				{Date: 1705759200, Locked: 32, Delayed: 7},
				{Date: 1705845600, Locked: 28, Delayed: 5},
				{Date: 1705932000, Locked: 41, Delayed: 11},
				{Date: 1706018400, Locked: 35, Delayed: 8},
				{Date: 1706104800, Locked: 38, Delayed: 9},
			},
		},
		{
			Name:  "NGFW-03",
			State: "warning",
			Statistics: []models.ResourcePoint{
				{Date: 1705759200, Locked: 18, Delayed: 3},
				{Date: 1705845600, Locked: 22, Delayed: 4},
				{Date: 1705932000, Locked: 15, Delayed: 2},
				{Date: 1706018400, Locked: 25, Delayed: 6},
				{Date: 1706104800, Locked: 20, Delayed: 5},
			},
		},
		{
			Name:  "NGFW-04",
			State: "inactive",
			Statistics: []models.ResourcePoint{
				{Date: 1705759200, Locked: 0, Delayed: 0},
				{Date: 1705845600, Locked: 0, Delayed: 0},
				{Date: 1705932000, Locked: 0, Delayed: 0},
				{Date: 1706018400, Locked: 0, Delayed: 0},
				{Date: 1706104800, Locked: 0, Delayed: 0},
			},
		},
	}

	return devices, nil
}

// GetTrafficStat возвращает статистику трафика за указанный период
func (u *Usecase) GetTrafficStat(ctx context.Context, start time.Time, end time.Time, filter string) ([]models.TrafficPoint, error) {
	u.l.Info("GetTrafficStat called",
		"start", start,
		"end", end,
		"filter", filter)

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	// Заглушка с тестовыми данными трафика
	traffic := []models.TrafficPoint{
		{Date: 1705759200, Input: 1450, Output: 890},
		{Date: 1705845600, Input: 1620, Output: 1020},
		{Date: 1705932000, Input: 1380, Output: 760},
		{Date: 1706018400, Input: 1780, Output: 1150},
		{Date: 1706104800, Input: 1520, Output: 940},
		{Date: 1706191200, Input: 1950, Output: 1280},
		{Date: 1706277600, Input: 1420, Output: 810},
		{Date: 1706364000, Input: 1680, Output: 990},
		{Date: 1706450400, Input: 1830, Output: 1120},
		{Date: 1706536800, Input: 1570, Output: 870},
	}

	return traffic, nil
}

// GetRequestsStat возвращает статистику запросов за указанный период
func (u *Usecase) GetRequestsStat(ctx context.Context, start time.Time, end time.Time, filter string) (models.RequestsStat, error) {
	u.l.Info("GetRequestsStat called",
		"start", start,
		"end", end,
		"filter", filter)

	if err := ctx.Err(); err != nil {
		return models.RequestsStat{}, err
	}

	// Заглушка с тестовыми данными статистики запросов
	stat := models.RequestsStat{
		Accepted: []models.RequestPoint{
			{Date: 1705759200, Count: 1240},
			{Date: 1705845600, Count: 1320},
			{Date: 1705932000, Count: 1180},
			{Date: 1706018400, Count: 1450},
			{Date: 1706104800, Count: 1280},
		},
		Blocked: []models.RequestPoint{
			{Date: 1705759200, Count: 89},
			{Date: 1705845600, Count: 76},
			{Date: 1705932000, Count: 102},
			{Date: 1706018400, Count: 68},
			{Date: 1706104800, Count: 94},
		},
		BeforeBlocked: []models.RequestPoint{
			{Date: 1705759200, Count: 245},
			{Date: 1705845600, Count: 198},
			{Date: 1705932000, Count: 312},
			{Date: 1706018400, Count: 187},
			{Date: 1706104800, Count: 276},
		},
		Delayed: []models.RequestPoint{
			{Date: 1705759200, Count: 45},
			{Date: 1705845600, Count: 38},
			{Date: 1705932000, Count: 52},
			{Date: 1706018400, Count: 41},
			{Date: 1706104800, Count: 49},
		},
	}

	return stat, nil
}

// GetProhActSchedule возвращает расписание запрещенных активностей по дням за указанный период
func (u *Usecase) GetProhActSchedule(ctx context.Context, start time.Time, filter string) (map[int64]int, error) {
	u.l.Info("GetProhActSchedule called",
		"start", start,
		"filter", filter)

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	// Заглушка с тестовыми данными расписания запрещенных активностей по дням
	schedule := map[int64]int{
		1705708800: 156, // 2024-01-20
		1705795200: 142, // 2024-01-21
		1705881600: 198, // 2024-01-22
		1705968000: 124, // 2024-01-23
		1706054400: 167, // 2024-01-24
		1706140800: 189, // 2024-01-25
		1706227200: 145, // 2024-01-26
		1706313600: 176, // 2024-01-27
		1706400000: 132, // 2024-01-28
		1706486400: 154, // 2024-01-29
		1706572800: 168, // 2024-01-30
		1706659200: 142, // 2024-01-31
		1706745600: 195, // 2024-02-01
		1706832000: 178, // 2024-02-02
		1706918400: 123, // 2024-02-03
	}

	return schedule, nil
}

// GetAnomaly возвращает статистику аномалий за указанный период
func (u *Usecase) GetAnomalies(ctx context.Context, start time.Time, end time.Time) (models.Anomaly, error) {
	u.l.Info("GetAnomaly called",
		"start", start,
		"end", end)

	if err := ctx.Err(); err != nil {
		return models.Anomaly{}, err
	}

	// Заглушка с тестовыми данными аномалий
	anomaly := models.Anomaly{
		BlockedResourcesCount: 45,
		FirewallsCount:        3,
		Stat: []models.AnomalyResourse{
			{
				ResourseName:      "malicious-site.com",
				UnblockedRequests: 23,
			},
			{
				ResourseName:      "phishing-bank.com",
				UnblockedRequests: 18,
			},
			{
				ResourseName:      "torrent-tracker.org",
				UnblockedRequests: 15,
			},
			{
				ResourseName:      "suspicious-download.net",
				UnblockedRequests: 12,
			},
			{
				ResourseName:      "unknown-proxy.ru",
				UnblockedRequests: 8,
			},
			{
				ResourseName:      "gambling-casino.com",
				UnblockedRequests: 6,
			},
			{
				ResourseName:      "crypto-mining.pool",
				UnblockedRequests: 5,
			},
		},
	}

	return anomaly, nil
}
