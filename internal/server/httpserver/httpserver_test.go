package httpserver

import (
	"testing"

	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/config"
	l "github.com/shestooy/go-musthave-metrics-tpl.git/internal/logger"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/storage"
	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestGetRouter(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "TestGetRouter"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			MStorage := &storage.Storage{}
			logger, err := l.Initialize("info")
			require.NoError(t, err)
			s := initServer(MStorage, logger, &config.ServerCfg{ServerKey: ""})
			assert.NotEmpty(t, s)
		})
	}
}
