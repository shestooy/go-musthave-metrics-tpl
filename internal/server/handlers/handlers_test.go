package handlers

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/logger"

	"github.com/labstack/echo/v4"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/config"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/server/middlewares"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/storage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string) int {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)
	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	err = resp.Body.Close()
	require.NoError(t, err)

	return resp.StatusCode
}

func testServer(t *testing.T) *echo.Echo {

	MStorage := &storage.Storage{}
	l, err := logger.Initialize("info")
	require.NoError(t, err)
	err = MStorage.Init(context.Background(), l, &config.ServerCfg{Restore: false, StorageInterval: 5000})
	require.NoError(t, err)

	h := NewHandler(l, MStorage)
	e := echo.New()

	e.HideBanner = true
	e.HidePort = true

	e.Use(middlewares.Gzip)
	e.Use(middlewares.GetLogg(l))

	e.GET("/", h.GetAllMetrics)
	e.GET("/ping", h.PingHandler)

	e.POST("/updates/", h.UpdateSomeMetrics)

	updateGroup := e.Group("/update")
	updateGroup.POST("/", h.PostMetricsWithJSON)
	updateGroup.POST("/:type/:name/:value", h.PostMetrics)

	valueGroup := e.Group("/value")
	valueGroup.POST("/", h.GetMetricIDWithJSON)
	valueGroup.GET("/:type/:name", h.GetMetricID)
	return e
}

func TestChangeMetric(t *testing.T) {
	ts := httptest.NewServer(testServer(t))
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

func TestGetMetricID(t *testing.T) {
	ts := httptest.NewServer(testServer(t))
	defer ts.Close()
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
	ts := httptest.NewServer(testServer(t))
	defer ts.Close()
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

func TestPostMetricsWithJSON(t *testing.T) {
	ts := httptest.NewServer(testServer(t))
	defer ts.Close()

	tests := []struct {
		name         string
		method       string
		body         string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "method_get",
			method:       http.MethodGet,
			expectedCode: http.StatusMethodNotAllowed,
			expectedBody: "",
		},
		{
			name:         "method_put",
			method:       http.MethodPut,
			expectedCode: http.StatusMethodNotAllowed,
			expectedBody: "",
		},
		{
			name:         "method_delete",
			method:       http.MethodDelete,
			expectedCode: http.StatusMethodNotAllowed,
			expectedBody: "",
		},
		{
			name:         "method_post_without_body",
			method:       http.MethodPost,
			expectedCode: http.StatusBadRequest,
			expectedBody: "",
		},
		{
			name:         "method_post_unsupported_type",
			method:       http.MethodPost,
			body:         `{"request": {"type": "idunno", "command": "do something"}, "version": "1.0"}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: "",
		},
		{
			name:   "method_post_success",
			method: http.MethodPost,
			body:   `{"id": "temperature","type": "counter", "delta": 34}`,

			expectedCode: http.StatusOK,
			expectedBody: `{"id": "temperature","type": "counter", "delta": 34}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, ts.URL+"/update/", bytes.NewReader([]byte(tt.body)))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			resp, err := ts.Client().Do(req)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedCode, resp.StatusCode, "unexpected response code")
			err = resp.Body.Close()
			require.NoError(t, err)
		})
	}
}

func TestGetMetricIDWithJSON(t *testing.T) {
	ts := httptest.NewServer(testServer(t))
	defer ts.Close()

	prepReq, err := http.NewRequest(http.MethodPost, ts.URL+"/update/",
		bytes.NewReader([]byte(`{"id": "temperature","type": "counter", "delta": 34}`)))

	require.NoError(t, err)
	prepReq.Header.Set("Content-Type", "application/json")

	prepResp, err := ts.Client().Do(prepReq)
	require.NoError(t, err)
	err = prepResp.Body.Close()
	require.NoError(t, err)

	tests := []struct {
		name         string
		method       string
		body         string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "method_get",
			method:       http.MethodGet,
			expectedCode: http.StatusMethodNotAllowed,
			expectedBody: "",
		},
		{
			name:         "method_put",
			method:       http.MethodPut,
			expectedCode: http.StatusMethodNotAllowed,
			expectedBody: "",
		},
		{
			name:         "method_delete",
			method:       http.MethodDelete,
			expectedCode: http.StatusMethodNotAllowed,
			expectedBody: "",
		},
		{
			name:         "method_post_without_body",
			method:       http.MethodPost,
			expectedCode: http.StatusNotFound,
			expectedBody: "",
		},
		{
			name:         "method_post_unsupported_type",
			method:       http.MethodPost,
			body:         `{"request": {"type": "idunno", "command": "do something"}, "version": "1.0"}`,
			expectedCode: http.StatusNotFound,
			expectedBody: "",
		},
		{
			name:   "method_post_success",
			method: http.MethodPost,
			body:   `{"id": "temperature","type": "counter"}`,

			expectedCode: http.StatusOK,
			expectedBody: `{"id": "temperature","type": "counter","delta": 6}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, ts.URL+"/value/", bytes.NewReader([]byte(tt.body)))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			resp, err := ts.Client().Do(req)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedCode, resp.StatusCode, "unexpected response code")
			err = resp.Body.Close()
			require.NoError(t, err)
		})
	}
}
