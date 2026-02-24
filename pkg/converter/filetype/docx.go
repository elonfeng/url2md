package filetype

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/fumiama/go-docx"
)

// ConvertDOCX extracts text from DOCX bytes and returns markdown.
func ConvertDOCX(data []byte, filename string) (string, error) {
	r, err := docx.Parse(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return "", fmt.Errorf("docx parse: %w", err)
	}

	var md strings.Builder

	if filename == "" {
		filename = "document.docx"
	}
	md.WriteString(fmt.Sprintf("# %s\n\n", filename))

	for _, item := range r.Document.Body.Items {
		switch v := item.(type) {
		case *docx.Paragraph:
			text := extractParagraphText(v)
			if text == "" {
				continue
			}

			// detect heading style
			if v.Properties != nil && v.Properties.Style != nil {
				style := v.Properties.Style.Val
				switch {
				case strings.Contains(style, "Heading1") || strings.HasPrefix(style, "heading 1") || style == "1":
					md.WriteString("# " + text + "\n\n")
					continue
				case strings.Contains(style, "Heading2") || strings.HasPrefix(style, "heading 2") || style == "2":
					md.WriteString("## " + text + "\n\n")
					continue
				case strings.Contains(style, "Heading3") || strings.HasPrefix(style, "heading 3") || style == "3":
					md.WriteString("### " + text + "\n\n")
					continue
				}
			}

			md.WriteString(text + "\n\n")

		case *docx.Table:
			md.WriteString(convertTable(v))
		}
	}

	return strings.TrimSpace(md.String()), nil
}

func extractParagraphText(p *docx.Paragraph) string {
	var parts []string
	for _, child := range p.Children {
		switch v := child.(type) {
		case *docx.Run:
			text := extractRunText(v)
			if text == "" {
				continue
			}
			// apply inline formatting
			if v.RunProperties != nil {
				if v.RunProperties.Bold != nil {
					text = "**" + text + "**"
				}
				if v.RunProperties.Italic != nil {
					text = "_" + text + "_"
				}
				if v.RunProperties.Strike != nil {
					text = "~~" + text + "~~"
				}
			}
			parts = append(parts, text)
		case *docx.Hyperlink:
			var linkText string
			for _, hc := range v.Run.Children {
				if t, ok := hc.(*docx.Text); ok {
					linkText += t.Text
				}
			}
			if linkText != "" {
				parts = append(parts, linkText)
			}
		}
	}
	return strings.Join(parts, "")
}

func extractRunText(r *docx.Run) string {
	var parts []string
	for _, child := range r.Children {
		switch v := child.(type) {
		case *docx.Text:
			parts = append(parts, v.Text)
		case *docx.Tab:
			parts = append(parts, "\t")
		case *docx.BarterRabbet:
			parts = append(parts, "\n")
		}
	}
	return strings.Join(parts, "")
}

func convertTable(t *docx.Table) string {
	if len(t.TableRows) == 0 {
		return ""
	}

	var md strings.Builder
	md.WriteString("\n")

	for i, row := range t.TableRows {
		md.WriteString("|")
		for _, cell := range row.TableCells {
			cellText := ""
			for _, item := range cell.Paragraphs {
				cellText += extractParagraphText(item)
			}
			cellText = strings.ReplaceAll(strings.TrimSpace(cellText), "\n", " ")
			md.WriteString(" " + cellText + " |")
		}
		md.WriteString("\n")

		if i == 0 {
			md.WriteString("|")
			for range row.TableCells {
				md.WriteString(" --- |")
			}
			md.WriteString("\n")
		}
	}
	md.WriteString("\n")

	return md.String()
}
