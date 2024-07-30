package storage

import (
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestMemStorage_GetAllMetrics(t *testing.T) {
	type fields struct {
		Metrics map[string]model.Metrics
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]model.Metrics
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				Metrics: tt.fields.Metrics,
			}
			assert.Equalf(t, tt.want, m.GetAllMetrics(), "GetAllMetrics()")
		})
	}
}
