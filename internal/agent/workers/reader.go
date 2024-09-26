package workers

import (
	"time"

	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/agent/metrics"
	"go.uber.org/zap"
)

type ReadWorker struct {
	logger         *zap.SugaredLogger
	doneChan       chan struct{}
	dataChan       chan []metrics.Metric
	intervalTicker *time.Ticker
}

func NewReadWorker(logger *zap.SugaredLogger, dataChan chan []metrics.Metric, readInterval int) *ReadWorker {
	return &ReadWorker{
		logger:         logger,
		doneChan:       make(chan struct{}),
		dataChan:       dataChan,
		intervalTicker: time.NewTicker(time.Duration(readInterval) * time.Second),
	}
}

func (r *ReadWorker) Start() {
	r.logger.Info("Starting Read Worker")
	for {
		select {
		case <-r.doneChan:
			r.logger.Info("Read Worker is done")
			r.intervalTicker.Stop()
			close(r.dataChan)
			return
		case <-r.intervalTicker.C:
			r.dataChan <- metrics.GetRuntimeMetrics()
			r.dataChan <- metrics.GetMemoryMetrics(r.logger)
			r.dataChan <- []metrics.Metric{{MType: "counter", ID: "PollCount", Delta: func(val int64) *int64 {
				return &val
			}(1)}}
		}
	}
}

func (r *ReadWorker) Stop() {
	close(r.doneChan)
}
