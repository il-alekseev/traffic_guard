package grafanaclient

import (
	"context"
	"encoding/json"
	"userctrl/pkg/slogger/wsl"
	"fmt"
	"net/http"
)

// Organization представляет структуру организации Grafana
type Organization struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// GetAllOrganizations получает список всех организаций Grafana
//
// Возвращает:
//
//	[]Organization - список организаций
//	error - ошибку при выполнении запроса (nil если успешно)
func (gc *GrafanaClient) GetAllOrganizations() ([]Organization, error) {
	log := gc.logger.With(wsl.Label("method", "GetAllOrganizations"))
	// Создаем HTTP клиент с конфигурацией транспорта
	client := &http.Client{
		Transport: gc.createHTTPTransport(),
	}

	// Формируем URL запроса
	url := fmt.Sprintf("%s://%s/api/orgs", gc.config.Schemes[0], gc.config.Host)

	// Создаем новый запрос
	req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	pass, ok := gc.config.BasicAuth.Password()
	if !ok {
		return nil, fmt.Errorf("gc.config.BasicAuth.Password() - not set up")
	}
	// Добавляем Basic Auth
	req.SetBasicAuth(gc.config.BasicAuth.Username(), pass)

	// Выполняем запрос
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем статус код
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	client.CloseIdleConnections()

	// Декодируем ответ
	var orgs []Organization
	if err := json.NewDecoder(resp.Body).Decode(&orgs); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	log.Debug("count=%d", wsl.Int("count", len(orgs)))

	return orgs, nil
}

// createHTTPTransport создает HTTP транспорт с TLS конфигурацией
func (gc *GrafanaClient) createHTTPTransport() *http.Transport {
	return &http.Transport{
		TLSClientConfig: gc.config.TLSConfig,
	}
}
