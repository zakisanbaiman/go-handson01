package store

import (
	"errors"

	"github.com/zakisanbaiman/go-handson01/entity"
)

var (
	Tasks = &TaskStore{Tasks: map[entity.TaskID]*entity.Task{}}

	ErrNotFound = errors.New("not found")
)

type TaskStore struct {
	LastID entity.TaskID
	Tasks  map[entity.TaskID]*entity.Task
}

func (ts *TaskStore) Add(t *entity.Task) (entity.TaskID, error) {
	ts.LastID++
	t.ID = ts.LastID
	ts.Tasks[t.ID] = t
	return t.ID, nil
}

func (ts *TaskStore) All() entity.Tasks {
	tasks := make([]*entity.Task, 0, len(ts.Tasks))
	for _, t := range ts.Tasks {
		tasks = append(tasks, t)
	}
	return tasks
}
