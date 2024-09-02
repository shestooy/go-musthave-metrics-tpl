package httpserver

import (
	f "github.com/shestooy/go-musthave-metrics-tpl.git/internal/flags"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/handlers"
	l "github.com/shestooy/go-musthave-metrics-tpl.git/internal/logger"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/middlewares"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/storage"
	"go.uber.org/zap"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Start() error {
	storage.MStorage.Init()
	if err := l.Initialize(f.LogLevel); err != nil {
		return err
	}
	l.Log.Info("Running server", zap.String("address", f.ServerEndPoint))

	return http.ListenAndServe(f.ServerEndPoint, GetRouter())
}

func GetRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(middlewares.LoggingMiddleware)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	r.Get("/", handlers.GetAllMetrics)

	r.Route("/update", func(r chi.Router) {
		r.Post("/", handlers.PostMetricsWithJSON)
		r.Post("/{type}/{name}/{value}", handlers.PostMetrics)
	})

	r.Route("/value", func(r chi.Router) {
		r.Post("/", handlers.GetMetricIDWithJSON)
		r.Get("/{type}/{name}", handlers.GetMetricID)
	})

	return r
}
