package usecase

import (
	"cmd/etl/internal/models"
	"cmd/etl/internal/utils"
	"cmd/etl/pkg/slogger/wsl"
	"context"
	"fmt"
	"log/slog"
)

// processNewLogs обрабатывает новые записи из IdsLogs
func (uc *UseCase) ProcessNewLogs(ctx context.Context) error {
	uc.l.DebugContext(ctx, "checking for new IDS logs")

	// Получаем новые записи
	logs, err := uc.GetLogs(ctx, uc.lastLog, uc.maxCount)
	if err != nil {
		return fmt.Errorf("failed to get new logs: %w", err)
	}

	if len(logs) == 0 {
		uc.l.DebugContext(ctx, "no new logs found")
		return nil
	}

	uc.l.InfoContext(ctx, "processing new IDS logs", slog.Int("count", len(logs)))

	// Обрабатываем каждую запись
	for _, log := range logs {
		if err := uc.ProcessLog(ctx, log); err != nil {
			uc.l.ErrorContext(ctx, "failed to process log",
				wsl.Int64("log_id", log.ID), wsl.Err(err))
			// Продолжаем обработку следующих записей
			continue
		}
		// Обновляем lastLog
		if uc.lastLog == nil {
			uc.lastLog = &log
		} else if log.ID > uc.lastLog.ID && log.Timestamp.After(uc.lastLog.Timestamp) {
			uc.lastLog = &log
		}
	}

	uc.l.DebugContext(ctx, "successfully processed IDS logs",
		slog.Int("processed_count", len(logs)))
	if uc.lastLog != nil {
		uc.l.InfoContext(ctx, "last log",
			slog.Int("id", int(uc.lastLog.ID)))
	}
	return nil
}

// processLog обрабатывает одну запись лога
func (uc *UseCase) ProcessLog(ctx context.Context, log models.IdsLog) error {
	//uc.l.DebugContext(ctx, "processing log", wsl.Int64("log_id", log.ID))

	// Обрабатываем источник (Source)
	source, err := uc.ProcessSource(ctx, log)
	if err != nil {
		return fmt.Errorf("failed to process source: %w", err)
	}

	// Обрабатываем домен (Domain)
	domain, err := uc.ProcessDomain(ctx, log)
	if err != nil {
		return fmt.Errorf("failed to process domain: %w", err)
	}

	// Обрабатываем устройство (Device)
	device, err := uc.ProcessDevice(ctx, log)
	if err != nil {
		return fmt.Errorf("failed to process device: %w", err)
	}

	// Обрабатываем сессию (Session)
	if err := uc.ProcessSession(ctx, log, source, domain, device); err != nil {
		return fmt.Errorf("failed to process session: %w", err)
	}

	//uc.l.DebugContext(ctx, "successfully processed log", slog.Uint64("log_id", uint64(log.ID)))
	return nil
}

// processSource обрабатывает и создает/обновляет запись Source
func (uc *UseCase) ProcessSource(ctx context.Context, log models.IdsLog) (*models.Source, error) {
	// Ищем существующий источник
	existingSource, err := uc.GetSourceByAddr(ctx, log.SrcIP)
	if err != nil {
		return nil, fmt.Errorf("failed to get source: %w", err)
	}

	if existingSource != nil {
		// Источник уже существует, можно обновить счетчики и т.д.
		//uc.l.DebugContext(ctx, "source already exists",
		//	wsl.String("ip", log.SrcIP), wsl.Int("port", log.SrcPort))
		return existingSource, nil
	}

	// Создаем новый источник
	source := models.Source{
		IP:       log.SrcIP,
		Country:  log.SrcCountry,
		Username: log.Username, // TODO: пока username не извлекаетс системой
	}

	if err := uc.CreateSource(ctx, source); err != nil {
		return nil, fmt.Errorf("failed to create source: %w", err)
	}
	uc.l.DebugContext(ctx, "source created",
		wsl.String("ip", log.SrcIP), wsl.Int("port", log.SrcPort))
	// Получаем созданный источник чтобы вернуть с ID
	newSource, err := uc.GetSourceByAddr(ctx, log.SrcIP)
	if err != nil {
		return nil, fmt.Errorf("failed to get created source: %w", err)
	}
	return newSource, nil
}

// processDomain обрабатывает и создает/обновляет запись Domain
func (uc *UseCase) ProcessDomain(ctx context.Context, log models.IdsLog) (*models.Domain, error) {
	// алиасы для полей лога
	domain := log.DestDomain
	ip := log.DestIP
	port := log.DestPort
	rawURL := log.Payload
	// перемены для заполнения
	var existingDomain *models.Domain
	path := ""
	var err error
	// Находим path и URL из лога
	if rawURL != "" {
		path, err = utils.ExtractDomain(rawURL)
		if err != nil {
			return nil, fmt.Errorf("failed to parse URL: %w", err)
		}
		// Логируем, если домен, полученный из URL, не совпадает с доменом, полученным из поля  КСУ
		if domain != "" && domain != path {
			uc.l.WarnContext(ctx,
				"different domains in URL and log.DestDomain",
				wsl.String("url domain", path),
				wsl.String("log.DestDomain", domain),
			)
		}
	} else if domain != "" {
		path = domain
	}
	// Ищем существующий домен по имени или IP
	// по доменному имени
	if path != "" {
		existingDomain, err = uc.GetDomainByPath(ctx, path)
		if err != nil {
			return nil, fmt.Errorf("failed to get domain by path: %w", err)
		}
		// по IP
	} else if ip != "" {
		existingDomain, err = uc.GetDomainByAddr(ctx, ip, port)
		if err != nil {
			return nil, fmt.Errorf("failed to get domain by addr: %w", err)
		}
		// Иначе ошибка - в строке нет данных для работы с доменом
	} else {
		return nil, fmt.Errorf("empty log")
	}
	// Обработка полученных данных
	if existingDomain != nil {
		// Домен уже существует, обновляем счетчики
		if err := uc.UpdateDomainLastAccess(ctx, existingDomain.ID, log.Timestamp); err != nil {
			uc.l.ErrorContext(ctx, "failed to update domain last access", wsl.Err(err))
		}
		//uc.l.DebugContext(ctx, "domain already exists, updated counters",
		//	wsl.String("ip", log.DestIP), wsl.Int("port", log.DestPort))
		return existingDomain, nil
	}
	// Создаем новый домен
	newDomain := models.Domain{
		IP:                     ip,
		Port:                   port,
		Country:                log.DestCountry,
		Path:                   path,
		AccessCount:            1,
		AnalysisAttemptsCount:  0,
		ContentAnalysisCounter: 0,
		DecisionID:             0, // По умолчанию
		LastAccessDatetime:     log.Timestamp,
	}

	if err := uc.CreateDomain(ctx, newDomain); err != nil {
		return nil, fmt.Errorf("failed to create domain: %w", err)
	}

	uc.l.InfoContext(ctx, "created new domain",
		slog.String("path", path),
		slog.String("ip", ip), slog.Int("port", port))

	// Получаем созданный домен чтобы вернуть с ID
	var createdDomain *models.Domain
	if path != "" {
		createdDomain, err = uc.GetDomainByPath(ctx, path)
		if err != nil {
			return nil, fmt.Errorf("failed to get created domain: %w", err)
		}
	} else {
		createdDomain, err = uc.GetDomainByAddr(ctx, ip, port)
		if err != nil {
			return nil, fmt.Errorf("failed to get created domain: %w", err)
		}
	}
	return createdDomain, nil
}

// processDevice обрабатывает и создает/обновляет запись Device
func (uc *UseCase) ProcessDevice(ctx context.Context, log models.IdsLog) (*models.Device, error) {
	// Ищем устройство по идентификатору
	existingDevice, err := uc.GetDeviceByID(ctx, log.SensorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get device: %w", err)
	}
	if existingDevice != nil {
		//uc.l.DebugContext(ctx, "device already exists",
		//	wsl.String("hostname", log.Hostname))
		return existingDevice, nil
	}

	// Создаем новое устройство
	device := models.Device{
		ID:       uint(log.SensorID),
		HostName: log.Hostname,
	}

	if err := uc.CreateDevice(ctx, device); err != nil {
		return nil, fmt.Errorf("failed to create device: %w", err)
	}

	uc.l.InfoContext(ctx, "created new device", slog.String("device name", log.Hostname))

	// Получаем созданное устройство
	newDevice, err := uc.GetDeviceByID(ctx, log.SensorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get created device: %w", err)
	}
	return newDevice, nil
}

// processSession создает запись Session
func (uc *UseCase) ProcessSession(ctx context.Context, log models.IdsLog, source *models.Source, domain *models.Domain, device *models.Device) error {
	// Парсим URL
	url, err := utils.ExtractURL(log.Payload)
	if err != nil {
		return fmt.Errorf("failed to parse URL: %w", err)
	}
	// Создаем сессию
	session := models.Session{
		DatetimeUTC: log.Timestamp,
		DeviceID:    device.ID,
		Type:        log.EventType,
		// TODO: Добавить логику обработки и получения статуса
		Status:   log.Action,
		URL:      url,
		DomainID: domain.ID,
		SrcID:    source.ID,
		Proto:    log.Proto,
	}

	err = uc.CreateSession(ctx, session)
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	//uc.l.DebugContext(ctx, "created new session",
	//	wsl.String("session dst_ip", session.IP),
	//	wsl.String("proto", session.Proto))
	return nil
}
