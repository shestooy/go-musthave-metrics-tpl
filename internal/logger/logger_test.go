package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitialize(t *testing.T) {
	tests := []struct {
		name  string
		level string
	}{
		{name: "TestLogInitialize", level: "info"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Initialize(tt.level)
			assert.NoError(t, err)
		})
	}
}
