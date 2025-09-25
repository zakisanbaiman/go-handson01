package auth

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/zakisanbaiman/go-handson01/clock"
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
	sut, err := NewJWTer(moq, clock.RealClocker{})
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

func TestJWTer_GetToken(t *testing.T) {
	t.Parallel()

	c := clock.FixedClocker{}

	// payload作成
	want, err := jwt.NewBuilder().
		JwtID(uuid.New().String()).
		Issuer("zakisanbaiman").
		Subject("access_token").
		IssuedAt(c.Now()).
		Expiration(c.Now().Add(30*time.Minute)).
		Claim(RoleKey, "test").
		Claim(UserNameKey, "test_user").
		Build()

	if err != nil {
		t.Fatalf("want no error, but got %v", err)
	}
	pkey, err := jwk.ParseKey(rawPriKey, jwk.WithPEM(true))
	if err != nil {
		t.Fatalf("want no error, but got %v", err)
	}

	// 署名
	signed, err := jwt.Sign(want, jwt.WithKey(jwa.RS256, pkey))
	if err != nil {
		t.Fatalf("want no error, but got %v", err)
	}
	userID := entity.UserID(20)

	ctx := context.Background()
	moq := &StoreMock{}
	moq.LoadFunc = func(ctx context.Context, key string) (entity.UserID, error) {
		return userID, nil
	}
	sut, err := NewJWTer(moq, c)
	if err != nil {
		t.Fatalf("want no error, but got %v", err)
	}

	req := httptest.NewRequest(
		http.MethodGet,
		"https://github.com/zakisanbaiman",
		nil,
	)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", signed))

	got, err := sut.GetToken(ctx, req)
	if err != nil {
		t.Fatalf("want no error, but got %v", err)
	}
	if got == nil {
		t.Fatalf("want not nil, but got nil")
	}
}
