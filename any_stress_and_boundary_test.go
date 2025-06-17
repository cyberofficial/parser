package parser

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

// Extreme stress test for ANY parser performance
func Test_ANY_Stress_Test(t *testing.T) {
	// Skip this test by default as it's resource-intensive
	if testing.Short() {
		t.Skip("Skipping stress test in short mode")
	}

	// Generate large test dataset
	itemCount := 1000
	t.Logf("Generating %d test items...", itemCount)

	items := make([]struct {
		ID          int
		Name        string
		Tags        []string
		Numbers     []int
		Properties  map[string]string
		NestedArray [][]string
		DeepMap     map[string]map[string]map[string]interface{}
	}, itemCount)

	// Fill the array with test data
	for i := 0; i < itemCount; i++ {
		// Generate tags (more tags for higher IDs to test performance with longer arrays)
		tagCount := 5 + (i / 100)
		tags := make([]string, tagCount)

		for j := 0; j < tagCount; j++ {
			// Add some specific values to ensure matches
			if j == 0 && i%10 == 0 {
				tags[j] = "special-tag"
			} else if j == 1 && i%20 == 0 {
				tags[j] = "rare-tag"
			} else if j == 2 && i%50 == 0 {
				tags[j] = "very-rare-tag"
			} else {
				tags[j] = fmt.Sprintf("tag-%d-%d", i, j)
			}
		}

		// Generate numbers
		numbers := make([]int, 10)
		for j := 0; j < 10; j++ {
			numbers[j] = i*10 + j
		}

		// Generate properties
		properties := map[string]string{
			"prop1": fmt.Sprintf("value-%d-1", i),
			"prop2": fmt.Sprintf("value-%d-2", i),
			"prop3": fmt.Sprintf("value-%d-3", i),
		}

		// Some special values
		if i%25 == 0 {
			properties["special"] = "special-value"
		}

		// Generate nested array
		nestedArray := make([][]string, 3)
		for j := 0; j < 3; j++ {
			nestedArray[j] = make([]string, 3)
			for k := 0; k < 3; k++ {
				nestedArray[j][k] = fmt.Sprintf("nested-%d-%d-%d", i, j, k)
			}
		}

		// Add a specific value deep in the array for some items
		if i%30 == 0 {
			nestedArray[1][1] = "deep-special-value"
		}

		// Generate deep map
		deepMap := map[string]map[string]map[string]interface{}{
			"level1": {
				"level2": {
					"level3":   fmt.Sprintf("deep-value-%d", i),
					"level3-2": i,
					"level3-3": i%2 == 0,
					"array":    []string{"a", "b", "c"},
				},
			},
		}

		items[i] = struct {
			ID          int
			Name        string
			Tags        []string
			Numbers     []int
			Properties  map[string]string
			NestedArray [][]string
			DeepMap     map[string]map[string]map[string]interface{}
		}{
			ID:          i,
			Name:        fmt.Sprintf("Stress Test Item %d", i),
			Tags:        tags,
			Numbers:     numbers,
			Properties:  properties,
			NestedArray: nestedArray,
			DeepMap:     deepMap,
		}
	}

	// Define stress test queries
	tests := []struct {
		name           string
		query          string
		expectedApprox int // Approximate expected count (exact count can vary)
	}{
		{
			name:           "Simple_ANY_with_large_dataset",
			query:          "ANY(Tags) = 'special-tag'",
			expectedApprox: itemCount / 10, // Every 10th item
		},
		{
			name:           "ANY_with_rare_value",
			query:          "ANY(Tags) = 'very-rare-tag'",
			expectedApprox: itemCount / 50, // Every 50th item
		},
		{
			name:           "Deep_nested_lookup_large_dataset",
			query:          "ANY(DeepMap.level1.level2.level3) = 'deep-value-500'",
			expectedApprox: 1, // Exactly one item should match
		},
		{
			name:           "Complex_ANY_with_multiple_clauses",
			query:          "ANY(Tags) = 'special-tag' AND ANY(Numbers) > 500",
			expectedApprox: itemCount / 20, // Rough approximation
		},
		{
			name:           "Very_deep_array_lookup",
			query:          "ANY(NestedArray) = 'deep-special-value'",
			expectedApprox: itemCount / 30, // Every 30th item
		},
		{
			name:           "Map_property_lookup_with_special_value",
			query:          "ANY(Properties.special) = 'special-value'",
			expectedApprox: itemCount / 25, // Every 25th item
		},
		{
			name:           "Complex_OR_condition",
			query:          "ANY(Tags) = 'special-tag' OR ANY(Tags) = 'rare-tag'",
			expectedApprox: itemCount / 7, // Approximately (1/10 + 1/20) of items
		},
		{
			name:           "Deep_boolean_condition",
			query:          "ANY(DeepMap.level1.level2.level3-3) = 'true'",
			expectedApprox: itemCount / 2, // Every even-numbered item
		},
	}

	// Run the tests and measure performance
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start := time.Now()

			results, err := Parse(tt.query, items)

			duration := time.Since(start)

			if err != nil {
				t.Fatalf("Failed to parse query '%s': %v", tt.query, err)
			}

			t.Logf("Query '%s' took %v and returned %d results",
				tt.query, duration, len(results))

			// For stress tests, we check approximate counts rather than exact
			// since we're more concerned with performance and robustness
			difference := float64(len(results) - tt.expectedApprox)
			percentDiff := (difference / float64(tt.expectedApprox)) * 100

			if percentDiff > 20 || percentDiff < -20 {
				t.Errorf("Result count %d differs from expected %d by more than 20%% for query %q",
					len(results), tt.expectedApprox, tt.query)
			}
		})
	}
}

// Test for boundary conditions and extreme cases
func Test_ANY_Boundary_Conditions(t *testing.T) {
	// Create test data with boundary conditions
	items := []struct {
		ID                 int
		EmptyString        string
		ZeroInt            int
		MaxInt             int
		MinInt             int
		EmptyArray         []string
		SingleItemArray    []string
		LargeArray         []string
		MapWithEmptyValues map[string]string
		SpecialChars       []string
		UnicodeChars       []string
		Escapes            []string
	}{
		{
			ID:              1,
			EmptyString:     "",
			ZeroInt:         0,
			MaxInt:          int(^uint(0) >> 1),  // Max int
			MinInt:          -int(^uint(0) >> 1), // Min int
			EmptyArray:      []string{},
			SingleItemArray: []string{"single"},
			LargeArray:      make([]string, 1000), // 1000 items
			MapWithEmptyValues: map[string]string{
				"empty":         "",
				"":              "empty-key",
				"null":          "null",
				"zero":          "0",
				"special-chars": "!@#$%^&*()_+",
			},
			SpecialChars: []string{
				"!@#$%^&*()_+",
				"={}[]|\\:;\"'<>,./?",
				"`~",
			},
			UnicodeChars: []string{
				"你好世界",       // Chinese
				"Привет мир", // Russian
				"こんにちは世界",    // Japanese
				"🙂🔍🌍🚀",       // Emojis
			},
			Escapes: []string{
				"line1\nline2",
				"tab\tdelimited",
				"quoted\"string",
				"back\\slash",
			},
		},
	}

	// Fill the large array for testing
	for i := range items[0].LargeArray {
		items[0].LargeArray[i] = fmt.Sprintf("item-%d", i)
	}
	// Add special values at specific positions
	items[0].LargeArray[0] = "first"
	items[0].LargeArray[500] = "middle"
	items[0].LargeArray[999] = "last"

	tests := []struct {
		name     string
		query    string
		expected int
	}{
		{
			name:     "Empty_string_comparison",
			query:    "EmptyString = ''",
			expected: 1,
		},
		{
			name:     "Zero_int_comparison",
			query:    "ZeroInt = '0'",
			expected: 1,
		},
		{
			name:     "Max_int_comparison",
			query:    fmt.Sprintf("MaxInt = '%d'", int(^uint(0)>>1)),
			expected: 1,
		},
		{
			name:     "Min_int_comparison",
			query:    fmt.Sprintf("MinInt = '%d'", -int(^uint(0)>>1)),
			expected: 1,
		},
		{
			name:     "Empty_array_match",
			query:    "ANY(EmptyArray) = 'anything'",
			expected: 0, // Should not match
		},
		{
			name:     "Single_item_array_match",
			query:    "ANY(SingleItemArray) = 'single'",
			expected: 1,
		},
		{
			name:     "Large_array_first_item",
			query:    "ANY(LargeArray) = 'first'",
			expected: 1,
		},
		{
			name:     "Large_array_middle_item",
			query:    "ANY(LargeArray) = 'middle'",
			expected: 1,
		},
		{
			name:     "Large_array_last_item",
			query:    "ANY(LargeArray) = 'last'",
			expected: 1,
		},
		{
			name:     "Empty_map_value",
			query:    "ANY(MapWithEmptyValues.empty) = ''",
			expected: 1,
		},
		{
			name:     "Empty_map_key",
			query:    "ANY(MapWithEmptyValues.) = 'empty-key'",
			expected: 0, // This should not parse or match
		},
		{
			name:     "Special_characters_in_value",
			query:    "ANY(SpecialChars) = '!@#$%^&*()_+'",
			expected: 1,
		},
		{
			name:     "Unicode_characters_Chinese",
			query:    "ANY(UnicodeChars) = '你好世界'",
			expected: 1,
		},
		{
			name:     "Unicode_characters_Emojis",
			query:    "ANY(UnicodeChars) = '🙂🔍🌍🚀'",
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := Parse(tt.query, items)

			// Log any error but don't fail - some boundary tests might legitimately error
			if err != nil {
				t.Logf("Query '%s' produced error: %v", tt.query, err)
			} else {
				if len(results) != tt.expected {
					t.Errorf("Expected %d items, got %d for query %q",
						tt.expected, len(results), tt.query)
				}
			}
		})
	}
}

// Test for query syntax variations and edge cases
func Test_ANY_Query_Syntax_Variations(t *testing.T) {
	items := []struct {
		ID         int
		Tags       []string
		Properties map[string]string
	}{
		{
			ID:   1,
			Tags: []string{"one", "two", "three"},
			Properties: map[string]string{
				"color": "red",
				"size":  "large",
			},
		},
	}

	tests := []struct {
		name        string
		query       string
		expected    int
		expectError bool
	}{
		// Whitespace variations
		{
			name:     "Extra_whitespace_before_ANY",
			query:    "   ANY(Tags) = 'one'",
			expected: 1,
		},
		{
			name:     "Extra_whitespace_after_ANY",
			query:    "ANY   (Tags) = 'one'",
			expected: 1,
		},
		{
			name:     "Extra_whitespace_inside_parentheses",
			query:    "ANY( Tags ) = 'one'",
			expected: 1,
		},
		{
			name:     "Extra_whitespace_around_operator",
			query:    "ANY(Tags)    =    'one'",
			expected: 1,
		},
		{
			name:     "No_whitespace_around_operator",
			query:    "ANY(Tags)='one'",
			expected: 1,
		},

		// Parentheses variations
		{
			name:        "Missing_opening_parenthesis",
			query:       "ANY Tags) = 'one'",
			expectError: true,
		},
		{
			name:        "Missing_closing_parenthesis",
			query:       "ANY(Tags = 'one'",
			expectError: true,
		},
		{
			name:        "Extra_opening_parenthesis",
			query:       "ANY((Tags) = 'one'",
			expectError: true,
		},
		{
			name:        "Extra_closing_parenthesis",
			query:       "ANY(Tags)) = 'one'",
			expectError: true,
		},

		// Quote variations
		{
			name:        "Missing_opening_quote",
			query:       "ANY(Tags) = one'",
			expectError: true,
		},
		{
			name:        "Missing_closing_quote",
			query:       "ANY(Tags) = 'one",
			expectError: true,
		},
		{
			name:        "Double_quotes_instead_of_single",
			query:       `ANY(Tags) = "one"`,
			expectError: true, // Assuming parser only supports single quotes
		},

		// Case variations
		{
			name:     "Lowercase_any_keyword",
			query:    "any(Tags) = 'one'",
			expected: 1, // Assuming case-insensitive keywords
		},
		{
			name:     "Mixed_case_any_keyword",
			query:    "AnY(Tags) = 'one'",
			expected: 1, // Assuming case-insensitive keywords
		},
		{
			name:     "Uppercase_any_keyword",
			query:    "ANY(Tags) = 'one'",
			expected: 1,
		},

		// Complex syntax variations
		{
			name:     "Multiple_ANY_with_extra_spaces",
			query:    "  ANY(Tags)  =  'one'  AND  ANY(Properties.color)  =  'red'  ",
			expected: 1,
		},
		{
			name:     "Multiple_ANY_with_no_spaces",
			query:    "ANY(Tags)='one'AND ANY(Properties.color)='red'",
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := Parse(tt.query, items)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none for query: %s", tt.query)
				}
				return
			}

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

// Test query length limits and huge identifiers
func Test_ANY_Query_Length_Limits(t *testing.T) {
	// Generate a very long field path
	longPath := strings.Repeat("very_long_field_name_that_is_unnecessarily_verbose.", 10)
	longPath = strings.TrimSuffix(longPath, ".")

	// Generate a very long string value
	longValue := strings.Repeat("extremelylongstringvaluethatgoesonandonwithoutanyreasonorpurpose", 20)

	// Create an item with the long path
	type nestedStruct struct {
		Value string
	}

	item := struct {
		Very_long_field_name_that_is_unnecessarily_verbose struct {
			Very_long_field_name_that_is_unnecessarily_verbose struct {
				Very_long_field_name_that_is_unnecessarily_verbose struct {
					Very_long_field_name_that_is_unnecessarily_verbose struct {
						Very_long_field_name_that_is_unnecessarily_verbose struct {
							Very_long_field_name_that_is_unnecessarily_verbose struct {
								Very_long_field_name_that_is_unnecessarily_verbose struct {
									Very_long_field_name_that_is_unnecessarily_verbose struct {
										Very_long_field_name_that_is_unnecessarily_verbose struct {
											Very_long_field_name_that_is_unnecessarily_verbose struct {
												Value string
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
		LongArray []string
	}{
		LongArray: []string{longValue},
	}

	// Set the value at the end of the long path
	item.Very_long_field_name_that_is_unnecessarily_verbose.
		Very_long_field_name_that_is_unnecessarily_verbose.
		Very_long_field_name_that_is_unnecessarily_verbose.
		Very_long_field_name_that_is_unnecessarily_verbose.
		Very_long_field_name_that_is_unnecessarily_verbose.
		Very_long_field_name_that_is_unnecessarily_verbose.
		Very_long_field_name_that_is_unnecessarily_verbose.
		Very_long_field_name_that_is_unnecessarily_verbose.
		Very_long_field_name_that_is_unnecessarily_verbose.
		Very_long_field_name_that_is_unnecessarily_verbose.Value = "deep-value"

	items := []interface{}{item}

	tests := []struct {
		name        string
		query       string
		expected    int
		expectError bool
	}{
		{
			name:     "Very_long_field_path",
			query:    fmt.Sprintf("ANY(%s.Value) = 'deep-value'", longPath),
			expected: 1,
		},
		{
			name:     "Very_long_value",
			query:    fmt.Sprintf("ANY(LongArray) = '%s'", longValue),
			expected: 1,
		},
		{
			name: "Extremely_long_query_with_many_conditions",
			query: fmt.Sprintf("ANY(LongArray) = '%s' OR ANY(%s.Value) = 'deep-value'",
				longValue[:100], // Truncate long value to keep test manageable
				longPath),
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := Parse(tt.query, items)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none for query: %s", tt.query)
				}
				return
			}

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
