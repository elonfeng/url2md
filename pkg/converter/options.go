package converter

import "time"

// Options configures the conversion behavior.
type Options struct {
	RetainImages  bool
	RetainLinks   bool
	Frontmatter   bool // prepend YAML frontmatter (title, description, image)
	Timeout       time.Duration
	EnableBrowser bool
	UserAgent     string
	Method        string // "auto" | "negotiate" | "static" | "browser"
}

// DefaultOptions returns sensible defaults for conversion.
func DefaultOptions() *Options {
	return &Options{
		RetainImages:  false,
		RetainLinks:   true,
		Frontmatter:   true,
		Timeout:       30 * time.Second,
		EnableBrowser: false,
		UserAgent:     "url2md/1.0",
		Method:        "auto",
	}
}

// Result holds the conversion output and associated metadata.
type Result struct {
	URL         string
	Markdown    string
	Title       string
	Description string
	TokenCount  int
	Method      string
	Metadata    map[string]string
	FetchTime   time.Duration
	ConvertTime time.Duration
}
