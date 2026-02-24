package filetype

import "strings"

// ConvertMD passes through markdown content as-is.
func ConvertMD(data []byte) (string, error) {
	return strings.TrimSpace(string(data)), nil
}
