package storage

import (
	"errors"

	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/model"
)

var Storage = MemStorage{}

type IStorage interface {
	UpdateMetric(t, k, v string) error
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
