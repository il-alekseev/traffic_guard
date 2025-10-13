package parser

import (
	"cmd/etl/internal/models"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

// ParseCSVFileToIdsLog парсит CSV файл по указанному пути и возвращает слайс структур IdsLog
func ParseCSVFileToIdsLog(filePath string) ([]models.IdsLog, error) {
	// Открываем файл
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// Создаем CSV reader
	csvReader := csv.NewReader(file)
	csvReader.Comma = ','
	csvReader.FieldsPerRecord = -1 // Разрешаем разное количество полей в строках

	var logs []models.IdsLog

	// Читаем и парсим все записи
	for {
		record, err := csvReader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, fmt.Errorf("error reading CSV record: %v", err)
		}

		// Пропускаем пустые строки
		if len(record) == 0 {
			continue
		}

		log, err := parseRecordToIdsLog(record)
		if err != nil {
			// Можно добавить логирование ошибок парсинга и продолжить обработку
			fmt.Printf("record: %v", record)
			fmt.Printf("Warning: skipping record due to error: %v\n", err)
			continue
		}

		logs = append(logs, log)
	}

	return logs, nil
}

// ParseCSVToIdsLog парсит CSV файл и возвращает слайс структур IdsLog
func ParseCSVToIdsLog(reader io.Reader) ([]models.IdsLog, error) {
	csvReader := csv.NewReader(reader)
	csvReader.Comma = ','

	var logs []models.IdsLog

	// Пропускаем заголовок если есть
	// Если заголовка нет, закомментируйте следующие две строки
	if _, err := csvReader.Read(); err != nil {
		return nil, fmt.Errorf("error reading header: %v", err)
	}

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading record: %v", err)
		}

		log, err := parseRecordToIdsLog(record)
		if err != nil {
			return nil, fmt.Errorf("error parsing record: %v", err)
		}

		logs = append(logs, log)
	}

	return logs, nil
}

// parseRecordToIdsLog парсит одну строку CSV в структуру IdsLog
func parseRecordToIdsLog(record []string) (models.IdsLog, error) {
	var log models.IdsLog
	var err error

	if len(record) < 30 {
		return log, fmt.Errorf("invalid record length: expected at least 30 fields, got %d", len(record))
	}

	// ID
	log.ID, err = strconv.ParseInt(record[0], 10, 64)
	if err != nil {
		return log, fmt.Errorf("error parsing ID: %v", err)
	}

	// Timestamp
	log.Timestamp, err = time.Parse("2006-01-02 15:04:05.999", record[1])
	if err != nil {
		return log, fmt.Errorf("error parsing timestamp: %v", err)
	}

	// TimestampStr
	log.TimestampStr = record[2]

	// SensorID
	log.SensorID, err = strconv.Atoi(record[3])
	if err != nil {
		log.SensorID = 0 // или значение по умолчанию
	}

	// EventType
	log.EventType = record[4]

	// SrcIP
	log.SrcIP = record[11]

	// SrcPort
	if record[18] != "" {
		log.SrcPort, err = strconv.Atoi(record[18])
		if err != nil {
			log.SrcPort = 0
		}
	}

	// DestIP
	log.DestIP = record[13]

	// DestPort
	if record[19] != "" {
		log.DestPort, err = strconv.Atoi(record[19])
		if err != nil {
			log.DestPort = 0
		}
	}

	// Proto
	log.Proto = record[17]

	// InIface
	log.InIface = record[20]

	// Action
	log.Action = record[8]

	// SignatureID
	if record[9] != "" {
		log.SignatureID, err = strconv.Atoi(record[9])
		if err != nil {
			log.SignatureID = 0
		}
	}

	// Signature
	log.Signature = record[5]

	// Category
	log.Category = record[6]

	// Severity
	if record[10] != "" {
		severity, err := strconv.Atoi(record[10])
		if err != nil {
			log.Severity = 0
		} else {
			log.Severity = int16(severity)
		}
	}

	// Payload
	log.Payload = record[21]

	// Packet
	log.Packet = record[22]

	// Msgrepeatcount
	if record[23] != "" {
		log.Msgrepeatcount, err = strconv.Atoi(record[23])
		if err != nil {
			log.Msgrepeatcount = 0
		}
	}

	// DestDomain
	log.DestDomain = record[24]

	// Revision
	if record[25] != "" {
		log.Revision, err = strconv.Atoi(record[25])
		if err != nil {
			log.Revision = 0
		}
	}

	// SignatureBody
	log.SignatureBody = record[26]

	// Sourcename
	log.Sourcename = record[27]

	// Hostname
	log.Hostname = record[28]

	// DestCountry
	log.DestCountry = record[14]

	// SrcCountry
	log.SrcCountry = record[12]

	// Username
	log.Username = record[29]

	// Vrf
	if len(record) > 30 {
		log.Vrf = record[30]
	} else {
		log.Vrf = "default" // значение по умолчанию, так как поле not null
	}

	return log, nil
}
