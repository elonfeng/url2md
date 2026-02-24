package converter

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/elonfeng/url2md/internal/metadata"
	"github.com/elonfeng/url2md/internal/token"
)

// Converter converts a URL to Markdown.
type Converter interface {
	Convert(ctx context.Context, url string, opts *Options) (*Result, error)
}

// Layer is a single conversion strategy.
type Layer interface {
	Name() string
	Convert(ctx context.Context, url string, opts *Options) (markdown string, rawHTML string, err error)
}

type converter struct {
	negotiate Layer
	static    Layer
	browser   Layer
}

// New creates a Converter with the three-layer fallback pipeline.
func New() Converter {
	return &converter{
		negotiate: &NegotiateLayer{},
		static:    &StaticLayer{},
		browser:   &BrowserLayer{},
	}
}

func (c *converter) Convert(ctx context.Context, rawURL string, opts *Options) (*Result, error) {
	if opts == nil {
		opts = DefaultOptions()
	}

	layers := c.buildLayers(opts)
	if len(layers) == 0 {
		return nil, fmt.Errorf("no conversion layers configured")
	}

	var lastErr error
	for _, layer := range layers {
		fetchStart := time.Now()
		md, rawHTML, err := layer.Convert(ctx, rawURL, opts)
		fetchTime := time.Since(fetchStart)

		if err != nil {
			lastErr = fmt.Errorf("[%s] %w", layer.Name(), err)
			continue
		}

		convertStart := time.Now()
		meta := metadata.Extract(rawHTML)

		// build final markdown with optional frontmatter and title
		var final strings.Builder

		if opts.Frontmatter && (meta.Title != "" || meta.Description != "" || meta.OG["og:image"] != "") {
			final.WriteString("---\n")
			if meta.Title != "" {
				final.WriteString(fmt.Sprintf("title: %s\n", meta.Title))
			}
			if meta.Description != "" {
				final.WriteString(fmt.Sprintf("description: %s\n", meta.Description))
			}
			if img := meta.OG["og:image"]; img != "" {
				final.WriteString(fmt.Sprintf("image: %s\n", img))
			}
			final.WriteString("---\n\n")
		}

		// auto-prepend # Title if markdown doesn't already start with one
		if meta.Title != "" && !strings.HasPrefix(md, "# ") {
			final.WriteString("# ")
			final.WriteString(meta.Title)
			final.WriteString("\n\n")
		}

		final.WriteString(md)
		finalMd := final.String()

		tokenCount := token.Estimate(finalMd)
		convertTime := time.Since(convertStart)

		result := &Result{
			URL:         rawURL,
			Markdown:    finalMd,
			Title:       meta.Title,
			Description: meta.Description,
			TokenCount:  tokenCount,
			Method:      layer.Name(),
			Metadata:    meta.OG,
			FetchTime:   fetchTime,
			ConvertTime: convertTime,
		}
		return result, nil
	}

	return nil, fmt.Errorf("all layers failed: %w", lastErr)
}

func (c *converter) buildLayers(opts *Options) []Layer {
	switch opts.Method {
	case "negotiate":
		return []Layer{c.negotiate}
	case "static":
		return []Layer{c.static}
	case "browser":
		return []Layer{c.browser}
	default: // "auto"
		layers := []Layer{c.negotiate, c.static}
		if opts.EnableBrowser {
			layers = append(layers, c.browser)
		}
		return layers
	}
}
