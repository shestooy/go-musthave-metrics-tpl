package storage

import (
	"errors"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/storage/types"
)

var MStorage = Storage{}

type IModel interface {
	Init()
	UpdateMetric(t, k, v string) error
	GetMetricID(t, n string) (interface{}, error)
	GetAllMetrics() map[string]types.Metrics
}

type Storage struct {
	Metrics map[string]types.Metrics
}

func (m *Storage) Init() {
	m.Metrics = make(map[string]types.Metrics)
	m.Metrics["gauge"] = &types.Gauge{}
	m.Metrics["gauge"].Init()
	m.Metrics["counter"] = &types.Counter{}
	m.Metrics["counter"].Init()
}

func (m *Storage) UpdateMetric(t, k, v string) error {
	if _, ok := m.Metrics[t]; !ok {
		return errors.New("non correct type of metric")
	}
	err := m.Metrics[t].AddValue(k, v)
	if err != nil {
		return err
	}
	return nil
}

func (m *Storage) GetMetricID(t, n string) (interface{}, error) {
	if _, ok := m.Metrics[t]; !ok {
		return 0, errors.New("non correct type of metric")
	}
	value, err := m.Metrics[t].GetValueID(n)
	if err != nil {
		return value, err
	}
	return value, nil
}

func (m *Storage) GetAllMetrics() map[string]types.Metrics {
	return m.Metrics
}
