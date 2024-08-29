package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/server/model"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/storage"
	"html/template"
	"log"
	"net/http"
	"strings"
)

func PostMetricsWithJSON(res http.ResponseWriter, req *http.Request) {
	if !strings.Contains(req.Header.Get("Content-Type"), "application/json") {
		http.Error(res, "bad request", http.StatusBadRequest)
		return
	}

	var m = model.Metrics{}
	if err := json.NewDecoder(req.Body).Decode(&m); err != nil {
		http.Error(res, "bad request", http.StatusBadRequest)
		return
	}

	err := storage.MStorage.UpdateMetric(m)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	m, err = storage.MStorage.GetMetricID(m.ID)
	if err != nil {
		http.Error(res, err.Error(), http.StatusNotFound)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	resp, err := json.Marshal(&m)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
	_, err = res.Write(resp)
	if err != nil {
		log.Println(err.Error())
	}
	err = req.Body.Close()
	if err != nil {
		log.Println(err.Error())
	}
}

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
	var m = model.Metrics{
		MType: params[0],
		ID:    params[1],
	}
	err := m.SetValue(params[2])
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	err = storage.MStorage.UpdateMetric(m)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	res.Header().Set("Content-Type", "text/plain")
}

func GetMetricIDWithJSON(res http.ResponseWriter, req *http.Request) {
	if !strings.Contains(req.Header.Get("Content-Type"), "application/json") {
		http.Error(res, "bad request", http.StatusBadRequest)
		return
	}
	var m = model.Metrics{}
	if err := json.NewDecoder(req.Body).Decode(&m); err != nil {
		http.Error(res, "bad request", http.StatusBadRequest)
		return
	}
	m, err := storage.MStorage.GetMetricID(m.ID)
	if err != nil {
		http.Error(res, err.Error(), http.StatusNotFound)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	resp, err := json.Marshal(&m)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
	_, err = res.Write(resp)
	if err != nil {
		log.Println(err.Error())
	}
	err = req.Body.Close()
	if err != nil {
		log.Println(err.Error())
	}
}

func GetMetricID(res http.ResponseWriter, req *http.Request) {
	params := make([]string, 2)
	params[0] = chi.URLParam(req, "type")
	params[1] = chi.URLParam(req, "name")

	for _, param := range params {
		if param == "" {
			http.Error(res, "invalid params", http.StatusNotFound)
			return
		}
	}
	m, err := storage.MStorage.GetMetricID(params[1])
	if err != nil {
		http.Error(res, err.Error(), http.StatusNotFound)
		return
	}
	res.Header().Set("Content-Type", "text/plain")
	_, err = res.Write([]byte(m.GetValueAsString()))
	if err != nil {
		log.Println(err.Error())
	}
}

func GetAllMetrics(res http.ResponseWriter, _ *http.Request) {
	metrics := storage.MStorage.GetAllMetrics()

	counters := make(map[string]model.Metrics)
	gauges := make(map[string]model.Metrics)

	for id, metric := range metrics {
		if metric.MType == "counter" {
			counters[id] = metric
		}
		if metric.MType == "gauge" {
			gauges[id] = metric
		}
	}

	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmp := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Metrics</title>
		</head>
		<body>
			<h1>Metrics</h1>

			<h2>Counters</h2>
			<table border="1">
				<tr>
					<th>Name</th>
					<th>Value</th>
				</tr>
				{{ range $name, $metric := .Counters }}
				<tr>
					<td>{{ $name }}</td>
					<td>{{ GetValueAsString $metric }}</td>
				</tr>
				{{ end }}
			</table>

			<h2>Gauges</h2>
			<table border="1">
				<tr>
					<th>Name</th>
					<th>Value</th>
				</tr>
				{{ range $name, $metric := .Gauges }}
				<tr>
					<td>{{ $name }}</td>
					<td>{{ GetValueAsString $metric }}</td>
				</tr>
				{{ end }}
			</table>

		</body>
		</html>
		`

	t, err := template.New("metrics").
		Funcs(template.FuncMap{"GetValueAsString": func(m model.Metrics) string {
			return m.GetValueAsString()
		},
		}).Parse(tmp)
	if err != nil {
		http.Error(res, "не удалось создать шаблон", http.StatusInternalServerError)
		return
	}

	data := struct {
		Counters map[string]model.Metrics
		Gauges   map[string]model.Metrics
	}{
		Counters: counters,
		Gauges:   gauges,
	}

	err = t.Execute(res, data)
	if err != nil {
		http.Error(res, "не удалось выполнить шаблон", http.StatusInternalServerError)
	}
}
