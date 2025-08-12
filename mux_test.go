package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewMux(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/health", nil)
	sut := NewMux()
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
