package usecase

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

type AntiBotDetector struct {
	phrases   []phraseEntry
	minLength int
}

type phraseEntry struct {
	raw        string
	normalized string
}

func NewAntiBotDetector(phrases []string, minLength int) *AntiBotDetector {
	if len(phrases) == 0 {
		return nil
	}

	entries := make([]phraseEntry, 0, len(phrases))
	seen := make(map[string]struct{}, len(phrases))
	for _, phrase := range phrases {
		raw := strings.TrimSpace(phrase)
		if raw == "" {
			continue
		}
		norm := normalizeForMatch(raw)
		if norm == "" {
			continue
		}
		if _, ok := seen[norm]; ok {
			continue
		}
		seen[norm] = struct{}{}
		entries = append(entries, phraseEntry{raw: raw, normalized: norm})
	}

	if len(entries) == 0 {
		return nil
	}

	if minLength <= 0 {
		minLength = 500
	}

	return &AntiBotDetector{
		phrases:   entries,
		minLength: minLength,
	}
}

func (d *AntiBotDetector) Check(text string) string {
	if d == nil {
		return ""
	}

	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return ""
	}

	if utf8.RuneCountInString(trimmed) >= d.minLength {
		return ""
	}

	normalized := normalizeForMatch(trimmed)
	if normalized == "" {
		return ""
	}

	for _, entry := range d.phrases {
		if strings.Contains(normalized, entry.normalized) {
			return entry.raw
		}
	}

	return ""
}

func normalizeForMatch(value string) string {
	value = strings.ToLower(value)
	var b strings.Builder
	b.Grow(len(value))

	lastWasSpace := true
	for _, r := range value {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			b.WriteRune(r)
			lastWasSpace = false
		case unicode.IsSpace(r):
			if !lastWasSpace {
				b.WriteRune(' ')
				lastWasSpace = true
			}
		default:
			if !lastWasSpace {
				b.WriteRune(' ')
				lastWasSpace = true
			}
		}
	}

	return strings.TrimSpace(b.String())
}
