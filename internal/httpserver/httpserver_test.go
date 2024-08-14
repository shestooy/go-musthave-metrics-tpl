package httpserver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRouter(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "TestGetRouter"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := GetRouter()
			assert.NotEmpty(t, r)
		})
	}
}
