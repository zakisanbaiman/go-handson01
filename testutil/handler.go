package testutil

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func AssertJSON(t *testing.T, expected []byte, actual []byte) {
	t.Helper()

	var jsonExpected, jsonActual any
	if err := json.Unmarshal([]byte(expected), &jsonExpected); err != nil {
		t.Fatalf("cannot unmarshal expected: %v", err)
	}
	if err := json.Unmarshal([]byte(actual), &jsonActual); err != nil {
		t.Fatalf("cannot unmarshal actual: %v", err)
	}
	if diff := cmp.Diff(jsonExpected, jsonActual); diff != "" {
		t.Fatalf("expected and actual are not equal: %s", diff)
	}
}

func AssertResponse(t *testing.T, got *http.Response, wantStatus int, body []byte) {
	t.Helper()
	t.Cleanup(func() { _ = got.Body.Close() })

	gb, err := io.ReadAll(got.Body)
	if err != nil {
		t.Fatal(err)
	}
	if got.StatusCode != wantStatus {
		t.Fatalf("expected status code %d, but got %d", wantStatus, got.StatusCode)
	}

	if len(gb) == 0 && len(body) == 0 {
		return
	}
	AssertJSON(t, gb, body)
}

func LoadFile(t *testing.T, path string) []byte {
	t.Helper()

	bt, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("cannot read file %s: %v", path, err)
	}
	return bt
}
