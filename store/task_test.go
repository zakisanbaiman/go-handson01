package store

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/go-cmp/cmp"
	"github.com/jmoiron/sqlx"
	"github.com/zakisanbaiman/go-handson01/clock"
	"github.com/zakisanbaiman/go-handson01/entity"
	"github.com/zakisanbaiman/go-handson01/testutil"
	"github.com/zakisanbaiman/go-handson01/testutil/fixture"
)

// RDBMSを使ったテスト
func TestRepository_ListTasks(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	tx, err := testutil.OpenDBForTest(t).BeginTxx(ctx, nil)

	t.Cleanup(func() {
		_ = tx.Rollback()
	})

	if err != nil {
		t.Fatalf("failed to begin tx: %s", err)
	}

	wantUserID, wants := prepareTasks(ctx, t, tx)

	sut := &Repository{}
	gots, err := sut.ListTasks(ctx, tx, wantUserID)
	if err != nil {
		t.Fatalf("failed to list tasks: %s", err)
	}

	if diff := cmp.Diff(wants, gots); diff != "" {
		t.Errorf("ListTasks() mismatch (-want +got):\n%s", diff)
	}
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

func prepareUser(ctx context.Context, t *testing.T, db Execer) entity.UserID {
	t.Helper()

	user := fixture.User(nil)
	result, err := db.ExecContext(ctx,
		`INSERT INTO users (name, email, password, created_at, modified_at) VALUES
		(?, ?, ?, ?, ?);`,
		user.ID, user.Name, user.Role, user.CreatedAt, user.ModifiedAt,
	)
	if err != nil {
		t.Fatalf("failed to insert user: %s", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("failed to get last insert id: %s", err)
	}

	return entity.UserID(id)
}

func prepareTasks(ctx context.Context, t *testing.T, con Execer) (entity.UserID, entity.Tasks) {
	t.Helper()

	userID := prepareUser(ctx, t, con)
	ohterUserID := prepareUser(ctx, t, con)
	c := clock.FixedClocker{}
	wants := entity.Tasks{
		{
			UserID:     userID,
			Title:      "test1",
			Status:     entity.TaskStatusTodo,
			CreatedAt:  c.Now(),
			ModifiedAt: c.Now(),
		},
		{
			UserID:     userID,
			Title:      "test2",
			Status:     entity.TaskStatusDone,
			CreatedAt:  c.Now(),
			ModifiedAt: c.Now(),
		},
	}

	tasks := entity.Tasks{
		wants[0],
		{
			UserID:     ohterUserID,
			Title:      "test3 not want",
			Status:     entity.TaskStatusTodo,
			CreatedAt:  c.Now(),
			ModifiedAt: c.Now(),
		},
		wants[1],
	}

	result, err := con.ExecContext(ctx,
		`INSERT INTO tasks (user_id, title, status, created_at, modified_at) VALUES
		(?, ?, ?, ?, ?),
		(?, ?, ?, ?, ?),
		(?, ?, ?, ?, ?);`,
		tasks[0].UserID, tasks[0].Title, tasks[0].Status, tasks[0].CreatedAt, tasks[0].ModifiedAt,
		tasks[1].UserID, tasks[1].Title, tasks[1].Status, tasks[1].CreatedAt, tasks[1].ModifiedAt,
		tasks[2].UserID, tasks[2].Title, tasks[2].Status, tasks[2].CreatedAt, tasks[2].ModifiedAt,
	)
	if err != nil {
		t.Fatalf("failed to insert tasks: %s", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("failed to get last insert id: %s", err)
	}

	tasks[0].ID = entity.TaskID(id)
	tasks[1].ID = entity.TaskID(id + 1)
	tasks[2].ID = entity.TaskID(id + 2)
	return userID, wants
}
