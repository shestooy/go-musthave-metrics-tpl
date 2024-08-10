package handlers

import (
	"fmt"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/storage"
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
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
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "text/plain")
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
	value, err := storage.Storage.GetMetricID(params[0], params[1])
	if err != nil {
		http.Error(res, err.Error(), http.StatusNotFound)
		return
	}
	ans := fmt.Sprintf("%v", value)
	res.Header().Set("Content-Type", "text/plain")
	_, err = res.Write([]byte(ans))
	if err != nil {
		log.Println(err.Error())
	}
}

func GetAllMetrics(res http.ResponseWriter, req *http.Request) {
	metrics := storage.Storage.GetAllMetrics()

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
				{{ range $name, $value := .Counters }}
				<tr>
					<td>{{ $name }}</td>
					<td>{{ $value }}</td>
				</tr>
				{{ end }}
			</table>

			<h2>Gauges</h2>
			<table border="1">
				<tr>
					<th>Name</th>
					<th>Value</th>
				</tr>
				{{ range $name, $value := .Gauges }}
				<tr>
					<td>{{ $name }}</td>
					<td>{{ $value }}</td>
				</tr>
				{{ end }}
			</table>

		</body>
		</html>
		`

	t, err := template.New("metrics").Parse(tmp)
	if err != nil {
		http.Error(res, "не удалось создать шаблон", http.StatusInternalServerError)
		return
	}

	data := struct {
		Counters interface{}
		Gauges   interface{}
	}{
		Counters: metrics["counter"].GetAllValue(),
		Gauges:   metrics["gauge"].GetAllValue(),
	}

	err = t.Execute(res, data)
	if err != nil {
		http.Error(res, "не удалось выполнить шаблон", http.StatusInternalServerError)
	}
}
