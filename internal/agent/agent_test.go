package agent

import (
	m "github.com/shestooy/go-musthave-metrics-tpl.git/internal/agent/metrics"
	r "github.com/shestooy/go-musthave-metrics-tpl.git/internal/httpserver"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/storage"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostMetrics(t *testing.T) {
	s := httptest.NewServer(r.GetRouter())
	storage.MStorage.Init()
	defer s.Close()
	tests := []struct {
		name   string
		values []m.Metric
	}{
		{name: "TestWithMetrics", values: m.GetAllMetrics()},
		{name: "TestEmptyValues", values: []m.Metric{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NoError(t, postMetrics(s.URL, tt.values))
		})
	}

}
