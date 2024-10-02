package agent

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/logger"

	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/agent/metrics"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/agent/workers"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/config"
)

func Start() error {
	cfg, err := config.GetAgentCfg()
	if err != nil {
		return err
	}

	l, err := logger.Initialize("info")
	if err != nil {
		return err
	}

	l.Infow("starting agent",
		"address", cfg.AgentEndPoint,
		"pollInterval", cfg.PollInterval,
		"reportInterval", cfg.ReportInterval,
	)

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM)
	dataCh := make(chan []metrics.Metric, 2)

	readWorker := workers.NewReadWorker(l, dataCh, int(cfg.PollInterval))
	sendWorker := workers.NewSender(l, int(cfg.ReportInterval), int(cfg.RateLimit), cfg.AgentEndPoint,
		cfg.AgentKey, dataCh)

	go func() {
		readWorker.Start()
	}()
	go func() {
		sendWorker.Start()
	}()
	<-stopCh
	l.Info("shutting down agent")
	readWorker.Stop()
	sendWorker.Stop()
	return err
}
