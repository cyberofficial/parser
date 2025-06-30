package parser

import (
	"fmt"
	"testing"
)

func TestNormalizeHumanizedValues(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Byte sizes - unquoted",
			input:    "Drive.Size > 10GiB",
			expected: "Drive.Size > 10000000000",
		},
		{
			name:     "Byte sizes with decimal - unquoted",
			input:    "Drive.Size > 1.5GB",
			expected: "Drive.Size > 1500000000",
		},
		{
			name:     "SI prefix numbers - unquoted",
			input:    "Count > 2.5K",
			expected: "Count > 2500",
		},
		{
			name:     "Comma separated numbers - unquoted",
			input:    "Population > 1,000,000",
			expected: "Population > 1000000",
		},
		{
			name:     "Mixed query with strings and humanized numbers",
			input:    "Person.Name = 'alice' AND Drive.Size > 10GB",
			expected: "Person.Name = 'alice' AND Drive.Size > 10000000000",
		},
		{
			name:     "Quoted strings should remain unchanged",
			input:    "Description CONTAINS 'has 10GB of storage'",
			expected: "Description CONTAINS 'has 10GB of storage'",
		},
		{
			name:     "Complex query",
			input:    "Drive.Size > 1.5TB AND Count < 5K AND Name = 'test'",
			expected: "Drive.Size > 1500000000000 AND Count < 5000 AND Name = 'test'",
		},
		{
			name:     "No humanized values",
			input:    "Age > 25 AND Name = 'john'",
			expected: "Age > 25 AND Name = 'john'",
		},
		{
			name:     "Binary byte sizes - unquoted",
			input:    "Drive.Size > 10GiB",
			expected: "Drive.Size > 10737418240",
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

func TestHumanizedValuesIntegration(t *testing.T) {
	// Test with actual structs
	type Drive struct {
		Name string
		Size int64
	}

	drives := []Drive{
		{Name: "SSD", Size: 500000000000},  // 500GB
		{Name: "HDD", Size: 2000000000000}, // 2TB
	}

	tests := []struct {
		name     string
		query    string
		expected int
	}{
		{
			name:     "Find drives larger than 1TB using humanized value",
			query:    "Size > 1TB",
			expected: 1,
		},
		{
			name:     "Find drives larger than 100GB using humanized value",
			query:    "Size > 100GB",
			expected: 2,
		},
		{
			name:     "Find SSD drive by name",
			query:    "Name = 'SSD'",
			expected: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Parse(tt.query, drives)
			if err != nil {
				t.Fatalf("Parse(%q) returned error: %v", tt.query, err)
			}
			if len(result) != tt.expected {
				t.Errorf("Parse(%q) returned %d results, want %d", tt.query, len(result), tt.expected)
				fmt.Printf("Query: %s\n", tt.query)
				fmt.Printf("Normalized: %s\n", normalizeHumanizedValues(tt.query))
				fmt.Printf("Results: %+v\n", result)
			}
		})
	}
}
