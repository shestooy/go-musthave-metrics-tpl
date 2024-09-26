package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getAllMetrics(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"TestMetricsCollector"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotEmpty(t, GetRuntimeMetrics(), "failed to collect metrics")
		})
	}
}
