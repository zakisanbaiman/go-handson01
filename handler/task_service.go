package handler

import (
	"context"

	"github.com/zakisanbaiman/go-handson01/entity"
)

//go:generate go run github.com/matryer/moq -out moq_test.go . ListTaskService AddTaskService
type ListTaskService interface {
	ListTasks(ctx context.Context) (entity.Tasks, error)
}

type AddTaskService interface {
	AddTask(ctx context.Context, title string) (*entity.Task, error)
}
