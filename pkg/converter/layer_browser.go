package converter

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/chromedp/chromedp"
	md "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/go-shiori/go-readability"
)

// BrowserLayer uses headless Chrome to render JavaScript-heavy pages.
type BrowserLayer struct{}

func (l *BrowserLayer) Name() string { return "browser" }

func (l *BrowserLayer) Convert(ctx context.Context, rawURL string, opts *Options) (string, string, error) {
	allocCtx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	var html string
	err := chromedp.Run(allocCtx,
		chromedp.Navigate(rawURL),
		chromedp.WaitReady("body"),
		chromedp.OuterHTML("html", &html),
	)
	if err != nil {
		return "", "", fmt.Errorf("chromedp: %w", err)
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
