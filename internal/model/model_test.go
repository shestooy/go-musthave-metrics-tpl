package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetrics_GetValue(t *testing.T) {
	tests := []struct {
		name     string
		metric   Metrics
		expected string
	}{
		{
			name: "gaugeWithValue",
			metric: Metrics{
				ID:    "test_gauge",
				MType: "gauge",
				Value: func() *float64 { v := 42.42; return &v }(),
			},
			expected: "42.42",
		},
		{
			name: "gaugeWithOutValue",
			metric: Metrics{
				ID:    "test_gauge_nil",
				MType: "gauge",
				Value: nil,
			},
			expected: "",
		},
		{
			name: "counterWithDelta",
			metric: Metrics{
				ID:    "test_counter",
				MType: "counter",
				Delta: func() *int64 { d := int64(10); return &d }(),
			},
			expected: "10",
		},
		{
			name: "counterWithoutDelta",
			metric: Metrics{
				ID:    "test_counter_nil",
				MType: "counter",
				Delta: nil,
			},
			expected: "",
		},
		{
			name: "unsupported metric type",
			metric: Metrics{
				ID:    "test_unsupported",
				MType: "unsupported",
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.metric.GetValue())
		})
	}
}

func TestMetrics_SetValue(t *testing.T) {
	tests := []struct {
		name          string
		metric        Metrics
		valueToSet    interface{}
		expectedVal   *float64
		expectedDelta *int64
	}{
		{
			name: "setGaugeValue",
			metric: Metrics{
				ID:    "test_gauge",
				MType: "gauge",
			},
			valueToSet:  42.42,
			expectedVal: func() *float64 { v := 42.42; return &v }(),
		},
		{
			name: "setCounterDelta",
			metric: Metrics{
				ID:    "test_counter",
				MType: "counter",
			},
			valueToSet:    int64(10),
			expectedDelta: func() *int64 { d := int64(10); return &d }(),
		},
		{
			name: "unsupportedMetricType",
			metric: Metrics{
				ID:    "test_unsupported",
				MType: "unsupported",
			},
			valueToSet: 42.42,
		},
		{
			name: "setGaugeValueNil",
			metric: Metrics{
				ID:    "test_gauge_nil",
				MType: "gauge",
				Value: nil,
			},
			valueToSet:  42.42,
			expectedVal: func() *float64 { v := 42.42; return &v }(),
		},
		{
			name: "setCounterDeltaNil",
			metric: Metrics{
				ID:    "test_counter_nil",
				MType: "counter",
				Delta: nil,
			},
			valueToSet:    int64(10),
			expectedDelta: func() *int64 { d := int64(10); return &d }(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.metric.SetValue(tt.valueToSet)

			switch tt.metric.MType {
			case "gauge":
				if tt.expectedVal != nil {
					assert.NotNil(t, tt.metric.Value)
					assert.Equal(t, *tt.expectedVal, *tt.metric.Value)
				} else {
					assert.Nil(t, tt.metric.Value)
				}

			case "counter":
				if tt.expectedDelta != nil {
					assert.NotNil(t, tt.metric.Delta)
					assert.Equal(t, *tt.expectedDelta, *tt.metric.Delta)
				} else {
					assert.Nil(t, tt.metric.Delta)
				}
			}

		})
	}
}
