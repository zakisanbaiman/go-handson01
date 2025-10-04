package store

import (
	"context"

	"github.com/zakisanbaiman/go-handson01/entity"
)

func (r *Repository) ListTasks(
	ctx context.Context, db Queryer, userID entity.UserID,
) (entity.Tasks, error) {
	tasks := entity.Tasks{}
	sql := `SELECT
		id,
		user_id,
		title,
		status,
		created_at,
		modified_at
	FROM tasks
	WHERE user_id = ?;`

	if err := db.SelectContext(ctx, &tasks, sql, userID); err != nil {
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
		(user_id, title, status, created_at, modified_at) VALUES (?, ?, ?, ?, ?);`

	result, err := db.ExecContext(
		ctx, sql, t.UserID, t.Title, t.Status, t.CreatedAt, t.ModifiedAt,
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
