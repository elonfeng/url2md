package converter

import (
	"strings"
	"testing"
)

func TestCleanHTML_RemovesScript(t *testing.T) {
	html := `<html><body><script>alert("xss")</script><p>Hello</p></body></html>`
	result := CleanHTML(html)
	if strings.Contains(result, "script") {
		t.Error("expected script tag to be removed")
	}
	if !strings.Contains(result, "Hello") {
		t.Error("expected content to be preserved")
	}
}

func TestCleanHTML_RemovesNav(t *testing.T) {
	html := `<html><body><nav>Menu</nav><p>Content</p></body></html>`
	result := CleanHTML(html)
	if strings.Contains(result, "Menu") {
		t.Error("expected nav content to be removed")
	}
	if !strings.Contains(result, "Content") {
		t.Error("expected main content to be preserved")
	}
}

func TestCleanHTML_RemovesCookieBanner(t *testing.T) {
	html := `<html><body><div class="cookie-consent">Accept</div><p>Article</p></body></html>`
	result := CleanHTML(html)
	if strings.Contains(result, "Accept") {
		t.Error("expected cookie banner to be removed")
	}
}

func TestCleanHTML_RemovesAdElements(t *testing.T) {
	html := `<html><body><div class="ad-container">Buy now</div><p>Real content</p></body></html>`
	result := CleanHTML(html)
	if strings.Contains(result, "Buy now") {
		t.Error("expected ad content to be removed")
	}
	if !strings.Contains(result, "Real content") {
		t.Error("expected real content to be preserved")
	}
}

func TestCleanHTML_RemovesAriaHidden(t *testing.T) {
	html := `<html><body><div aria-hidden="true">Hidden</div><p>Visible</p></body></html>`
	result := CleanHTML(html)
	if strings.Contains(result, "Hidden") {
		t.Error("expected aria-hidden content to be removed")
	}
}

func TestCleanHTML_PreservesArticleContent(t *testing.T) {
	html := `<html><body><article><h1>Title</h1><p>Body text</p></article></body></html>`
	result := CleanHTML(html)
	if !strings.Contains(result, "Title") || !strings.Contains(result, "Body text") {
		t.Error("expected article content to be preserved")
	}
}

func TestCleanHTML_InvalidHTML(t *testing.T) {
	html := "not valid html at all <><><>"
	result := CleanHTML(html)
	if result == "" {
		t.Error("expected some output even for invalid HTML")
	}
}
