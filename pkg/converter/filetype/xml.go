package filetype

import (
	"fmt"
	"strings"
)

// ConvertXML converts an XML file to markdown with syntax highlighting.
func ConvertXML(data []byte, filename string) (string, error) {
	if filename == "" {
		filename = "data.xml"
	}

	var md strings.Builder
	md.WriteString(fmt.Sprintf("# %s\n\n", filename))
	md.WriteString("```xml\n")
	md.WriteString(strings.TrimSpace(string(data)))
	md.WriteString("\n```\n")

	return strings.TrimSpace(md.String()), nil
}
