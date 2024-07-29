package storage

import (
	"errors"

	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/model"
)

var Storage = MemStorage{}

type IStorage interface {
	Init()
	UpdateMetric(t, k, v string) error
	GetMetricId(t, n string) (interface{}, error)
	GetAllMetrics() map[string]model.Metrics
}

type MemStorage struct {
	Metrics map[string]model.Metrics
}

func (m *MemStorage) Init() {
	m.Metrics = make(map[string]model.Metrics)
	m.Metrics["gauge"] = &model.Gauge{}
	m.Metrics["gauge"].Init()
	m.Metrics["counter"] = &model.Counter{}
	m.Metrics["counter"].Init()
}

func (m *MemStorage) UpdateMetric(t, k, v string) error {
	if _, ok := m.Metrics[t]; !ok {
		return errors.New("non correct type of metric")
	}
	err := m.Metrics[t].AddValue(k, v)
	if err != nil {
		return err
	}
	return nil
}

func (m *MemStorage) GetMetricId(t, n string) (interface{}, error) {
	if _, ok := m.Metrics[t]; !ok {
		return 0, errors.New("non correct type of metric")
	}
	value, err := m.Metrics[t].GetValueId(n)
	if err != nil {
		return value, err
	}
	return value, nil
}

func (m *MemStorage) GetAllMetrics() map[string]model.Metrics {
	return m.Metrics
}
