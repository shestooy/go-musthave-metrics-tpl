package main

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
			assert.NotEmpty(t, getAllMetrics(), "failed to collect metrics")
		})
	}
}

// TODO: Спросить у ментора как тут лучше сделать
//func Test_postMetrics(t *testing.T) {
//	tests := []struct {
//		name   string
//		values []Metric
//	}{
//		{name: "TestWithMetrics", values: getAllMetrics()},
//		{name: "TestEmptyValues", values: []Metric{}},
//	}
//	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		handlers.ChangeMetric(w, r)
//	}))
//	defer server.Close()
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			assert.NoError(t, postMetrics(tt.values))
//		})
//	}
//}

//func Test_start(t *testing.T) {
//	tests := []struct {
//		name string
//	}{
//
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			start()
//		})
//	}
//}
