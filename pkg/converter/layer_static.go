package converter

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/go-shiori/go-readability"
)

// StaticLayer fetches HTML via standard HTTP, extracts with go-readability, converts with html-to-markdown.
type StaticLayer struct{}

func (l *StaticLayer) Name() string { return "static" }

func (l *StaticLayer) Convert(ctx context.Context, rawURL string, opts *Options) (string, string, error) {
	html, err := l.fetch(ctx, rawURL, opts)
	if err != nil {
		return "", "", err
	}

	cleaned := CleanHTML(html)

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", "", fmt.Errorf("parse url: %w", err)
	}

	article, err := readability.FromReader(strings.NewReader(cleaned), parsedURL)
	if err != nil {
		return "", "", fmt.Errorf("readability: %w", err)
	}

	markdown, err := md.ConvertString(article.Content)
	if err != nil {
		return "", "", fmt.Errorf("html-to-markdown: %w", err)
	}

	if !opts.RetainImages {
		markdown = stripImages(markdown)
	}

	return strings.TrimSpace(markdown), html, nil
}

func (l *StaticLayer) fetch(ctx context.Context, rawURL string, opts *Options) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	if opts.UserAgent != "" {
		req.Header.Set("User-Agent", opts.UserAgent)
	}

	client := &http.Client{Timeout: opts.Timeout}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("fetch: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 10<<20))
	if err != nil {
		return "", fmt.Errorf("read body: %w", err)
	}

	return string(body), nil
}

// stripImages removes markdown image syntax ![alt](url).
func stripImages(text string) string {
	var result strings.Builder
	i := 0
	for i < len(text) {
		if i < len(text)-1 && text[i] == '!' && text[i+1] == '[' {
			// skip ![...](...) pattern
			j := i + 2
			depth := 1
			for j < len(text) && depth > 0 {
				if text[j] == '[' {
					depth++
				} else if text[j] == ']' {
					depth--
				}
				j++
			}
			if j < len(text) && text[j] == '(' {
				j++
				for j < len(text) && text[j] != ')' {
					j++
				}
				if j < len(text) {
					j++ // skip closing )
				}
			}
			i = j
			continue
		}
		result.WriteByte(text[i])
		i++
	}
	return result.String()
}
