package middlewares

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/hex"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/config"
	l "github.com/shestooy/go-musthave-metrics-tpl.git/internal/logger"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/server/handlers"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoggingMiddleware(t *testing.T) {
	logger, err := l.Initialize("info")
	require.NoError(t, err)

	testHandler := func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	}
	require.NotEmpty(t, testHandler)

	e := echo.New()
	logg := GetLogg(logger)
	middleware := logg(testHandler)
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
	MStorage := &storage.Storage{}
	logger, err := l.Initialize("info")
	require.NoError(t, err)
	err = MStorage.Init(context.Background(), logger, &config.ServerCfg{StorageInterval: 5000, Restore: false})
	require.NoError(t, err)

	h := handlers.NewHandler(logger, MStorage)
	e := echo.New()
	e.POST("/update/", Gzip(h.PostMetricsWithJSON))

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

func TestHashMiddleware(t *testing.T) {
	tests := []struct {
		name               string
		key                string
		requestBody        string
		modifyHash         bool
		expectedStatus     int
		expectedResponse   string
		expectResponseHash bool
	}{
		{
			name:               "Valid key and correct hash",
			key:                "TEST",
			requestBody:        "test body",
			modifyHash:         false,
			expectedStatus:     http.StatusOK,
			expectedResponse:   "Success",
			expectResponseHash: true,
		},
		{
			name:               "Valid key and incorrect hash",
			key:                "TEST",
			requestBody:        "test body",
			modifyHash:         true,
			expectedStatus:     http.StatusBadRequest,
			expectedResponse:   "the hash checksum did not match",
			expectResponseHash: false,
		},
		{
			name:               "Empty key, hash check skipped",
			key:                "",
			requestBody:        "test body",
			modifyHash:         false,
			expectedStatus:     http.StatusOK,
			expectedResponse:   "Success",
			expectResponseHash: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()
			e.Use(Hash(tc.key))

			e.POST("/", func(c echo.Context) error {
				return c.String(http.StatusOK, "Success")
			})

			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tc.requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMETextPlain)

			if tc.key != "" {
				bodyHash := hash([]byte(tc.requestBody), tc.key)
				if tc.modifyHash {
					bodyHash = []byte("invalidHash")
				}
				req.Header.Set("HashSHA256", hex.EncodeToString(bodyHash))
			}

			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			responseBody := rec.Body.String()
			assert.Contains(t, responseBody, tc.expectedResponse)

			if tc.expectResponseHash {
				resHash := rec.Header().Get("HashSHA256")
				assert.NotEmpty(t, resHash)
				expectedResHash := hash([]byte(responseBody), tc.key)
				expectedResHashHex := hex.EncodeToString(expectedResHash)
				assert.Equal(t, expectedResHashHex, resHash)
			} else {
				resHash := rec.Header().Get("HashSHA256")
				assert.Empty(t, resHash)
			}
		})
	}
}
