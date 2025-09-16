package auth

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/yumemi-inc/go-oidc/pkg/jwt"
	"github.com/zakisanbaiman/go-handson01/clock"
	"github.com/zakisanbaiman/go-handson01/entity"
)

//go:embed cert/secret.pem
var rawPriKey []byte

//go:embed cert/public.pem
var rawPubKey []byte

type JWTer struct {
	PrivateKey, PublicKey jwt.Key
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

func parse(rawKey []byte) (jwt.Key, error) {
	key, err := jwt.ParseKey(rawKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse key: %w", err)
	}
	return key, nil
}
