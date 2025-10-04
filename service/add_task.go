package service

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/zakisanbaiman/go-handson01/auth"
	"github.com/zakisanbaiman/go-handson01/entity"
	"github.com/zakisanbaiman/go-handson01/store"
)

type AddTask struct {
	DB   *sqlx.DB
	Repo *store.Repository
}

func (a *AddTask) AddTask(ctx context.Context, title string) (*entity.Task, error) {
	userID, ok := auth.GetUserID(ctx)
	if !ok {
		return nil, fmt.Errorf("user_id not found")
	}

	task := &entity.Task{
		UserID: userID,
		Title:  title,
		Status: entity.TaskStatusTodo,
	}
	err := a.Repo.AddTask(ctx, a.DB, task)
	if err != nil {
		return nil, fmt.Errorf("failed to register: %w", err)
	}
	return task, nil
}
