package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/zakisanbaiman/go-handson01/entity"
	"github.com/zakisanbaiman/go-handson01/store"
)

type AddTask struct {
	Store     *store.TaskStore
	Validator *validator.Validate
}

func (h *AddTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var b struct {
		Title string `json:"title" validate:"required,max=100"`
	}
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: "failed to decode request",
			Details: []string{err.Error()},
		}, http.StatusBadRequest)
		return
	}
	err := h.Validator.Struct(b)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: "failed to validate request",
			Details: []string{err.Error()},
		}, http.StatusBadRequest)
		return
	}

	task := &entity.Task{
		Title:     b.Title,
		Status:    entity.TaskStatusTodo,
		CreatedAt: time.Now(),
	}
	id, err := h.Store.Add(task)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: "failed to add task",
			Details: []string{err.Error()},
		}, http.StatusInternalServerError)
		return
	}
	rsp := struct {
		ID entity.TaskID `json:"id"`
	}{ID: id}
	RespondJSON(ctx, w, rsp, http.StatusCreated)
}
