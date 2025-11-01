package utils

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

type URLInfo struct {
	Proto  string
	Domain string
	URL    string
}

func ParseRawURL(input string) (URLInfo, error) {
	info := URLInfo{}

	// Извлекаем полный URL
	fullURL, err := extractURL(input)
	if err != nil {
		return info, fmt.Errorf("failed to extract URL: %w", err)
	}

	// Очищаем URL от недопустимых символов
	cleanedURL := cleanURL(fullURL)
	info.URL = cleanedURL

	// Извлекаем протокол
	proto, err := extractProto(cleanedURL)
	if err != nil {
		return info, fmt.Errorf("failed to extract protocol: %w", err)
	}
	info.Proto = proto

	// Извлекаем домен
	domain, err := extractDomain(cleanedURL)
	if err != nil {
		return info, fmt.Errorf("failed to extract domain: %w", err)
	}
	info.Domain = domain

	return info, nil
}

// cleanURL удаляет недопустимые символы из URL
func cleanURL(rawURL string) string {
	// Удаляем управляющие символы (ASCII 0-31 и 127)
	var result strings.Builder
	for _, r := range rawURL {
		if r > 31 && r != 127 && r != ' ' {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// extractURL извлекает полный URL из строки
func extractURL(input string) (string, error) {
	// Регулярное выражение для поиска URL (учитывает GET/POST и другие методы)
	re := regexp.MustCompile(`(?:GET|POST|PUT|DELETE|PATCH|HEAD|OPTIONS)\s+(https?://[^\s\)]+)`)
	matches := re.FindStringSubmatch(input)

	if len(matches) < 2 {
		// Если не нашли с методом, ищем просто URL
		re = regexp.MustCompile(`(https?://[^\s\)]+)`)
		matches = re.FindStringSubmatch(input)
		if len(matches) == 0 {
			return "", fmt.Errorf("URL not found in string")
		}
	}

	if len(matches) >= 2 {
		return matches[1], nil // берем первую группу (URL)
	} else {
		return matches[0], nil // берем полное совпадение
	}
}

// extractDomain извлекает домен из URL строки
func extractDomain(fullURL string) (string, error) {
	// Парсим URL для извлечения домена
	parsedURL, err := url.Parse(fullURL)
	if err != nil {
		return "", fmt.Errorf("URL parsing error: %w", err)
	}

	domain := parsedURL.Hostname()
	if domain == "" {
		return "", fmt.Errorf("failed to extract domain from URL")
	}

	return domain, nil
}

// extractProto определяет тип протокола из URL
func extractProto(fullURL string) (string, error) {
	// Извлекаем протокол из URL
	if strings.HasPrefix(fullURL, "https://") {
		return "https", nil
	} else if strings.HasPrefix(fullURL, "http://") {
		return "http", nil
	} else {
		return "", fmt.Errorf("unknown protocol in URL: %s", fullURL)
	}
}
