package storage

import (
	"bufio"
	"encoding/json"
	"errors"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/flags"
	l "github.com/shestooy/go-musthave-metrics-tpl.git/internal/logger"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/server/model"
	"go.uber.org/zap"
	"os"
	"sync"
)

var MStorage = Storage{}

type Storage struct {
	Metrics map[string]model.Metrics
	mu      sync.RWMutex
}

func (m *Storage) Init() error {
	m.Metrics = make(map[string]model.Metrics)
	return m.restore()
}

func (m *Storage) UpdateMetric(metric model.Metrics) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if metric.MType != "gauge" && metric.MType != "counter" {
		return errors.New("non correct type of metric")
	}

	if _, ok := m.Metrics[metric.ID]; metric.MType == "counter" && ok {
		*m.Metrics[metric.ID].Delta += *metric.Delta
		return nil
	}
	m.Metrics[metric.ID] = metric
	if flags.StorageInterval == 0 {
		return m.WriteInFile()
	}
	return nil
}

func (m *Storage) GetMetricID(id string) (model.Metrics, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if _, ok := m.Metrics[id]; !ok {
		l.Log.Info("ok", zap.Bool("ok", ok))
		l.Log.Info("storage", zap.Any("stor", m.GetAllMetrics()))
		return model.Metrics{}, errors.New("non correct type of metric")
	}
	return m.Metrics[id], nil
}

func (m *Storage) GetAllMetrics() map[string]model.Metrics {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.Metrics
}

func (m *Storage) restore() error {
	if !flags.Restore {
		return nil
	}

	f, err := os.OpenFile(flags.FileStoragePath, os.O_RDONLY|os.O_CREATE, 0666)
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
		err = m.UpdateMetric(*metric)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Storage) WriteInFile() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	f, err := os.OpenFile(flags.FileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
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
