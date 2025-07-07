package parser

import (
	"testing"
)

func TestLowercaseConflicts(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			// This should be parsed as time (minutes), not SI prefix (million)
			name:     "10m should be 600 seconds (time), not 10,000,000 (SI million)",
			input:    "Duration > 10m",
			expected: "Duration > 600",
		},
		{
			// This should be parsed as SI prefix (million), not time
			name:     "10M should be 10,000,000 (SI million), not time",
			input:    "Count > 10M",
			expected: "Count > 10000000",
		},
		{
			// Test the conflict: what happens with "10g"?
			name:     "10g should be 10,000,000,000 (SI giga)",
			input:    "Count > 10g",
			expected: "Count > 10000000000",
		},
		{
			// Test the conflict: what happens with "10t"?
			name:     "10t should be 10,000,000,000,000 (SI tera)",
			input:    "Count > 10t",
			expected: "Count > 10000000000000",
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
