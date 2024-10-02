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
			l, err := Initialize(tt.level)
			assert.NotEmpty(t, l)
			assert.NoError(t, err)
		})
	}
}
