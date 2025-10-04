package service

import (
	"context"

	"github.com/zakisanbaiman/go-handson01/entity"
	"github.com/zakisanbaiman/go-handson01/store"
)

//go:generate go run github.com/matryer/moq -out moq_test.go . TaskAdder TaskLister UserGetter TokenGenerator
type TaskAdder interface {
	AddTask(ctx context.Context, db store.Execer, t *entity.Task) error
}

type TaskLister interface {
	ListTasks(ctx context.Context, db store.Queryer, userID entity.UserID) (entity.Tasks, error)
}

type UserGetter interface {
	GetUser(ctx context.Context, db store.Queryer, userName string) (*entity.User, error)
}

type TokenGenerator interface {
	GenerateToken(ctx context.Context, user entity.User) ([]byte, error)
}
