package model

import "strconv"

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

func (m *Metrics) GetValueAsString() string {
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

func (m *Metrics) SetValue(v string) error {
	switch m.MType {
	case "gauge":
		val, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return err
		}
		if m.Value == nil {
			m.Value = new(float64)
		}
		*m.Value = val
	case "counter":
		val, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return err
		}
		if m.Delta == nil {
			m.Delta = new(int64)
		}
		*m.Delta = val
	}
	return nil
}
