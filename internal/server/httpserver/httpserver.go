package httpserver

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	f "github.com/shestooy/go-musthave-metrics-tpl.git/internal/flags"
	l "github.com/shestooy/go-musthave-metrics-tpl.git/internal/logger"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/server/handlers"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/server/middlewares"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/storage"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := l.Initialize(f.LogLevel); err != nil {
		return err
	}

	if err := initializeStorage(ctx); err != nil {
		l.Log.Info("Error init Storage", zap.Error(err))
	}

	l.Log.Info("Running server", zap.String("address", f.ServerEndPoint))

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- http.ListenAndServe(f.ServerEndPoint, GetRouter())
	}()

	select {
	case err := <-serverErr:
		if err != nil {
			return err
		}
	case <-stop:
		l.Log.Info("Shutting down...")
		cancel()
	}

	if err := storage.MStorage.Close(); err != nil {
		l.Log.Error("Error closing storage", zap.Error(err))
	}

	l.Log.Info("Server shutdown complete")
	return nil
}

func GetRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(middlewares.LoggingMiddleware)
	r.Use(middlewares.GzipMiddleware)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	r.Get("/", handlers.GetAllMetrics)
	r.Get("/ping", handlers.PingHandler)

	r.Post("/updates/", handlers.UpdateSomeMetrics)

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

func initializeStorage(ctx context.Context) error {
	if f.AddrDB == "" {
		storage.MStorage = &storage.Storage{}
	} else {
		storage.MStorage = &storage.DB{}
	}
	err := storage.MStorage.Init(ctx)
	if err != nil {
		return err
	}
	return nil
}
