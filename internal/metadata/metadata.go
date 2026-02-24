package metadata

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Meta holds extracted page metadata.
type Meta struct {
	Title       string
	Description string
	OG          map[string]string
}

// Extract parses HTML and returns page metadata including title, description, and Open Graph tags.
func Extract(html string) *Meta {
	m := &Meta{OG: make(map[string]string)}

	if html == "" {
		return m
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return m
	}

	m.Title = strings.TrimSpace(doc.Find("title").First().Text())

	doc.Find("meta").Each(func(_ int, s *goquery.Selection) {
		name, _ := s.Attr("name")
		property, _ := s.Attr("property")
		content, _ := s.Attr("content")

		if strings.EqualFold(name, "description") {
			m.Description = content
		}

		if strings.HasPrefix(property, "og:") {
			m.OG[property] = content
		}
	})

	if m.Title == "" {
		if ogTitle, ok := m.OG["og:title"]; ok {
			m.Title = ogTitle
		}
	}

	if m.Description == "" {
		if ogDesc, ok := m.OG["og:description"]; ok {
			m.Description = ogDesc
		}
	}

	return m
}
