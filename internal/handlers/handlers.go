package handlers

import (
	"net/http"
	"strings"

	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/storage"
)

func ChangeMetric(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Invalid Method", http.StatusMethodNotAllowed)
		return
	}
	parts := strings.Split(req.URL.String(), "/")
	if len(parts) != 5 {
		http.NotFound(res, req)
		return
	}
	err := storage.Storage.UpdateMetric(parts[2], parts[3], parts[4])
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}
	res.Header().Set("Content-Type", "text/plain")
}
