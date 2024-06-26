package model

import (
	"fmt"
	"strconv"
)

type Metrics interface {
	AddValue(k, v string) error
	Init()
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
	g.Values = make(map[string]float64)
}

func (r *Counter) AddValue(k, v string) error {
	value, err := strconv.Atoi(v)
	if err != nil {
		return fmt.Errorf("expected int64, got %T", v)
	}
	r.Values[k] = int64(value)
	return nil
}

func (c *Counter) Init() {
	c.Values = make(map[string]int64)
}
