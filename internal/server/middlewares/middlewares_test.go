package middlewares

import (
	"bytes"
	"compress/gzip"
	"context"
	"github.com/labstack/echo/v4"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/flags"
	l "github.com/shestooy/go-musthave-metrics-tpl.git/internal/logger"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/server/handlers"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoggingMiddleware(t *testing.T) {
	err := l.Initialize("info")
	require.NoError(t, err)

	testHandler := func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	}
	require.NotEmpty(t, testHandler)

	e := echo.New()

	middleware := Logging(testHandler)
	require.NotEmpty(t, middleware)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	err = middleware(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "test", rec.Body.String())
}

func TestGzipCompression(t *testing.T) {
	storage.MStorage = &storage.Storage{}
	err := storage.MStorage.Init(context.Background())
	require.NoError(t, err)

	flags.Restore = false
	flags.StorageInterval = 5000

	e := echo.New()
	e.POST("/update/", Gzip(handlers.PostMetricsWithJSON))

	srv := httptest.NewServer(e.Server.Handler)
	defer srv.Close()

	requestBody := `{
					  "id": "Alloc",
					  "type": "gauge",
					  "value": 100
					}`

	successBody := `{
					  "id": "Alloc",
					  "type": "gauge",
					  "value": 100
					}`

	t.Run("sends_gzip", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		zb := gzip.NewWriter(buf)
		_, err := zb.Write([]byte(requestBody))
		require.NoError(t, err)
		err = zb.Close()
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, srv.URL+"/update/", buf)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Content-Encoding", "gzip")
		req.RequestURI = ""

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		b, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		err = resp.Body.Close()
		require.NoError(t, err)
		require.JSONEq(t, successBody, string(b))
	})

	t.Run("accepts_gzip", func(t *testing.T) {
		buf := bytes.NewBufferString(requestBody)
		req := httptest.NewRequest(http.MethodPost, srv.URL+"/update/", buf)
		req.Header.Set("Accept-Encoding", "gzip")
		req.Header.Set("Content-Type", "application/json")
		req.RequestURI = ""

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		zr, err := gzip.NewReader(resp.Body)
		assert.NoError(t, err)
		err = resp.Body.Close()
		require.NoError(t, err)
		b, err := io.ReadAll(zr)
		require.NoError(t, err)

		require.JSONEq(t, successBody, string(b))
	})
}
