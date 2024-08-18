package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	gauge   = "gauge"
	counter = "counter"
)

type Metric struct {
	Type  string
	Name  string
	Value interface{}
}

var pollCount int64

func getAllMetrics() []Metric {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	pollCount++
	return []Metric{
		{Type: gauge, Name: "Alloc", Value: m.Alloc},
		{Type: gauge, Name: "BuckHashSys", Value: m.BuckHashSys},
		{Type: gauge, Name: "Frees", Value: m.Frees},
		{Type: gauge, Name: "GCCPUFraction", Value: m.GCCPUFraction},
		{Type: gauge, Name: "GCSys", Value: m.GCSys},
		{Type: gauge, Name: "HeapAlloc", Value: m.HeapAlloc},
		{Type: gauge, Name: "HeapIdle", Value: m.HeapIdle},
		{Type: gauge, Name: "HeapInuse", Value: m.HeapInuse},
		{Type: gauge, Name: "HeapObjects", Value: m.HeapObjects},
		{Type: gauge, Name: "HeapReleased", Value: m.HeapReleased},
		{Type: gauge, Name: "HeapSys", Value: m.HeapSys},
		{Type: gauge, Name: "LastGC", Value: m.LastGC},
		{Type: gauge, Name: "Lookups", Value: m.Lookups},
		{Type: gauge, Name: "MCacheInuse", Value: m.MCacheInuse},
		{Type: gauge, Name: "MCacheSys", Value: m.MCacheSys},
		{Type: gauge, Name: "MSpanInuse", Value: m.MCacheInuse},
		{Type: gauge, Name: "MSpanSys", Value: m.MSpanSys},
		{Type: gauge, Name: "Mallocs", Value: m.Mallocs},
		{Type: gauge, Name: "NextGC", Value: m.NextGC},
		{Type: gauge, Name: "NumForcedGC", Value: m.NumForcedGC},
		{Type: gauge, Name: "NumGC", Value: m.NumGC},
		{Type: gauge, Name: "OtherSys", Value: m.GCSys},
		{Type: gauge, Name: "PauseTotalNs", Value: m.PauseTotalNs},
		{Type: gauge, Name: "StackInuse", Value: m.StackInuse},
		{Type: gauge, Name: "StackSys", Value: m.GCSys},
		{Type: gauge, Name: "Sys", Value: m.GCSys},
		{Type: gauge, Name: "TotalAlloc", Value: m.TotalAlloc},
		{Type: counter, Name: "PollCount", Value: pollCount},
		{Type: gauge, Name: "RandomValue", Value: rand.Float64()},
	}
}

func postMetrics(url string, metrics []Metric) error {
	client := resty.New()
	url, _ = strings.CutPrefix(url, "http://")

	for _, metric := range metrics {
		resp, err := client.R().SetPathParams(map[string]string{
			"type":  metric.Type,
			"name":  metric.Name,
			"value": fmt.Sprintf("%v", metric.Value),
		}).SetHeader("Content-Type", "text/plain").Post("http://" + url + "/update/{type}/{name}/{value}")

		if err != nil {
			log.Printf("error send request: %s. Name metric: %s", err, metric.Name)
			return err
		}

		if resp.StatusCode() != http.StatusOK {
			log.Printf("unexpected status code. Expected code 200, got %d. Name metric: %s", resp.StatusCode(), metric.Name)
			return errors.New("unexpected status code")
		}
	}
	pollCount = 0
	return nil
}

func start() error {
	pollTicker := time.NewTicker(time.Duration(pollInterval) * time.Second)
	defer pollTicker.Stop()

	reportTicker := time.NewTicker(time.Duration(reportInterval) * time.Second)
	defer reportTicker.Stop()

	metrics := make([]Metric, 0)

	for {
		select {
		case <-pollTicker.C:
			metrics = append(metrics, getAllMetrics()...)

		case <-reportTicker.C:
			err := postMetrics(agentEndPoint, metrics)
			if err != nil {
				return err
			}
			metrics = make([]Metric, 0)
		}
	}
}

func main() {
	parseFlag()
	err := start()
	if err != nil {
		log.Fatal("send metrics failed")
	}
}
