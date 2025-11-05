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
	report.MainActivityPage.Traffic = trf

	// Получаем данные для второй страницы
	dv, err := u.db.GetDevicesAnalytics(ctx, tr, "")
	if err != nil {
		err = fmt.Errorf("%s: failed to get devices stat for second page: %w", method, err)
		u.l.ErrorContext(ctx, "Database operation failed",
			wsl.String("method", method),
			wsl.String("error", err.Error()),
		)
	}
	report.DeviceAnalyticsPage = dv

	// Получаем данные для третьей страницы
	ans, err := u.db.GetAnomaliesList(ctx, tr, "")
	if err != nil {
		err = fmt.Errorf("%s: failed to get anomalies list for third page: %w", method, err)
		u.l.ErrorContext(ctx, "Database operation failed",
			wsl.String("method", method),
			wsl.String("error", err.Error()),
		)
	}
	report.AnomaliesListPage = ans

	// Получаем данные для четвертой страницы
	anr, err := u.db.GetTopAnomalies(ctx, tr)
	if err != nil {
		err = fmt.Errorf("%s: failed to get top anomalies for forth page: %w", method, err)
		u.l.ErrorContext(ctx, "Database operation failed",
			wsl.String("method", method),
			wsl.String("error", err.Error()),
		)
	}
	report.TopAnomaliesPage = anr

	// Получаем данные для пятой страницы
	cs, err := u.db.GetTopCategoriesForReport(ctx, tr, "")
	if err != nil {
		err = fmt.Errorf("%s: failed to get top categories for fifth page: %w", method, err)
		u.l.ErrorContext(ctx, "Database operation failed",
			wsl.String("method", method),
			wsl.String("error", err.Error()),
		)
	}
	report.TopCategoriesPage = cs

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

	// Получаем данные для первой страницы
	trf, err := u.mdb.GetTrafficStat(ctx, tr, hostname, 20)
	if err != nil {
		err = fmt.Errorf("%s: failed to get traffic stat for main page: %w", method, err)
		u.l.ErrorContext(ctx, "Database operation failed",
			wsl.String("method", method),
			wsl.String("error", err.Error()),
		)
	}
	report.DeviceAnalyticsPage.Traffic = trf

	allowed, err := u.db.GetRequestStat(ctx, tr, hostname, string(models.RequestStatusAllowed), 20)
	if err != nil {
		err = fmt.Errorf("%s: failed to get allowed requests statistics: %w", method, err)
		u.l.ErrorContext(ctx, "Database operation failed",
			wsl.String("method", method),
			wsl.String("error", err.Error()),
		)
	}
	report.DeviceAnalyticsPage.Allowed = allowed

	blocked, err := u.db.GetRequestStat(ctx, tr, hostname, string(models.RequestStatusBlocked), 20)
	if err != nil {
		err = fmt.Errorf("%s: failed to get blocked requests statistics: %w", method, err)
		u.l.ErrorContext(ctx, "Database operation failed",
			wsl.String("method", method),
			wsl.String("error", err.Error()),
		)
	}
	report.DeviceAnalyticsPage.Blocked = blocked

	pending, err := u.db.GetRequestStat(ctx, tr, hostname, string(models.RequestStatusPending), 20)
	if err != nil {
		err = fmt.Errorf("%s: failed to get pending requests statistics: %w", method, err)
		u.l.ErrorContext(ctx, "Database operation failed",
			wsl.String("method", method),
			wsl.String("error", err.Error()),
		)
	}
	report.DeviceAnalyticsPage.Pending = pending

	dv, err := u.db.GetDevicesAnalytics(ctx, tr, hostname)
	if err != nil {
		err = fmt.Errorf("%s: failed to get device stat: %w", method, err)
		u.l.ErrorContext(ctx, "Database operation failed",
			wsl.String("method", method),
			wsl.String("error", err.Error()),
		)
	}

	// Безопасное извлечение данных для конкретного устройства
	if deviceAnalytics, exists := dv.Analytics[hostname]; exists {
		report.DeviceAnalyticsPage.Anomalies = deviceAnalytics.Anomalies
		report.DeviceAnalyticsPage.Blocks = deviceAnalytics.Blocks
		report.DeviceAnalyticsPage.All = deviceAnalytics.All
	}

	// Получаем данные для второй страницы
	ans, err := u.db.GetAnomaliesList(ctx, tr, hostname)
	if err != nil {
		err = fmt.Errorf("%s: failed to get anomalies list: %w", method, err)
		u.l.ErrorContext(ctx, "Database operation failed",
			wsl.String("method", method),
			wsl.String("error", err.Error()),
		)
	}

	// Безопасное извлечение аномалий для конкретного устройства
	if deviceAnomalies, exists := ans[hostname]; exists {
		report.AnomaliesListPage = deviceAnomalies.Anomalies // Теперь типы совпадают
	}

	// Получаем данные для третьей страницы
	cs, err := u.db.GetTopCategoriesForReport(ctx, tr, hostname)
	if err != nil {
		err = fmt.Errorf("%s: failed to get top categories: %w", method, err)
		u.l.ErrorContext(ctx, "Database operation failed",
			wsl.String("method", method),
			wsl.String("error", err.Error()),
		)
	}

	// Безопасное извлечение категорий для конкретного устройства
	if deviceCategories, exists := cs[hostname]; exists {
		report.CategoriesPage = deviceCategories.Categories
	}

	u.l.InfoContext(ctx, "report for device created",
		slog.String("method", method),
	)
	return report, nil
}
