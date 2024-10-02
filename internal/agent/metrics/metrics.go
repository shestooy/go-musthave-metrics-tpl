package metrics

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"log"
	"math/rand"
	"runtime"

	"github.com/shirou/gopsutil/v4/mem"
	"go.uber.org/zap"
)

const (
	Gauge = "gauge"
)

type Metric struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

func GetRuntimeMetrics() []Metric {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	float64Ptr := func(val float64) *float64 {
		return &val
	}
	return []Metric{
		{MType: Gauge, ID: "Alloc", Value: float64Ptr(float64(m.Alloc))},
		{MType: Gauge, ID: "BuckHashSys", Value: float64Ptr(float64(m.BuckHashSys))},
		{MType: Gauge, ID: "Frees", Value: float64Ptr(float64(m.Frees))},
		{MType: Gauge, ID: "GCCPUFraction", Value: float64Ptr(m.GCCPUFraction)},
		{MType: Gauge, ID: "GCSys", Value: float64Ptr(float64(m.GCSys))},
		{MType: Gauge, ID: "HeapAlloc", Value: float64Ptr(float64(m.HeapAlloc))},
		{MType: Gauge, ID: "HeapIdle", Value: float64Ptr(float64(m.HeapIdle))},
		{MType: Gauge, ID: "HeapInuse", Value: float64Ptr(float64(m.HeapInuse))},
		{MType: Gauge, ID: "HeapObjects", Value: float64Ptr(float64(m.HeapObjects))},
		{MType: Gauge, ID: "HeapReleased", Value: float64Ptr(float64(m.HeapReleased))},
		{MType: Gauge, ID: "HeapSys", Value: float64Ptr(float64(m.HeapSys))},
		{MType: Gauge, ID: "LastGC", Value: float64Ptr(float64(m.LastGC))},
		{MType: Gauge, ID: "Lookups", Value: float64Ptr(float64(m.Lookups))},
		{MType: Gauge, ID: "MCacheInuse", Value: float64Ptr(float64(m.MCacheInuse))},
		{MType: Gauge, ID: "MCacheSys", Value: float64Ptr(float64(m.MCacheSys))},
		{MType: Gauge, ID: "MSpanInuse", Value: float64Ptr(float64(m.MSpanInuse))},
		{MType: Gauge, ID: "MSpanSys", Value: float64Ptr(float64(m.MSpanSys))},
		{MType: Gauge, ID: "Mallocs", Value: float64Ptr(float64(m.Mallocs))},
		{MType: Gauge, ID: "NextGC", Value: float64Ptr(float64(m.NextGC))},
		{MType: Gauge, ID: "NumForcedGC", Value: float64Ptr(float64(m.NumForcedGC))},
		{MType: Gauge, ID: "NumGC", Value: float64Ptr(float64(m.NumGC))},
		{MType: Gauge, ID: "OtherSys", Value: float64Ptr(float64(m.GCSys))},
		{MType: Gauge, ID: "PauseTotalNs", Value: float64Ptr(float64(m.PauseTotalNs))},
		{MType: Gauge, ID: "StackInuse", Value: float64Ptr(float64(m.StackInuse))},
		{MType: Gauge, ID: "StackSys", Value: float64Ptr(float64(m.StackSys))},
		{MType: Gauge, ID: "Sys", Value: float64Ptr(float64(m.Sys))},
		{MType: Gauge, ID: "TotalAlloc", Value: float64Ptr(float64(m.TotalAlloc))},
		{MType: Gauge, ID: "RandomValue", Value: float64Ptr(rand.Float64())},
	}
}

func GetMemoryMetrics(log *zap.SugaredLogger) []Metric {
	m, err := mem.VirtualMemory()
	if err != nil {
		log.Error("Error getting memory metrics: ", err)
		return nil
	}
	float64Ptr := func(val float64) *float64 {
		return &val
	}
	return []Metric{
		{MType: Gauge, ID: "TotalMemory", Value: float64Ptr(float64(m.Total))},
		{MType: Gauge, ID: "FreeMemory", Value: float64Ptr(float64(m.Free))},
		{MType: Gauge, ID: "UsedMemory", Value: float64Ptr(m.UsedPercent)},
	}
}

func Compress(m []Metric) ([]byte, error) {
	var buf bytes.Buffer

	w := gzip.NewWriter(&buf)

	err := json.NewEncoder(w).Encode(m)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	if err = w.Close(); err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return buf.Bytes(), nil
}

func GetMetricsAsBody(m []Metric) ([]byte, error) {
	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)
	err := json.NewEncoder(w).Encode(m)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return buf.Bytes(), nil
}
