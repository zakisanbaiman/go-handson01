package handler

import (
	"context"

	"github.com/zakisanbaiman/go-handson01/entity"
)

//go:generate go run github.com/matryer/moq -out moq_test.go . ListTaskService AddTaskService RegisterUserService LoginService
type ListTaskService interface {
	ListTasks(ctx context.Context) (entity.Tasks, error)
}

type AddTaskService interface {
	AddTask(ctx context.Context, title string) (*entity.Task, error)
}

type RegisterUserService interface {
	RegisterUser(ctx context.Context, name string, password string, role string) (*entity.User, error)
}

type LoginService interface {
	Login(ctx context.Context, name string, password string) (string, error)
}
