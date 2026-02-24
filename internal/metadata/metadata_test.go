package metadata

import (
	"testing"
)

func TestExtract_TitleAndDescription(t *testing.T) {
	html := `<html><head>
		<title>Test Page</title>
		<meta name="description" content="A test description">
	</head><body></body></html>`

	m := Extract(html)
	if m.Title != "Test Page" {
		t.Errorf("expected title 'Test Page', got %q", m.Title)
	}
	if m.Description != "A test description" {
		t.Errorf("expected description 'A test description', got %q", m.Description)
	}
}

func TestExtract_OGTags(t *testing.T) {
	html := `<html><head>
		<meta property="og:title" content="OG Title">
		<meta property="og:description" content="OG Desc">
		<meta property="og:image" content="https://example.com/img.png">
	</head><body></body></html>`

	m := Extract(html)
	if m.OG["og:title"] != "OG Title" {
		t.Errorf("expected og:title 'OG Title', got %q", m.OG["og:title"])
	}
	if m.OG["og:image"] != "https://example.com/img.png" {
		t.Errorf("expected og:image, got %q", m.OG["og:image"])
	}
}

func TestExtract_FallbackToOG(t *testing.T) {
	html := `<html><head>
		<meta property="og:title" content="Fallback Title">
		<meta property="og:description" content="Fallback Desc">
	</head><body></body></html>`

	m := Extract(html)
	if m.Title != "Fallback Title" {
		t.Errorf("expected fallback title from og:title, got %q", m.Title)
	}
	if m.Description != "Fallback Desc" {
		t.Errorf("expected fallback description from og:description, got %q", m.Description)
	}
}

func TestExtract_EmptyHTML(t *testing.T) {
	m := Extract("")
	if m.Title != "" || m.Description != "" {
		t.Error("expected empty meta for empty HTML")
	}
	if m.OG == nil {
		t.Error("expected non-nil OG map")
	}
}

func TestExtract_InvalidHTML(t *testing.T) {
	m := Extract("<not valid>>>")
	if m.OG == nil {
		t.Error("expected non-nil OG map even for invalid HTML")
	}
}
