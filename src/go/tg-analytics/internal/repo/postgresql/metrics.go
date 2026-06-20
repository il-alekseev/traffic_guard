package postgresql

import (
	"context"
	"fmt"
	"tg-an/internal/models"
	"tg-an/pkg/trparser"
	"time"
)

// GetTrafficStat возвращает статистику по трафику
func (r *RepoPG) GetTrafficStat(ctx context.Context, tr *trparser.TimeRange, hostname string, count uint) (models.TrafficStat, error) {
	result := models.TrafficStat{
		Time:   make([]time.Time, count),
		Input:  make([]uint, count),
		Output: make([]uint, count),
	}

	if !tr.IsValid() {
		return result, fmt.Errorf("invalid time range: %s", tr.String())
	}

	if count == 0 {
		return result, fmt.Errorf("count must be greater than 0")
	}

	totalDuration := tr.Duration()
	intervalDuration := totalDuration / time.Duration(count)

	// Создаем временные интервалы
	for i := uint(0); i < count; i++ {
		start := tr.From.Add(time.Duration(i) * intervalDuration)
		result.Time[i] = start.Add(intervalDuration / 2)
	}

	// Используем сырой SQL запрос с JSON функциями PostgreSQL для ускорения запроса
	// TODO: приспособить под различные сетевые интерфейсы (нужно менять структуру БД метрик)
	query := `
		SELECT 
			FLOOR(EXTRACT(EPOCH FROM (date - $1)) / $2)::integer as interval_index,
			SUM(
				(stat->'statistics'->'network'->'interfaces'->'ge-0-0'->'input'->'bytes'->>'count')::bigint
			) as input_bytes,
			SUM(
				(stat->'statistics'->'network'->'interfaces'->'ge-0-0'->'output'->'bytes'->>'count')::bigint
			) as output_bytes
		FROM stats_json 
		WHERE date BETWEEN $1 AND $3
		AND ($4 = '' OR stat->'statistics'->'common'->>'hostname' = $4)
		GROUP BY interval_index
		ORDER BY interval_index
	`

	intervalSeconds := int(intervalDuration.Seconds())
	var rows []struct {
		IntervalIndex int   `gorm:"column:interval_index"`
		InputBytes    int64 `gorm:"column:input_bytes"`
		OutputBytes   int64 `gorm:"column:output_bytes"`
	}

	err := r.db.GetDB().WithContext(ctx).Raw(query,
		tr.From, intervalSeconds, tr.To, hostname).Find(&rows).Error

	if err != nil {
		return result, fmt.Errorf("failed to execute query: %w", err)
	}

	// Заполняем результат
	for _, row := range rows {
		if row.IntervalIndex >= 0 && row.IntervalIndex < int(count) {
			result.Input[row.IntervalIndex] = uint(row.InputBytes / 1024)
			result.Output[row.IntervalIndex] = uint(row.OutputBytes / 1024)
		}
	}

	return result, nil
}
