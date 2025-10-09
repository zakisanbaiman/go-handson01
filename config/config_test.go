package config

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	expectedPort := 3333
	t.Setenv("PORT", fmt.Sprintf("%d", expectedPort))

	got, err := New()
	if err != nil {
		t.Fatalf("failed to get config: %v", err)
	}

	if got.Port != expectedPort {
		t.Errorf("expected port %d, but got %d", expectedPort, got.Port)
	}

	if got.Env != "dev" {
		t.Errorf("expected env %s, but got %s", "dev", got.Env)
	}
}

func TestGetCORSOptions(t *testing.T) {
	cfg := &Config{
		CORSAllowedOrigins: "https://example.com,https://app.example.com",
		CORSAllowedMethods: "GET,POST,PUT,DELETE,OPTIONS",
		CORSAllowedHeaders: "Content-Type,Authorization",
		CORSMaxAge:         3600,
	}

	options := cfg.GetCORSOptions()

	expectedOrigins := []string{"https://example.com", "https://app.example.com"}
	if len(options.AllowedOrigins) != len(expectedOrigins) {
		t.Errorf("expected %d origins, got %d", len(expectedOrigins), len(options.AllowedOrigins))
	}

	for i, expected := range expectedOrigins {
		if options.AllowedOrigins[i] != expected {
			t.Errorf("expected origin %s, got %s", expected, options.AllowedOrigins[i])
		}
	}

	expectedMethods := []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	if len(options.AllowedMethods) != len(expectedMethods) {
		t.Errorf("expected %d methods, got %d", len(expectedMethods), len(options.AllowedMethods))
	}

	if options.MaxAge != 3600 {
		t.Errorf("expected MaxAge 3600, got %d", options.MaxAge)
	}
}
