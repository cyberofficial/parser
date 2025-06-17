package parser

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

// Define test structures with unusual field types and edge cases
type EdgeCaseItem struct {
	ID           int
	Name         string
	EmptyArray   []string                          // Always empty
	NilArray     []int                             // Always nil
	MixedArray   []interface{}                     // Contains mixed types
	NestedArrays [][]string                        // Array of arrays
	MapOfArrays  map[string][]string               // Map containing arrays
	ArrayOfMaps  []map[string]string               // Array containing maps
	NestedMaps   map[string]map[string]interface{} // Nested maps
	ZeroValues   ZeroValueStruct                   // Struct with zero values
	Channels     []chan int                        // Unsupported type
	Functions    []func()                          // Unsupported type
	Times        []time.Time                       // Time type
	TimePointers []*time.Time                      // Pointers to time
}

type ZeroValueStruct struct {
	Int    int
	String string
	Bool   bool
	Float  float64
	Slice  []string
}

func generateEdgeCaseTestData() []EdgeCaseItem {
	items := make([]EdgeCaseItem, 10)

	for i := 0; i < 10; i++ {
		now := time.Now().Add(time.Duration(i) * time.Hour)
		nowPtr := now

		items[i] = EdgeCaseItem{
			ID:         i,
			Name:       fmt.Sprintf("EdgeCase-%d", i),
			EmptyArray: []string{},
			MixedArray: []interface{}{
				"string",
				123,
				true,
				3.14,
				[]string{"nested", "array"},
				map[string]int{"nested": 42},
			},
			NestedArrays: [][]string{
				{"a", "b", "c"},
				{"d", "e", "f"},
				{"special", "test", "value"},
			},
			MapOfArrays: map[string][]string{
				"colors":    {"red", "blue", "green"},
				"sizes":     {"small", "medium", "large"},
				"empty":     {},
				"withSpace": {"has space", "another space"},
			},
			ArrayOfMaps: []map[string]string{
				{"key1": "value1", "key2": "value2"},
				{"key3": "value3", "keySpecial": "specialValue"},
				{"empty": ""},
			},
			NestedMaps: map[string]map[string]interface{}{
				"user": {
					"name":   "User " + fmt.Sprint(i),
					"active": i%2 == 0,
					"score":  float64(i) * 10.5,
					"tags":   []string{"user", "profile"},
				},
				"settings": {
					"theme":  "dark",
					"notify": true,
					"limits": map[string]int{"daily": i * 10, "weekly": i * 50},
				},
			},
			Times:        []time.Time{now, now.AddDate(0, 0, -1), now.AddDate(0, 0, -2)},
			TimePointers: []*time.Time{&nowPtr},
		}

		// Add special test cases for different indices
		if i == 3 {
			items[i].MixedArray = append(items[i].MixedArray, "needle")
		}

		if i == 5 {
			items[i].MixedArray = nil
		}

		if i == 7 {
			items[i].NilArray = nil
			items[i].EmptyArray = nil
		} else {
			items[i].NilArray = []int{}
		}
	}

	return items
}

func Test_ANY_EdgeCases_Advanced(t *testing.T) {
	items := generateEdgeCaseTestData()

	tests := []struct {
		name        string
		query       string
		expected    int
		expectError bool
	}{
		{
			name:     "ANY_with_empty_array_field",
			query:    "ANY(EmptyArray) = 'anything'",
			expected: 0,
		},
		{
			name:     "ANY_with_nil_array_field",
			query:    "ANY(NilArray) = '5'",
			expected: 0,
		},
		{
			name:     "ANY_equals_with_mixed_types_array",
			query:    "ANY(MixedArray) = '123'",
			expected: 9, // Should match string representation of numbers
		},
		{
			name:     "ANY_with_nested_arrays",
			query:    "ANY(NestedArrays) CONTAINS 'special'",
			expected: 10, // Should match all items with nested array containing "special"
		},
		{
			name:     "ANY_with_map_of_arrays",
			query:    "ANY(MapOfArrays.colors) = 'red'",
			expected: 10,
		},
		{
			name:     "ANY_with_spaces_in_values",
			query:    "ANY(MapOfArrays.withSpace) = 'has space'",
			expected: 10,
		},
		{
			name:     "ANY_with_array_of_maps",
			query:    "ANY(ArrayOfMaps) CONTAINS 'specialValue'",
			expected: 10,
		},
		{
			name:     "ANY_with_multiple_nested_maps",
			query:    "ANY(NestedMaps.user.tags) = 'profile'",
			expected: 10,
		},
		{
			name:     "ANY_with_numeric_in_nested_maps",
			query:    "ANY(NestedMaps.settings.limits.daily) > 30",
			expected: 7, // Items with ID 4-9
		},
		{
			name:     "ANY_with_boolean_in_nested_maps",
			query:    "ANY(NestedMaps.user.active) = 'true'",
			expected: 5, // Half the items
		},
		{
			name:     "ANY_with_multiple_comparison",
			query:    "ANY(MixedArray) = ANY('string', '123', 'needle')",
			expected: 9, // All non-nil MixedArray items
		},
		{
			name:     "ANY_with_multiple_fields",
			query:    "ANY(MapOfArrays.colors) = ANY('red', 'blue') OR ANY(MapOfArrays.sizes) = 'small'",
			expected: 10,
		},
		{
			name:     "Complex_nested_ANY_query",
			query:    "ANY(NestedMaps.user.tags) = 'profile' AND ANY(NestedMaps.settings.limits.weekly) > 200",
			expected: 5, // Items with ID 5-9
		},
		{
			name:     "ANY_with_nested_special_characters",
			query:    "ANY(ArrayOfMaps) CONTAINS 'key'",
			expected: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := Parse(tt.query, items)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, but got none for query: %s", tt.query)
				}
				return
			}

			if err != nil {
				t.Fatalf("Failed to parse query '%s': %v", tt.query, err)
			}

			if len(results) != tt.expected {
				t.Errorf("Expected %d items, got %d for query %q",
					tt.expected, len(results), tt.query)

				// Print the IDs of matching items to help debug
				var ids []int
				for _, item := range results {
					ids = append(ids, item.ID)
				}
				t.Logf("Matching item IDs: %v", ids)
			}
		})
	}
}

func Test_ANY_ErrorHandling_Extended(t *testing.T) {
	items := generateEdgeCaseTestData()

	tests := []struct {
		name        string
		query       string
		expectError bool
		errorSubstr string // Expected substring in error
	}{
		{
			name:        "ANY_with_malformed_syntax",
			query:       "ANY(Tags = 'test'",
			expectError: true,
			errorSubstr: "parse",
		},
		{
			name:        "ANY_missing_closed_parenthesis",
			query:       "ANY(Tags",
			expectError: true,
			errorSubstr: "parse",
		},
		{
			name:        "ANY_with_non_existent_field",
			query:       "ANY(NonExistentField) = 'test'",
			expectError: false, // Should not error, just return no matches
		},
		{
			name:        "ANY_with_unsupported_type",
			query:       "ANY(Channels) = 'anything'",
			expectError: false, // Should not error, just return no matches
		},
		{
			name:        "ANY_invalid_comparison_operator",
			query:       "ANY(Name) INVALID 'test'",
			expectError: true,
			errorSubstr: "token",
		},
		{
			name:        "ANY_unclosed_quotes",
			query:       "ANY(Tags) = 'unclosed",
			expectError: true,
			errorSubstr: "parse",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := Parse(tt.query, items)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none for query: %s", tt.query)
				} else if tt.errorSubstr != "" && !strings.Contains(err.Error(), tt.errorSubstr) {
					t.Errorf("Error message %q doesn't contain expected substring %q",
						err.Error(), tt.errorSubstr)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				// For valid queries, just log the results count
				t.Logf("Query '%s' returned %d results", tt.query, len(results))
			}
		})
	}
}

func Test_ANY_TimeHandling(t *testing.T) {
	items := generateEdgeCaseTestData()

	// Add a specific item with known time values
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	lastWeek := now.AddDate(0, 0, -7)

	specialItem := EdgeCaseItem{
		ID:   999,
		Name: "TimeTestItem",
		Times: []time.Time{
			now,
			yesterday,
			lastWeek,
		},
	}
	items = append(items, specialItem)

	// Format times as strings in the same format the parser would use
	nowStr := now.Format(time.RFC3339)
	yesterdayStr := yesterday.Format(time.RFC3339)

	tests := []struct {
		name     string
		query    string
		expected int
	}{
		{
			name:     "ANY_exact_time_match",
			query:    fmt.Sprintf("ANY(Times) = '%s'", nowStr),
			expected: 1, // Only the special item should match exactly
		},
		{
			name:     "ANY_with_multiple_time_values",
			query:    fmt.Sprintf("ANY(Times) = ANY('%s', '%s')", nowStr, yesterdayStr),
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := Parse(tt.query, items)

			if err != nil {
				t.Fatalf("Failed to parse query '%s': %v", tt.query, err)
			}

			if len(results) != tt.expected {
				t.Errorf("Expected %d items, got %d for query %q",
					tt.expected, len(results), tt.query)
			}
		})
	}
}

// Test to validate that unrelated test suites and functionality still work
func Test_ANYBackwardCompatibility(t *testing.T) {
	// Use a simple standard test structure to ensure basic functionality
	type TestItem struct {
		ID    int
		Name  string
		Value int
		Tags  []string
	}

	// Create simple test data
	items := []TestItem{
		{ID: 1, Name: "Item One", Value: 100, Tags: []string{"important"}},
		{ID: 2, Name: "Item Two", Value: 200, Tags: []string{"normal"}},
		{ID: 3, Name: "Item Three", Value: 300, Tags: []string{"low"}},
	}

	// Test basic non-ANY queries still work
	basicTests := []struct {
		name     string
		query    string
		expected int
	}{
		{
			name:     "Basic_equal_comparison_still_works",
			query:    "Name = 'Item One'",
			expected: 1,
		},
		{
			name:     "Basic_greater_than_still_works",
			query:    "Value > 150",
			expected: 2,
		},
		{
			name:     "Basic_AND_still_works",
			query:    "Value >= 100 AND Value <= 200",
			expected: 2,
		},
	}

	// Test basic ANY queries
	anyTests := []struct {
		name     string
		query    string
		expected int
	}{
		{
			name:     "Simple_ANY_tag_match",
			query:    "ANY(Tags) = 'important'",
			expected: 1,
		},
		{
			name:     "ANY_with_numerical_comparison",
			query:    "ANY(Tags) = 'normal' AND Value = 200",
			expected: 1,
		},
	}

	// Run basic tests
	for _, tt := range basicTests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := Parse(tt.query, items)

			if err != nil {
				t.Fatalf("Failed to parse query '%s': %v", tt.query, err)
			}

			if len(results) != tt.expected {
				t.Errorf("Expected %d items, got %d for query %q",
					tt.expected, len(results), tt.query)
			}
		})
	}

	// Run ANY tests
	for _, tt := range anyTests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := Parse(tt.query, items)

			if err != nil {
				t.Fatalf("Failed to parse query '%s': %v", tt.query, err)
			}

			if len(results) != tt.expected {
				t.Errorf("Expected %d items, got %d for query %q",
					tt.expected, len(results), tt.query)
			}
		})
	}
}
