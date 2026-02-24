package converter

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// NegotiateLayer attempts content negotiation by requesting text/markdown directly.
type NegotiateLayer struct{}

func (l *NegotiateLayer) Name() string { return "negotiate" }

func (l *NegotiateLayer) Convert(ctx context.Context, url string, opts *Options) (string, string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", "", fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Accept", "text/markdown, text/x-markdown")
	if opts.UserAgent != "" {
		req.Header.Set("User-Agent", opts.UserAgent)
	}

	client := &http.Client{Timeout: opts.Timeout}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("fetch: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	ct := resp.Header.Get("Content-Type")
	if !strings.Contains(ct, "text/markdown") && !strings.Contains(ct, "text/x-markdown") {
		return "", "", fmt.Errorf("server returned %q, not markdown", ct)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 10<<20)) // 10 MB limit
	if err != nil {
		return "", "", fmt.Errorf("read body: %w", err)
	}

	md := string(body)
	return md, "", nil // no raw HTML in negotiate path
}
