package types

import (
	"errors"
	"fmt"
	"strconv"
)

type Metrics interface {
	AddValue(k, v string) error
	Init()
	GetValueID(n string) (interface{}, error)
	GetAllValue() interface{}
}

type Gauge struct {
	Values map[string]float64
}

type Counter struct {
	Values map[string]int64
}

func (g *Gauge) AddValue(k, v string) error {
	value, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return fmt.Errorf("expected float64, got %T", v)
	}
	g.Values[k] = value
	return nil
}

func (g *Gauge) Init() {
	g.Values = map[string]float64{
		"Alloc":         0.0,
		"BuckHashSys":   0.0,
		"Frees":         0.0,
		"GCCPUFraction": 0.0,
		"GCSys":         0.0,
		"HeapAlloc":     0.0,
		"HeapIdle":      0.0,
		"HeapInuse":     0.0,
		"HeapObjects":   0.0,
		"HeapReleased":  0.0,
		"HeapSys":       0.0,
		"LastGC":        0.0,
		"Lookups":       0.0,
		"MCacheInuse":   0.0,
		"MCacheSys":     0.0,
		"MSpanInuse":    0.0,
		"MSpanSys":      0.0,
		"Mallocs":       0.0,
		"NextGC":        0.0,
		"NumForcedGC":   0.0,
		"NumGC":         0.0,
		"OtherSys":      0.0,
		"PauseTotalNs":  0.0,
		"StackInuse":    0.0,
		"StackSys":      0.0,
		"Sys":           0.0,
		"TotalAlloc":    0.0,
		"RandomValue":   0.0,
	}
}

func (g *Gauge) GetValueID(n string) (interface{}, error) {
	if value, ok := g.Values[n]; !ok {
		return value, errors.New("unknown metric name")
	}
	return g.Values[n], nil
}

func (g *Gauge) GetAllValue() interface{} {
	return g.Values
}

func (c *Counter) AddValue(k, v string) error {
	value, err := strconv.Atoi(v)
	if err != nil {
		return fmt.Errorf("expected int64, got %T", v)
	}
	c.Values[k] += int64(value)
	return nil
}

func (c *Counter) Init() {
	c.Values = map[string]int64{
		"PollCount": 0,
	}
}

func (c *Counter) GetValueID(n string) (interface{}, error) {
	if value, ok := c.Values[n]; !ok {
		return value, errors.New("unknown metric name")
	}
	return c.Values[n], nil
}

func (c *Counter) GetAllValue() interface{} {
	return c.Values
}
