package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCounter_Init(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "TestInitCounter"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Counter{}
			c.Init()
			assert.NotEmpty(t, tests, "init failed")
		})
	}
}

func TestCounter_AddValue(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{name: "TestFloat64", value: "434.32", wantErr: true},
		{name: "TestInt", value: "3543", wantErr: false},
		{name: "TestString", value: "eoka", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Counter{}
			c.Init()
			err := c.AddValue(tt.name, tt.value)
			assert.Equal(t, tt.wantErr, err != nil, "unexpected error")
		})
	}
}

func TestGauge_Init(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "TestInitCounter"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Gauge{}
			c.Init()
			assert.NotEmpty(t, tests, "init failed")
		})
	}
}

func TestGauge_AddValue(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{name: "TestFloat64", value: "992348.32", wantErr: false},
		{name: "TestInt", value: "23512", wantErr: false},
		{name: "TestString", value: "test", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Gauge{}
			c.Init()
			err := c.AddValue(tt.name, tt.value)
			assert.Equal(t, tt.wantErr, err != nil, "unexpected error")
		})
	}
}
