package httpserver

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	f "github.com/shestooy/go-musthave-metrics-tpl.git/internal/flags"
	l "github.com/shestooy/go-musthave-metrics-tpl.git/internal/logger"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/server/handlers"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/server/middlewares"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/storage"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func Start() error {
	if err := l.Initialize(f.LogLevel); err != nil {
		return err
	}
	if err := storage.MStorage.Init(); err != nil {
		l.Log.Info("Error init storage", zap.Error(err))
		return err
	}
	go startSaveMetrics()
	l.Log.Info("Running server", zap.String("address", f.ServerEndPoint))

	err := http.ListenAndServe(f.ServerEndPoint, GetRouter())
	if err != nil {
		return err
	}
	return storage.MStorage.WriteInFile()
}

func GetRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(middlewares.LoggingMiddleware)
	r.Use(middlewares.GzipMiddleware)
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

func startSaveMetrics() {
	if f.StorageInterval > 0 {
		ticker := time.NewTicker(time.Duration(f.StorageInterval) * time.Second)
		defer ticker.Stop()
		go func() {
			for range ticker.C {
				if err := storage.MStorage.WriteInFile(); err != nil {
					l.Log.Info("error saving metrics", zap.Error(err))
				}
			}
		}()
	}
}
