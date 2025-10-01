package parser

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/net/html"
	htmlstd "html"
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
	node, err := html.Parse(strings.NewReader(raw))
	if err != nil {
		return "", err
	}

	var builder strings.Builder

	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "script", "style", "noscript", "template":
				return
			}
		}

		if n.Type == html.TextNode {
			text := strings.TrimSpace(n.Data)
			if text != "" {
				if builder.Len() > 0 {
					builder.WriteByte(' ')
				}
				builder.WriteString(text)
			}
		}

		for child := n.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
		}
	}

	walk(node)

	cleaned := strings.Join(strings.Fields(builder.String()), " ")
	if cleaned == "" {
		return cleaned, nil
	}

	cleaned = htmlstd.UnescapeString(cleaned)
	cleaned = sanitizeControlRunes(cleaned)
	cleaned = tagPattern.ReplaceAllString(cleaned, "")
	cleaned = urlPattern.ReplaceAllString(cleaned, "")
	cleaned = strings.Join(strings.Fields(cleaned), " ")

	return strings.TrimSpace(cleaned), nil
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
