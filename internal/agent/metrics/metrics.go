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
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

var PollCount int64

func GetAllMetrics() []Metric {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	float64Ptr := func(val float64) *float64 {
		return &val
	}
	return []Metric{
		{MType: gauge, ID: "Alloc", Value: float64Ptr(float64(m.Alloc))},
		{MType: gauge, ID: "BuckHashSys", Value: float64Ptr(float64(m.BuckHashSys))},
		{MType: gauge, ID: "Frees", Value: float64Ptr(float64(m.Frees))},
		{MType: gauge, ID: "GCCPUFraction", Value: float64Ptr(m.GCCPUFraction)},
		{MType: gauge, ID: "GCSys", Value: float64Ptr(float64(m.GCSys))},
		{MType: gauge, ID: "HeapAlloc", Value: float64Ptr(float64(m.HeapAlloc))},
		{MType: gauge, ID: "HeapIdle", Value: float64Ptr(float64(m.HeapIdle))},
		{MType: gauge, ID: "HeapInuse", Value: float64Ptr(float64(m.HeapInuse))},
		{MType: gauge, ID: "HeapObjects", Value: float64Ptr(float64(m.HeapObjects))},
		{MType: gauge, ID: "HeapReleased", Value: float64Ptr(float64(m.HeapReleased))},
		{MType: gauge, ID: "HeapSys", Value: float64Ptr(float64(m.HeapSys))},
		{MType: gauge, ID: "LastGC", Value: float64Ptr(float64(m.LastGC))},
		{MType: gauge, ID: "Lookups", Value: float64Ptr(float64(m.Lookups))},
		{MType: gauge, ID: "MCacheInuse", Value: float64Ptr(float64(m.MCacheInuse))},
		{MType: gauge, ID: "MCacheSys", Value: float64Ptr(float64(m.MCacheSys))},
		{MType: gauge, ID: "MSpanInuse", Value: float64Ptr(float64(m.MSpanInuse))},
		{MType: gauge, ID: "MSpanSys", Value: float64Ptr(float64(m.MSpanSys))},
		{MType: gauge, ID: "Mallocs", Value: float64Ptr(float64(m.Mallocs))},
		{MType: gauge, ID: "NextGC", Value: float64Ptr(float64(m.NextGC))},
		{MType: gauge, ID: "NumForcedGC", Value: float64Ptr(float64(m.NumForcedGC))},
		{MType: gauge, ID: "NumGC", Value: float64Ptr(float64(m.NumGC))},
		{MType: gauge, ID: "OtherSys", Value: float64Ptr(float64(m.GCSys))},
		{MType: gauge, ID: "PauseTotalNs", Value: float64Ptr(float64(m.PauseTotalNs))},
		{MType: gauge, ID: "StackInuse", Value: float64Ptr(float64(m.StackInuse))},
		{MType: gauge, ID: "StackSys", Value: float64Ptr(float64(m.StackSys))},
		{MType: gauge, ID: "Sys", Value: float64Ptr(float64(m.Sys))},
		{MType: gauge, ID: "TotalAlloc", Value: float64Ptr(float64(m.TotalAlloc))},
		{MType: gauge, ID: "RandomValue", Value: float64Ptr(rand.Float64())},
	}
}
