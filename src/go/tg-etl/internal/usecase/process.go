package usecase

import (
	"cmd/etl/internal/models"
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
	uc.l.DebugContext(ctx, "processing log", wsl.Int64("log_id", log.ID))

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

	uc.l.DebugContext(ctx, "successfully processed log", slog.Uint64("log_id", uint64(log.ID)))
	return nil
}

// processSource обрабатывает и создает/обновляет запись Source
func (uc *UseCase) ProcessSource(ctx context.Context, log models.IdsLog) (*models.Source, error) {
	// Ищем существующий источник
	existingSource, err := uc.GetSourceByAddr(ctx, log.SrcIP, log.SrcPort)
	if err != nil {
		return nil, fmt.Errorf("failed to get source: %w", err)
	}

	if existingSource != nil {
		// Источник уже существует, можно обновить счетчики и т.д.
		uc.l.DebugContext(ctx, "source already exists",
			wsl.String("ip", log.SrcIP), wsl.Int("port", log.SrcPort))
		return existingSource, nil
	}

	// Создаем новый источник
	source := models.Source{
		IP:       log.SrcIP,
		Port:     log.SrcPort,
		Country:  log.SrcCountry,
		Username: log.Username, // TODO: пока username не извлекаетс системой
	}

	if err := uc.CreateSource(ctx, source); err != nil {
		return nil, fmt.Errorf("failed to create source: %w", err)
	}

	uc.l.InfoContext(ctx, "created new source",
		slog.String("ip", log.SrcIP), slog.Int("port", log.SrcPort))

	// Получаем созданный источник чтобы вернуть с ID
	newSource, err := uc.GetSourceByAddr(ctx, log.SrcIP, log.SrcPort)
	if err != nil {
		return nil, fmt.Errorf("failed to get created source: %w", err)
	}
	return newSource, nil
}

// processDomain обрабатывает и создает/обновляет запись Domain
func (uc *UseCase) ProcessDomain(ctx context.Context, log models.IdsLog) (*models.Domain, error) {
	// Ищем существующий домен
	existingDomain, err := uc.GetDomainByAddr(ctx, log.DestIP, log.DestPort)
	if err != nil {
		return nil, fmt.Errorf("failed to get domain: %w", err)
	}

	if existingDomain != nil {
		// Домен уже существует, обновляем счетчики
		if err := uc.UpdateDomainLastAccess(ctx, existingDomain.ID, log.Timestamp); err != nil {
			uc.l.ErrorContext(ctx, "failed to update domain last access", wsl.Err(err))
		}

		uc.l.DebugContext(ctx, "domain already exists, updated counters",
			wsl.String("ip", log.DestIP), wsl.Int("port", log.DestPort))
		return existingDomain, nil
	}

	// Создаем новый домен
	domain := models.Domain{
		IP:                     log.DestIP,
		Port:                   log.DestPort,
		Country:                log.DestCountry,
		Path:                   log.DestDomain, // TODO: пока система не выдает домена
		AccessCount:            1,
		AnalysisAttemptsCount:  0,
		ContentAnalysisCounter: 0,
		DecisionID:             0, // По умолчанию
		LastAccessDatetime:     log.Timestamp,
	}

	if err := uc.CreateDomain(ctx, domain); err != nil {
		return nil, fmt.Errorf("failed to create domain: %w", err)
	}

	uc.l.InfoContext(ctx, "created new domain",
		slog.String("ip", log.DestIP), slog.Int("port", log.DestPort))

	// Получаем созданный домен чтобы вернуть с ID
	newDomain, err := uc.GetDomainByAddr(ctx, log.DestIP, log.DestPort)
	if err != nil {
		return nil, fmt.Errorf("failed to get created domain: %w", err)
	}

	return newDomain, nil
}

// processDevice обрабатывает и создает/обновляет запись Device
func (uc *UseCase) ProcessDevice(ctx context.Context, log models.IdsLog) (*models.Device, error) {
	// Ищем устройство по идентификатору
	existingDevice, err := uc.GetDeviceByID(ctx, log.SensorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get device: %w", err)
	}
	if existingDevice != nil {
		uc.l.DebugContext(ctx, "device already exists",
			wsl.String("hostname", log.Hostname))
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

	// Создаем сессию
	session := models.Session{
		DatetimeUTC: log.Timestamp,
		DeviceID:    device.ID,
		Type:        log.EventType,
		// TODO: Добавить логику обработки и получения статуса
		Status: log.Action,
		// TODO: пока система не выдает URL
		URL: "",
		// TODO: проверить корректность получения ID домена
		DomainID: domain.ID,
		SrcID:    source.ID,
		Proto:    log.Proto,
	}

	err := uc.CreateSession(ctx, session)
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	uc.l.DebugContext(ctx, "created new session",
		wsl.String("session dst_ip", session.IP),
		wsl.String("proto", session.Proto))
	return nil
}
