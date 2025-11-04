package trparser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// TimeRangeParser парсит строковые представления временных диапазонов
type TimeRangeParser struct{}

// Parse парсит строку в TimeRange
func (p *TimeRangeParser) Parse(timeRangeStr string, now time.Time) (*TimeRange, error) {
	if timeRangeStr == "" {
		return nil, fmt.Errorf("empty time range")
	}

	// Обработка формата "from=X&to=Y"
	if strings.Contains(timeRangeStr, "from=") && strings.Contains(timeRangeStr, "to=") {
		return p.parseQueryStringFormat(timeRangeStr, now)
	}

	// Обработка относительных временных диапазонов (например, "now-1h", "now-24h")
	if strings.HasPrefix(timeRangeStr, "now") {
		return p.parseRelativeTimeRange(timeRangeStr, now)
	}

	// Обработка абсолютных временных диапазонов
	return p.parseAbsoluteTimeRange(timeRangeStr)
}

// parseQueryStringFormat парсит строку в формате "from=X&to=Y"
func (p *TimeRangeParser) parseQueryStringFormat(queryStr string, now time.Time) (*TimeRange, error) {
	var fromStr, toStr string

	// Парсим параметры from и to
	params := strings.Split(queryStr, "&")
	for _, param := range params {
		if strings.HasPrefix(param, "from=") {
			fromStr = strings.TrimPrefix(param, "from=")
		} else if strings.HasPrefix(param, "to=") {
			toStr = strings.TrimPrefix(param, "to=")
		}
	}

	if fromStr == "" || toStr == "" {
		return nil, fmt.Errorf("missing 'from' or 'to' parameter in query string: %s", queryStr)
	}

	// Парсим время начала (from)
	var fromTime time.Time
	if fromStr == "now" {
		fromTime = now
	} else if strings.HasPrefix(fromStr, "now-") {
		durationStr := strings.TrimPrefix(fromStr, "now-")
		duration, err := p.parseDuration(durationStr)
		if err != nil {
			return nil, fmt.Errorf("invalid 'from' duration: %v", err)
		}
		fromTime = now.Add(-duration)
	} else {
		// Пытаемся парсить как абсолютное время
		parsedTime, err := time.Parse(time.RFC3339, fromStr)
		if err != nil {
			return nil, fmt.Errorf("invalid 'from' time format: %v", err)
		}
		fromTime = parsedTime
	}

	// Парсим время окончания (to)
	var toTime time.Time
	if toStr == "now" {
		toTime = now
	} else if strings.HasPrefix(toStr, "now-") {
		durationStr := strings.TrimPrefix(toStr, "now-")
		duration, err := p.parseDuration(durationStr)
		if err != nil {
			return nil, fmt.Errorf("invalid 'to' duration: %v", err)
		}
		toTime = now.Add(-duration)
	} else {
		// Пытаемся парсить как абсолютное время
		parsedTime, err := time.Parse(time.RFC3339, toStr)
		if err != nil {
			return nil, fmt.Errorf("invalid 'to' time format: %v", err)
		}
		toTime = parsedTime
	}

	// Проверяем, что from < to
	if !fromTime.Before(toTime) {
		return nil, fmt.Errorf("'from' time must be before 'to' time: from=%s, to=%s",
			fromTime.Format(time.RFC3339), toTime.Format(time.RFC3339))
	}

	return &TimeRange{From: fromTime, To: toTime}, nil
}

// parseRelativeTimeRange парсит относительные временные диапазоны
func (p *TimeRangeParser) parseRelativeTimeRange(timeRangeStr string, now time.Time) (*TimeRange, error) {
	// Обработка формата "now-X" (только начало указано)
	if strings.Contains(timeRangeStr, "-") && !strings.Contains(timeRangeStr, " to ") {
		parts := strings.Split(timeRangeStr, "-")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid relative time range format: %s", timeRangeStr)
		}

		durationStr := strings.TrimSpace(parts[1])
		duration, err := p.parseDuration(durationStr)
		if err != nil {
			return nil, err
		}

		from := now.Add(-duration)
		return &TimeRange{From: from, To: now}, nil
	}

	// Обработка формата "now-X to now-Y"
	parts := strings.Split(timeRangeStr, " to ")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid relative time range format: %s", timeRangeStr)
	}

	fromPart := strings.TrimSpace(parts[0])
	toPart := strings.TrimSpace(parts[1])

	// Парсим время начала
	var from time.Time
	if fromPart == "now" {
		from = now
	} else if strings.HasPrefix(fromPart, "now-") {
		durationStr := strings.TrimPrefix(fromPart, "now-")
		duration, err := p.parseDuration(durationStr)
		if err != nil {
			return nil, fmt.Errorf("invalid from duration: %v", err)
		}
		from = now.Add(-duration)
	} else {
		return nil, fmt.Errorf("invalid from time format: %s", fromPart)
	}

	// Парсим время окончания
	var to time.Time
	if toPart == "now" {
		to = now
	} else if strings.HasPrefix(toPart, "now-") {
		durationStr := strings.TrimPrefix(toPart, "now-")
		duration, err := p.parseDuration(durationStr)
		if err != nil {
			return nil, fmt.Errorf("invalid to duration: %v", err)
		}
		to = now.Add(-duration)
	} else {
		return nil, fmt.Errorf("invalid to time format: %s", toPart)
	}

	// Проверяем валидность диапазона
	if !from.Before(to) {
		return nil, fmt.Errorf("from time must be before to time")
	}

	return &TimeRange{From: from, To: to}, nil
}

// parseDuration парсит строку продолжительности
func (p *TimeRangeParser) parseDuration(durationStr string) (time.Duration, error) {
	// Регулярное выражение для парсинга продолжительности
	re := regexp.MustCompile(`^(\d+)([smhdwMy])$`)
	matches := re.FindStringSubmatch(durationStr)

	if matches == nil {
		return 0, fmt.Errorf("invalid duration format: %s", durationStr)
	}

	value, _ := strconv.Atoi(matches[1])
	unit := matches[2]

	switch unit {
	case "s": // seconds
		return time.Duration(value) * time.Second, nil
	case "m": // minutes
		return time.Duration(value) * time.Minute, nil
	case "h": // hours
		return time.Duration(value) * time.Hour, nil
	case "d": // days
		return time.Duration(value) * 24 * time.Hour, nil
	case "w": // weeks
		return time.Duration(value) * 7 * 24 * time.Hour, nil
	case "M": // months (приблизительно)
		return time.Duration(value) * 30 * 24 * time.Hour, nil
	case "y": // years (приблизительно)
		return time.Duration(value) * 365 * 24 * time.Hour, nil
	default:
		return 0, fmt.Errorf("unknown duration unit: %s", unit)
	}
}

// parseAbsoluteTimeRange парсит абсолютные временные диапазоны
func (p *TimeRangeParser) parseAbsoluteTimeRange(timeRangeStr string) (*TimeRange, error) {
	// Формат: "2024-01-01T00:00:00Z to 2024-01-02T00:00:00Z"
	parts := strings.Split(timeRangeStr, " to ")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid absolute time range format: %s", timeRangeStr)
	}

	from, err := time.Parse(time.RFC3339, strings.TrimSpace(parts[0]))
	if err != nil {
		return nil, fmt.Errorf("invalid from time: %v", err)
	}

	to, err := time.Parse(time.RFC3339, strings.TrimSpace(parts[1]))
	if err != nil {
		return nil, fmt.Errorf("invalid to time: %v", err)
	}

	return &TimeRange{From: from, To: to}, nil
}
