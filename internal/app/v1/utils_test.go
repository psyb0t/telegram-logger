package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLogLevelEmoji(t *testing.T) {
	tests := []struct {
		level    string
		expected string
	}{
		{"debug", "🐛"},
		{"dbg", "🐛"},
		{"debugging", "🐛"},
		{"info", "💬"},
		{"inf", "💬"},
		{"information", "💬"},
		{"warn", "⚠️"},
		{"wrn", "⚠️"},
		{"warning", "⚠️"},
		{"error", "❌"},
		{"err", "❌"},
		{"fatal", "💣"},
		{"critical", "💣"},
		{"invalid", ""},
	}

	for _, test := range tests {
		actual := getLogLevelEmoji(test.level)
		assert.Equal(t, actual, test.expected)
	}
}
