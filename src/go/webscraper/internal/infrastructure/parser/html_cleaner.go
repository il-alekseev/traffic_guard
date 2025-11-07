package parser

import (
	"regexp"
	"strings"
	"unicode"

	htmlstd "html"

	"github.com/PuerkitoBio/goquery"
)

var (
	urlPattern = regexp.MustCompile(`https?://\S+`)
	tagPattern = regexp.MustCompile(`<[^>]+>`)
)

type HTMLCleaner struct{}

func NewHTMLCleaner() *HTMLCleaner {
	return &HTMLCleaner{}
}

func (c *HTMLCleaner) Clean(raw string) (string, error) {
	if strings.TrimSpace(raw) == "" {
		return "", nil
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(raw))
	if err != nil {
		fallback := c.postProcess(raw)
		if fallback == "" {
			return fallback, err
		}
		return fallback, nil
	}

	doc.Find("script, style, noscript, template").Remove()

	cleaned := c.postProcess(doc.Text())
	if cleaned == "" {
		return cleaned, nil
	}

	return cleaned, nil
}

func (c *HTMLCleaner) CleanFragment(fragment string) (string, error) {
	fragment = strings.TrimSpace(fragment)
	if fragment == "" {
		return "", nil
	}

	reader := strings.NewReader("<div>" + fragment + "</div>")
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		fallback := c.postProcess(fragment)
		if fallback == "" {
			return fallback, err
		}
		return fallback, nil
	}

	doc.Find("script, style, noscript, template").Remove()

	cleaned := c.postProcess(doc.Text())
	if cleaned == "" {
		return cleaned, nil
	}

	return cleaned, nil
}

func (c *HTMLCleaner) postProcess(input string) string {
	cleaned := strings.Join(strings.Fields(input), " ")
	if cleaned == "" {
		return ""
	}

	cleaned = htmlstd.UnescapeString(cleaned)
	cleaned = sanitizeControlRunes(cleaned)
	cleaned = tagPattern.ReplaceAllString(cleaned, " ")
	cleaned = urlPattern.ReplaceAllString(cleaned, " ")
	cleaned = strings.Join(strings.Fields(cleaned), " ")

	return strings.TrimSpace(cleaned)
}

func sanitizeControlRunes(s string) string {
	return strings.Map(func(r rune) rune {
		if r == '\n' || r == '\r' || r == '\t' {
			return ' '
		}
		if unicode.IsControl(r) {
			return -1
		}
		return r
	}, s)
}
