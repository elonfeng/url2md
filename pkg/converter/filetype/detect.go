package filetype

import (
	"bytes"
	"net/http"
	"path"
	"strings"
)

// Type represents a supported file type.
type Type string

const (
	TypeHTML Type = "html"
	TypePDF  Type = "pdf"
	TypeDOCX Type = "docx"
	TypeXLSX Type = "xlsx"
	TypeCSV  Type = "csv"
	TypePNG  Type = "png"
	TypeJPEG Type = "jpeg"
	TypeGIF  Type = "gif"
	TypeWEBP Type = "webp"
)

// IsImage returns true if the type is an image format.
func (t Type) IsImage() bool {
	switch t {
	case TypePNG, TypeJPEG, TypeGIF, TypeWEBP:
		return true
	}
	return false
}

// DetectFromURL guesses the file type from the URL path extension.
func DetectFromURL(rawURL string) Type {
	ext := strings.ToLower(path.Ext(strings.Split(rawURL, "?")[0]))
	switch ext {
	case ".pdf":
		return TypePDF
	case ".docx":
		return TypeDOCX
	case ".xlsx":
		return TypeXLSX
	case ".csv":
		return TypeCSV
	case ".png":
		return TypePNG
	case ".jpg", ".jpeg":
		return TypeJPEG
	case ".gif":
		return TypeGIF
	case ".webp":
		return TypeWEBP
	}
	return TypeHTML
}

// DetectFromContentType determines the file type from an HTTP Content-Type header.
func DetectFromContentType(ct string) Type {
	ct = strings.ToLower(ct)
	switch {
	case strings.Contains(ct, "application/pdf"):
		return TypePDF
	case strings.Contains(ct, "application/vnd.openxmlformats-officedocument.wordprocessingml"):
		return TypeDOCX
	case strings.Contains(ct, "application/vnd.openxmlformats-officedocument.spreadsheetml"):
		return TypeXLSX
	case strings.Contains(ct, "application/vnd.ms-excel"):
		return TypeXLSX
	case strings.Contains(ct, "text/csv"):
		return TypeCSV
	case strings.Contains(ct, "image/png"):
		return TypePNG
	case strings.Contains(ct, "image/jpeg"):
		return TypeJPEG
	case strings.Contains(ct, "image/gif"):
		return TypeGIF
	case strings.Contains(ct, "image/webp"):
		return TypeWEBP
	case strings.Contains(ct, "application/octet-stream"):
		// generic binary — defer to magic bytes detection
		return TypeHTML
	case strings.Contains(ct, "text/html"), strings.Contains(ct, "application/xhtml"):
		return TypeHTML
	}
	return TypeHTML
}

// Magic byte signatures for binary file detection.
var (
	magicPDF  = []byte("%PDF")
	magicZIP  = []byte("PK\x03\x04") // DOCX and XLSX are ZIP-based
	magicPNG  = []byte("\x89PNG\r\n\x1a\n")
	magicJPEG = []byte("\xff\xd8\xff")
	magicGIF  = []byte("GIF8")
	magicWEBP = []byte("RIFF") // WEBP starts with RIFF....WEBP
)

// DetectFromBytes uses magic byte signatures to identify file types.
// For ZIP-based formats (DOCX/XLSX), it peeks inside the archive markers.
func DetectFromBytes(data []byte) Type {
	if len(data) < 8 {
		return TypeHTML
	}
	switch {
	case bytes.HasPrefix(data, magicPDF):
		return TypePDF
	case bytes.HasPrefix(data, magicZIP):
		return detectZIPType(data)
	case bytes.HasPrefix(data, magicPNG):
		return TypePNG
	case bytes.HasPrefix(data, magicJPEG):
		return TypeJPEG
	case bytes.HasPrefix(data, magicGIF):
		return TypeGIF
	case bytes.HasPrefix(data, magicWEBP) && len(data) >= 12 && string(data[8:12]) == "WEBP":
		return TypeWEBP
	}
	return TypeHTML
}

// detectZIPType distinguishes DOCX from XLSX by looking for known paths in the ZIP.
func detectZIPType(data []byte) Type {
	// Look for Office Open XML markers in the first 4KB
	peek := data
	if len(peek) > 4096 {
		peek = peek[:4096]
	}
	if bytes.Contains(peek, []byte("word/")) || bytes.Contains(peek, []byte("word\\")) {
		return TypeDOCX
	}
	if bytes.Contains(peek, []byte("xl/")) || bytes.Contains(peek, []byte("xl\\")) {
		return TypeXLSX
	}
	// Unknown ZIP type
	return TypeHTML
}

// Detect determines file type using URL extension, final redirect URL, Content-Type,
// and magic bytes (in that order of priority).
func Detect(rawURL string, resp *http.Response, data []byte) Type {
	// 1. Original URL extension
	t := DetectFromURL(rawURL)
	if t != TypeHTML {
		return t
	}

	// 2. Final redirect URL extension (after following redirects)
	if resp != nil && resp.Request != nil && resp.Request.URL != nil {
		finalURL := resp.Request.URL.String()
		if finalURL != rawURL {
			t = DetectFromURL(finalURL)
			if t != TypeHTML {
				return t
			}
		}
	}

	// 3. Content-Type header
	if resp != nil {
		ct := resp.Header.Get("Content-Type")
		if ct != "" {
			t = DetectFromContentType(ct)
			if t != TypeHTML {
				return t
			}
		}
	}

	// 4. Magic bytes
	if len(data) > 0 {
		t = DetectFromBytes(data)
		if t != TypeHTML {
			return t
		}
	}

	return TypeHTML
}
