package utils

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// ExtractDomain извлекает домен из строки
func ExtractDomain(input string) (string, error) {
	fullURL, err := ExtractURL(input)
	if err != nil {
		return "", err
	}

	// Парсим URL для извлечения домена
	parsedURL, err := url.Parse(fullURL)
	if err != nil {
		return "", fmt.Errorf("ошибка парсинга URL: %w", err)
	}

	domain := parsedURL.Hostname()
	if domain == "" {
		return "", fmt.Errorf("не удалось извлечь домен из URL")
	}

	return domain, nil
}

// ExtractURL извлекает полный URL из строки
func ExtractURL(input string) (string, error) {
	// Регулярное выражение для поиска URL (учитывает GET/POST и другие методы)
	re := regexp.MustCompile(`(?:GET|POST|PUT|DELETE|PATCH|HEAD|OPTIONS)\s+(https?://[^\s\)]+)`)
	matches := re.FindStringSubmatch(input)

	if len(matches) < 2 {
		// Если не нашли с методом, ищем просто URL
		re = regexp.MustCompile(`(https?://[^\s\)]+)`)
		matches = re.FindStringSubmatch(input)
		if len(matches) == 0 {
			return "", fmt.Errorf("URL не найден в строке")
		}
	}

	if len(matches) >= 2 {
		return matches[1], nil // берем первую группу (URL)
	} else {
		return matches[0], nil // берем полное совпадение
	}
}

// Определяет тип протокола из записи URL
func ExtractProto(input string) (*string, error) {
	proto := ""
	url, err := ExtractURL(input)
	if err != nil {
		return nil, fmt.Errorf("не удалось извлечь протокол: %w", err)
	}

	// Извлекаем протокол из URL
	if strings.HasPrefix(url, "https://") {
		proto = "https"
	} else if strings.HasPrefix(url, "http://") {
		proto = "http"
	} else {
		return nil, fmt.Errorf("неизвестный протокол в URL: %s", url)
	}
	return &proto, nil
}
