package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/zakisanbaiman/go-handson01/config"
)

func TestNewMux(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/health", nil)
	cfg := &config.Config{
		Env:        "test",
		Port:       8080,
		DBHost:     "localhost",
		DBPort:     3306,
		DBUser:     "test",
		DBPassword: "test",
		DBName:     "test",
	}
	sut, cleanup, err := NewMux(context.Background(), cfg)
	if err != nil {
		t.Fatalf("failed to create mux: %v", err)
	}
	defer cleanup()
	sut.ServeHTTP(w, r)
	resp := w.Result()
	t.Cleanup(func() {
		resp.Body.Close()
	})

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, but got %d", http.StatusOK, resp.StatusCode)
	}

	got, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	if string(got) != `{"status": "ok"}` {
		t.Errorf("expected body %s, but got %s", `{"status": "ok"}`, string(got))
	}
}
