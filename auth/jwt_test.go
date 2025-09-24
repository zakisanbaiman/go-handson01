package auth

import (
	"bytes"
	"context"
	"testing"

	"github.com/zakisanbaiman/go-handson01/entity"
	"github.com/zakisanbaiman/go-handson01/testutil/fixture"
)

func TestEmbed(t *testing.T) {
	want := []byte("-----BEGIN PUBLIC KEY-----")
	if !bytes.Contains(rawPubKey, want) {
		t.Errorf("rawPubKey = %v, want %v", rawPubKey, want)
	}
	want = []byte("-----BEGIN RSA PRIVATE KEY-----")
	if !bytes.Contains(rawPriKey, want) {
		t.Errorf("rawPriKey = %v, want %v", rawPriKey, want)
	}
}

func TestJWTer_GenerateToken(t *testing.T) {
	cts := context.Background()
	moq := &StoreMock{}
	wantID := entity.UserID(20)
	user := fixture.User(func(u *entity.User) {
		u.ID = wantID
	})
	moq.SaveFunc = func(ctx context.Context, key string, userID entity.UserID) error {
		if userID != wantID {
			t.Errorf("want %d, but got %d", wantID, userID)
		}
		return nil
	}
	sut, err := NewJWTer(moq)
	if err != nil {
		t.Fatalf("want no error, but got %v", err)
	}
	got, err := sut.GenerateToken(cts, *user)
	if err != nil {
		t.Fatalf("want no error, but got %v", err)
	}
	if len(got) == 0 {
		t.Fatalf("want not empty, but got empty")
	}
}
