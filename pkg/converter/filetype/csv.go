package filetype

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"strings"
)

// ConvertCSV converts CSV bytes to a markdown table.
func ConvertCSV(data []byte, filename string) (string, error) {
	reader := csv.NewReader(bytes.NewReader(data))
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true

	var rows [][]string
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			// try to continue on parse errors
			continue
		}
		rows = append(rows, record)
	}

	if len(rows) == 0 {
		return "", fmt.Errorf("csv: no data")
	}

	var md strings.Builder

	if filename == "" {
		filename = "data.csv"
	}
	md.WriteString(fmt.Sprintf("# %s\n\n", filename))

	// determine max columns
	maxCols := 0
	for _, row := range rows {
		if len(row) > maxCols {
			maxCols = len(row)
		}
	}

	// build markdown table
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

		// separator after first row (header)
		if i == 0 {
			md.WriteString("|")
			for j := 0; j < maxCols; j++ {
				md.WriteString(" --- |")
			}
			md.WriteString("\n")
		}
	}

	return strings.TrimSpace(md.String()), nil
}
