package converter

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var removeSelectors = []string{
	"script", "style", "noscript", "iframe", "svg",
	"nav", "header", "footer", "aside",
	"[class*='cookie']", "[class*='popup']", "[class*='modal']",
	"[class*='ad-']", "[class*='advert']", "[class*='advertisement']",
	"[class*='sidebar']", "[class*='comment']",
	"[class*='share']", "[class*='social']",
	"[role='navigation']", "[role='banner']", "[role='complementary']",
	"[aria-hidden='true']",
}

// CleanHTML removes noisy elements from HTML before conversion.
func CleanHTML(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return html
	}

	for _, sel := range removeSelectors {
		doc.Find(sel).Remove()
	}

	result, err := doc.Html()
	if err != nil {
		return html
	}
	return result
}
