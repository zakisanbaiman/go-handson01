package auth

import (
	"context"
	_ "embed"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/zakisanbaiman/go-handson01/clock"
	"github.com/zakisanbaiman/go-handson01/entity"
)

//go:embed cert/secret.pem
var rawPriKey []byte

//go:embed cert/public.pem
var rawPubKey []byte

type JWTer struct {
	PrivateKey, PublicKey jwk.Key
	Store                 Store
	Clocker               clock.Clocker
}

//go:generate go run github.com/matryer/moq -out moq_test.go . Store
type Store interface {
	Save(ctx context.Context, key string, userID entity.UserID) error
	Load(ctx context.Context, key string) (entity.UserID, error)
}

func NewJWTer(s Store) (*JWTer, error) {
	j := &JWTer{Store: s}
	privateKey, err := parse(rawPriKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}
	publicKey, err := parse(rawPubKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	j.PrivateKey = privateKey
	j.PublicKey = publicKey
	j.Clocker = clock.RealClocker{}

	return j, nil
}

func parse(rawKey []byte) (jwk.Key, error) {
	key, err := jwk.ParseKey(rawKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse key: %w", err)
	}
	return key, nil
}

const (
	RoleKey     = "role"
	UserNameKey = "name"
)

func (j *JWTer) GenerateToken(ctx context.Context, user entity.User) ([]byte, error) {
	tok, err := jwt.NewBuilder().
		JwtID(uuid.New().String()).
		Issuer("access_token").
		IssuedAt(j.Clocker.Now()).
		Expiration(j.Clocker.Now().Add(30*time.Minute)).
		Claim(RoleKey, user.Role).
		Claim(UserNameKey, user.Name).
		Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build token: %w", err)
	}
	if err := j.Store.Save(ctx, tok.JwtID(), user.ID); err != nil {
		return nil, fmt.Errorf("failed to save token: %w", err)
	}

	signed, err := jwt.Sign(tok, jwt.WithKey(jwa.RS256, j.PrivateKey))
	if err != nil {
		return nil, fmt.Errorf("failed to sign token: %w", err)
	}
	return signed, nil
}
