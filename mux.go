package main

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"github.com/zakisanbaiman/go-handson01/auth"
	"github.com/zakisanbaiman/go-handson01/clock"
	"github.com/zakisanbaiman/go-handson01/config"
	"github.com/zakisanbaiman/go-handson01/handler"
	"github.com/zakisanbaiman/go-handson01/service"
	"github.com/zakisanbaiman/go-handson01/store"
)

func NewMux(ctx context.Context, cfg *config.Config) (http.Handler, func(), error) {
	mux := chi.NewRouter()

	// CORSミドルウェアを適用
	mux.Use(handler.CORSMiddleware(cfg.GetCORSOptions()))

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write([]byte(`{"status": "ok"}`))
	})
	v := validator.New()

	db, cleanup, err := store.New(ctx, cfg)
	if err != nil {
		return nil, cleanup, err
	}

	r := store.Repository{Clocker: clock.RealClocker{}}

	// jwt関連
	clocker := clock.RealClocker{}
	repo := store.Repository{Clocker: clocker}
	rcli, err := store.NewKVS(ctx, cfg)
	if err != nil {
		return nil, cleanup, err
	}
	jwter, err := auth.NewJWTer(rcli, clocker)
	if err != nil {
		return nil, cleanup, err
	}

	// login
	login := &handler.Login{
		Service:   &service.Login{DB: db, Repo: &repo, TokenGenerator: jwter},
		Validator: v,
	}
	mux.Post("/login", login.ServeHTTP)

	// task
	at := &handler.AddTask{
		Service:   &service.AddTask{DB: db, Repo: &r},
		Validator: v,
	}
	lt := &handler.ListTask{
		Service: &service.ListTask{DB: db, Repo: &r},
	}
	mux.Route("/tasks", func(r chi.Router) {
		r.Use(handler.AuthMiddleware(jwter))
		r.Post("/", at.ServeHTTP)
		r.Get("/", lt.ServeHTTP)
	})

	// user
	ru := &handler.RegisterUser{
		Service:   &service.RegisterUser{DB: db, Repo: &r},
		Validator: v,
	}
	mux.Post("/users", ru.ServeHTTP)

	// admin
	mux.Route("/admin", func(r chi.Router) {
		r.Use(handler.AuthMiddleware(jwter))
		r.Use(handler.AdminMiddleware)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			_, _ = w.Write([]byte(`{"status": "ok"}`))
		})
	})

	return mux, cleanup, nil
}
