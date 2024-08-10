package main

import (
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/handlers"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/storage"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func GetRouter() chi.Router {
	storage.Storage.Init()
	r := chi.NewRouter()

	r.Post("/update/{type}/{name}/{value}", handlers.PostMetrics)
	r.Get("/value/{type}/{name}", handlers.GetMetricID)
	r.Get("/", handlers.GetAllMetrics)
	return r
}

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

func TestPostMetrics(t *testing.T) {
	s := httptest.NewServer(GetRouter())
	defer s.Close()
	tests := []struct {
		name   string
		values []Metric
	}{
		{name: "TestWithMetrics", values: getAllMetrics()},
		{name: "TestEmptyValues", values: []Metric{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NoError(t, postMetrics(s.URL, tt.values))
		})
	}

}

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
