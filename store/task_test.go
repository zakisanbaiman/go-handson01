package store

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/go-cmp/cmp"
	"github.com/jmoiron/sqlx"
	"github.com/zakisanbaiman/go-handson01/clock"
	"github.com/zakisanbaiman/go-handson01/entity"
	"github.com/zakisanbaiman/go-handson01/testutil"
)

// RDBMSを使ったテスト
func TestRepository_ListTasks(t *testing.T) {
	ctx := context.Background()

	tx, err := testutil.OpenDBForTest(t).BeginTxx(ctx, nil)

	t.Cleanup(func() {
		_ = tx.Rollback()
	})

	if err != nil {
		t.Fatalf("failed to begin tx: %s", err)
	}

	wants := prepareTasks(ctx, t, tx)

	sut := &Repository{}
	gots, err := sut.ListTasks(ctx, tx)
	if err != nil {
		t.Fatalf("failed to list tasks: %s", err)
	}

	if diff := cmp.Diff(wants, gots); diff != "" {
		t.Errorf("ListTasks() mismatch (-want +got):\n%s", diff)
	}
}

func prepareTasks(ctx context.Context, t *testing.T, con Execer) entity.Tasks {
	t.Helper()

	if _, err := con.ExecContext(ctx, "DELETE FROM tasks"); err != nil {
		t.Fatalf("failed to delete tasks: %s", err)
	}

	fixedTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	wants := entity.Tasks{
		{
			Title:      "test1",
			Status:     entity.TaskStatusTodo,
			CreatedAt:  fixedTime,
			ModifiedAt: fixedTime,
		},
		{
			Title:      "test2",
			Status:     entity.TaskStatusDone,
			CreatedAt:  fixedTime,
			ModifiedAt: fixedTime,
		},
		{
			Title:      "test3",
			Status:     entity.TaskStatusTodo,
			CreatedAt:  fixedTime,
			ModifiedAt: fixedTime,
		},
	}

	result, err := con.ExecContext(ctx,
		`INSERT INTO tasks (title, status, created_at, modified_at) VALUES
		(?, ?, ?, ?),
		(?, ?, ?, ?),
		(?, ?, ?, ?);`,
		wants[0].Title, wants[0].Status, wants[0].CreatedAt, wants[0].ModifiedAt,
		wants[1].Title, wants[1].Status, wants[1].CreatedAt, wants[1].ModifiedAt,
		wants[2].Title, wants[2].Status, wants[2].CreatedAt, wants[2].ModifiedAt,
	)
	if err != nil {
		t.Fatalf("failed to insert tasks: %s", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("failed to get last insert id: %s", err)
	}

	wants[0].ID = entity.TaskID(id)
	wants[1].ID = entity.TaskID(id + 1)
	wants[2].ID = entity.TaskID(id + 2)

	return wants
}

// sqlmockを使ったテスト(RDBMSを使わない)
func TestRepository_AddTask(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	c := clock.FixedClocker{}

	var wantID int64 = 20

	okTask := &entity.Task{
		Title:      "ok test",
		Status:     entity.TaskStatusTodo,
		CreatedAt:  c.Now(),
		ModifiedAt: c.Now(),
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })

	mock.ExpectExec("INSERT INTO tasks \\(title, status, created_at, modified_at\\) VALUES \\(\\?, \\?, \\?, \\?\\);").
		WithArgs(okTask.Title, okTask.Status, okTask.CreatedAt, okTask.ModifiedAt).
		WillReturnResult(sqlmock.NewResult(wantID, 1))

	xdb := sqlx.NewDb(db, "mysql")
	r := &Repository{Clocker: c}
	if err := r.AddTask(ctx, xdb, okTask); err != nil {
		t.Errorf("failed to add task: %s", err)
	}
}
