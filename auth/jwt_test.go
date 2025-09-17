package auth

import (
	"bytes"
	"testing"
)

func TestEmbed(t *testing.T) {
	t.Skip("Skipping TestEmbed as it references local PEM files")
	want := []byte("-----BEGIN PUBLIC KEY-----")
	if !bytes.Contains(rawPubKey, want) {
		t.Errorf("rawPubKey = %v, want %v", rawPubKey, want)
	}
	want = []byte("-----BEGIN RSA PRIVATE KEY-----")
	if !bytes.Contains(rawPriKey, want) {
		t.Errorf("rawPriKey = %v, want %v", rawPriKey, want)
	}
}
