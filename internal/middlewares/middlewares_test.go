package middlewares

import (
	l "github.com/shestooy/go-musthave-metrics-tpl.git/internal/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoggingMiddleware(t *testing.T) {
	err := l.Initialize("info")
	require.NoError(t, err)

	testHandler := http.HandlerFunc(func(res http.ResponseWriter, _ *http.Request) {
		_, err = res.Write([]byte("test"))
		require.NoError(t, err)
	})
	require.NotEmpty(t, testHandler)

	middleware := LoggingMiddleware(testHandler)
	require.NotEmpty(t, middleware)

	req, err := http.NewRequest(http.MethodGet, "/test", nil)
	require.NoError(t, err)
	require.NotEmpty(t, req)

	rr := httptest.NewRecorder()
	require.NotEmpty(t, rr)

	middleware.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, len("test"), rr.Body.Len())
}
