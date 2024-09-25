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
			s := initServer()
			assert.NotEmpty(t, s)
		})
	}
}
