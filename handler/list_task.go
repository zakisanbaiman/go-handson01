package handler

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/zakisanbaiman/go-handson01/entity"
	"github.com/zakisanbaiman/go-handson01/store"
)

type ListTask struct {
	DB      *sqlx.DB
	Repo    *store.Repository
	Service ListTaskService
}

type task struct {
	ID     entity.TaskID     `json:"id"`
	Title  string            `json:"title"`
	Status entity.TaskStatus `json:"status"`
}

func (lt *ListTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tasks, err := lt.Service.ListTasks(ctx)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: "failed to list tasks",
			Details: []string{err.Error()},
		}, http.StatusInternalServerError)
		return
	}
	rsp := []task{}
	for _, t := range tasks {
		rsp = append(rsp, task{
			ID:     t.ID,
			Title:  t.Title,
			Status: t.Status,
		})
	}
	RespondJSON(ctx, w, rsp, http.StatusOK)
}
