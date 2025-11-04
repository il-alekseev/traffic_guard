package usecase

import (
	"context"
	"log/slog"
	"tg-an/internal/models"
	"time"
)

// GetDevicesStat возвращает статистику по устройствам за указанный период
func (u *Usecase) GetDevicesStat(ctx context.Context, start time.Time, end time.Time) ([]models.DeviceStat, error) {
	u.l.InfoContext(ctx,
		"GetDevicesStat called",
		slog.String("start", start.String()),
		slog.String("end", end.String()),
	)

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

// GetProhActSchedule возвращает расписание запрещенных активностей по дням за указанный период
func (u *Usecase) GetProhActSchedule(ctx context.Context, start time.Time, filter string) (map[int64]int, error) {
	u.l.InfoContext(ctx,
		"GetProhActSchedule called",
		slog.String("start", start.String()),
		slog.String("filter", filter),
	)

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
	u.l.InfoContext(ctx,
		"GetAnomaly called",
		slog.String("start", start.String()),
		slog.String("end", end.String()),
	)

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
