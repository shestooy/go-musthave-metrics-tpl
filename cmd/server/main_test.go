package main

import "testing"

func Test_start(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"the test is temporary"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}
