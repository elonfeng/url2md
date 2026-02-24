package filetype

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/extrame/xls"
)

// ConvertXLS converts legacy .xls (BIFF) bytes to markdown tables.
func ConvertXLS(data []byte, filename string) (string, error) {
	// extrame/xls requires a file on disk
	tmp, err := os.CreateTemp("", "url2md-*.xls")
	if err != nil {
		return "", fmt.Errorf("xls temp file: %w", err)
	}
	defer os.Remove(tmp.Name())

	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return "", fmt.Errorf("xls write temp: %w", err)
	}
	tmp.Close()

	wb, err := xls.Open(tmp.Name(), "utf-8")
	if err != nil {
		return "", fmt.Errorf("xls open: %w", err)
	}

	if filename == "" {
		filename = "spreadsheet.xls"
	}

	var md strings.Builder
	md.WriteString(fmt.Sprintf("# %s\n\n", filename))

	numSheets := wb.NumSheets()
	for i := 0; i < numSheets; i++ {
		sheet := wb.GetSheet(i)
		if sheet == nil {
			continue
		}

		maxRow := int(sheet.MaxRow)
		if maxRow == 0 {
			continue
		}

		if numSheets > 1 {
			md.WriteString(fmt.Sprintf("## %s\n\n", sheet.Name))
		}

		// Determine max columns from all rows
		maxCols := 0
		for r := 0; r <= maxRow; r++ {
			row := sheet.Row(r)
			if row == nil {
				continue
			}
			cols := row.LastCol()
			if cols > maxCols {
				maxCols = cols
			}
		}

		if maxCols == 0 {
			continue
		}

		for r := 0; r <= maxRow; r++ {
			row := sheet.Row(r)
			md.WriteString("|")
			for c := 0; c < maxCols; c++ {
				cell := ""
				if row != nil {
					cell = xlsCellValue(strings.TrimSpace(row.Col(c)))
					cell = strings.ReplaceAll(cell, "|", "\\|")
				}
				md.WriteString(" " + cell + " |")
			}
			md.WriteString("\n")

			if r == 0 {
				md.WriteString("|")
				for c := 0; c < maxCols; c++ {
					md.WriteString(" --- |")
				}
				md.WriteString("\n")
			}
		}
		md.WriteString("\n")
	}

	return strings.TrimSpace(md.String()), nil
}

// xlsCellValue fixes a known issue in extrame/xls where numeric cells with
// user-defined formats get incorrectly converted to date strings.
func xlsCellValue(raw string) string {
	// If the value looks like an RFC3339 date near the Excel epoch (1900-01-01),
	// it's likely a small number misinterpreted as a date.
	if t, err := time.Parse(time.RFC3339, raw); err == nil {
		if t.Year() == 1900 && t.Month() <= 12 {
			// Excel serial: Jan 1, 1900 = 1, Jan 2 = 2, etc.
			// Excel has a leap year bug treating 1900 as leap year, so Feb 29 = serial 60.
			epoch := time.Date(1899, 12, 30, 0, 0, 0, 0, time.UTC)
			days := int(t.Sub(epoch).Hours() / 24)
			return strconv.Itoa(days)
		}
	}
	return raw
}
