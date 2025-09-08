package service

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/zakisanbaiman/go-handson01/entity"
	"github.com/zakisanbaiman/go-handson01/store"
)

type AddTask struct {
	DB   *sqlx.DB
	Repo *store.Repository
}

func (a *AddTask) AddTask(ctx context.Context, title string) (*entity.Task, error) {
	task := &entity.Task{
		Title:  title,
		Status: entity.TaskStatusTodo,
	}
	err := a.Repo.AddTask(ctx, a.DB, task)
	if err != nil {
		return nil, fmt.Errorf("failed to register: %w", err)
	}
	return task, nil
}
