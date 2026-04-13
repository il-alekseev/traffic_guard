package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"tg-etl/internal/controllers/values"
	"tg-etl/internal/models"
	"tg-etl/internal/utils"
	"tg-etl/pkg/blog/operations"
	"tg-etl/pkg/slogger/wsl"
	"time"

	pkg "tg-etl/pkg/models"

	"github.com/google/uuid"
)

// processNewLogs обрабатывает новые записи из IdsLogs
func (uc *UseCase) ProcessNewLogs(ctx context.Context) (bool, error) {
	uc.processingLock.Lock()
	defer uc.processingLock.Unlock()
	var logsAfterExists = false
	// Получаем новые записи
	logs, err := uc.q.GetLogs(ctx, uc.lastLog, uc.batchSize)
	if err != nil {
		return logsAfterExists, fmt.Errorf("failed to get new logs: %w", err)
	}
	if len(logs) == 0 {
		uc.l.InfoContext(ctx, "new logs for process not found")
		return logsAfterExists, nil
	}
	var lastIDSLog *models.IdsLog
	// Обрабатываем каждую запись
	for _, log := range logs {
		if err := uc.processLog(ctx, log); err != nil {
			uc.l.ErrorContext(ctx, "failed to process log",
				wsl.Int64("log_id", log.ID), wsl.Err(err))
			// Продолжаем обработку следующих записей
			continue
		}
		// Обновляем lastLog
		if lastIDSLog == nil {
			lastIDSLog = &log
		} else if log.ID > lastIDSLog.ID && log.Timestamp.After(lastIDSLog.Timestamp) {
			lastIDSLog = &log
		}
	}
	id := 0

	if lastIDSLog != nil {
		id = int(lastIDSLog.ID)
		// Сохраняем в базу последний обработанны лог
		log := models.LastLog{
			ID:        uint(lastIDSLog.ID),
			Timestamp: lastIDSLog.Timestamp,
		}
		if err := uc.q.CreateOrUpdateLastLog(ctx, log); err != nil {
			uc.l.ErrorContext(ctx, "error with create or update last log", wsl.Err(err))
		}
		uc.lastLog = &log
	}
	uc.l.InfoContext(ctx, "successfully processed IDS logs",
		wsl.Int("processed_count", len(logs)),
		wsl.Int("last_log id", id),
	)
	// Если получили логов столько же сколько батчсайз, значит логи еще есть, выполняем без перерыва
	if len(logs) == int(uc.batchSize) {
		logsAfterExists = true
	}
	return logsAfterExists, nil
}

// processLog обрабатывает одну запись лога
func (uc *UseCase) processLog(ctx context.Context, log models.IdsLog) error {
	// Обрабатываем устройство (Device)
	device, err := uc.getDevice(ctx, log)
	if err != nil {
		return fmt.Errorf("failed to process device: %w", err)
	}

	// Обрабатываем источник (Source)
	source, err := uc.getSource(ctx, log)
	if err != nil {
		return fmt.Errorf("failed to process source: %w", err)
	}

	// Обрабатываем URL (URL)
	domain, url, err := uc.getDomainAndURL(ctx, log)
	if err != nil {
		return fmt.Errorf("failed to process domain: %w", err)
	}

	// Добавляем сессию (Session)
	if err := uc.createSession(ctx, log, source, domain, device, url); err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	return nil
}

// getDevice обрабатывает и создает/обновляет запись Device
func (uc *UseCase) getDevice(ctx context.Context, log models.IdsLog) (*models.Device, error) {
	method := "getDevice"
	// Ищем устройство по идентификатору
	existingDevice, err := uc.q.GetDeviceByID(ctx, log.SensorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get device: %w", err)
	}
	if existingDevice != nil {
		return existingDevice, nil
	}

	// Создаем новое устройство
	device := models.Device{
		ID:       uint(log.SensorID),
		HostName: log.Hostname,
	}

	if err := uc.q.CreateDevice(ctx, device); err != nil {
		return nil, fmt.Errorf("failed to create device: %w", err)
	}

	// Создание групп и ролей в keycloak
	authInfo, err := uc.getServiceAuthInfo(ctx)
	if err != nil {
		uc.l.WarnContext(ctx, "failed to get service auth, device created without keycloak roles",
			slog.String("device", log.Hostname),
			slog.String("error", err.Error()))
		// Продолжаем без ролей в Keycloak
		return &device, nil
	}
	// Проверяем, создавалась ли уже роль для данного сетевого узла
	roleExists, err := uc.hasDeviceInRoles(&authInfo, log.Hostname)
	if err != nil {
		return nil, fmt.Errorf("failed to get roles for device: %w", err)
	}
	// Если не создавалась, то создаем
	if !roleExists {
		if err = uc.addRolesToKeyCloak(ctx, &authInfo, log.Hostname); err != nil {
			return nil, fmt.Errorf("failed to create role for device: %w", err)
		}
	}
	uc.l.InfoContext(ctx, "created new device", slog.String("device name", log.Hostname))

	// Запись события в бизнес-лог
	uuidStr := uuid.New().String()

	newValue := make(map[string]any)
	newValue["id"] = device.ID
	newValue["hostname"] = device.HostName

	record := pkg.DtoBusinessLog{
		Description: fmt.Sprintf("создание нового устройства %s", device.HostName),
		Entity:      "device",
		EntityID:    strconv.FormatUint(uint64(device.ID), 10),
		UserName:    values.ThisServiceName,
		NewValue:    newValue,
		OldValue:    "",
		EventType:   "CREATE",
	}

	_, err = uc.blclient.Operations.PostAPIV1Add(&operations.PostAPIV1AddParams{
		XCallerService: values.ThisServiceName,
		XRequestID:     uuidStr,
		Record:         &record,
		Context:        ctx,
	},
	)
	if err != nil {
		err = fmt.Errorf("%s: failed to blog event: %w", method, err)
		uc.l.ErrorContext(ctx, "failed to blog event",
			wsl.String("method", method),
			wsl.String("error", err.Error()),
		)
		return nil, err
	}

	// Получаем созданное устройство
	newDevice, err := uc.q.GetDeviceByID(ctx, log.SensorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get created device: %w", err)
	}
	return newDevice, nil
}

// getSource обрабатывает и создает запись Source
func (uc *UseCase) getSource(ctx context.Context, log models.IdsLog) (*models.Source, error) {
	// Ищем существующий источник
	existingSource, err := uc.q.GetSourceByAddr(ctx, log.SrcIP)
	if err != nil {
		return nil, fmt.Errorf("failed to get source: %w", err)
	}

	if existingSource != nil {
		return existingSource, nil
	}

	// Создаем новый источник
	source := models.Source{
		IP:       log.SrcIP,
		Country:  log.SrcCountry,
		Username: log.Username, // TODO: пока username не извлекается системой
	}

	if err := uc.q.CreateSource(ctx, source); err != nil {
		return nil, fmt.Errorf("failed to create source: %w", err)
	}
	uc.l.DebugContext(ctx, "source created",
		wsl.String("ip", log.SrcIP), wsl.Int("port", log.SrcPort))
	// Получаем созданный источник чтобы вернуть с ID
	newSource, err := uc.q.GetSourceByAddr(ctx, log.SrcIP)
	if err != nil {
		return nil, fmt.Errorf("failed to get created source: %w", err)
	}
	return newSource, nil
}

func (uc *UseCase) getDomainAndURL(ctx context.Context, log models.IdsLog) (*models.Domain, *models.URL, error) {
	uc.l.DebugContext(ctx, "got domain from KSU: ",
		wsl.String("url", log.Payload),
	)
	// Находим path и URL из лога
	if log.Payload != "" {
		return uc.handleURLWithDomain(ctx, log, log.DestDomain, log.Payload, log.DestIP, log.DestPort)
	}

	if log.DestDomain != "" {
		return uc.handleDomainOnly(ctx, log, log.DestDomain, log.DestIP, log.DestPort)
	}

	if log.DestIP != "" {
		return uc.handleIPOnly(ctx, log, log.DestIP, log.DestPort)
	}

	return nil, nil, fmt.Errorf("empty log: no URL, domain or IP provided")
}

func (uc *UseCase) handleURLWithDomain(ctx context.Context, log models.IdsLog, domain, rawURL, ip string, port int) (*models.Domain, *models.URL, error) {
	urlInfo, err := utils.ParseRawURL(rawURL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse raw URL: %w", err)
	}
	// Логируем несовпадение доменов
	if domain != "" && domain != urlInfo.Domain {
		uc.l.WarnContext(ctx,
			"different domains in URL and log.DestDomain",
			wsl.String("url_domain", urlInfo.Domain),
			wsl.String("log_dest_domain", domain),
		)
	}
	existingDomain, err := uc.q.GetDomainByPath(ctx, urlInfo.Domain)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get domain by path: %w", err)
	}

	if existingDomain != nil {
		return uc.handleExistingDomainWithURL(ctx, log, existingDomain, urlInfo)
	}
	// Создаем новый домен и URL
	return uc.createNewDomainAndURL(ctx, log, urlInfo, ip, port)
}

// handleExistingDomainWithURL обрабатывает существующий домен с URL
func (uc *UseCase) handleExistingDomainWithURL(ctx context.Context, log models.IdsLog, domain *models.Domain, urlInfo utils.URLInfo) (*models.Domain, *models.URL, error) {
	// Ищем существующий URL
	existingURL, err := uc.q.GetURLByPath(ctx, urlInfo.URL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get URL by path: %w", err)
	}

	if existingURL != nil {
		return domain, existingURL, nil
	}

	// Создаем новый URL
	return uc.createURLForDomain(ctx, log, domain, urlInfo)
}

// createURLForDomain создает новый URL для домена
func (uc *UseCase) createURLForDomain(ctx context.Context, log models.IdsLog, domain *models.Domain, urlInfo utils.URLInfo) (*models.Domain, *models.URL, error) {
	newURL := models.URL{
		Path:      urlInfo.URL,
		Proto:     urlInfo.Proto,
		DomainID:  domain.ID,
		IDSLogsAt: log.Timestamp,
	}

	// Проверяем списки и отправляем в Kafka при необходимости
	list, err := uc.q.GetListByDomainID(ctx, domain.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get list type for domain: %w", err)
	}

	// Если домен не в списках - отправляем на анализ
	if list == nil {
		ts := time.Now()
		uuid := uuid.New()
		newURL.RequestID = uuid
		newURL.PutKafkaAt = ts

		req := models.AnalysisRequest{
			RequestID: newURL.RequestID,
			Src:       models.Src{IP: log.SrcIP},
			Dst: models.Destination{
				Type:     models.DestinationURL,
				Resource: newURL.Path,
				Port:     &log.DestPort,
				Proto:    (*models.Protocol)(&urlInfo.Proto),
			},
			Timestamp: ts,
		}

		if err := uc.kc.SendAnalysisRequest(ctx, req); err != nil {
			return nil, nil, fmt.Errorf("failed to send URL request to Kafka: %w", err)
		}
		//uc.l.Debug("url content analysis sent", wsl.String("url", newURL.Path), wsl.String("request_id", newURL.RequestID.String()))
	}

	if err := uc.q.CreateURL(ctx, newURL); err != nil {
		return nil, nil, fmt.Errorf("failed to create new URL: %w", err)
	}

	// Снова получаем URL, чтобы получить его ID
	url, err := uc.q.GetURLByPath(ctx, newURL.Path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get URL: %w", err)
	}

	return domain, url, nil
}

// createNewDomainAndURL создает новый домен и URL
func (uc *UseCase) createNewDomainAndURL(ctx context.Context, log models.IdsLog, urlInfo utils.URLInfo, ip string, port int) (*models.Domain, *models.URL, error) {
	newDomain := models.Domain{
		IP:       ip,
		Port:     port,
		Country:  log.DestCountry,
		Path:     urlInfo.Domain,
		ActionID: 0,
	}
	if err := uc.q.CreateDomain(ctx, newDomain); err != nil {
		return nil, nil, fmt.Errorf("failed to create new domain: %w", err)
	}
	uc.l.DebugContext(ctx, "domain created",
		wsl.String("path", newDomain.Path),
	)
	// Получаем домен из базы, чтобы найти его ID
	domain, err := uc.q.GetDomainByPath(ctx, newDomain.Path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get domain: %w", err)
	}

	newURL := models.URL{
		Path:      urlInfo.URL,
		Proto:     urlInfo.Proto,
		DomainID:  domain.ID,
		RequestID: uuid.New(),
		IDSLogsAt: log.Timestamp,
	}
	// Создаем запрос для Kafka
	ts := time.Now()
	req := models.AnalysisRequest{
		RequestID: newURL.RequestID,
		Src:       models.Src{IP: log.SrcIP},
		Dst: models.Destination{
			Type:     models.DestinationURL,
			Resource: newURL.Path,
			Port:     &log.DestPort,
			Proto:    (*models.Protocol)(&urlInfo.Proto),
		},
		Timestamp: ts,
	}
	newURL.PutKafkaAt = ts
	// Сохраняем в БД
	if err := uc.q.CreateURL(ctx, newURL); err != nil {
		return nil, nil, fmt.Errorf("failed to create new URL: %w", err)
	}
	// Снова получаем URL, чтобы получить его ID
	url, err := uc.q.GetURLByPath(ctx, newURL.Path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get URL: %w", err)
	}
	// Отправляем в Kafka
	if err := uc.kc.SendAnalysisRequest(ctx, req); err != nil {
		return nil, nil, fmt.Errorf("failed to send URL request to Kafka: %w", err)
	}
	//uc.l.Debug("url content analysis sent", wsl.String("url", url.Path), wsl.String("request_id", url.RequestID.String()))
	return domain, url, nil
}

func (uc *UseCase) handleDomainOnly(ctx context.Context, log models.IdsLog, domain, ip string, port int) (*models.Domain, *models.URL, error) {
	existingDomain, err := uc.q.GetDomainByPath(ctx, domain)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get domain by path: %w", err)
	}

	if existingDomain != nil {
		return existingDomain, &models.URL{}, nil
	}

	return uc.createDomainWithAnalysis(ctx, log, domain, ip, port, models.DestinationDomain, domain)
}

// handleIPOnly обрабатывает случай только с IP
func (uc *UseCase) handleIPOnly(ctx context.Context, log models.IdsLog, ip string, port int) (*models.Domain, *models.URL, error) {
	existingDomain, err := uc.q.GetDomainByAddr(ctx, ip, port)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get domain by addr: %w", err)
	}

	if existingDomain != nil {
		return existingDomain, &models.URL{}, nil
	}

	return uc.createDomainWithAnalysis(ctx, log, "", ip, port, models.DestinationIP, ip)
}

// createDomainWithAnalysis создает домен и отправляет запрос на анализ
func (uc *UseCase) createDomainWithAnalysis(ctx context.Context, log models.IdsLog, domainPath, ip string, port int, destType models.DestinationType, resource string) (*models.Domain, *models.URL, error) {
	newDomain := models.Domain{
		IP:      ip,
		Port:    port,
		Country: log.DestCountry,
		Path:    domainPath,
	}

	if err := uc.q.CreateDomain(ctx, newDomain); err != nil {
		return nil, nil, fmt.Errorf("failed to create new domain: %w", err)
	}
	uc.l.DebugContext(ctx, "domain created",
		wsl.String("path", newDomain.Path),
	)
	// Создаем URL для анализа
	createdDomain, err := uc.q.GetDomainByAddr(ctx, newDomain.IP, newDomain.Port)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get created domain: %w", err)
	}

	newURL := models.URL{
		DomainID:  createdDomain.ID,
		RequestID: uuid.New(),
		IDSLogsAt: log.Timestamp,
	}
	// Создаем запрос в Kafka
	ts := time.Now()
	req := models.AnalysisRequest{
		RequestID: newURL.RequestID,
		Src:       models.Src{IP: log.SrcIP},
		Dst: models.Destination{
			Type:     destType,
			Resource: resource,
			Port:     &log.DestPort,
		},
		Timestamp: ts,
	}
	// Сохраняем URL в БД
	newURL.PutKafkaAt = ts
	if err := uc.q.CreateURL(ctx, newURL); err != nil {
		return nil, nil, fmt.Errorf("failed to create URL: %w", err)
	}
	// Снова получаем URL, чтобы получить его ID
	// Здесь используем ID домена, т.к поиск по только по пути выдаст множество URL,
	// т.к они все пустые
	url, err := uc.q.GetURLByPathDomain(ctx, newURL.Path, createdDomain.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get URL: %w", err)
	}
	// Отправляем в Kafka
	if err := uc.kc.SendAnalysisRequest(ctx, req); err != nil {
		return nil, nil, fmt.Errorf("failed to send analysis request to Kafka: %w", err)
	}
	uc.l.Debug("content analysis sent", wsl.String("destType", string(destType)), wsl.String("url", url.Path), wsl.String("request_id", url.RequestID.String()))
	return createdDomain, url, nil
}

// createSession создает запись Session
func (uc *UseCase) createSession(ctx context.Context, log models.IdsLog, source *models.Source, domain *models.Domain, device *models.Device, url *models.URL) error {
	if device == nil {
		return fmt.Errorf("device is required")
	}
	if source == nil {
		return fmt.Errorf("source is required")
	}
	if domain == nil {
		return fmt.Errorf("source is required")
	}
	if url == nil {
		return fmt.Errorf("source is required")
	}

	status, err := uc.processStatus(ctx, log, domain)
	if err != nil {
		return fmt.Errorf("failed to process status")
	}

	sessionType, err := uc.processSessionType(domain)
	if err != nil {
		return fmt.Errorf("failed to process session type")
	}

	// Создаем сессию
	session := models.Session{
		DatetimeUTC: log.Timestamp,
		DeviceID:    device.ID,
		Type:        sessionType.String(),
		Status:      status.String(),
		URLID:       uint(url.ID),
		DomainID:    uint(domain.ID),
		SrcID:       source.ID,
	}

	err = uc.q.CreateSession(ctx, session)
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	return nil
}

func (uc *UseCase) processStatus(ctx context.Context, log models.IdsLog, domain *models.Domain) (models.Status, error) {
	// Алгоритм определения статуса:
	// 1. Если log.Action == "blocked" -> Заблокирован
	// 2. Если категория домена положительная или нейтральная -> Разрешен
	// 3. Если категория негативная:
	//    - domain.ActionID == 0 (решение не принято) -> Ожидает
	//    - domain.ActionID != 0 (решение принято):
	//      - action.Action == "Заблокировано" -> Аномалия (доступ к заблокированному домену)
	//      - иначе -> Разрешен
	if log.Action == "blocked" {
		return models.StatusBlocked, nil
	} else {
		// Проверяем категорию домена
		cat, err := uc.q.GetCategoryByID(ctx, uint(domain.CategoryID))
		if err != nil {
			return models.StatusAllowed, fmt.Errorf("failed to get category for domain: %w", err)
		}
		// если категория положительная или нейтральная -> разрешенный
		if cat.Type != models.CategoryTypeNegative {
			return models.StatusAllowed, nil
		} else { // Если категория негативная, смотрим, было ли совершено действия над доменом
			if domain.ActionID == 0 { // решение еще не принято -> ожидает
				return models.StatusPending, nil
			} else {
				// Проверяем, какое действие было выбрано для домена
				action, err := uc.q.GetActionByDomainID(ctx, domain.ID)
				if err != nil {
					return models.StatusAllowed, fmt.Errorf("failed to get action for domain: %w", err)
				}
				// если принято решение заблокировать домен, но все равно происходит обращение к домену -> аномалия
				if action != nil && action.Action == "Заблокировано" { // TODO: добавить структуру сюда вместо жестко прописанного поля
					return models.StatusAnomaly, nil
				} else {
					return models.StatusAllowed, nil
				}
			}
		}

	}
}

// TODO: получение типа сессии VPN
func (uc *UseCase) processSessionType(domain *models.Domain) (models.SessionType, error) {
	// Получаем категорию домена
	domainCategory, err := models.ParseContentCategoryByID(uint(domain.CategoryID))
	if err != nil {
		return models.SessionTypeAllowed, err
	}
	// SessionType = blocked - только тогда, когда у домена отрицательная категория и по нему принято решение заблокировать его
	if domainCategory.Type == models.CategoryTypeNegative && domain.ActionID != 0 {
		return models.SessionTypeBlocked, nil
	}
	// В остальных случаях - будет allowed
	return models.SessionTypeAllowed, nil
}
