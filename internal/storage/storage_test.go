package storage

import (
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemStorage_Init(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "TestInitCounter"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &MemStorage{}
			c.Init()
			assert.NotEmpty(t, tests)
		})
	}
}

func TestMemStorage_UpdateMetric(t *testing.T) {
	type args struct {
		t string
		k string
		v string
	}
	tests := []struct {
		name    string
		values  args
		wantErr bool
	}{
		{name: "TestTypeGauge", values: args{t: "gauge", k: "name", v: "4354"}, wantErr: false},
		{name: "TestTypeCounter", values: args{t: "counter", k: "test", v: "4324"}, wantErr: false},
		{name: "TestTypeStr", values: args{t: "str", k: "test", v: "4324"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{}
			m.Init()
			err := m.UpdateMetric(tt.values.t, tt.values.k, tt.values.v)
			assert.Equal(t, tt.wantErr, err != nil, "unexpected error")
		})
	}
}

func TestMemStorage_GetMetricID(t *testing.T) {

	tests := []struct {
		name   string
		values map[string]model.Metrics
		args   []string
		want   interface{}
	}{
		{name: "TestGetMetricByID", values: map[string]model.Metrics{"gauge": &model.Gauge{Values: map[string]float64{"test": 1.5}}},
			args: []string{"gauge", "test"}, want: 1.5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				Metrics: tt.values,
			}
			v, err := m.GetMetricID(tt.args[0], tt.args[1])
			require.NoError(t, err)
			assert.Equal(t, tt.want, v)
		})
	}
}
