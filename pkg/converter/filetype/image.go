package filetype

import (
	"fmt"
	"strings"
)

// ConvertImage creates a markdown representation for an image file.
// Since we don't have AI vision, we output the image as a markdown image embed
// with file metadata.
func ConvertImage(data []byte, filename string, rawURL string, contentType string) (string, error) {
	if filename == "" {
		filename = "image"
	}

	var md strings.Builder
	md.WriteString(fmt.Sprintf("# %s\n\n", filename))

	// image metadata
	md.WriteString("## Metadata\n\n")
	md.WriteString(fmt.Sprintf("- **File**: %s\n", filename))
	md.WriteString(fmt.Sprintf("- **Size**: %d bytes\n", len(data)))
	if contentType != "" {
		md.WriteString(fmt.Sprintf("- **Type**: %s\n", contentType))
	}
	md.WriteString("\n")

	// embed as markdown image
	md.WriteString("## Image\n\n")
	md.WriteString(fmt.Sprintf("![%s](%s)\n", filename, rawURL))

	return strings.TrimSpace(md.String()), nil
}
