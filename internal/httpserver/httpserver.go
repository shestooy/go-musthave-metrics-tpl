package httpserver

import (
	"fmt"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/flags"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/handlers"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/storage"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Start() error {
	storage.Storage.Init()
	fmt.Printf("Server start on %s\n", flags.ServerEndPoint)
	return http.ListenAndServe(flags.ServerEndPoint, GetRouter())
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
