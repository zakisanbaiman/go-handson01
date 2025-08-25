package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"github.com/zakisanbaiman/go-handson01/handler"
	"github.com/zakisanbaiman/go-handson01/store"
)

// func NewMux() *http.ServeMux {
// 	mux := http.NewServeMux()
// 	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Content-Type", "application/json; charset=utf-8")
// 		_, _ = w.Write([]byte(`{"status": "ok"}`))
// 	})
// 	return mux
// }

func NewMux() http.Handler {
	mux := chi.NewRouter()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write([]byte(`{"status": "ok"}`))
	})
	v := validator.New()

	at := &handler.AddTask{Store: store.Tasks, Validator: v}
	mux.Post("/tasks", at.ServeHTTP)

	lt := &handler.ListTask{Store: store.Tasks}
	mux.Get("/tasks", lt.ServeHTTP)

	return mux
}
