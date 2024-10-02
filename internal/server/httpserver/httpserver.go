package httpserver

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/config"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/logger"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/server/handlers"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/server/middlewares"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/storage"
	"go.uber.org/zap"
)

type Server struct {
	server *echo.Echo
	h      *handlers.Handler
}

func Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.GetServerCfg()
	if err != nil {
		return err
	}

	l, err := logger.Initialize(cfg.LogLevel)
	if err != nil {
		return err
	}
	s, err := initializeStorage(ctx, l, cfg)
	if err != nil {
		l.Info("Error init Storage", zap.Error(err))
	}

	server := initServer(s, l, cfg)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- server.Start(cfg.ServerEndPoint)
	}()

	select {
	case err = <-serverErr:
		if err != nil {
			return err
		}
	case <-stop:
		l.Info("Shutting down...")
		if err = server.Stop(ctx); err != nil {
			l.Info("Error shutting down", zap.Error(err))
		}
		cancel()
	}
	l.Info("Server shutdown complete")
	return nil
}

func initializeStorage(ctx context.Context, l *zap.SugaredLogger, cfg *config.ServerCfg) (s storage.IStorage, err error) {
	if cfg.AddrDB == "" {
		s = &storage.Storage{}
	} else {
		s = &storage.DB{}
	}
	err = s.Init(ctx, l, cfg)
	if err != nil {
		return nil, err
	}
	return s, err
}

func initServer(s storage.IStorage, l *zap.SugaredLogger, cfg *config.ServerCfg) *Server {
	h := handlers.NewHandler(l, s)

	e := echo.New()

	e.HideBanner = true
	e.HidePort = true

	e.Use(middlewares.Gzip)
	e.Use(middlewares.GetLogg(l))
	e.Use(middlewares.Hash(cfg.ServerKey))

	e.GET("/", h.GetAllMetrics)
	e.GET("/ping", h.PingHandler)

	e.POST("/updates/", h.UpdateSomeMetrics)

	updateGroup := e.Group("/update")
	updateGroup.POST("/", h.PostMetricsWithJSON)
	updateGroup.POST("/:type/:name/:value", h.PostMetrics)

	valueGroup := e.Group("/value")
	valueGroup.POST("/", h.GetMetricIDWithJSON)
	valueGroup.GET("/:type/:name", h.GetMetricID)

	return &Server{server: e, h: h}
}

func (s *Server) Start(endPoint string) error {
	s.h.Logger.Info("Server starting on: ", endPoint)
	return s.server.Start(endPoint)
}

func (s *Server) Stop(ctx context.Context) error {
	err := s.server.Server.Shutdown(ctx)
	if err != nil {
		s.h.Logger.Info("Error shutting down", zap.Error(err))
	}
	return s.h.DB.Close()
}
