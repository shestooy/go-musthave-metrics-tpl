package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/storage"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string) int {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	return resp.StatusCode
}

func testServer() chi.Router {
	storage.Storage.Init()
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/update/{type}/{name}/{value}", PostMetrics)
	r.Get("/value/{type}/{name}", GetMetricID)
	r.Get("/", GetAllMetrics)
	return r
}

func TestChangeMetric(t *testing.T) {
	ts := httptest.NewServer(testServer())
	defer ts.Close()
	tests := []struct {
		name         string
		method       string
		path         string
		expectedCode int
	}{
		{name: "TestPost_WithParam", method: http.MethodPost, expectedCode: http.StatusOK, path: "/update/gauge/testGet/1.0"},
		{name: "TestNoParam", method: http.MethodPost, expectedCode: http.StatusNotFound, path: "/update/gauge/testNoParam/"},
		{name: "TestAnotherMethod", method: http.MethodDelete, expectedCode: http.StatusMethodNotAllowed, path: "/update/gauge/testDelete/1.0"},
		{name: "PutTest", method: http.MethodPut, expectedCode: http.StatusMethodNotAllowed, path: "/update/gauge/testPut/1.0"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code := testRequest(t, ts, tt.method, tt.path)
			assert.Equal(t, tt.expectedCode, code, "unexpected response code")
		})
	}
}

func TestGetMetricId(t *testing.T) {
	ts := httptest.NewServer(testServer())
	tests := []struct {
		name         string
		method       string
		path         string
		expectedCode int
	}{
		{name: "TestGet_WithParam", method: http.MethodGet, expectedCode: http.StatusOK, path: "/value/gauge/testGet"},
		{name: "TestPostMethod", method: http.MethodPost, expectedCode: http.StatusMethodNotAllowed, path: "/value/gauge/testGet"},
		{name: "TestWithInvalidParam", method: http.MethodGet, expectedCode: http.StatusNotFound, path: "/value/gauge/"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = testRequest(t, ts, http.MethodPost, "/update/gauge/testGet/534.23")
			code := testRequest(t, ts, tt.method, tt.path)
			assert.Equal(t, tt.expectedCode, code, "unexpected response code")
		})
	}
}

func TestGetAllMetrics(t *testing.T) {
	ts := httptest.NewServer(testServer())
	tests := []struct {
		name         string
		method       string
		path         string
		expectedCode int
	}{
		{name: "TestGet", method: http.MethodGet, expectedCode: http.StatusOK, path: "/"},
		{name: "TestPost", method: http.MethodPost, expectedCode: http.StatusMethodNotAllowed, path: "/"},
		{name: "TestPut", method: http.MethodPut, expectedCode: http.StatusMethodNotAllowed, path: "/"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code := testRequest(t, ts, tt.method, tt.path)
			assert.Equal(t, tt.expectedCode, code, "unexpected response code")
		})
	}
}
