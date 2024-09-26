package workers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/avast/retry-go"
	"github.com/go-resty/resty/v2"
	m "github.com/shestooy/go-musthave-metrics-tpl.git/internal/agent/metrics"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/agent/semaphore"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/utils"
	"go.uber.org/zap"
)

type Sender struct {
	logger     *zap.SugaredLogger
	sendTicker *time.Ticker
	address    string
	agentKey   string
	doneChan   chan struct{}
	dataChan   chan []m.Metric
	semaphore  *semaphore.Semaphore
	wg         *sync.WaitGroup
	mu         sync.RWMutex
}

func NewSender(
	logger *zap.SugaredLogger,
	sendInterval, limit int,
	address, agentKey string,
	dataChan chan []m.Metric,
) *Sender {
	return &Sender{
		logger:     logger,
		sendTicker: time.NewTicker(time.Duration(sendInterval) * time.Second),
		address:    address,
		agentKey:   agentKey,
		dataChan:   dataChan,
		doneChan:   make(chan struct{}),
		semaphore:  semaphore.New(limit),
		wg:         &sync.WaitGroup{},
	}
}

func (s *Sender) Start() {
	s.logger.Info("Starting sender...")
	var metrics []m.Metric
	for {
		select {
		case <-s.doneChan:
			s.sendTicker.Stop()
			s.wg.Wait()
			return
		case data := <-s.dataChan:
			metrics = append(metrics, data...)
		case <-s.sendTicker.C:
			s.wg.Add(1)
			go s.sendMetrics(metrics)
			metrics = make([]m.Metric, 0)
		}
	}
}

func (s *Sender) Stop() {
	close(s.doneChan)
}

func (s *Sender) sendMetrics(metrics []m.Metric) {
	s.semaphore.Acquire()
	defer s.semaphore.Release()
	defer s.wg.Done()
	s.mu.RLock()
	url := s.address
	s.mu.RUnlock()

	client := resty.New()
	url, _ = strings.CutPrefix(url, "http://")

	body, err := m.Compress(metrics)
	if err != nil {
		s.logger.Error("error compress procedure", zap.Error(err))
		return
	}
	uncompressedBody, err := m.GetMetricsAsBody(metrics)
	if err != nil {
		s.logger.Error("error get metrics", zap.Error(err))
		return
	}
	hash := getHash(uncompressedBody, s.getAgentKey())

	err = retry.Do(func() error {
		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Content-Encoding", "gzip").
			SetHeader("HashSHA256", hash).
			SetBody(body).
			Post("http://" + url + "/updates/")

		if err != nil {
			s.logger.Error("error sending metrics to agent", zap.String("url", url), zap.Error(err))
			if !utils.IsRetriableError(err) {
				return retry.Unrecoverable(err)
			}
			return err
		}
		if resp.StatusCode() != http.StatusOK {
			s.logger.Error("unexpected status code. Expected code 200, got %d.", resp.StatusCode())
		}
		return nil
	},
		retry.Attempts(4),
		retry.DelayType(utils.RetryDelay))
	if err != nil {
		s.logger.Error(err)
	}
}

func (s *Sender) getAgentKey() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.agentKey
}

func getHash(body []byte, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(body)
	return hex.EncodeToString(h.Sum(nil))
}
