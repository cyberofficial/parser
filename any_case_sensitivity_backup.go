package parser

import (
	"testing"
)

// Test structure with fields having different cases
type CaseSensitiveItem struct {
	ID              int
	Name            string
	Tags            []string            // Normal case field 
	TagsLower       []string            // For testing lowercase
	TagsUpper       []string            // For testing uppercase
	MixedCaseTags   []string            // camelCase field
	MixedCaseTagsSnake []string         // snake_case field
	Properties      map[string]string   // normal case map
	PropertiesLower map[string]string   // lowercase map
	PropertiesUpper map[string]string   // uppercase map
	NestedCase      NestedCaseItem
	NestedCaseLower NestedCaseItem      // lowercase nested
}

type NestedCaseItem struct {
	Name      string
	Tags      []string
	SubTags   []string
	SUBTAGS   []string
	Properties map[string]string
	// Maps with case-sensitive keys
	CaseSensitiveMap map[string]interface{}
}

func Test_ANY_Case_Sensitivity_Comprehensive(t *testing.T) {
	// Create test data with various case variations
	items := []CaseSensitiveItem{
		{
			ID:            1,
			Name:          "Case Test 1",
			Tags:          []string{"normal", "mixed", "test"},
			TagsLower:     []string{"lowercase", "mixed", "test"},
			TagsUpper:     []string{"UPPERCASE", "MIXED", "TEST"},
			MixedCaseTags: []string{"CamelCase", "MixedCase", "testCase"},
			MixedCaseTagsSnake: []string{"snake_case", "under_score", "test_case"},
			Properties: map[string]string{
				"Color": "Red",
				"Size":  "Large",
				"Type":  "Main",
			},
			PropertiesLower: map[string]string{
				"color": "red",
				"size":  "large",
				"type":  "main",
			},
			PropertiesUpper: map[string]string{
				"COLOR": "RED",
				"SIZE":  "LARGE",
				"TYPE":  "MAIN",
			},
			NestedCase: NestedCaseItem{
				Name:    "Nested Normal Case",
				Tags:    []string{"nested", "normal", "case"},
				SubTags: []string{"sub", "tags", "normal"},
				SUBTAGS: []string{"SUB", "TAGS", "UPPER"},
				Properties: map[string]string{
					"Color": "Blue",
					"Size":  "Medium",
				},
				CaseSensitiveMap: map[string]interface{}{
					"Red":      "FF0000",
					"red":      "red value",
					"GREEN":    "00FF00",
					"green":    "green value",
					"Blue":     "0000FF",
					"BLUE":     "BLUE VALUE",
					"RedGreen": "Mixed",
					"redgreen": "mixed",
				},
			},
			NestedCaseLower: NestedCaseItem{
				Name:    "nested lowercase case",
				Tags:    []string{"nested", "lowercase", "field"},
				SubTags: []string{"sub", "tags", "lower"},
				Properties: map[string]string{
					"color": "yellow",
					"size":  "small",
				},
			},
		},
	}

	// Test cases for field name case sensitivity
	fieldTests := []struct {
		name     string
		query    string
		expected int
	}{
		{
			name:     "Case_sensitive_field_exact_match",
			query:    "ANY(TagsLower) = 'lowercase'",
			expected: 1, // Should match lowercase field
		},
		{
			name:     "Case_sensitive_field_uppercase_exact",
			query:    "ANY(TagsUpper) = 'UPPERCASE'",
			expected: 1, // Should match uppercase field
		},
		{
			name:     "Case_insensitive_field_lookup",
			query:    "ANY(Tags) = 'lowercase'",
			expected: 0, // Should not match because Tags doesn't contain this value
		},
		{
			name:     "Case_insensitive_field_lookup_TagsLower",
			query:    "ANY(TagsLower) = 'LOWERCASE'",
			expected: 0, // Should not match due to case sensitivity in values
		},
		{
			name:     "Mixed_case_field_lookup",
			query:    "ANY(MixedCaseTags) = 'CamelCase'", 
			expected: 1, // Exact case match
		},
		{
			name:     "Snake_case_field_lookup",
			query:    "ANY(MixedCaseTagsSnake) = 'snake_case'", 
			expected: 1,
		},
		{
			name:     "Case_sensitive_map_exact_field",
			query:    "ANY(Properties.Color) = 'Red'",
			expected: 1, // Exact match
		},
		{
			name:     "Case_insensitive_map_field_lookup",
			query:    "ANY(PropertiesLower.color) = 'red'",
			expected: 1, // Exact case match
		},
		{
			name:     "Case_insensitive_map_value_lookup",
			query:    "ANY(Properties.color) = 'RED'",
			expected: 0, // Should not match case-sensitive value
		},
		{
			name:     "Nested_case_insensitive_lookup",
			query:    "ANY(NestedCaseLower.Tags) = 'lowercase'",
			expected: 1, // Exact case match
		},
		{
			name:     "Deep_case_sensitive_map_lookup",
			query:    "ANY(NestedCase.CaseSensitiveMap.Red) = 'FF0000'",
			expected: 1, // Exact case map key
		},
		{
			name:     "Deep_case_sensitive_map_lookup_exact",
			query:    "ANY(NestedCase.CaseSensitiveMap.red) = 'red value'",
			expected: 1, // Different case map key
		},
		{
			name:     "Deep_case_insensitive_map_lookup",
			query:    "ANY(NestedCase.CaseSensitiveMap.RED) = 'FF0000'",
			expected: 0, // Case-sensitive map key lookup
		},
	}

	// Run the tests	for _, tt := range fieldTests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := Parse[CaseSensitiveItem](tt.query, items)
			
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

func Test_ANY_Case_Sensitivity_Values(t *testing.T) {
	// Test specific case for case-sensitive value matching
	items := []struct {
		ID    int
		Name  string
		Tags  []string
		Codes map[string]string
	}{
		{
			ID:   1,
			Name: "Item One",
			Tags: []string{"Red", "Blue", "Green"},
			Codes: map[string]string{
				"Color": "Red",
				"SIZE":  "Large",
				"type":  "Main",
			},
		},
		{
			ID:   2,
			Name: "Item Two",
			Tags: []string{"red", "blue", "green"},
			Codes: map[string]string{
				"color": "red",
				"size":  "large",
				"Type":  "Secondary",
			},
		},
		{
			ID:   3,
			Name: "Item Three",
			Tags: []string{"RED", "BLUE", "GREEN"},
			Codes: map[string]string{
				"COLOR": "RED",
				"Size":  "LARGE",
				"Type":  "MAIN",
			},
		},
	}

	valueTests := []struct {
		name     string
		query    string
		expected int
	}{
		{
			name:     "Case_sensitive_exact_value_match",
			query:    "ANY(Tags) = 'Red'",
			expected: 1, // Should only match item 1
		},
		{
			name:     "Case_sensitive_lowercase_value_match",
			query:    "ANY(Tags) = 'red'",
			expected: 1, // Should only match item 2
		},
		{
			name:     "Case_sensitive_uppercase_value_match",
			query:    "ANY(Tags) = 'RED'",
			expected: 1, // Should only match item 3
		},
		{
			name:     "Case_insensitive_ANY_value_comparison",
			query:    "ANY(Tags) = ANY('red', 'blue')",
			expected: 1, // Should match item 2 only (exact case)
		},
		{
			name:     "Case_comparison_not_equals",
			query:    "ANY(Tags) != 'Red'",
			expected: 2, // Should match items 2 and 3
		},
		{
			name:     "Case_mixed_value_matches",
			query:    "ANY(Codes.Type) = 'MAIN'",
			expected: 1, // Should match item 3 only
		},
	}
	// Run the tests
	for _, tt := range valueTests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := Parse[struct {
				ID    int
				Name  string
				Tags  []string
				Codes map[string]string
			}](tt.query, items)
			
			if err != nil {
				t.Fatalf("Failed to parse query '%s': %v", tt.query, err)
			}
			
			if len(results) != tt.expected {
				t.Errorf("Expected %d items, got %d for query %q", 
					tt.expected, len(results), tt.query)
				
				// Print IDs to help debugging
				var ids []int
				for _, item := range results {
					ids = append(ids, item.ID)
				}
				t.Logf("Matching IDs: %v", ids)
			}
		})
	}
}

func Test_ANY_Multiple_Value_Sets(t *testing.T) {
	// Tests with multiple ANY value sets
	items := []struct {
		ID       int
		Name     string
		Category string
		Tags     []string
		Scores   []int
		Active   bool
	}{
		{
			ID:       1,
			Name:     "Item One",
			Category: "Electronics",
			Tags:     []string{"popular", "featured", "sale"},
			Scores:   []int{90, 85, 95},
			Active:   true,
		},
		{
			ID:       2, 
			Name:     "Item Two",
			Category: "Clothing",
			Tags:     []string{"clearance", "winter", "discount"},
			Scores:   []int{70, 65, 75},
			Active:   true,
		},
		{
			ID:       3,
			Name:     "Item Three",
			Category: "Books",
			Tags:     []string{"bestseller", "fiction", "featured"},
			Scores:   []int{95, 90, 98},
			Active:   false,
		},
	}

	tests := []struct {
		name     string
		query    string
		expected int
	}{
		{
			name:     "ANY_with_multiple_match_values",
			query:    "ANY(Tags) = ANY('featured', 'bestseller')",
			expected: 2, // Items 1 and 3
		},
		{
			name:     "ANY_with_multiple_values_and_AND",
			query:    "ANY(Tags) = ANY('featured', 'sale') AND Active = true",
			expected: 1, // Item 1
		},
		{
			name:     "ANY_with_complex_value_combinations",
			query:    "(ANY(Tags) = ANY('featured', 'winter') OR ANY(Scores) > 90) AND Active = true",
			expected: 1, // Item 1
		},
		{
			name:     "ANY_with_not_equals_multiple_values",
			query:    "ANY(Tags) != ANY('clearance', 'winter') AND Category = 'Books'",
			expected: 1, // Item 3
		},
		{
			name:     "ANY_with_empty_value_set",
			query:    "ANY(Tags) = ANY()",
			expected: 0, // Should match none
		},
		{
			name:     "ANY_with_numeric_value_set",
			query:    "ANY(Scores) = ANY(90, 95)",
			expected: 2, // Items 1 and 3
		},
	}
	// Run the tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := Parse[struct {
				ID       int
				Name     string
				Category string
				Tags     []string
				Scores   []int
				Active   bool
			}](tt.query, items)
			
			if err != nil {
				t.Fatalf("Failed to parse query '%s': %v", tt.query, err)
			}
			
			if len(results) != tt.expected {
				t.Errorf("Expected %d items, got %d for query %q", 
					tt.expected, len(results), tt.query)
				
				// Print IDs to help debugging
				var ids []int
				for _, item := range results {
					ids = append(ids, item.ID)
				}
				t.Logf("Matching IDs: %v", ids)
			}
		})
	}
}
