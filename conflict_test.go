package parser

import (
	"testing"
)

func TestPotentialConflicts(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "10m should be minutes (600), not SI prefix million",
			input:    "Duration > 10m",
			expected: "Duration > 600",
		},
		{
			name:     "10M should be SI prefix million (10000000), not minutes",
			input:    "Count > 10M",
			expected: "Count > 10000000",
		},
		{
			name:     "10g should be SI prefix giga (10000000000), not a time unit",
			input:    "Count > 10g",
			expected: "Count > 10000000000",
		},
		{
			name:     "10G should be SI prefix giga (10000000000), not bytes",
			input:    "Count > 10G",
			expected: "Count > 10000000000",
		},
		{
			name:     "10GB should be bytes (10000000000), not SI prefix",
			input:    "Size > 10GB",
			expected: "Size > 10000000000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeHumanizedValues(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeHumanizedValues(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
