package parser

import (
	"testing"
)

// Test_unicode_support tests Unicode characters in string literals
func Test_unicode_support(t *testing.T) {
	type Item struct {
		Name string
	}

	items := []Item{
		{Name: "こんにちは"}, // Hello in Japanese
		{Name: "你好"},    // Hello in Chinese
		{Name: "안녕하세요"}, // Hello in Korean
		{Name: "Hello"}, // Hello in English
	}

	tests := []struct {
		name     string
		query    string
		expected int
	}{
		{
			name:     "Japanese string",
			query:    "Name = 'こんにちは'",
			expected: 1,
		},
		{
			name:     "Chinese string",
			query:    "Name = '你好'",
			expected: 1,
		},
		{
			name:     "Korean string",
			query:    "Name = '안녕하세요'",
			expected: 1,
		},
		{
			name:     "English string",
			query:    "Name = 'Hello'",
			expected: 1,
		},
		{
			name:     "Non-matching unicode",
			query:    "Name = 'Привет'", // Russian hello
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Parse(tt.query, items)
			if err != nil {
				t.Errorf("Failed to parse %q: %v", tt.query, err)
			}
			if len(result) != tt.expected {
				t.Errorf("Expected %d results, got %d for query %q", tt.expected, len(result), tt.query)
			}
		})
	}
}

// Test_deeply_nested_fields tests traversal through deeply nested struct fields
func Test_deeply_nested_fields(t *testing.T) {
	type Level4 struct {
		Value string
	}

	type Level3 struct {
		Level4 Level4
		Value  string
	}

	type Level2 struct {
		Level3 Level3
		Value  string
	}

	type Level1 struct {
		Level2 Level2
		Value  string
	}

	items := []Level1{
		{
			Value: "level1",
			Level2: Level2{
				Value: "level2",
				Level3: Level3{
					Value: "level3",
					Level4: Level4{
						Value: "level4",
					},
				},
			},
		},
	}

	tests := []struct {
		name     string
		query    string
		expected int
	}{
		{
			name:     "Level 1 field",
			query:    "Value = 'level1'",
			expected: 1,
		},
		{
			name:     "Level 2 field",
			query:    "Level2.Value = 'level2'",
			expected: 1,
		},
		{
			name:     "Level 3 field",
			query:    "Level2.Level3.Value = 'level3'",
			expected: 1,
		},
		{
			name:     "Level 4 field",
			query:    "Level2.Level3.Level4.Value = 'level4'",
			expected: 1,
		},
		{
			name:     "Incorrect value at deep nesting",
			query:    "Level2.Level3.Level4.Value = 'wrong'",
			expected: 0,
		},
		{
			name:     "Non-existent deep field",
			query:    "Level2.Level3.Level4.NonExistent = 'level4'",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Parse(tt.query, items)
			if err != nil {
				t.Errorf("Failed to parse %q: %v", tt.query, err)
			}
			if len(result) != tt.expected {
				t.Errorf("Expected %d results, got %d for query %q", tt.expected, len(result), tt.query)
			}
		})
	}
}

// Test_is_null_operator tests the IS NULL and IS NOT NULL operators
func Test_is_null_operator(t *testing.T) {
	type ItemWithNulls struct {
		ID          int
		Name        string
		Description *string
		Tags        []string
		Metadata    map[string]string
	}

	desc1 := "Item 1"
	desc2 := "Item 2"

	items := []ItemWithNulls{
		{ID: 1, Name: "Item 1", Description: &desc1, Tags: []string{"tag1"}, Metadata: map[string]string{"key": "value"}},
		{ID: 2, Name: "Item 2", Description: &desc2, Tags: []string{}, Metadata: nil},
		{ID: 3, Name: "Item 3", Description: nil, Tags: nil, Metadata: map[string]string{}},
	}

	tests := []struct {
		name     string
		query    string
		expected int
	}{
		{
			name:     "IS NULL on nil pointer",
			query:    "Description IS NULL",
			expected: 1, // Only Item 3
		},
		{
			name:     "IS NOT NULL on pointer",
			query:    "Description IS NOT NULL",
			expected: 2, // Item 1 and Item 2
		},
		{
			name:     "IS NULL on empty slice",
			query:    "Tags IS NULL",
			expected: 2, // Item 2 (empty slice) and Item 3 (nil slice)
		},
		{
			name:     "IS NOT NULL on non-empty slice",
			query:    "Tags IS NOT NULL",
			expected: 1, // Only Item 1
		},
		{
			name:     "IS NULL on nil map",
			query:    "Metadata IS NULL",
			expected: 1, // Only Item 2
		},
		{
			name:     "IS NULL on empty map",
			query:    "Metadata IS NULL",
			expected: 1, // Item 2 (nil map) considered null
		},
		{
			name:     "IS NOT NULL on map",
			query:    "Metadata IS NOT NULL",
			expected: 2, // Item 1 and Item 3
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Parse(tt.query, items)
			if err != nil {
				t.Errorf("Failed to parse %q: %v", tt.query, err)
			}
			if len(result) != tt.expected {
				t.Errorf("Expected %d results, got %d for query %q", tt.expected, len(result), tt.query)
			}
		})
	}
}

// Test_quotes_in_strings tests handling of quotes in string literals
func Test_quotes_in_strings(t *testing.T) {
	type StringItem struct {
		Value string
	}

	items := []StringItem{
		{Value: "No special chars"},
		{Value: "Contains 'single quotes'"},
		{Value: "Contains \"double quotes\""},
		{Value: "Contains 'mixed \"quotes\"'"},
	}

	tests := []struct {
		name     string
		query    string
		expected int
	}{
		{
			name:     "Simple string",
			query:    "Value = 'No special chars'",
			expected: 1,
		},
		{
			name:     "String with double quotes",
			query:    "Value = 'Contains \"double quotes\"'",
			expected: 1,
		},
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
