package agent

import (
	"context"
	m "github.com/shestooy/go-musthave-metrics-tpl.git/internal/agent/metrics"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/flags"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/server/httpserver"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/storage"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostMetrics(t *testing.T) {
	s := httptest.NewServer(httpserver.GetRouter())
	defer s.Close()

	flags.Restore = false
	flags.StorageInterval = 5000
	storage.MStorage = &storage.Storage{}
	err := storage.MStorage.Init(context.Background())
	require.NoError(t, err)

	tests := []struct {
		name   string
		values []m.Metric
	}{
		{name: "TestWithMetrics", values: m.GetAllMetrics()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := postMetrics(s.URL, tt.values)
			require.NoError(t, err)
			req, err := http.NewRequest(http.MethodGet, s.URL+"/value/gauge/Alloc", http.NoBody)
			require.NoError(t, err)
			resp, err := s.Client().Do(req)
			require.Equal(t, http.StatusOK, resp.StatusCode)
			require.NoError(t, err)
			b, err := io.ReadAll(resp.Body)
			defer func() {
				err := resp.Body.Close()
				require.NoError(t, err)
			}()

			require.NoError(t, err)
			assert.NotEqual(t, "", string(b))
		})
	}

}
