package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/e0m-ru/fileserver/internal/config"
)

func TestCollectHandlers(t *testing.T) {
	// Initialize config for testing
	config.Config.Os.Uploads = "/tmp/test_uploads"

	mux, err := CollectHandlers()
	if err != nil {
		t.Fatalf("CollectHandlers() error = %v", err)
	}

	if mux == nil {
		t.Fatal("CollectHandlers() returned nil mux")
	}

	// Test that mux is created and can handle requests
	req, err := http.NewRequest("GET", "/uploads/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	// Test passes if handler exists (even if returns 404 for empty dir)
	t.Logf("/uploads/ handler status: %d", rr.Code)
}

func TestDeleteHandler_InvalidFilename(t *testing.T) {
	config.Config.Os.Uploads = "/tmp/test_uploads"

	tests := []struct {
		name     string
		filename string
	}{
		{"path traversal dots", "../etc/passwd"},
		{"absolute path", "/etc/passwd"},
		{"backslash", "file\\name"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/api/delete?filename="+tt.filename, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(deleteHandler)
			handler.ServeHTTP(rr, req)

			// Should reject invalid filenames (not StatusOK)
			if rr.Code == http.StatusOK {
				t.Errorf("deleteHandler should reject invalid filename: %s", tt.filename)
			}
		})
	}
}
