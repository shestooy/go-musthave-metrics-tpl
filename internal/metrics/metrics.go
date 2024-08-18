package metrics

import (
	"math/rand"
	"runtime"
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

var PollCount int64

func GetAllMetrics() []Metric {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	PollCount++
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
		{Type: counter, Name: "PollCount", Value: PollCount},
		{Type: gauge, Name: "RandomValue", Value: rand.Float64()},
	}
}
