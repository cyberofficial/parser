package parser

import (
	"testing"
)

// Test_number_formats tests various number formats in queries
func Test_number_formats(t *testing.T) {
	type NumberItem struct {
		IntValue    int
		FloatValue  float64
		UintValue   uint
		StringValue string
	}
	
	items := []NumberItem{
		{IntValue: 123, FloatValue: 123.45, UintValue: 123, StringValue: "123"},
		{IntValue: 456, FloatValue: 456.78, UintValue: 456, StringValue: "456"},
		{IntValue: 0, FloatValue: 0.0, UintValue: 0, StringValue: "0"},
		{IntValue: -123, FloatValue: -123.45, UintValue: 0, StringValue: "-123"},
	}
	
	tests := []struct {
		name     string
		query    string
		expected int
	}{
		// Integer comparisons
		{name: "Integer equality", query: "IntValue = 123", expected: 1},
		{name: "Integer inequality", query: "IntValue != 123", expected: 3},
		{name: "Integer greater than", query: "IntValue > 123", expected: 1},
		{name: "Integer less than", query: "IntValue < 123", expected: 2},
		{name: "Integer greater than or equal", query: "IntValue >= 123", expected: 2},
		{name: "Integer zero", query: "IntValue = 0", expected: 1},
		
		// Float comparisons
		{name: "Float equality", query: "FloatValue = 123.45", expected: 1},
		{name: "Float inequality", query: "FloatValue != 123.45", expected: 3},
		{name: "Float greater than", query: "FloatValue > 123.45", expected: 1},
		{name: "Float less than", query: "FloatValue < 123.45", expected: 2},
		{name: "Float zero", query: "FloatValue = 0.0", expected: 1},
		
		// String number comparisons
		{name: "String number equality", query: "StringValue = '123'", expected: 1},
		{name: "String number as string", query: "StringValue = '0'", expected: 1},
		
		// Mixed comparisons
		{name: "Int to float comparison", query: "IntValue = 123.0", expected: 1},
		{name: "Float to int comparison", query: "FloatValue = 123", expected: 0}, // This should fail as 123.45 != 123
		
		// Complex expressions with numbers
		{name: "Number range", query: "IntValue > 0 AND IntValue < 200", expected: 1},
		{name: "Multiple number conditions", query: "IntValue = 123 AND FloatValue = 123.45 AND StringValue = '123'", expected: 1},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Parse(tt.query, items)
			if err != nil {
				t.Errorf("Failed to parse %q: %v", tt.query, err)
				return
			}
			if len(result) != tt.expected {
				t.Errorf("Expected %d results, got %d for query %q", tt.expected, len(result), tt.query)
			}
		})
	}
}

// Test_negative_numbers specifically tests handling of negative numbers in queries
func Test_negative_numbers(t *testing.T) {
	type NumberItem struct {
		IntValue    int
		FloatValue  float64
		StringValue string
	}
	
	items := []NumberItem{
		{IntValue: 100, FloatValue: 100.0, StringValue: "100"},
		{IntValue: 0, FloatValue: 0.0, StringValue: "0"},
		{IntValue: -100, FloatValue: -100.0, StringValue: "-100"},
		{IntValue: 50, FloatValue: -50.5, StringValue: "-50"},
	}
	
	tests := []struct {
		name     string
		query    string
		expected int
	}{
		// Negative integer comparisons
		{name: "Negative integer equality", query: "IntValue = -100", expected: 1},
		{name: "Negative integer inequality", query: "IntValue != -100", expected: 3},
		{name: "Negative integer greater than", query: "IntValue > -100", expected: 3},
		{name: "Negative integer less than", query: "IntValue < -100", expected: 0},
		
		// Negative float comparisons
		{name: "Negative float equality", query: "FloatValue = -100.0", expected: 1},
		{name: "Negative float inequality", query: "FloatValue != -50.5", expected: 3},
		{name: "Negative float greater than", query: "FloatValue > -100.0", expected: 3},
		{name: "Negative float less than", query: "FloatValue < -50.5", expected: 1},
		
		// Mixed comparisons with negative numbers
		{name: "Compare negative and positive", query: "IntValue < 0", expected: 1},
		{name: "Range with negative lower bound", query: "IntValue > -50 AND IntValue < 75", expected: 2},
		
		// String representations of negative numbers
		{name: "String negative number", query: "StringValue = '-100'", expected: 1},
		{name: "String with negative number comparison", query: "StringValue = '-50'", expected: 1},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Parse(tt.query, items)
			if err != nil {
				t.Errorf("Failed to parse %q: %v", tt.query, err)
				return
			}
			if len(result) != tt.expected {
				t.Errorf("Expected %d results, got %d for query %q", tt.expected, len(result), tt.query)
			}
		})
	}
}
