package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthEndpoint(t *testing.T) {
	srv := New(0)
	mux := http.NewServeMux()
	mux.HandleFunc("/health", srv.handleHealth)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var body map[string]string
	json.NewDecoder(w.Body).Decode(&body)
	if body["status"] != "ok" {
		t.Errorf("expected status ok, got %q", body["status"])
	}
}

func TestConvertEndpoint_POST(t *testing.T) {
	// create a target server that returns HTML
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `<html><head><title>Test</title></head><body>
		<article><p>Test article content that is long enough for readability extraction.</p>
		<p>Second paragraph of content for the readability algorithm to work.</p>
		<p>Third paragraph of substantial content to ensure extraction succeeds.</p></article></body></html>`)
	}))
	defer target.Close()

	srv := New(0)
	mux := http.NewServeMux()
	mux.HandleFunc("/", srv.handleConvert)

	payload := fmt.Sprintf(`{"url":"%s","method":"static"}`, target.URL)
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d; body: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["method"] != "static" {
		t.Errorf("expected method static, got %v", resp["method"])
	}
	if resp["markdown"] == "" {
		t.Error("expected non-empty markdown")
	}

	if w.Header().Get("X-Convert-Method") != "static" {
		t.Errorf("expected X-Convert-Method header, got %q", w.Header().Get("X-Convert-Method"))
	}
}

func TestConvertEndpoint_MissingURL(t *testing.T) {
	srv := New(0)
	mux := http.NewServeMux()
	mux.HandleFunc("/", srv.handleConvert)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestConvertEndpoint_MethodNotAllowed(t *testing.T) {
	srv := New(0)
	mux := http.NewServeMux()
	mux.HandleFunc("/", srv.handleConvert)

	req := httptest.NewRequest(http.MethodDelete, "/test", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", w.Code)
	}
}
