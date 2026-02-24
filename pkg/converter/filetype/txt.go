package filetype

import (
	"fmt"
	"strings"
)

// ConvertTXT converts a plain text file to markdown.
func ConvertTXT(data []byte, filename string) (string, error) {
	if filename == "" {
		filename = "document.txt"
	}

	var md strings.Builder
	md.WriteString(fmt.Sprintf("# %s\n\n", filename))
	md.WriteString(strings.TrimSpace(string(data)))
	md.WriteString("\n")

	return strings.TrimSpace(md.String()), nil
}
