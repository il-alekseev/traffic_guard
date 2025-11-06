package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"tg-an/internal/models"
	"tg-an/pkg/slogger/wsl"
	"tg-an/pkg/trparser"
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

	// Получаем данные для первой страницы
	cats, err := u.db.GetCategories(ctx, tr)
	if err != nil {
		err = fmt.Errorf("%s: failed to get categories for main page: %w", method, err)
		u.l.ErrorContext(ctx, "Database operation failed",
			wsl.String("method", method),
			wsl.String("error", err.Error()),
		)
	}
	report.MainActivityPage.TopCategories = cats

	rs, err := u.db.GetResourses(ctx, tr)
	if err != nil {
		err = fmt.Errorf("%s: failed to get resourses for main page: %w", method, err)
		u.l.ErrorContext(ctx, "Database operation failed",
			wsl.String("method", method),
			wsl.String("error", err.Error()),
		)
	}
	report.MainActivityPage.TopResources = rs

	trf, err := u.mdb.GetTrafficStat(ctx, tr, "", 20)
	if err != nil {
		err = fmt.Errorf("%s: failed to get traffic stat for main page: %w", method, err)
		u.l.ErrorContext(ctx, "Database operation failed",
			wsl.String("method", method),
			wsl.String("error", err.Error()),
		)
	}
	report.MainActivityPage.Traffic.Data = trf
	report.MainActivityPage.Traffic.Count = 20

	// Получаем данные для второй страницы
	dv, err := u.db.GetDevicesAnalytics(ctx, tr, "")
	if err != nil {
		err = fmt.Errorf("%s: failed to get devices stat for second page: %w", method, err)
		u.l.ErrorContext(ctx, "Database operation failed",
			wsl.String("method", method),
			wsl.String("error", err.Error()),
		)
	}
	report.DeviceAnalyticsPage.Analytics = dv

	// Получаем данные для третьей страницы
	ans, err := u.db.GetAnomaliesList(ctx, tr, "")
	if err != nil {
		err = fmt.Errorf("%s: failed to get anomalies list for third page: %w", method, err)
		u.l.ErrorContext(ctx, "Database operation failed",
			wsl.String("method", method),
			wsl.String("error", err.Error()),
		)
	}
	report.AnomaliesListPage.Anomalies = ans

	// Получаем данные для четвертой страницы
	anr, err := u.db.GetTopAnomalies(ctx, tr)
	if err != nil {
		err = fmt.Errorf("%s: failed to get top anomalies for forth page: %w", method, err)
		u.l.ErrorContext(ctx, "Database operation failed",
			wsl.String("method", method),
			wsl.String("error", err.Error()),
		)
	}
	report.TopAnomaliesPage.DeviceAnomaly = anr

	// Получаем данные для пятой страницы
	cs, err := u.db.GetTopCategoriesForReport(ctx, tr, "")
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

	// Получаем данные для первой страницы (DeviceAnalyticsPage)
	trf, err := u.mdb.GetTrafficStat(ctx, tr, hostname, 20)
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
	allowed, err := u.db.GetRequestStat(ctx, tr, hostname, string(models.RequestStatusAllowed), 20)
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

	blocked, err := u.db.GetRequestStat(ctx, tr, hostname, string(models.RequestStatusBlocked), 20)
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

	pending, err := u.db.GetRequestStat(ctx, tr, hostname, string(models.RequestStatusPending), 20)
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
	dv, err := u.db.GetDevicesAnalytics(ctx, tr, hostname)
	if err != nil {
		err = fmt.Errorf("%s: failed to get device stat: %w", method, err)
		u.l.ErrorContext(ctx, "Database operation failed",
			wsl.String("method", method),
			wsl.String("error", err.Error()),
		)
	}

	// Ищем данные для конкретного устройства
	for _, device := range dv {
		if device.HostName == hostname {
			report.DeviceAnalyticsPage.AnomalyBlockStat = device.AnomalyBlockStat
			break
		}
	}

	// Получаем данные для второй страницы (AnomaliesListPage)
	ans, err := u.db.GetAnomaliesList(ctx, tr, hostname)
	if err != nil {
		err = fmt.Errorf("%s: failed to get anomalies list: %w", method, err)
		u.l.ErrorContext(ctx, "Database operation failed",
			wsl.String("method", method),
			wsl.String("error", err.Error()),
		)
	}

	// Ищем аномалии для конкретного устройства
	for _, anomaly := range ans {
		if anomaly.HostName == hostname {
			report.AnomaliesListPage.DeviceAnomaly = []models.DeviceAnomalyAnalytics{
				{
					HostName:         anomaly.HostName,
					Traffic:          models.Traffic{},          // TODO: заполнить из данных
					Requests:         0,                         // TODO: заполнить из данных
					AnomalyBlockStat: models.AnomalyBlockStat{}, // TODO: заполнить из данных
					Detections:       models.DetectionReport{},  // TODO: заполнить из данных
				},
			}
			break
		}
	}

	// Получаем данные для третьей страницы (CategoriesPage)
	cs, err := u.db.GetTopCategoriesForReport(ctx, tr, hostname)
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
	)
	return report, nil
}
