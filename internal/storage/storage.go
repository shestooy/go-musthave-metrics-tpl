package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"

	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/config"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/server/model"
	"go.uber.org/zap"
)

const (
	gauge   = "gauge"
	counter = "counter"
)

type IStorage interface {
	Init(_ context.Context, l *zap.SugaredLogger, cfg *config.ServerCfg) error
	SaveMetric(ctx context.Context, m model.Metrics) (model.Metrics, error)
	GetAllMetrics(ctx context.Context) (map[string]model.Metrics, error)
	SaveMetrics(ctx context.Context, metrics []model.Metrics) ([]model.Metrics, error)
	GetByID(ctx context.Context, id string) (model.Metrics, error)
	Ping(ctx context.Context) error
	Close() error
}

type Storage struct {
	Metrics map[string]model.Metrics
	mu      sync.RWMutex
	logger  *zap.SugaredLogger
	cfg     *config.ServerCfg
}

func (m *Storage) Init(ctx context.Context, l *zap.SugaredLogger, cfg *config.ServerCfg) error {
	m.mu.Lock()
	m.Metrics = make(map[string]model.Metrics)
	m.mu.Unlock()
	m.logger = l
	m.cfg = cfg
	go m.startSaveMetrics(ctx)
	return m.restore(ctx)
}

func (m *Storage) SaveMetric(ctx context.Context, metric model.Metrics) (model.Metrics, error) {
	select {
	case <-ctx.Done():
		return model.Metrics{}, ctx.Err()
	default:
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	if metric.MType != gauge && metric.MType != counter {
		return model.Metrics{}, errors.New("non correct type of metric")
	}

	if _, ok := m.Metrics[metric.ID]; metric.MType == counter && ok {
		*m.Metrics[metric.ID].Delta += *metric.Delta
		metric.Delta = m.Metrics[metric.ID].Delta
	} else {
		m.Metrics[metric.ID] = metric
	}

	if m.cfg.StorageInterval == 0 {
		return metric, m.writeInFile(ctx)
	}
	return metric, nil
}

func (m *Storage) SaveMetrics(ctx context.Context, metrics []model.Metrics) ([]model.Metrics, error) {
	ans := make([]model.Metrics, 0)
	for _, metric := range metrics {
		newMetric, err := m.SaveMetric(ctx, metric)
		if err != nil {
			return nil, err
		}
		ans = append(ans, newMetric)
	}
	return ans, nil
}

func (m *Storage) GetByID(ctx context.Context, id string) (model.Metrics, error) {
	select {
	case <-ctx.Done():
		return model.Metrics{}, ctx.Err()
	default:
	}

	m.mu.RLock()
	defer m.mu.RUnlock()
	if _, ok := m.Metrics[id]; !ok {
		return model.Metrics{}, errors.New("non correct type of metric")
	}
	return m.Metrics[id], nil
}

func (m *Storage) GetAllMetrics(_ context.Context) (map[string]model.Metrics, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.Metrics, nil
}

func (m *Storage) restore(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if !m.cfg.Restore {
		return nil
	}

	f, err := os.OpenFile(m.cfg.FileStoragePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var metric = &model.Metrics{}
		if err = json.Unmarshal(scanner.Bytes(), &metric); err != nil {
			return err
		}
		_, err = m.SaveMetric(ctx, *metric)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Storage) writeInFile(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	f, err := os.OpenFile(m.cfg.FileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	for _, metric := range m.Metrics {
		if err = enc.Encode(metric); err != nil {
			return err
		}
	}
	return nil
}

func (m *Storage) Ping(_ context.Context) error {
	return errors.New("not supported")
}

func (m *Storage) Close() error {
	if err := m.writeInFile(context.Background()); err != nil {
		m.logger.Info("error saving metrics", zap.Error(err))
		return err
	}
	m.logger.Info("Last save in file complete")
	return nil
}

func (m *Storage) startSaveMetrics(ctx context.Context) {
	if m.cfg.StorageInterval > 0 {
		ticker := time.NewTicker(time.Duration(m.cfg.StorageInterval) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := m.writeInFile(ctx); err != nil {
					m.logger.Info("error saving metrics", zap.Error(err))
				}
			case <-ctx.Done():
				return
			}
		}
	}
}
