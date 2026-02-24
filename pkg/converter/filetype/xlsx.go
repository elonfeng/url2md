package filetype

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/xuri/excelize/v2"
)

// ConvertXLSX converts XLSX bytes to markdown tables (one per sheet).
func ConvertXLSX(data []byte, filename string) (string, error) {
	f, err := excelize.OpenReader(bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("xlsx open: %w", err)
	}
	defer f.Close()

	if filename == "" {
		filename = "spreadsheet.xlsx"
	}

	var md strings.Builder
	md.WriteString(fmt.Sprintf("# %s\n\n", filename))

	sheets := f.GetSheetList()
	for _, sheet := range sheets {
		rows, err := f.GetRows(sheet)
		if err != nil || len(rows) == 0 {
			continue
		}

		if len(sheets) > 1 {
			md.WriteString(fmt.Sprintf("## %s\n\n", sheet))
		}

		// determine max columns
		maxCols := 0
		for _, row := range rows {
			if len(row) > maxCols {
				maxCols = len(row)
			}
		}

		if maxCols == 0 {
			continue
		}

		for i, row := range rows {
			md.WriteString("|")
			for j := 0; j < maxCols; j++ {
				cell := ""
				if j < len(row) {
					cell = strings.ReplaceAll(strings.TrimSpace(row[j]), "|", "\\|")
				}
				md.WriteString(" " + cell + " |")
			}
			md.WriteString("\n")

			if i == 0 {
				md.WriteString("|")
				for j := 0; j < maxCols; j++ {
					md.WriteString(" --- |")
				}
				md.WriteString("\n")
			}
		}
		md.WriteString("\n")
	}

	return strings.TrimSpace(md.String()), nil
}
