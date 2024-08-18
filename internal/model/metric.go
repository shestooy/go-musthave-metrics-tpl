package model

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
	g.Values = make(map[string]float64)
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
	c.Values = make(map[string]int64)
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
