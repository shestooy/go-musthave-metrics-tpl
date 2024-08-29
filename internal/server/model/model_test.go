package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			assert.Equal(t, tt.expected, tt.metric.GetValueAsString())
		})
	}
}

func TestMetrics_SetValue(t *testing.T) {
	tests := []struct {
		name       string
		metric     Metrics
		valueToSet string
		expected   string
	}{
		{
			name: "setGaugeValue",
			metric: Metrics{
				ID:    "test_gauge",
				MType: "gauge",
			},
			valueToSet: "42.42",
			expected:   "42.42",
		},
		{
			name: "setCounterDelta",
			metric: Metrics{
				ID:    "test_counter",
				MType: "counter",
			},
			valueToSet: "10",
			expected:   "10",
		},
		{
			name: "unsupportedMetricType",
			metric: Metrics{
				ID:    "test_unsupported",
				MType: "unsupported",
			},
			valueToSet: "42.42",
		},
		{
			name: "setGaugeValueNil",
			metric: Metrics{
				ID:    "test_gauge_nil",
				MType: "gauge",
				Value: nil,
			},
			valueToSet: "42.42",
			expected:   "42.42",
		},
		{
			name: "setCounterDeltaNil",
			metric: Metrics{
				ID:    "test_counter_nil",
				MType: "counter",
				Delta: nil,
			},
			valueToSet: "10",
			expected:   "10",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.metric.SetValue(tt.valueToSet)
			require.NoError(t, err)

			assert.Equal(t, tt.expected, tt.metric.GetValueAsString())

		})
	}
}
