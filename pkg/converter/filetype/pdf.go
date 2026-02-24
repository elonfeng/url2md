package filetype

import (
	"bytes"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/ledongthuc/pdf"
)

// ConvertPDF extracts text from PDF bytes and returns markdown.
func ConvertPDF(data []byte, filename string) (string, error) {
	reader, err := pdf.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return "", fmt.Errorf("pdf open: %w", err)
	}

	var md strings.Builder

	// filename as title
	if filename == "" {
		filename = "document.pdf"
	}
	md.WriteString(fmt.Sprintf("# %s\n\n", filename))

	numPages := reader.NumPage()

	md.WriteString("## Contents\n\n")
	for i := 1; i <= numPages; i++ {
		p := reader.Page(i)
		if p.V.IsNull() {
			continue
		}

		text, err := p.GetPlainText(nil)
		if err != nil {
			continue
		}

		text = strings.TrimSpace(text)
		if text == "" {
			continue
		}

		md.WriteString(fmt.Sprintf("### Page %d\n\n", i))
		md.WriteString(text)
		md.WriteString("\n\n")
	}

	return strings.TrimSpace(md.String()), nil
}

// FilenameFromURL extracts a filename from a URL, decoding percent-encoding.
func FilenameFromURL(rawURL string) string {
	p := strings.Split(rawURL, "?")[0]
	name := path.Base(p)
	if decoded, err := url.PathUnescape(name); err == nil {
		return decoded
	}
	return name
}
