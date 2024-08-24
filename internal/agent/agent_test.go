package agent

import (
	m "github.com/shestooy/go-musthave-metrics-tpl.git/internal/agent/metrics"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/handlers"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/storage"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostMetrics(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(handlers.GetMetricIDWithJSON))
	storage.MStorage.Init()
	defer s.Close()

	tests := []struct {
		name   string
		values []m.Metric
	}{
		{name: "TestWithMetrics", values: m.GetAllMetrics()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			postMetrics(s.URL, tt.values)
			req, err := http.NewRequest(http.MethodGet, s.URL+"/update/gauge/Alloc", http.NoBody)
			require.NoError(t, err)
			resp, err := s.Client().Do(req)
			require.NoError(t, err)
			b, err := io.ReadAll(resp.Body)
			defer resp.Body.Close()
			require.NoError(t, err)
			assert.NotEqual(t, "", string(b))
		})
	}

}
