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
