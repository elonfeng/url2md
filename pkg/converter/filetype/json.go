package filetype

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

// ConvertJSON converts a JSON file to markdown with syntax highlighting.
func ConvertJSON(data []byte, filename string) (string, error) {
	if filename == "" {
		filename = "data.json"
	}

	// Prettify JSON
	var pretty bytes.Buffer
	if err := json.Indent(&pretty, bytes.TrimSpace(data), "", "  "); err != nil {
		// If invalid JSON, use raw content
		pretty.Reset()
		pretty.Write(data)
	}

	var md strings.Builder
	md.WriteString(fmt.Sprintf("# %s\n\n", filename))
	md.WriteString("```json\n")
	md.WriteString(pretty.String())
	md.WriteString("\n```\n")

	return strings.TrimSpace(md.String()), nil
}
