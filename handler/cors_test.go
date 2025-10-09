package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zakisanbaiman/go-handson01/config"
)

func TestCORSMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		origin         string
		method         string
		allowedOrigins []string
		expectedOrigin string
		expectedStatus int
	}{
		{
			name:           "allowed origin",
			origin:         "https://example.com",
			method:         "GET",
			allowedOrigins: []string{"https://example.com"},
			expectedOrigin: "https://example.com",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "wildcard origin",
			origin:         "https://anydomain.com",
			method:         "GET",
			allowedOrigins: []string{"*"},
			expectedOrigin: "https://anydomain.com",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "wildcard subdomain",
			origin:         "https://app.example.com",
			method:         "GET",
			allowedOrigins: []string{"*.example.com"},
			expectedOrigin: "https://app.example.com",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "disallowed origin",
			origin:         "https://malicious.com",
			method:         "GET",
			allowedOrigins: []string{"https://example.com"},
			expectedOrigin: "",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "OPTIONS preflight request",
			origin:         "https://example.com",
			method:         "OPTIONS",
			allowedOrigins: []string{"https://example.com"},
			expectedOrigin: "https://example.com",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := &config.CORSOptions{
				AllowedOrigins: tt.allowedOrigins,
				AllowedMethods: []string{"GET", "POST", "OPTIONS"},
				AllowedHeaders: []string{"Content-Type", "Authorization"},
				MaxAge:         3600,
			}

			handler := CORSMiddleware(options)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest(tt.method, "/test", nil)
			if tt.origin != "" {
				req.Header.Set("Origin", tt.origin)
			}

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedOrigin != "" {
				actualOrigin := w.Header().Get("Access-Control-Allow-Origin")
				if actualOrigin != tt.expectedOrigin {
					t.Errorf("expected origin %s, got %s", tt.expectedOrigin, actualOrigin)
				}
			}

			// OPTIONSリクエストの場合、追加のヘッダーをチェック
			if tt.method == "OPTIONS" {
				if w.Header().Get("Access-Control-Allow-Methods") == "" {
					t.Error("expected Access-Control-Allow-Methods header")
				}
				if w.Header().Get("Access-Control-Allow-Headers") == "" {
					t.Error("expected Access-Control-Allow-Headers header")
				}
				if w.Header().Get("Access-Control-Max-Age") == "" {
					t.Error("expected Access-Control-Max-Age header")
				}
			}
		})
	}
}

func TestIsAllowedOrigin(t *testing.T) {
	tests := []struct {
		name           string
		origin         string
		allowedOrigins []string
		expected       bool
	}{
		{
			name:           "exact match",
			origin:         "https://example.com",
			allowedOrigins: []string{"https://example.com"},
			expected:       true,
		},
		{
			name:           "wildcard match",
			origin:         "https://anydomain.com",
			allowedOrigins: []string{"*"},
			expected:       true,
		},
		{
			name:           "wildcard subdomain match",
			origin:         "https://app.example.com",
			allowedOrigins: []string{"*.example.com"},
			expected:       true,
		},
		{
			name:           "wildcard subdomain no match",
			origin:         "https://app.other.com",
			allowedOrigins: []string{"*.example.com"},
			expected:       false,
		},
		{
			name:           "no match",
			origin:         "https://malicious.com",
			allowedOrigins: []string{"https://example.com"},
			expected:       false,
		},
		{
			name:           "empty origin",
			origin:         "",
			allowedOrigins: []string{"*"},
			expected:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isAllowedOrigin(tt.origin, tt.allowedOrigins)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestDefaultCORSOptions(t *testing.T) {
	options := DefaultCORSOptions()

	if len(options.AllowedOrigins) == 0 {
		t.Error("expected non-empty AllowedOrigins")
	}

	if len(options.AllowedMethods) == 0 {
		t.Error("expected non-empty AllowedMethods")
	}

	if len(options.AllowedHeaders) == 0 {
		t.Error("expected non-empty AllowedHeaders")
	}

	if options.MaxAge <= 0 {
		t.Error("expected positive MaxAge")
	}
}
