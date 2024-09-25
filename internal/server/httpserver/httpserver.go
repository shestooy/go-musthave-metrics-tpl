package httpserver

import (
	"context"
	"github.com/labstack/echo/v4"
	f "github.com/shestooy/go-musthave-metrics-tpl.git/internal/flags"
	l "github.com/shestooy/go-musthave-metrics-tpl.git/internal/logger"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/server/handlers"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/server/middlewares"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/storage"
	"go.uber.org/zap"
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

	server := initServer()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- server.Start(f.ServerEndPoint)
	}()

	select {
	case err := <-serverErr:
		if err != nil {
			return err
		}
	case <-stop:
		l.Log.Info("Shutting down...")
		if err := server.Stop(ctx); err != nil {
			l.Log.Info("Error shutting down", zap.Error(err))
		}
		cancel()
	}

	if err := storage.MStorage.Close(); err != nil {
		l.Log.Error("Error closing storage", zap.Error(err))
	}

	l.Log.Info("Server shutdown complete")
	return nil
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

type Server struct {
	server *echo.Echo
}

func initServer() *Server {
	e := echo.New()

	e.HideBanner = true
	e.HidePort = true

	e.Use(middlewares.Gzip)
	e.Use(middlewares.Logging)
	e.Use(middlewares.Hash(f.ServerKey))

	e.GET("/", handlers.GetAllMetrics)
	e.GET("/ping", handlers.PingHandler)

	e.POST("/updates/", handlers.UpdateSomeMetrics)

	updateGroup := e.Group("/update")
	updateGroup.POST("/", handlers.PostMetricsWithJSON)
	updateGroup.POST("/:type/:name/:value", handlers.PostMetrics)

	valueGroup := e.Group("/value")
	valueGroup.POST("/", handlers.GetMetricIDWithJSON)
	valueGroup.GET("/:type/:name", handlers.GetMetricID)

	return &Server{server: e}
}

func (s *Server) Start(endPoint string) error {
	return s.server.Start(endPoint)
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Server.Shutdown(ctx)
}
