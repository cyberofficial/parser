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
			name:     "Decimal byte sizes - unquoted",
			input:    "Drive.Size > 10GB",
			expected: "Drive.Size > 10000000000",
		},
		{
			name:     "Binary byte sizes - unquoted",
			input:    "Drive.Size > 10GiB",
			expected: "Drive.Size > 10737418240",
		},
		{
			name:     "Binary MiB units - unquoted",
			input:    "Memory > 512MiB",
			expected: "Memory > 536870912",
		},
		{
			name:     "Binary KiB units - unquoted",
			input:    "Cache < 100KiB",
			expected: "Cache < 102400",
		},
		{
			name:     "Decimal MB units - unquoted",
			input:    "Storage > 500MB",
			expected: "Storage > 500000000",
		},
		{
			name:     "Decimal KB units - unquoted",
			input:    "Buffer < 100KB",
			expected: "Buffer < 100000",
		},
		{
			name:     "Decimal TB units - unquoted",
			input:    "Backup > 2TB",
			expected: "Backup > 2000000000000",
		},
		{
			name:     "Binary TiB units - unquoted",
			input:    "Archive > 1TiB",
			expected: "Archive > 1099511627776",
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
			name:     "Find drives larger than 1TB using decimal units",
			query:    "Size > 1TB",
			expected: 1,
		},
		{
			name:     "Find drives larger than 1TiB using binary units",
			query:    "Size > 1TiB",
			expected: 1,
		},
		{
			name:     "Find drives larger than 100GB using decimal units",
			query:    "Size > 100GB",
			expected: 2,
		},
		{
			name:     "Find drives larger than 100GiB using binary units",
			query:    "Size > 100GiB",
			expected: 2,
		},
		{
			name:     "Find drives smaller than 600GB using decimal units",
			query:    "Size < 600GB",
			expected: 1, // Only the 500GB SSD
		},
		{
			name:     "Find drives smaller than 600GiB using binary units",
			query:    "Size < 600GiB",
			expected: 1, // Only the 500GB SSD (500GB < 600GiB)
		},
		{
			name:     "Find SSD drive by name",
			query:    "Name = 'SSD'",
			expected: 1,
		},
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

func TestDecimalVsBinaryUnits(t *testing.T) {
	// Test the difference between decimal and binary units
	testCases := []struct {
		query    string
		expected string
	}{
		{"Size > 1TB", "Size > 1000000000000"},   // Decimal: 1,000,000,000,000
		{"Size > 1TiB", "Size > 1099511627776"},  // Binary: 1,099,511,627,776
		{"Memory > 1GB", "Memory > 1000000000"},  // Decimal: 1,000,000,000
		{"Memory > 1GiB", "Memory > 1073741824"}, // Binary: 1,073,741,824
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Convert %s", tc.query), func(t *testing.T) {
			result := normalizeHumanizedValues(tc.query)
			if result != tc.expected {
				t.Errorf("normalizeHumanizedValues(%q) = %q, want %q", tc.query, result, tc.expected)
			}
		})
	}
}
