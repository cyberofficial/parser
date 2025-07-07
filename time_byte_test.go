package parser

import (
	"testing"
)

func TestTimeVsByteDistinction(t *testing.T) {
	// Test the specific case mentioned in the user's request
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "10m should be parsed as 600 seconds (time), not megabytes",
			input:    "Duration > 10m",
			expected: "Duration > 600",
		},
		{
			name:     "10MB should be parsed as megabytes (bytes), not minutes",
			input:    "Size > 10MB",
			expected: "Size > 10000000",
		},
		{
			name:     "1h should be parsed as 3600 seconds (time)",
			input:    "Timeout > 1h",
			expected: "Timeout > 3600",
		},
		{
			name:     "1GB should be parsed as gigabytes (bytes)",
			input:    "Storage > 1GB",
			expected: "Storage > 1000000000",
		},
		{
			name:     "Mixed time and byte units in same query",
			input:    "Duration > 30m AND Size < 500MB",
			expected: "Duration > 1800 AND Size < 500000000",
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

func TestTimeDurationParsing(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
		hasError bool
	}{
		{
			name:     "10m should parse as 600 seconds",
			input:    "10m",
			expected: 600,
			hasError: false,
		},
		{
			name:     "1h should parse as 3600 seconds",
			input:    "1h",
			expected: 3600,
			hasError: false,
		},
		{
			name:     "30s should parse as 30 seconds",
			input:    "30s",
			expected: 30,
			hasError: false,
		},
		{
			name:     "1d should parse as 86400 seconds",
			input:    "1d",
			expected: 86400,
			hasError: false,
		},
		{
			name:     "1.5h should parse as 5400 seconds",
			input:    "1.5h",
			expected: 5400,
			hasError: false,
		},
		{
			name:     "Invalid time unit should error",
			input:    "10x",
			expected: 0,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseTimeDuration(tt.input)
			if tt.hasError {
				if err == nil {
					t.Errorf("parseTimeDuration(%q) expected error but got none", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("parseTimeDuration(%q) unexpected error: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("parseTimeDuration(%q) = %d, want %d", tt.input, result, tt.expected)
				}
			}
		})
	}
}

func TestByteSizeParsing(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
		hasError bool
	}{
		{
			name:     "10MB should parse as 10000000 bytes",
			input:    "10MB",
			expected: 10000000,
			hasError: false,
		},
		{
			name:     "1GB should parse as 1000000000 bytes",
			input:    "1GB",
			expected: 1000000000,
			hasError: false,
		},
		{
			name:     "1GiB should parse as 1073741824 bytes",
			input:    "1GiB",
			expected: 1073741824,
			hasError: false,
		},
		{
			name:     "1KB should parse as 1000 bytes",
			input:    "1KB",
			expected: 1000,
			hasError: false,
		},
		{
			name:     "1KiB should parse as 1024 bytes",
			input:    "1KiB",
			expected: 1024,
			hasError: false,
		},
		{
			name:     "Invalid byte unit should error",
			input:    "10XB",
			expected: 0,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseByteSize(tt.input)
			if tt.hasError {
				if err == nil {
					t.Errorf("parseByteSize(%q) expected error but got none", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("parseByteSize(%q) unexpected error: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("parseByteSize(%q) = %d, want %d", tt.input, result, tt.expected)
				}
			}
		})
	}
}
