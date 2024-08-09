package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"

	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/handlers"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/storage"
)

func main() {
	parseFlag()

	err := start()
	if err != nil {
		panic(err)
	}
}

func start() error {
	storage.Storage.Init()
	fmt.Printf("Server start on %s\n", serverEndPoint)
	return http.ListenAndServe(serverEndPoint, GetRouter())
}

func GetRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/update/{type}/{name}/{value}", handlers.PostMetrics)
	r.Get("/value/{type}/{name}", handlers.GetMetricID)
	r.Get("/", handlers.GetAllMetrics)
	return r
}
