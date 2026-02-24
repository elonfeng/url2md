package converter

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestStripImages(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple image",
			input:    "before ![alt](http://img.png) after",
			expected: "before  after",
		},
		{
			name:     "no images",
			input:    "just text",
			expected: "just text",
		},
		{
			name:     "multiple images",
			input:    "![a](1.png) text ![b](2.png)",
			expected: " text ",
		},
		{
			name:     "image with nested brackets",
			input:    "![alt [text]](http://img.png)",
			expected: "",
		},
		{
			name:     "link not image",
			input:    "[link](http://example.com)",
			expected: "[link](http://example.com)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stripImages(tt.input)
			if result != tt.expected {
				t.Errorf("got %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestConverter_NegotiateLayer(t *testing.T) {
	mdContent := "# Hello\n\nThis is markdown."
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accept := r.Header.Get("Accept")
		if strings.Contains(accept, "text/markdown") {
			w.Header().Set("Content-Type", "text/markdown")
			fmt.Fprint(w, mdContent)
		} else {
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprint(w, "<html><body><p>HTML fallback</p></body></html>")
		}
	}))
	defer srv.Close()

	c := New()
	opts := DefaultOptions()
	opts.Method = "negotiate"

	result, err := c.Convert(context.Background(), srv.URL, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Method != "negotiate" {
		t.Errorf("expected method negotiate, got %s", result.Method)
	}

	if result.Markdown != mdContent {
		t.Errorf("expected markdown %q, got %q", mdContent, result.Markdown)
	}
}

func TestConverter_StaticLayer(t *testing.T) {
	htmlContent := `<!DOCTYPE html>
<html><head><title>Test</title></head>
<body>
<article>
<h1>Test Article</h1>
<p>This is a test article with enough content to pass readability extraction threshold for the go-readability library to work properly.</p>
<p>Second paragraph with additional meaningful content that helps the extraction algorithm determine this is real article content.</p>
<p>Third paragraph providing even more substance to the article body for reliable extraction.</p>
</article>
</body></html>`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, htmlContent)
	}))
	defer srv.Close()

	c := New()
	opts := DefaultOptions()
	opts.Method = "static"

	result, err := c.Convert(context.Background(), srv.URL, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Method != "static" {
		t.Errorf("expected method static, got %s", result.Method)
	}

	if result.Title != "Test" {
		t.Errorf("expected title 'Test', got %q", result.Title)
	}

	if result.Markdown == "" {
		t.Error("expected non-empty markdown")
	}

	if result.TokenCount == 0 {
		t.Error("expected non-zero token count")
	}
}

func TestConverter_AutoFallback(t *testing.T) {
	htmlContent := `<!DOCTYPE html>
<html><head><title>Fallback Test</title></head>
<body>
<article>
<h1>Fallback Article</h1>
<p>This article content should be extracted after negotiate fails. It has enough text to pass the readability threshold for extraction.</p>
<p>Adding more paragraphs to ensure the readability algorithm can properly identify this as article content.</p>
<p>Third paragraph with substantial content for reliable extraction by the go-readability library.</p>
</article>
</body></html>`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// always return HTML, negotiate should fail
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, htmlContent)
	}))
	defer srv.Close()

	c := New()
	opts := DefaultOptions()
	opts.Method = "auto"

	result, err := c.Convert(context.Background(), srv.URL, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// negotiate fails, should fall back to static
	if result.Method != "static" {
		t.Errorf("expected fallback to static, got %s", result.Method)
	}
}

func TestConverter_AllLayersFail(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	c := New()
	opts := DefaultOptions()
	opts.Method = "auto"

	_, err := c.Convert(context.Background(), srv.URL, opts)
	if err == nil {
		t.Fatal("expected error when all layers fail")
	}

	if !strings.Contains(err.Error(), "all layers failed") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestConverter_NilOptions(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `<html><head><title>Default</title></head><body>
		<article><p>Content with enough text for readability to detect it as a real article body.</p>
		<p>Second paragraph of content to help readability algorithm.</p>
		<p>Third paragraph of content for reliable extraction.</p></article></body></html>`)
	}))
	defer srv.Close()

	c := New()
	result, err := c.Convert(context.Background(), srv.URL, nil)
	if err != nil {
		t.Fatalf("unexpected error with nil options: %v", err)
	}
	if result.Markdown == "" {
		t.Error("expected non-empty markdown with nil options")
	}
}
