package storage

import (
	"context"
	"testing"

	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/logger"

	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/config"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/server/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemStorage_Init(t *testing.T) {
	MStorage := &Storage{}
	l, err := logger.Initialize("info")
	require.NoError(t, err)
	err = MStorage.Init(context.Background(), l, &config.ServerCfg{StorageInterval: 5000, Restore: false})
	require.NoError(t, err)
}

func TestStorage_UpdateMetric(t *testing.T) {
	MStorage := &Storage{}
	l, err := logger.Initialize("info")
	require.NoError(t, err)
	err = MStorage.Init(context.Background(), l, &config.ServerCfg{StorageInterval: 5000, Restore: false})
	require.NoError(t, err)

	metric := model.Metrics{
		ID:    "test_gauge",
		MType: "gauge",
		Value: func(v float64) *float64 { return &v }(10),
	}

	m, err := MStorage.SaveMetric(context.Background(), metric)
	require.NotEmpty(t, m)
	require.NoError(t, err)

	storedMetric, err := MStorage.GetByID(context.Background(), metric.ID)
	require.NoError(t, err)
	assert.Equal(t, metric, storedMetric)
}

func TestMemStorage_GetMetricID(t *testing.T) {
	MStorage := &Storage{}
	l, err := logger.Initialize("info")
	require.NoError(t, err)
	err = MStorage.Init(context.Background(), l, &config.ServerCfg{StorageInterval: 5000, Restore: false})
	require.NoError(t, err)

	metric := model.Metrics{
		ID:    "test_gauge",
		MType: "gauge",
		Value: func(v float64) *float64 { return &v }(20),
	}

	m, err := MStorage.SaveMetric(context.Background(), metric)
	require.NotEmpty(t, m)
	require.NoError(t, err)

	retrievedMetric, err := MStorage.GetByID(context.Background(), metric.ID)
	require.NoError(t, err)
	require.Equal(t, metric, retrievedMetric)

	_, err = MStorage.GetByID(context.Background(), "testErr")
	assert.Error(t, err)
}

func TestStorage_GetAllMetrics(t *testing.T) {
	MStorage := &Storage{}
	l, err := logger.Initialize("info")
	require.NoError(t, err)
	err = MStorage.Init(context.Background(), l, &config.ServerCfg{StorageInterval: 5000, Restore: false})
	require.NoError(t, err)

	metric1 := model.Metrics{
		ID:    "test_gauge",
		MType: "gauge",
		Value: func(v float64) *float64 { return &v }(450.2),
	}
	metric2 := model.Metrics{
		ID:    "test_counter",
		MType: "counter",
		Delta: func(v int64) *int64 { return &v }(432),
	}

	m, err := MStorage.SaveMetric(context.Background(), metric1)
	require.NotEmpty(t, m)
	require.NoError(t, err)
	m, err = MStorage.SaveMetric(context.Background(), metric2)
	require.NotEmpty(t, m)
	require.NoError(t, err)

	allMetrics, err := MStorage.GetAllMetrics(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 2, len(allMetrics))
	assert.Contains(t, allMetrics, metric1.ID)
	assert.Contains(t, allMetrics, metric2.ID)
}
