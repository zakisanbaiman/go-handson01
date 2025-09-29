package service

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/zakisanbaiman/go-handson01/auth"
	"github.com/zakisanbaiman/go-handson01/entity"
	"github.com/zakisanbaiman/go-handson01/store"
)

type Login struct {
	DB    *sqlx.DB
	Repo  *store.Repository
	JWTer *auth.JWTer
}

func (l *Login) Login(ctx context.Context, userName, password string) (string, error) {
	// ユーザー認証のロジックをここに実装
	// 1. ユーザー名でユーザーを検索
	// 2. パスワードを検証
	// 3. JWTトークンを生成して返す

	// TODO: 実際の認証ロジックを実装
	user := entity.User{
		ID:   1,
		Name: userName,
		Role: "user",
	}

	token, err := l.JWTer.GenerateToken(ctx, user)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return string(token), nil
}
