package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/elonfeng/url2md/pkg/converter"
)

type convertRequest struct {
	URL          string `json:"url"`
	Method       string `json:"method,omitempty"`
	RetainImages bool   `json:"retain_images,omitempty"`
	RetainLinks  *bool  `json:"retain_links,omitempty"`
}

type convertResponse struct {
	URL         string            `json:"url"`
	Markdown    string            `json:"markdown"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	TokenCount  int               `json:"token_count"`
	Method      string            `json:"method"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	FetchMs     int64             `json:"fetch_ms"`
	ConvertMs   int64             `json:"convert_ms"`
}

type errorResponse struct {
	Error string `json:"error"`
}

// Server is the url2md HTTP server.
type Server struct {
	conv converter.Converter
	port int
}

// New creates a new Server.
func New(port int) *Server {
	return &Server{
		conv: converter.New(),
		port: port,
	}
}

// ListenAndServe starts the HTTP server.
func (s *Server) ListenAndServe() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleConvert)
	mux.HandleFunc("/health", s.handleHealth)

	addr := fmt.Sprintf(":%d", s.port)
	fmt.Printf("url2md server listening on %s\n", addr)
	return http.ListenAndServe(addr, mux)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *Server) handleConvert(w http.ResponseWriter, r *http.Request) {
	var opts converter.Options
	*(&opts) = *converter.DefaultOptions()

	var targetURL string

	switch r.Method {
	case http.MethodGet:
		// GET /https://example.com?method=auto&retain_images=false
		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "" || path == "health" {
			return
		}
		rawQuery := r.URL.RawQuery
		targetURL = path
		if rawQuery != "" {
			// separate our params from the target URL's query
			targetURL = path
		}

		if m := r.URL.Query().Get("method"); m != "" {
			opts.Method = m
		}
		if r.URL.Query().Get("retain_images") == "true" {
			opts.RetainImages = true
		}
		if r.URL.Query().Get("retain_links") == "false" {
			opts.RetainLinks = false
		}
		if r.URL.Query().Get("enable_browser") == "true" {
			opts.EnableBrowser = true
		}

	case http.MethodPost:
		var req convertRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
			return
		}
		targetURL = req.URL
		if req.Method != "" {
			opts.Method = req.Method
		}
		opts.RetainImages = req.RetainImages
		if req.RetainLinks != nil {
			opts.RetainLinks = *req.RetainLinks
		}

	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	if targetURL == "" {
		writeError(w, http.StatusBadRequest, "url is required")
		return
	}

	// ensure scheme
	if !strings.HasPrefix(targetURL, "http://") && !strings.HasPrefix(targetURL, "https://") {
		targetURL = "https://" + targetURL
	}

	result, err := s.conv.Convert(r.Context(), targetURL, &opts)
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Markdown-Tokens", fmt.Sprintf("%d", result.TokenCount))
	w.Header().Set("X-Convert-Method", result.Method)
	w.Header().Set("X-Fetch-Time", fmt.Sprintf("%dms", result.FetchTime.Milliseconds()))

	resp := convertResponse{
		URL:         result.URL,
		Markdown:    result.Markdown,
		Title:       result.Title,
		Description: result.Description,
		TokenCount:  result.TokenCount,
		Method:      result.Method,
		Metadata:    result.Metadata,
		FetchMs:     result.FetchTime.Milliseconds(),
		ConvertMs:   result.ConvertTime.Milliseconds(),
	}

	json.NewEncoder(w).Encode(resp)
}

func writeError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(errorResponse{Error: msg})
}

// SetTimeout configures the converter timeout (used for testing).
func (s *Server) SetTimeout(d time.Duration) {
	// Not directly exposed, but could be extended.
}
