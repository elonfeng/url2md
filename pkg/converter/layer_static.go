package converter

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/elonfeng/url2md/pkg/converter/filetype"
	"github.com/go-shiori/go-readability"
)

// StaticLayer fetches content via standard HTTP. For HTML pages, it extracts
// with go-readability and converts with html-to-markdown. For other file types
// (PDF, DOCX, XLSX, CSV, images), it uses specialized parsers.
type StaticLayer struct{}

func (l *StaticLayer) Name() string { return "static" }

func (l *StaticLayer) Convert(ctx context.Context, rawURL string, opts *Options) (string, string, error) {
	data, resp, err := l.fetchRaw(ctx, rawURL, opts)
	if err != nil {
		return "", "", err
	}

	ft := filetype.Detect(rawURL, resp, data)
	ct := ""
	if resp != nil {
		ct = resp.Header.Get("Content-Type")
	}

	// Use the final redirect URL for filename extraction when available.
	fnURL := rawURL
	if resp != nil && resp.Request != nil && resp.Request.URL != nil {
		fnURL = resp.Request.URL.String()
	}
	filename := filetype.FilenameFromURL(fnURL)

	switch ft {
	case filetype.TypePDF:
		markdown, err := filetype.ConvertPDF(data, filename)
		return markdown, "", err

	case filetype.TypeDOCX:
		markdown, err := filetype.ConvertDOCX(data, filename)
		return markdown, "", err

	case filetype.TypeXLSX:
		markdown, err := filetype.ConvertXLSX(data, filename)
		return markdown, "", err

	case filetype.TypeXLS:
		markdown, err := filetype.ConvertXLS(data, filename)
		return markdown, "", err

	case filetype.TypeODT:
		markdown, err := filetype.ConvertODT(data, filename)
		return markdown, "", err

	case filetype.TypeCSV:
		markdown, err := filetype.ConvertCSV(data, filename)
		return markdown, "", err

	case filetype.TypeJSON:
		markdown, err := filetype.ConvertJSON(data, filename)
		return markdown, "", err

	case filetype.TypeXML:
		markdown, err := filetype.ConvertXML(data, filename)
		return markdown, "", err

	case filetype.TypeTXT:
		markdown, err := filetype.ConvertTXT(data, filename)
		return markdown, "", err

	case filetype.TypeMD:
		markdown, err := filetype.ConvertMD(data)
		return markdown, "", err

	case filetype.TypeSVG:
		markdown, err := filetype.ConvertImage(ctx, data, filename, rawURL, "image/svg+xml", opts.Vision)
		return markdown, "", err

	case filetype.TypePNG, filetype.TypeJPEG, filetype.TypeGIF, filetype.TypeWEBP:
		markdown, err := filetype.ConvertImage(ctx, data, filename, rawURL, ct, opts.Vision)
		return markdown, "", err
	}

	// default: HTML
	return l.convertHTML(data, rawURL, opts)
}

func (l *StaticLayer) convertHTML(data []byte, rawURL string, opts *Options) (string, string, error) {
	html := string(data)
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

func (l *StaticLayer) fetchRaw(ctx context.Context, rawURL string, opts *Options) ([]byte, *http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Accept", "*/*")
	if opts.UserAgent != "" {
		req.Header.Set("User-Agent", opts.UserAgent)
	}

	client := &http.Client{Timeout: opts.Timeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("fetch: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, resp, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 50<<20)) // 50 MB limit for files
	if err != nil {
		return nil, resp, fmt.Errorf("read body: %w", err)
	}

	return body, resp, nil
}

// stripImages removes markdown image syntax ![alt](url).
func stripImages(text string) string {
	var result strings.Builder
	i := 0
	for i < len(text) {
		if i < len(text)-1 && text[i] == '!' && text[i+1] == '[' {
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
					j++
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
