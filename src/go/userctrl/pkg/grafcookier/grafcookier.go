package grafcookier

import (
	"crypto/tls"
	"userctrl/pkg/slogger/wsl"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/publicsuffix"
)

type GrafCookierInterface interface {
	GetCookies(username, password string) (*GrafCookies, error)
}

type GrafCookier struct {
	grafanaURL      string
	grafanaAuthPath string
	logger          slog.Logger
}

type GrafCookies struct {
	GrafanaSession       string `json:"grafana_session"`
	GrafanaSessionExpiry string `json:"grafana_session_expiry"`
}

func New(url, path string, l slog.Logger) (*GrafCookier, error) {
	return &GrafCookier{
		grafanaURL:      url,
		grafanaAuthPath: path,
		logger:          l,
	}, nil
}

func (g *GrafCookier) newClient() (*http.Client, error) {
	tr := http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	// Создаем клиент с поддержкой cookies
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		g.logger.Error(err.Error())
		return nil, err
	}
	client := &http.Client{
		Jar:       jar,
		Timeout:   30 * time.Second,
		Transport: &tr,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Не следуем автоматически за редиректами
			return http.ErrUseLastResponse
		},
	}
	return client, nil
}

// GetCookies выполняет полный OAuth-флоу для получения сессионных кук Grafana
// через промежуточную авторизацию в Keycloak.
//
// Метод осуществляет последовательную авторизацию в системе, начиная с запроса
// к Grafana OAuth endpoint, перехода на страницу Keycloak, отправки учетных данных
// и завершения OAuth-флоу с получением финальных сессионных кук.
//
// Параметры:
//
//	username: имя пользователя для авторизации
//	password: пароль пользователя
//
// Возвращаемые значения:
//
//	*GrafCookies: структура с полученными сессионными куками
//	error: ошибка при выполнении процесса авторизации
func (g *GrafCookier) GetCookies(username, password string) (*GrafCookies, error) {
	log := g.logger.With(wsl.Label("method", "GetCookies"))
	// 1. Начальный запрос к Grafana OAuth endpoint
	grafanaOAuthURL := g.grafanaURL + g.grafanaAuthPath
	client, err := g.newClient()
	if err != nil {
		err = fmt.Errorf("new client: %w", err)
		log.Error("unable create client", wsl.Err(err))
		return nil, err
	}
	defer client.CloseIdleConnections()

	resp, err := client.Get(grafanaOAuthURL)
	if err != nil {
		log.Error("Ошибка при запросе к Grafana OAuth:", wsl.Err(err))
		return nil, err
	}
	defer resp.Body.Close()

	// 2. Переходим по редиректу на страницу авторизации Keycloak
	keycloakAuthURL := resp.Header.Get("Location")
	if keycloakAuthURL == "" {
		log.Error("Не получен URL для редиректа на Keycloak")
		return nil, err
	}

	resp, err = client.Get(keycloakAuthURL)
	if err != nil {
		log.Error("Ошибка при запросе к Keycloak:", wsl.Err(err))
		return nil, err
	}
	defer resp.Body.Close()

	// 3. Извлекаем форму авторизации и необходимые параметры
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("Ошибка чтения тела ответа:", wsl.Err(err))
		return nil, err
	}

	// Парсим форму для получения URL и скрытых полей
	formAction, formData, err := parseKeycloakLoginForm(string(body))
	if err != nil {
		log.Error("Ошибка парсинга формы:", wsl.Err(err))
		return nil, err
	}

	// 4. Отправляем данные для авторизации
	formData.Set("username", username)
	formData.Set("password", password)
	formData.Set("credentialId", "")

	req, err := http.NewRequest("POST", formAction, strings.NewReader(formData.Encode()))
	if err != nil {
		log.Error("Ошибка создания запроса:", wsl.Err(err))
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err = client.Do(req)
	if err != nil {
		log.Error("Ошибка при авторизации:", wsl.Err(err))
		return nil, err
	}
	defer resp.Body.Close()

	// 5. Получаем код авторизации от Keycloak
	if resp.StatusCode != http.StatusFound {
		err = fmt.Errorf("неожиданный статус код после авторизации: StatusCode=%d", resp.StatusCode)
		log.Error(err.Error())
		return nil, err
	}

	// 6. Завершаем OAuth flow в Grafana
	grafanaCallbackURL := resp.Header.Get("Location")
	if grafanaCallbackURL == "" {
		err = fmt.Errorf("не получен URL для обратного вызова Grafana: StatusCode=%d", resp.StatusCode)
		log.Error(err.Error())
		return nil, err
	}

	resp, err = client.Get(grafanaCallbackURL)
	if err != nil {
		log.Error("Ошибка при завершении OAuth", wsl.Err(err))
		return nil, err
	}
	defer resp.Body.Close()

	// 7. Получаем сессионную куку Grafana
	var cookies GrafCookies
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "grafana_session" {
			cookies.GrafanaSession = cookie.Value
		} else if cookie.Name == "grafana_session_expiry" {
			cookies.GrafanaSessionExpiry = cookie.Value
		}
	}

	if cookies.GrafanaSession == "" {
		err = fmt.Errorf("не удалось получить сессионную куку Grafana")
		log.Error(err.Error())
		return &cookies, err
	}
	log.Debug("Успешная авторизация! Grafana session cookie", wsl.Info(cookies.GrafanaSession))

	return &cookies, nil
}

// parseKeycloakLoginForm анализирует HTML-форму входа Keycloak и извлекает необходимые данные для авторизации.
//
// Функция парсит HTML-содержимое формы входа, находя URL действия формы и все скрытые поля,
// необходимые для успешной авторизации в системе.
//
// Параметры:
//
//	html - HTML-содержимое страницы формы входа
//
// Возвращаемые значения:
//
//	formAction - URL, на который отправляется форма
//	formData - карта скрытых полей формы (url.Values)
//	error - ошибка при парсинге (если произошла)
func parseKeycloakLoginForm(html string) (string, url.Values, error) {
	// Ищем URL формы
	actionStart := strings.Index(html, `action="`)
	if actionStart == -1 {
		return "", nil, fmt.Errorf("не найден атрибут action в форме")
	}
	actionEnd := strings.Index(html[actionStart+8:], `"`)
	if actionEnd == -1 {
		return "", nil, fmt.Errorf("не найден закрывающий атрибут action")
	}
	formAction := html[actionStart+8 : actionStart+8+actionEnd]

	// Парсим скрытые поля формы
	formData := url.Values{}
	start := 0
	for {
		// Ищем input fields
		inputStart := strings.Index(html[start:], `<input type="hidden"`)
		if inputStart == -1 {
			break
		}
		inputStart += start

		// Извлекаем name и value
		nameStart := strings.Index(html[inputStart:], `name="`)
		if nameStart == -1 {
			start = inputStart + 1
			continue
		}
		nameEnd := strings.Index(html[inputStart+nameStart+6:], `"`)
		if nameEnd == -1 {
			start = inputStart + 1
			continue
		}
		name := html[inputStart+nameStart+6 : inputStart+nameStart+6+nameEnd]

		valueStart := strings.Index(html[inputStart:], `value="`)
		if valueStart == -1 {
			start = inputStart + 1
			continue
		}
		valueEnd := strings.Index(html[inputStart+valueStart+7:], `"`)
		if valueEnd == -1 {
			start = inputStart + 1
			continue
		}
		value := html[inputStart+valueStart+7 : inputStart+valueStart+7+valueEnd]

		formData.Set(name, value)
		start = inputStart + valueStart + 7 + valueEnd
	}

	return formAction, formData, nil
}
