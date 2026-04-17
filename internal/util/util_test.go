package util

import (
	"flag"
	"testing"
)

func TestOpenURL(t *testing.T) {
	// Skip actual URL opening in tests
	if testing.Short() {
		t.Skip("Skipping OpenURL test in short mode")
	}
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "valid http url",
			url:     "http://localhost:8080",
			wantErr: false,
		},
		{
			name:    "valid https url",
			url:     "https://example.com",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We can't fully test OpenURL since it's OS-dependent
			// But we can verify it doesn't panic with valid URLs
			err := OpenURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("OpenURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func init() {
	// Suppress flag warnings in tests
	flag.String("p", "8080", "port for serve files")
	flag.String("u", "", "folder to save files")
}
