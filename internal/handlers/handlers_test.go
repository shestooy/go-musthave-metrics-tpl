package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChangeMetric(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		expectedCode int
		expectedBody string
	}{
		{name: "GetTest", method: http.MethodGet, expectedCode: http.StatusMethodNotAllowed, expectedBody: ""},
		{name: "PostTest", method: http.MethodPost, expectedCode: http.StatusNotFound, expectedBody: ""},
		{name: "DeleteTest", method: http.MethodDelete, expectedCode: http.StatusMethodNotAllowed, expectedBody: ""},
		{name: "PutTest", method: http.MethodPut, expectedCode: http.StatusMethodNotAllowed, expectedBody: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(tt.method, "/update/", nil)
			require.NoError(t, err)
			w := httptest.NewRecorder()

			PostMetrics(w, r)
			assert.Equal(t, tt.expectedCode, w.Code, "unexpected response code")
		})
	}
}

func TestGetAllMetrics(t *testing.T) {
	type args struct {
		res http.ResponseWriter
		req *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetAllMetrics(tt.args.res, tt.args.req)
		})
	}
}

func TestGetMetricId(t *testing.T) {
	type args struct {
		res http.ResponseWriter
		req *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetMetricId(tt.args.res, tt.args.req)
		})
	}
}
