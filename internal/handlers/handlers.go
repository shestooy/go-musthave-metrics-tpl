package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/storage"
	"net/http"
)

func PostMetrics(res http.ResponseWriter, req *http.Request) {
	params := make([]string, 3)
	params[0] = chi.URLParam(req, "type")
	params[1] = chi.URLParam(req, "name")
	params[2] = chi.URLParam(req, "value")

	for _, param := range params {
		if param == "" {
			http.Error(res, "invalid params", http.StatusBadRequest)
			return
		}
	}

	err := storage.Storage.UpdateMetric(params[0], params[1], params[2])
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	res.Header().Set("Content-Type", "text/plain")
}

func GetMetricId(res http.ResponseWriter, req *http.Request) {
	params := make([]string, 2)
	params[0] = chi.URLParam(req, "type")
	params[1] = chi.URLParam(req, "name")

	for _, param := range params {
		if param == "" {
			http.Error(res, "invalid params", http.StatusNotFound)
			return
		}
	}
	value, err := storage.Storage.GetMetricId(params[0], params[1])
	if err != nil {
		http.Error(res, err.Error(), http.StatusNotFound)
		return
	}
	ans := fmt.Sprintf("%s is equal to %v\n", params[1], value)
	res.Header().Set("Content-Type", "text/plain")
	res.Write([]byte(ans))
}

func GetAllMetrics(res http.ResponseWriter, req *http.Request) {

}
