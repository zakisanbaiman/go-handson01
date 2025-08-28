package store

import (
	"context"

	"github.com/zakisanbaiman/go-handson01/entity"
)

func (r *Repository) ListTasks(
	ctx context.Context, db Queryer,
) (entity.Tasks, error) {
	tasks := entity.Tasks{}
	sql := `SELECT
		id,
		title,
		status,
		created_at,
		modified_at
	FROM tasks;`

	if err := db.SelectContext(ctx, &tasks, sql); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *Repository) AddTask(
	ctx context.Context, db Execer, t *entity.Task,
) error {
	t.CreatedAt = r.Clocker.Now()
	t.ModifiedAt = r.Clocker.Now()

	sql := `INSERT INTO tasks
		(title, status, created_at, modified_at) VALUES (?, ?, ?, ?);`

	result, err := db.ExecContext(
		ctx, sql, t.Title, t.Status, t.CreatedAt, t.ModifiedAt,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	t.ID = entity.TaskID(id)

	return nil
}
