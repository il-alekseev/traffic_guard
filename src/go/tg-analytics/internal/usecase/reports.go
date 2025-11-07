package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"tg-an/internal/models"
	"tg-an/pkg/slogger/wsl"
	"tg-an/pkg/trparser"
	"time"
)

func (u *Usecase) CreateReport(ctx context.Context, tr *trparser.TimeRange) (models.Report, error) {
	method := "CreateReport"
	u.l.InfoContext(ctx,
		method,
		slog.Any("time_range", tr),
	)

	var report models.Report
	report.From = tr.From
	report.To = tr.To

	// Общее время выполнения
	overallStart := time.Now()
	defer func() {
		totalDuration := time.Since(overallStart)
		u.l.InfoContext(ctx, "Report generation completed",
			slog.String("method", method),
			slog.Duration("total_duration", totalDuration),
		)
	}()

	// Функция для логирования времени выполнения запроса
	logDBTime := func(operation string, start time.Time, err error) {
		duration := time.Since(start)

		// Создаем аргументы для логирования
		args := []any{
			wsl.String("method", method),
			wsl.String("operation", operation),
			slog.Duration("duration", duration),
		}

		// Добавляем ошибку только если она не nil
		if err != nil {
			args = append(args, wsl.Err(err))
			u.l.ErrorContext(ctx, "Database operation failed", args...)
		} else {
			u.l.InfoContext(ctx, "Database operation completed", args...)
		}
	}

	// Получаем данные для первой страницы (MainActivityPage)
	start := time.Now()
	cats, err := u.db.GetCategories(ctx, tr)
	logDBTime("GetCategories", start, err)
	if err != nil {
		err = fmt.Errorf("%s: failed to get categories for main page: %w", method, err)
		u.l.ErrorContext(ctx, "Database operation failed",
			wsl.String("method", method),
			wsl.String("error", err.Error()),
		)
	}
	report.MainActivityPage.TopCategories = cats

	start = time.Now()
	rs, err := u.db.GetResourses(ctx, tr)
	logDBTime("GetResourses", start, err)
	if err != nil {
		err = fmt.Errorf("%s: failed to get resourses for main page: %w", method, err)
		u.l.ErrorContext(ctx, "Database operation failed",
			wsl.String("method", method),
			wsl.String("error", err.Error()),
		)
	}
	report.MainActivityPage.TopResources = rs

	start = time.Now()
	trf, err := u.mdb.GetTrafficStat(ctx, tr, "", 20)
	logDBTime("GetTrafficStat", start, err)
	if err != nil {
		err = fmt.Errorf("%s: failed to get traffic stat for main page: %w", method, err)
		u.l.ErrorContext(ctx, "Database operation failed",
			wsl.String("method", method),
			wsl.String("error", err.Error()),
		)
	}
	report.MainActivityPage.Traffic.Data = trf
	report.MainActivityPage.Traffic.Count = 20

	// Получаем данные для второй страницы (DevicesAnalyticsPage)
	start = time.Now()
	dv, err := u.db.GetDevicesAnalytics(ctx, tr, "")
	logDBTime("GetDevicesAnalytics", start, err)
	if err != nil {
		err = fmt.Errorf("%s: failed to get devices stat for second page: %w", method, err)
		u.l.ErrorContext(ctx, "Database operation failed",
			wsl.String("method", method),
			wsl.String("error", err.Error()),
		)
	}
	report.DeviceAnalyticsPage.Analytics = dv

	// Получаем данные для третьей страницы (AnomaliesListPage)
	start = time.Now()
	ans, err := u.db.GetAnomaliesList(ctx, tr, "")
	logDBTime("GetAnomaliesList", start, err)
	if err != nil {
		err = fmt.Errorf("%s: failed to get anomalies list for third page: %w", method, err)
		u.l.ErrorContext(ctx, "Database operation failed",
			wsl.String("method", method),
			wsl.String("error", err.Error()),
		)
	}
	report.AnomaliesListPage.Anomalies = ans

	// Получаем данные для четвертой страницы (TopAnomaliesPage)
	start = time.Now()
	anr, err := u.db.GetTopAnomalies(ctx, tr)
	logDBTime("GetTopAnomalies", start, err)
	if err != nil {
		err = fmt.Errorf("%s: failed to get top anomalies for forth page: %w", method, err)
		u.l.ErrorContext(ctx, "Database operation failed",
			wsl.String("method", method),
			wsl.String("error", err.Error()),
		)
	}
	report.TopAnomaliesPage.DeviceAnomaly = anr

	// Получаем данные для пятой страницы (TopCategoriesPage)
	start = time.Now()
	cs, err := u.db.GetTopCategoriesForReport(ctx, tr, "")
	logDBTime("GetTopCategoriesForReport", start, err)
	if err != nil {
		err = fmt.Errorf("%s: failed to get top categories for fifth page: %w", method, err)
		u.l.ErrorContext(ctx, "Database operation failed",
			wsl.String("method", method),
			wsl.String("error", err.Error()),
		)
	}
	report.TopCategoriesPage.Categories = cs

	u.l.InfoContext(ctx, "report created",
		slog.String("method", method),
		slog.Int("categories_count", len(cats)),
		slog.Int("resources_count", len(rs)),
		slog.Int("devices_count", len(dv)),
		slog.Int("anomalies_count", len(ans)),
		slog.Int("top_anomalies_count", len(anr)),
		slog.Int("top_categories_count", len(cs)),
	)
	return report, nil
}

func (u *Usecase) CreateReportForDevice(ctx context.Context, tr *trparser.TimeRange, hostname string) (models.ReportForDevice, error) {
	method := "CreateReportForDevice"
	u.l.InfoContext(ctx,
		method,
		slog.Any("time_range", tr),
		wsl.String("hostname", hostname),
	)

	var report models.ReportForDevice
	report.From = tr.From
	report.To = tr.To
	report.HostName = hostname

	// Общее время выполнения
	overallStart := time.Now()
	defer func() {
		totalDuration := time.Since(overallStart)
		u.l.InfoContext(ctx, "Report for device generation completed",
			slog.String("method", method),
			slog.String("hostname", hostname),
			slog.Duration("total_duration", totalDuration),
		)
	}()

	// Функция для логирования времени выполнения запроса
	logDBTime := func(operation string, start time.Time, err error) {
		duration := time.Since(start)

		// Создаем аргументы для логирования
		args := []any{
			wsl.String("method", method),
			wsl.String("operation", operation),
			slog.Duration("duration", duration),
		}

		// Добавляем ошибку только если она не nil
		if err != nil {
			args = append(args, wsl.Err(err))
			u.l.ErrorContext(ctx, "Database operation failed", args...)
		} else {
			u.l.InfoContext(ctx, "Database operation completed", args...)
		}
	}

	// Получаем данные для первой страницы (DeviceAnalyticsPage)
	start := time.Now()
	trf, err := u.mdb.GetTrafficStat(ctx, tr, hostname, 20)
	logDBTime("GetTrafficStat", start, err)
	if err != nil {
		err = fmt.Errorf("%s: failed to get traffic stat for main page: %w", method, err)
		u.l.ErrorContext(ctx, "Database operation failed",
			wsl.String("method", method),
			wsl.String("error", err.Error()),
		)
	}
	report.DeviceAnalyticsPage.Traffic.Data = trf
	report.DeviceAnalyticsPage.Traffic.Count = 20

	// Получаем статистику запросов
	start = time.Now()
	allowed, err := u.db.GetRequestStat(ctx, tr, hostname, string(models.RequestStatusAllowed), 20)
	logDBTime("GetRequestStat (allowed)", start, err)
	if err != nil {
		err = fmt.Errorf("%s: failed to get allowed requests statistics: %w", method, err)
		u.l.ErrorContext(ctx, "Database operation failed",
			wsl.String("method", method),
			wsl.String("error", err.Error()),
		)
	}
	report.DeviceAnalyticsPage.RequestsAnalytics.Allowed = models.RequestStatData{
		Time:  allowed.Time,
		Data:  allowed.Data,
		Count: 20,
	}

	start = time.Now()
	blocked, err := u.db.GetRequestStat(ctx, tr, hostname, string(models.RequestStatusBlocked), 20)
	logDBTime("GetRequestStat (blocked)", start, err)
	if err != nil {
		err = fmt.Errorf("%s: failed to get blocked requests statistics: %w", method, err)
		u.l.ErrorContext(ctx, "Database operation failed",
			wsl.String("method", method),
			wsl.String("error", err.Error()),
		)
	}
	report.DeviceAnalyticsPage.RequestsAnalytics.Blocked = models.RequestStatData{
		Time:  blocked.Time,
		Data:  blocked.Data,
		Count: 20,
	}

	start = time.Now()
	pending, err := u.db.GetRequestStat(ctx, tr, hostname, string(models.RequestStatusPending), 20)
	logDBTime("GetRequestStat (pending)", start, err)
	if err != nil {
		err = fmt.Errorf("%s: failed to get pending requests statistics: %w", method, err)
		u.l.ErrorContext(ctx, "Database operation failed",
			wsl.String("method", method),
			wsl.String("error", err.Error()),
		)
	}
	report.DeviceAnalyticsPage.RequestsAnalytics.Pending = models.RequestStatData{
		Time:  pending.Time,
		Data:  pending.Data,
		Count: 20,
	}

	// Получаем статистику аномалий и блокировок
	start = time.Now()
	dv, err := u.db.GetDevicesAnalytics(ctx, tr, hostname)
	logDBTime("GetDevicesAnalytics", start, err)
	if err != nil {
		err = fmt.Errorf("%s: failed to get device stat: %w", method, err)
		u.l.ErrorContext(ctx, "Database operation failed",
			wsl.String("method", method),
			wsl.String("error", err.Error()),
		)
	}

	// Ищем данные для конкретного устройства
	var deviceFound bool
	for _, device := range dv {
		if device.HostName == hostname {
			report.DeviceAnalyticsPage.AnomalyBlockStat = device.AnomalyBlockStat
			deviceFound = true
			break
		}
	}

	// Если устройство не найдено, создаем пустую статистику
	if !deviceFound {
		report.DeviceAnalyticsPage.AnomalyBlockStat = models.AnomalyBlockStat{}
	}

	// Получаем данные для второй страницы (TopAnomaliesPage)
	start = time.Now()
	topAnomalies, err := u.db.GetTopAnomalies(ctx, tr)
	logDBTime("GetTopAnomalies", start, err)
	if err != nil {
		err = fmt.Errorf("%s: failed to get top anomalies: %w", method, err)
		u.l.ErrorContext(ctx, "Database operation failed",
			wsl.String("method", method),
			wsl.String("error", err.Error()),
		)
	}

	// Фильтруем топ аномалии для текущего устройства
	var deviceTopAnomalies []models.DeviceAnomalyAnalytics
	for _, anomaly := range topAnomalies {
		if anomaly.HostName == hostname {
			deviceTopAnomalies = append(deviceTopAnomalies, anomaly)
		}
	}
	report.AnomaliesListPage.DeviceAnomaly = deviceTopAnomalies

	// Получаем данные для третьей страницы (CategoriesPage)
	start = time.Now()
	cs, err := u.db.GetTopCategoriesForReport(ctx, tr, hostname)
	logDBTime("GetTopCategoriesForReport", start, err)
	if err != nil {
		err = fmt.Errorf("%s: failed to get top categories: %w", method, err)
		u.l.ErrorContext(ctx, "Database operation failed",
			wsl.String("method", method),
			wsl.String("error", err.Error()),
		)
	}
	report.CategoriesPage.Categories = cs

	u.l.InfoContext(ctx, "report for device created",
		slog.String("method", method),
		slog.String("hostname", hostname),
		slog.Int("categories_count", len(cs)),
		slog.Int("anomalies_count", len(deviceTopAnomalies)),
	)
	return report, nil
}
