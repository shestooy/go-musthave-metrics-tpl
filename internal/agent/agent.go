package agent

import (
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/agent/metrics"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/agent/workers"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/flags"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

func Start() error {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return err
	}
	defer logger.Sync()

	sugarLogger := logger.Sugar()

	sugarLogger.Infow("starting agent",
		"address", flags.AgentEndPoint,
		"pollInterval", flags.PollInterval,
		"reportInterval", flags.ReportInterval,
	)
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM)
	dataCh := make(chan []metrics.Metric, 2)

	readWorker := workers.NewReadWorker(sugarLogger, dataCh, int(flags.PollInterval))
	sendWorker := workers.NewSender(sugarLogger, int(flags.ReportInterval), int(flags.RateLimit), flags.AgentEndPoint,
		flags.AgentKey, dataCh)

	go func() {
		readWorker.Start()
	}()
	go func() {
		sendWorker.Start()
	}()
	<-stopCh
	logger.Info("shutting down agent")
	readWorker.Stop()
	sendWorker.Stop()
	return err
}
