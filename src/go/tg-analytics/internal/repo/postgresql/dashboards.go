package postgresql

import (
	"context"
	"fmt"
	"tg-an/internal/controllers/http/v1/dto"
	"tg-an/internal/models"
	"tg-an/internal/pkg/status"
	pkg "tg-an/pkg/models"
	"tg-an/pkg/slogger/wsl"
	"tg-an/pkg/trparser"
	"time"
)

// GetTopCategories возвращает топ категорий по количеству доступов
func (r *RepoPG) GetTopCategories(ctx context.Context, tr *trparser.TimeRange, f models.CategoryFilter, count int) ([]dto.Category, error) {
	var categories []dto.Category

	query := r.db.GetDB().WithContext(ctx).Table("sessions").
		Select("categories.name as name, COUNT(*) as access_count").
		Joins("LEFT JOIN domains ON sessions.domain_id = domains.id").
		Joins("LEFT JOIN categories ON domains.category_id = categories.id").
		Where("categories.name IS NOT NULL AND categories.name != 'Неизвестный класс'")

	// Применяем временной диапазон
	if tr != nil && !tr.From.IsZero() && !tr.To.IsZero() {
		query = query.Where("sessions.datetime_utc BETWEEN ? AND ?", tr.From, tr.To)
	}

	// Применяем фильтр по hostname
	if f.HostName != "" {
		query = query.Joins("JOIN devices ON sessions.device_id = devices.id").
			Where("devices.hostname = ?", f.HostName)
	}

	// Применяем фильтр по типу сессии
	if f.Type != "" {
		query = query.Where("sessions.type = ?", f.Type)
	}

	err := query.
		Group("categories.name").
		Order("access_count DESC").
		Limit(count).
		Find(&categories).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get top categories: %w", err)
	}

	return categories, nil
}

// GetRequestsStat возвращает статистику запросов по временным интервалам
func (r *RepoPG) GetRequestStat(ctx context.Context, tr *trparser.TimeRange, hostname, requestType string, count uint) (models.RequestStat, error) {
	result := models.RequestStat{
		Time: make([]time.Time, count),
		Data: make([]uint, count),
	}

	if !tr.IsValid() {
		return result, fmt.Errorf("invalid time range: %s", tr.String())
	}

	if count == 0 {
		return result, fmt.Errorf("count must be greater than 0")
	}

	totalDuration := tr.Duration()
	intervalDuration := totalDuration / time.Duration(count)

	// Создаем временные интервалы заранее
	for i := uint(0); i < count; i++ {
		start := tr.From.Add(time.Duration(i) * intervalDuration)
		result.Time[i] = start.Add(intervalDuration / 2)
	}

	// Получаем все записи используя GORM
	var statRecords []models.Session
	query := r.db.GetDB().WithContext(ctx).Table("sessions").
		Where("sessions.datetime_utc BETWEEN ? AND ?", tr.From, tr.To).
		Joins("LEFT JOIN domains ON sessions.domain_id = domains.id").
		Joins("LEFT JOIN categories ON domains.category_id = categories.id").
		Joins("LEFT JOIN actions ON domains.action_id = actions.id")

	// Применяем фильтр по hostname
	if hostname != "" {
		query = query.Joins("JOIN devices ON sessions.device_id = devices.id").
			Where("devices.hostname = ?", hostname)
	}

	rt, err := models.ParseRequestStatus(requestType)
	if err != nil {
		return result, fmt.Errorf("failed to parse request type: %w", err)
	}

	if rt == models.RequestStatusBeforeBlock {
		query = query.Where("sessions.status = ? or sessions.status = ?", models.RequestStatusPending.String(), status.StatusAnomaly.String())
	} else {
		query = query.Where("sessions.status = ?", rt.String())
	}
	if err := query.Find(&statRecords).Error; err != nil {
		return result, fmt.Errorf("failed to get request stat: %w", err)
	}

	// Обрабатываем каждую запись и распределяем по интервалам
	for _, record := range statRecords {
		// Определяем индекс интервала
		timeDiff := record.DatetimeUTC.Sub(tr.From)
		intervalIndex := int(timeDiff / intervalDuration)

		if intervalIndex >= 0 && intervalIndex < int(count) {
			result.Data[intervalIndex]++
		}
	}

	return result, nil
}

// GetTopUnresolvedDetections возвращает топ нерешенных выявлений по числу запросов
func (r *RepoPG) GetTopUnresolvedDetections(ctx context.Context, tr *trparser.TimeRange, hostName string, count int) ([]dto.UnresolvedDetection, error) {
	var detections []dto.UnresolvedDetection

	query := r.db.GetDB().WithContext(ctx).Table("sessions").
		Select(`
			domains.path as domain,
			COUNT(*) as requests_all,
			COUNT(CASE WHEN sessions.datetime_utc < domains.categorized_at THEN 1 END) as requests_before,
			COUNT(CASE WHEN sessions.datetime_utc >= domains.categorized_at THEN 1 END) as requests_after
		`).
		Joins("LEFT JOIN domains ON sessions.domain_id = domains.id").
		Joins("LEFT JOIN categories ON domains.category_id = categories.id").
		Where("categories.type = ?", pkg.CategoryTypeNegative).       // признак выявления
		Where("(domains.action_id IS NULL OR domains.action_id = 0)") // признак того, что выявление нерешенное

	// Применяем временной диапазон
	if tr != nil && !tr.To.IsZero() {
		query = query.Where("sessions.datetime_utc BETWEEN ? AND ?", tr.From, tr.To)
	}

	// Применяем фильтр по hostname
	if hostName != "" {
		query = query.Joins("JOIN devices ON sessions.device_id = devices.id").
			Where("devices.hostname = ?", hostName)
	}

	err := query.
		Group("domains.path").
		Order("requests_after DESC, requests_before DESC").
		Limit(count).
		Find(&detections).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get top unresolved detections: %w", err)
	}

	return detections, nil
}

func (r *RepoPG) GetDeviceStat(ctx context.Context, tr *trparser.TimeRange, hostname string, count uint) (dto.DeviceStatResponse, error) {
	result := dto.DeviceStatResponse{
		Time:  make([]time.Time, count),
		Data:  []models.DeviceRequestStat{},
		Count: count,
	}

	if !tr.IsValid() {
		return result, fmt.Errorf("invalid time range: %s", tr.String())
	}

	if count == 0 {
		return result, fmt.Errorf("count must be greater than 0")
	}

	totalDuration := tr.Duration()
	intervalDuration := totalDuration / time.Duration(count)

	// Создаем временные интервалы заранее
	for i := uint(0); i < count; i++ {
		start := tr.From.Add(time.Duration(i) * intervalDuration)
		result.Time[i] = start.Add(intervalDuration / 2)
	}

	// Получаем список хостов
	var hostnames []string
	var err error

	if hostname == "" {
		hostnames, err = r.GetDevices(ctx, hostname)
		if err != nil {
			return result, fmt.Errorf("failed to get devices: %w", err)
		}
	} else {
		hostnames = []string{hostname}
	}

	// Инициализируем структуры данных для всех хостов
	for _, h := range hostnames {
		// TODO: добавить определение текущего статуса сетевого узла
		deviceStat := models.DeviceRequestStat{
			HostName: h,
			Status:   "unknown", // Заглушка, нужно реализовать определение статуса
			Blocked:  make([]uint, count),
			Pending:  make([]uint, count),
		}
		result.Data = append(result.Data, deviceStat)
	}

	var sessionRecords []dto.Session
	if err := r.db.GetDB().WithContext(ctx).Table("sessions").
		Select(`
			sessions.id,
			sessions.datetime_utc,
			sessions.type,
			sessions.status,
			urls.path,
			urls.proto,
			devices.hostname as host_name,
			sources.ip as src_ip,
			sources.country as src_country,
			sources.username,
			domains.ip as dst_ip,
			domains.port as dst_port,
			domains.country as dst_country,
			categories.name as category
		`).
		Joins("LEFT JOIN devices ON sessions.device_id = devices.id").
		Joins("LEFT JOIN sources ON sessions.src_id = sources.id").
		Joins("LEFT JOIN domains ON sessions.domain_id = domains.id").
		Joins("LEFT JOIN urls ON domains.id = urls.domain_id").
		Joins("LEFT JOIN categories ON domains.category_id = categories.id").
		Where("sessions.datetime_utc BETWEEN ? AND ?", tr.From, tr.To).
		Find(&sessionRecords).Error; err != nil {
		return result, fmt.Errorf("failed to get sessions: %w", err)
	}

	// Создаем карту для быстрого доступа к данным по hostname
	deviceMap := make(map[string]*models.DeviceRequestStat)
	for i := range result.Data {
		deviceMap[result.Data[i].HostName] = &result.Data[i]
	}

	for _, record := range sessionRecords {
		// Находим устройство в результатах
		deviceStat, exists := deviceMap[record.HostName]
		if !exists {
			// Если устройства нет в списке, пропускаем запись
			continue
		}

		// Определяем индекс интервала
		timeDiff := record.DatetimeUTC.Sub(tr.From)
		intervalIndex := int(timeDiff / intervalDuration)

		if intervalIndex >= 0 && intervalIndex < int(count) {
			if record.Status == status.StatusBlocked.String() || record.Status == status.StatusAnomaly.String() {
				deviceStat.Blocked[intervalIndex]++
			} else if record.Status == status.StatusPending.String() {
				deviceStat.Pending[intervalIndex]++
			}
		}
	}

	return result, nil
}

func (r *RepoPG) GetAnomalies(ctx context.Context, tr *trparser.TimeRange, hostname string) (dto.GetAnomaliesResponse, error) {
	var result = dto.GetAnomaliesResponse{
		HostAnomalies: make([]models.HostAnomalies, 0),
	}

	// Получаем список хостов
	hostnames, err := r.GetDevices(ctx, hostname)
	if err != nil {
		return result, fmt.Errorf("failed to get devices: %w", err)
	}

	// Если указан конкретный hostname, проверяем его существование
	if hostname != "" {
		found := false
		for _, h := range hostnames {
			if h == hostname {
				found = true
				break
			}
		}
		if !found {
			return result, fmt.Errorf("hostname '%s' not found", hostname)
		}
		hostnames = []string{hostname} // Работаем только с указанным хостом
	}

	result.HostCount = uint(len(hostnames))
	// Инициализируем списки хостов
	for _, h := range hostnames {
		result.HostAnomalies = append(result.HostAnomalies, models.HostAnomalies{HostName: h, Domains: []string{}})
	}

	// Получаем число заблокированных ресурсов (доменов)
	var blockedDomainsCount int64
	if err := r.db.GetDB().WithContext(ctx).Table("sessions").
		Select(`
			COUNT(DISTINCT CASE WHEN actions.action = ? THEN domains.id END) as denied
		`, pkg.ActionTypeDenied.String()).
		Joins("LEFT JOIN devices ON sessions.device_id = devices.id").
		Joins("LEFT JOIN domains ON sessions.domain_id = domains.id").
		Joins("LEFT JOIN actions ON domains.action_id = actions.id").
		Joins("LEFT JOIN categories ON domains.category_id = categories.id").
		Where("categories.type = ?", pkg.CategoryTypeNegative.String()).
		Scan(&blockedDomainsCount).Error; err != nil {
		return result, fmt.Errorf("failed to get blocked domains count: %w", err)
	}

	result.BlockCount = uint(blockedDomainsCount)

	// Получаем аномалии для каждого хоста
	for i, host := range hostnames {
		var anomalyDomains []string

		// Получаем аномальные домены для конкретного хоста
		err := r.db.GetDB().WithContext(ctx).Table("sessions").
			Distinct("domains.path").
			Select("domains.path").
			Joins("LEFT JOIN domains ON sessions.domain_id = domains.id").
			Joins("LEFT JOIN devices ON sessions.device_id = devices.id").
			Where("sessions.status = ?", status.StatusAnomaly.String()).
			Where("devices.hostname = ?", host).
			Where("sessions.datetime_utc BETWEEN ? AND ?", tr.From, tr.To).
			Pluck("domains.path", &anomalyDomains).Error

		if err != nil {
			r.l.WarnContext(ctx, "warning with get domains", wsl.Err(err))
			continue
		}

		// Добавляем в result
		result.HostAnomalies[i].AnomalyCount = uint(len(anomalyDomains))
		result.HostAnomalies[i].Domains = anomalyDomains
	}

	return result, nil
}

// GetProhActivity возвращает статистику запрещенной активности по дням (оптимизированная версия)
func (r *RepoPG) GetProhActivity(ctx context.Context, tr *trparser.TimeRange, hostname string) (dto.GetProhActivityResponse, error) {
	// Проверяем временной диапазон
	if !tr.IsValid() {
		return dto.GetProhActivityResponse{}, fmt.Errorf("invalid time range: %s", tr.String())
	}

	// Вычисляем количество дней в диапазоне
	daysCount := calculateDaysCount(tr.From, tr.To)
	if daysCount == 0 {
		daysCount = 1
	}

	// Формируем структуру ответа
	result := dto.GetProhActivityResponse{
		TimeSince: tr.From,
		Count:     uint(daysCount),
		Data:      make([]uint, daysCount),
	}

	// Создаем карту для агрегации данных
	dayData := make(map[string]uint)

	// Используем группировку по дням в базе данных для большей эффективности
	type DayCount struct {
		Date  string `gorm:"column:date"`
		Count uint   `gorm:"column:count"`
	}

	var dayCounts []DayCount

	// Создаем запрос с группировкой по дням
	query := r.db.GetDB().WithContext(ctx).Table("sessions").
		Select(`
			DATE(sessions.datetime_utc) as date,
			COUNT(*) as count
		`).
		Where("sessions.datetime_utc BETWEEN ? AND ?", tr.From, tr.To).
		Where("(sessions.status = ? OR sessions.status = ?)",
			status.StatusAnomaly.String(),
			status.StatusBlocked.String())

	// Применяем фильтр по hostname, если указан
	if hostname != "" {
		query = query.Joins("JOIN devices ON sessions.device_id = devices.id").
			Where("devices.hostname = ?", hostname)
	}

	// Выполняем запрос с группировкой
	if err := query.Group("DATE(sessions.datetime_utc)").
		Order("DATE(sessions.datetime_utc)").
		Find(&dayCounts).Error; err != nil {
		return result, fmt.Errorf("failed to get prohibited activity: %w", err)
	}

	// Заполняем карту данными из БД
	for _, dc := range dayCounts {
		dayData[dc.Date] = dc.Count
	}

	// Заполняем массив Data в соответствии с днями
	currentDate := tr.From
	for i := 0; i < daysCount; i++ {
		dateStr := currentDate.Format(time.RFC3339)
		if count, exists := dayData[dateStr]; exists {
			result.Data[i] = count
		} else {
			result.Data[i] = 0
		}
		currentDate = currentDate.Add(24 * time.Hour)
	}

	return result, nil
}

// calculateDaysCount вычисляет количество дней в временном диапазоне
func calculateDaysCount(from, to time.Time) int {
	// Нормализуем время до начала дня
	fromDate := time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, from.Location())
	toDate := time.Date(to.Year(), to.Month(), to.Day(), 0, 0, 0, 0, to.Location())

	days := int(toDate.Sub(fromDate).Hours()/24) + 1 // +1 чтобы включить оба дня
	if days < 1 {
		return 1
	}
	return days
}
