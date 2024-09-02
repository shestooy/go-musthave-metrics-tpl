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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go startSaveMetrics(ctx)

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

	if err := storage.MStorage.WriteInFile(); err != nil {
		l.Log.Error("error write metrics in file", zap.Error(err))
		return err
	}

	l.Log.Info("server shutdown complete")
	return nil
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

func startSaveMetrics(ctx context.Context) {
	if f.StorageInterval > 0 {
		ticker := time.NewTicker(time.Duration(f.StorageInterval) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := storage.MStorage.WriteInFile(); err != nil {
					l.Log.Info("error saving metrics", zap.Error(err))
				}
			case <-ctx.Done():
				return
			}
		}
	}
}
