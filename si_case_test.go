package parser

import (
	"testing"
)

func TestSIPrefixCaseHandling(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Uppercase K should work (SI prefix)",
			input:    "Count > 2.5K",
			expected: "Count > 2500",
		},
		{
			name:     "Uppercase M should work (SI prefix)",
			input:    "Population > 10M",
			expected: "Population > 10000000",
		},
		{
			name:     "Lowercase m should be time unit (minutes), not SI prefix",
			input:    "Duration > 10m",
			expected: "Duration > 600",
		},
		{
			name:     "Lowercase g should NOT be parsed as SI prefix (should remain as-is)",
			input:    "Count > 10g",
			expected: "Count > 10g", // Should remain unchanged since lowercase g is not in SI map
		},
		{
			name:     "Mixed case: GB should be bytes, not SI prefix",
			input:    "Size > 1GB",
			expected: "Size > 1000000000",
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
