package main

import (
	"net/http"

	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/handlers"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/storage"
)

func main() {
	storage.Storage.Init()
	err := start()
	if err != nil {
		panic(err)
	}
}

func start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/update/{type}/{name}/{value}", handlers.ChangeMetric)
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		return err
	}
	return err
}
