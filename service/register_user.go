package service

import (
	"context"
	"fmt"

	"github.com/zakisanbaiman/go-handson01/entity"
	"github.com/zakisanbaiman/go-handson01/store"
)

type RegisterUser struct {
	DB   store.Execer
	Repo *store.Repository
}

func (r *RegisterUser) RegisterUser(ctx context.Context, name string, password string, role string) (*entity.User, error) {
	user := &entity.User{
		Name:     name,
		Password: password,
		Role:     role,
	}
	if err := r.Repo.RegisterUser(ctx, r.DB, user); err != nil {
		return nil, fmt.Errorf("failed to register user: %w", err)
	}
	return user, nil
}
