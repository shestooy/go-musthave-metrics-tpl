package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"

	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/handlers"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/storage"
)

func main() {
	storage.Storage.Init()
	err := start()
	if err != nil {
		panic(err)
	}
}

func start() error {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/update/{type}/{name}/{value}", handlers.PostMetrics)
	r.Get("/value/{type}/{name}", handlers.GetMetricID)
	r.Get("/", handlers.GetAllMetrics)
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		return err
	}
	return err
}
