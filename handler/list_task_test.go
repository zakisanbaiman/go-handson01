package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/zakisanbaiman/go-handson01/entity"
	"github.com/zakisanbaiman/go-handson01/store"
	"github.com/zakisanbaiman/go-handson01/testutil"
)

func TestListTaskHandler(t *testing.T) {
	t.Parallel()

	type want struct {
		status  int
		rspFile string
	}

	tests := map[string]struct {
		tasks []*entity.Task
		want  want
	}{
		"ok": {
			tasks: []*entity.Task{
				{ID: 1, Title: "test1", Status: entity.TaskStatusTodo},
				{ID: 2, Title: "test2", Status: entity.TaskStatusDone},
			},
			want: want{
				status:  http.StatusOK,
				rspFile: "testdata/list_task/ok_rsp.json",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/tasks", nil)

			moq := &ListTaskServiceMock{}
			moq.ListTasksFunc = func(ctx context.Context) (entity.Tasks, error) {
				if tt.tasks != nil {
					return tt.tasks, nil
				}
				return nil, errors.New("error from mock")
			}
			sut := ListTask{
				Service: moq,
				DB:      &sqlx.DB{},
				Repo:    &store.Repository{},
			}
			sut.ServeHTTP(w, r)
			testutil.AssertResponse(
				t, w.Result(), tt.want.status, testutil.LoadFile(t, tt.want.rspFile),
			)
		})
	}
}
