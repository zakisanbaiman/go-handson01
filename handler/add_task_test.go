package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/zakisanbaiman/go-handson01/clock"
	"github.com/zakisanbaiman/go-handson01/entity"
	"github.com/zakisanbaiman/go-handson01/store"
	"github.com/zakisanbaiman/go-handson01/testutil"
)

func TestAddTaskHandler(t *testing.T) {
	t.Parallel()

	type want struct {
		status  int
		rspFile string
	}

	tests := map[string]struct {
		reqFile string
		want    want
	}{
		"ok": {
			reqFile: "testdata/add_task/ok.json",
			want: want{
				status:  http.StatusCreated,
				rspFile: "testdata/add_task/ok_rsp.json",
			},
		},
		"bad_request": {
			reqFile: "testdata/add_task/bad_request.json",
			want: want{
				status:  http.StatusBadRequest,
				rspFile: "testdata/add_task/bad_request_rsp.json",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(testutil.LoadFile(t, tt.reqFile)))
			r.Header.Set("Content-Type", "application/json")

			moq := &AddTaskServiceMock{}
			moq.AddTaskFunc = func(ctx context.Context, title string) (*entity.Task, error) {
				if tt.want.status == http.StatusCreated {
					return &entity.Task{ID: 1}, nil
				}
				return nil, errors.New("error from mock")
			}
			sut := AddTask{
				Service: moq,
				Repo: &store.Repository{
					Clocker: clock.RealClocker{},
				},
				Validator: validator.New(),
			}
			sut.ServeHTTP(w, r)
			testutil.AssertResponse(
				t, w.Result(), tt.want.status, testutil.LoadFile(t, tt.want.rspFile),
			)
		})
	}
}
