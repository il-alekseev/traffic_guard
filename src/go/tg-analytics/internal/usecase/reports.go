package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"tg-an/internal/controllers/http/v1/values"
	"tg-an/internal/models"
	"tg-an/pkg/blog/operations"
	pkg "tg-an/pkg/models"
	"tg-an/pkg/slogger/wsl"
	"tg-an/pkg/trparser"
	"time"

	"github.com/google/uuid"
)

func (u *Usecase) CreateReport(ctx context.Context, userMeta *models.UserMeta, tr *trparser.TimeRange) (models.Report, error) {
	method := "CreateReport"
	u.l.InfoContext(ctx,
		method,
		slog.Any("time_range", tr),
	)

	// Инициализируем отчет с пустыми структурами вместо nil
	report := models.Report{
		From: tr.From,
		To:   tr.To,
		MainActivityPage: models.MainActivityPage{
			TopCategories: []models.CategoryStat{},
			TopResources:  []models.ResourceStat{},
			Traffic: models.TrafficStatData{
				Data:  models.TrafficStat{},
				Count: 20,
			},
		},
		DeviceAnalyticsPage: models.DevicesAnalyticsPage{
			Analytics: []models.DeviceReport{},
		},
		AnomaliesListPage: models.DevicesAnomaliesListPage{
			Anomalies: []models.DeviceAnomaly{},
		},
		TopAnomaliesPage: models.TopAnomaliesPage{
			DeviceAnomaly: []models.DeviceAnomalyAnalytics{},
		},
		TopCategoriesPage: models.TopCategoriesPage{
			Categories: []models.TopCategory{},
		},
	}

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
	} else {
		report.MainActivityPage.TopCategories = cats
	}

	start = time.Now()
	rs, err := u.db.GetResourses(ctx, tr)
	logDBTime("GetResourses", start, err)
	if err != nil {
		err = fmt.Errorf("%s: failed to get resourses for main page: %w", method, err)
		u.l.ErrorContext(ctx, "Database operation failed",
			wsl.String("method", method),
			wsl.String("error", err.Error()),
		)
	} else {
		report.MainActivityPage.TopResources = rs
	}

	start = time.Now()
	trf, err := u.mdb.GetTrafficStat(ctx, tr, "", 20)
	logDBTime("GetTrafficStat", start, err)
	if err != nil {
		err = fmt.Errorf("%s: failed to get traffic stat for main page: %w", method, err)
		u.l.ErrorContext(ctx, "Database operation failed",
			wsl.String("method", method),
			wsl.String("error", err.Error()),
		)
	} else {
		report.MainActivityPage.Traffic.Data = trf
		report.MainActivityPage.Traffic.Count = 20
	}

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
	} else {
		report.DeviceAnalyticsPage.Analytics = dv
	}

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
	} else {
		report.AnomaliesListPage.Anomalies = ans
	}

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
	} else {
		report.TopAnomaliesPage.DeviceAnomaly = anr
	}

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
	} else {
		report.TopCategoriesPage.Categories = cs
	}

	// Запись события в бизнес-лог
	uuidStr := uuid.New().String()

	newValue := make(map[string]any)
	newValue["from"] = tr.From
	newValue["to"] = tr.To

	record := pkg.DtoBusinessLog{
		Description: "Создание сводного отчета по всем устройствам",
		Entity:      "Report",
		EntityID:    "",
		NewValue:    newValue,
		EventType:   "CREATE",
		Context:     userMeta.ContextID,
		UserName:    userMeta.Username,
		UserRole:    userMeta.ClientRole,
	}

	_, err = u.blclient.Operations.PostAPIV1Add(&operations.PostAPIV1AddParams{
		XCallerService: values.ThisServiceName,
		XRequestID:     uuidStr,
		Record:         &record,
		Context:        ctx,
	},
	)
	if err != nil {
		err = fmt.Errorf("%s: failed to blog event: %w", method, err)
		u.l.ErrorContext(ctx, "failed to blog event",
			wsl.String("method", method),
			wsl.String("error", err.Error()),
		)
		return report, nil
	}

	u.l.InfoContext(ctx, "report created",
		slog.String("method", method),
		slog.String("blog uuid", uuidStr),
		slog.Int("categories_count", len(report.MainActivityPage.TopCategories)),
		slog.Int("resources_count", len(report.MainActivityPage.TopResources)),
		slog.Int("devices_count", len(report.DeviceAnalyticsPage.Analytics)),
		slog.Int("anomalies_count", len(report.AnomaliesListPage.Anomalies)),
		slog.Int("top_anomalies_count", len(report.TopAnomaliesPage.DeviceAnomaly)),
		slog.Int("top_categories_count", len(report.TopCategoriesPage.Categories)),
	)
	return report, nil
}

func (u *Usecase) CreateReportForDevice(ctx context.Context, userMeta *models.UserMeta, tr *trparser.TimeRange, hostname string) (models.ReportForDevice, error) {
	method := "CreateReportForDevice"
	u.l.InfoContext(ctx,
		method,
		slog.Any("time_range", tr),
		wsl.String("hostname", hostname),
	)

	// Инициализируем отчет с пустыми структурами вместо nil
	report := models.ReportForDevice{
		From:     tr.From,
		To:       tr.To,
		HostName: hostname,
		DeviceAnalyticsPage: models.DeviceAnalyticsPage{
			Traffic: models.TrafficStatData{
				Data:  models.TrafficStat{},
				Count: 20,
			},
			RequestsAnalytics: models.RequestsAnalytics{
				Allowed: models.RequestStatData{
					Time:  []time.Time{},
					Data:  []uint{},
					Count: 20,
				},
				Blocked: models.RequestStatData{
					Time:  []time.Time{},
					Data:  []uint{},
					Count: 20,
				},
				Pending: models.RequestStatData{
					Time:  []time.Time{},
					Data:  []uint{},
					Count: 20,
				},
			},
			AnomalyBlockStat: models.AnomalyBlockStat{},
		},
		AnomaliesListPage: models.TopAnomaliesPage{
			DeviceAnomaly: []models.DeviceAnomalyAnalytics{},
		},
		CategoriesPage: models.TopCategoriesPage{
			Categories: []models.TopCategory{},
		},
	}

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
	} else {
		report.DeviceAnalyticsPage.Traffic.Data = trf
		report.DeviceAnalyticsPage.Traffic.Count = 20
	}

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
	} else {
		report.DeviceAnalyticsPage.RequestsAnalytics.Allowed = models.RequestStatData{
			Time:  allowed.Time,
			Data:  allowed.Data,
			Count: 20,
		}
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
	} else {
		report.DeviceAnalyticsPage.RequestsAnalytics.Blocked = models.RequestStatData{
			Time:  blocked.Time,
			Data:  blocked.Data,
			Count: 20,
		}
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
	} else {
		report.DeviceAnalyticsPage.RequestsAnalytics.Pending = models.RequestStatData{
			Time:  pending.Time,
			Data:  pending.Data,
			Count: 20,
		}
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
	} else {
		// Ищем данные для конкретного устройства
		for _, device := range dv {
			if device.HostName == hostname {
				report.DeviceAnalyticsPage.AnomalyBlockStat = device.AnomalyBlockStat
				break
			}
		}
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
	} else {
		// Фильтруем топ аномалии для текущего устройства
		deviceTopAnomalies := make([]models.DeviceAnomalyAnalytics, 0)
		for _, anomaly := range topAnomalies {
			if anomaly.HostName == hostname {
				deviceTopAnomalies = append(deviceTopAnomalies, anomaly)
			}
		}
		report.AnomaliesListPage.DeviceAnomaly = deviceTopAnomalies
	}

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
	} else {
		report.CategoriesPage.Categories = cs
	}

	// Запись события в бизнес-лог
	uuidStr := uuid.New().String()

	newValue := make(map[string]any)
	newValue["from"] = tr.From
	newValue["to"] = tr.To
	newValue["hostname"] = hostname

	record := pkg.DtoBusinessLog{
		Description: "Создание отчета по устройству",
		Entity:      "Report",
		EntityID:    "",
		NewValue:    newValue,
		EventType:   "CREATE",
		Context:     userMeta.ContextID,
		UserName:    userMeta.Username,
		UserRole:    userMeta.ClientRole,
	}

	_, err = u.blclient.Operations.PostAPIV1Add(&operations.PostAPIV1AddParams{
		XCallerService: values.ThisServiceName,
		XRequestID:     uuidStr,
		Record:         &record,
		Context:        ctx,
	},
	)
	if err != nil {
		err = fmt.Errorf("%s: failed to blog event: %w", method, err)
		u.l.ErrorContext(ctx, "failed to blog event",
			wsl.String("method", method),
			wsl.String("blog uuid", uuidStr),
			wsl.String("error", err.Error()),
		)
		return report, nil
	}

	u.l.InfoContext(ctx, "report for device created",
		slog.String("method", method),
		slog.String("hostname", hostname),
		slog.Int("categories_count", len(report.CategoriesPage.Categories)),
		slog.Int("anomalies_count", len(report.AnomaliesListPage.DeviceAnomaly)),
	)
	return report, nil
}
