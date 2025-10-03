package service

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Login struct {
	DB             *sqlx.DB
	Repo           UserGetter
	TokenGenerator TokenGenerator
}

func (l *Login) Login(ctx context.Context, userName, password string) (string, error) {
	user, err := l.Repo.GetUser(ctx, l.DB, userName)
	if err != nil {
		return "", fmt.Errorf("failed to get user: %w", err)
	}
	if err := user.ComparePassword(password); err != nil {
		return "", fmt.Errorf("failed to compare password: %w", err)
	}

	token, err := l.TokenGenerator.GenerateToken(ctx, *user)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return string(token), nil
}
