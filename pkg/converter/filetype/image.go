package filetype

import (
	"context"
	"fmt"
	"strings"
)

// ConvertImage creates a markdown representation for an image file.
// If VisionConfig is provided and configured, it uses Cloudflare Workers AI
// to generate an AI description. Otherwise, it outputs metadata + image embed.
func ConvertImage(ctx context.Context, data []byte, filename string, rawURL string, contentType string, vision *VisionConfig) (string, error) {
	if filename == "" {
		filename = "image"
	}

	// Try AI vision description first
	if vision != nil && vision.IsConfigured() {
		desc, err := DescribeImage(ctx, vision, data, contentType)
		if err == nil && desc != "" {
			var md strings.Builder
			md.WriteString(fmt.Sprintf("# %s\n\n", filename))
			md.WriteString("## Description\n\n")
			md.WriteString(desc)
			md.WriteString("\n\n")
			md.WriteString(fmt.Sprintf("![%s](%s)\n", filename, rawURL))
			return strings.TrimSpace(md.String()), nil
		}
		// Fall through to metadata-only output on error
	}

	// Fallback: metadata + image embed
	var md strings.Builder
	md.WriteString(fmt.Sprintf("# %s\n\n", filename))

	md.WriteString("## Metadata\n\n")
	md.WriteString(fmt.Sprintf("- **File**: %s\n", filename))
	md.WriteString(fmt.Sprintf("- **Size**: %d bytes\n", len(data)))
	if contentType != "" {
		md.WriteString(fmt.Sprintf("- **Type**: %s\n", contentType))
	}
	md.WriteString("\n")

	md.WriteString("## Image\n\n")
	md.WriteString(fmt.Sprintf("![%s](%s)\n", filename, rawURL))

	return strings.TrimSpace(md.String()), nil
}
