package postresql

import (
	"context"
	"encoding/json"
	"fmt"
	"tg-dbd/internal/models"
	"tg-dbd/pkg/trparser"
	"time"
)

// Получение статистики по входящему и исходящему трафику за временной интервал в килобайтах
func (r *RepoPG) GetTrafficStat(ctx context.Context, tr *trparser.TimeRange, count uint) (models.TrafficStat, error) {
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

	// Создаем временные интервалы заранее
	for i := uint(0); i < count; i++ {
		start := tr.From.Add(time.Duration(i) * intervalDuration)
		result.Time[i] = start.Add(intervalDuration / 2)
	}

	// Получаем все записи используя GORM
	var statsRecords []models.StatsJSON
	err := r.db.GetDB().WithContext(ctx).
		Where("date BETWEEN ? AND ?", tr.From, tr.To).
		Order("date").
		Find(&statsRecords).Error

	if err != nil {
		return result, fmt.Errorf("failed to execute query: %w", err)
	}

	// Временные переменные для накопления байтов
	inputBytesTotal := make([]int64, count)
	outputBytesTotal := make([]int64, count)

	// Обрабатываем каждую запись и распределяем по интервалам
	for _, record := range statsRecords {
		// Преобразуем datatypes.JSON в []byte
		statJSON, err := record.Stat.MarshalJSON()
		if err != nil {
			// Пропускаем записи с ошибками JSON
			continue
		}

		// Определяем индекс интервала
		timeDiff := record.Date.Sub(tr.From)
		intervalIndex := int(timeDiff / intervalDuration)

		if intervalIndex >= 0 && intervalIndex < int(count) {
			inputBytes, outputBytes, err := getInterfaceBytes(statJSON, "ge-0-0")
			if err == nil {
				inputBytesTotal[intervalIndex] += inputBytes
				outputBytesTotal[intervalIndex] += outputBytes
			}
		}
	}

	// Преобразуем накопленные байты в килобайты
	for i := uint(0); i < count; i++ {
		result.Input[i] = uint(inputBytesTotal[i] / 1024)
		result.Output[i] = uint(outputBytesTotal[i] / 1024)
	}

	return result, nil
}

// Функция для получения числа байтов из JSON
func getInterfaceBytes(statsJSON []byte, interfaceName string) (inputBytes, outputBytes int64, err error) {
	var data models.Statistics

	if err := json.Unmarshal(statsJSON, &data); err != nil {
		return 0, 0, fmt.Errorf("ошибка парсинга JSON: %v", err)
	}

	iface, exists := data.Statistics.Network.Interfaces[interfaceName]
	if !exists {
		return 0, 0, fmt.Errorf("интерфейс %s не найден", interfaceName)
	}

	return iface.Input.Bytes.Count, iface.Output.Bytes.Count, nil
}
