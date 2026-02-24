package filetype

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

// ConvertODT converts OpenDocument Text (.odt) bytes to markdown.
func ConvertODT(data []byte, filename string) (string, error) {
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return "", fmt.Errorf("odt open zip: %w", err)
	}

	var contentXML []byte
	for _, f := range zr.File {
		if f.Name == "content.xml" {
			rc, err := f.Open()
			if err != nil {
				return "", fmt.Errorf("odt open content.xml: %w", err)
			}
			contentXML, err = io.ReadAll(rc)
			rc.Close()
			if err != nil {
				return "", fmt.Errorf("odt read content.xml: %w", err)
			}
			break
		}
	}

	if contentXML == nil {
		return "", fmt.Errorf("odt: content.xml not found")
	}

	text := extractODTText(contentXML)

	if filename == "" {
		filename = "document.odt"
	}

	var md strings.Builder
	md.WriteString(fmt.Sprintf("# %s\n\n", filename))
	md.WriteString(text)
	md.WriteString("\n")

	return strings.TrimSpace(md.String()), nil
}

// extractODTText walks the ODF XML and extracts paragraphs and headings.
func extractODTText(data []byte) string {
	decoder := xml.NewDecoder(bytes.NewReader(data))
	var result strings.Builder
	var inParagraph, inHeading bool
	var lineBuffer strings.Builder

	for {
		tok, err := decoder.Token()
		if err != nil {
			break
		}

		switch t := tok.(type) {
		case xml.StartElement:
			local := t.Name.Local
			switch local {
			case "p":
				inParagraph = true
				lineBuffer.Reset()
			case "h":
				inHeading = true
				lineBuffer.Reset()
			case "tab":
				lineBuffer.WriteString("\t")
			case "s":
				lineBuffer.WriteString(" ")
			case "line-break":
				lineBuffer.WriteString("\n")
			}
		case xml.EndElement:
			local := t.Name.Local
			switch local {
			case "p":
				if inParagraph {
					result.WriteString(lineBuffer.String())
					result.WriteString("\n\n")
					inParagraph = false
				}
			case "h":
				if inHeading {
					result.WriteString("## " + lineBuffer.String())
					result.WriteString("\n\n")
					inHeading = false
				}
			}
		case xml.CharData:
			if inParagraph || inHeading {
				lineBuffer.Write(t)
			}
		}
	}

	return strings.TrimSpace(result.String())
}
