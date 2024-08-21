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
	storage.Storage.Init()
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

	r.Post("/update/{type}/{name}/{value}", handlers.PostMetrics)
	r.Get("/value/{type}/{name}", handlers.GetMetricID)
	r.Get("/", handlers.GetAllMetrics)
	return r
}
