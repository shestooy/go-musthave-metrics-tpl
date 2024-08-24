package model

import "strconv"

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

func (m *Metrics) GetValue() string {
	switch m.MType {
	case "gauge":
		if m.Value == nil {
			return ""
		}
		return strconv.FormatFloat(*m.Value, 'f', -1, 64)
	case "counter":
		if m.Delta == nil {
			return ""
		}
		return strconv.FormatInt(*m.Delta, 10)
	}
	return ""
}

func (m *Metrics) SetValue(v interface{}) {
	switch m.MType {
	case "gauge":
		val, ok := v.(float64)
		if !ok {
			return
		}
		if m.Value == nil {
			m.Value = new(float64)
		}
		*m.Value = val
	case "counter":
		val, ok := v.(int64)
		if !ok {
			return
		}
		if m.Delta == nil {
			m.Delta = new(int64)
		}
		*m.Delta = val
	}
}
