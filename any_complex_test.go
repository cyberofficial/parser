package parser

import (
	"strings"
	"testing"
)

// Test specific error conditions and edge cases
func Test_ANY_ErrorHandling(t *testing.T) {
	// Simple test data
	data := []struct {
		ID    int
		Tags  []string
		Value string
	}{
		{ID: 1, Tags: []string{"a", "b", "c"}, Value: "test"},
	}

	tests := []struct {
		name        string
		query       string
		expectError bool
		errorSubstr string
	}{
		{
			name:        "Missing field name in ANY",
			query:       "ANY() = 'value'",
			expectError: true,
			errorSubstr: "field name",
		},
		{
			name:        "Missing parenthesis after ANY",
			query:       "ANY = 'value'",
			expectError: true,
			errorSubstr: "expected",
		},
		{
			name:        "Missing closing parenthesis",
			query:       "ANY(Tags = 'value'",
			expectError: true,
			errorSubstr: "expected",
		}, {
			name:        "Invalid operator with ANY",
			query:       "ANY(Tags) INVALID 'value'",
			expectError: true,
			errorSubstr: "expected",
		},
		{
			name:        "Missing value after ANY",
			query:       "ANY(Tags) =",
			expectError: true,
			errorSubstr: "expected",
		},
		{
			name:        "Missing values in ANY() = ANY()",
			query:       "ANY(Tags) = ANY()",
			expectError: true,
			errorSubstr: "expected",
		},
		{
			name:        "Invalid comma usage",
			query:       "ANY(Tags) = ANY('a', , 'b')",
			expectError: true,
			errorSubstr: "expected",
		},
		{
			name:        "Missing value after comma",
			query:       "ANY(Tags) = ANY('a',)",
			expectError: true,
			errorSubstr: "expected",
		},
		{
			name:        "Valid query with no syntax errors",
			query:       "ANY(Tags) = 'a'",
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse(tc.query, data)
			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error for query %q but got none", tc.query)
				} else if tc.errorSubstr != "" && !strings.Contains(err.Error(), tc.errorSubstr) {
					t.Errorf("Error message %q does not contain %q", err.Error(), tc.errorSubstr)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for query %q: %v", tc.query, err)
				}
			}
		})
	}
}

// Test complex query combinations
func Test_ANY_ComplexQueries(t *testing.T) {
	type ComplexItem struct {
		ID         int
		Name       string
		Tags       []string
		Categories []string
		Scores     map[string]int
		Details    struct {
			Status    string
			Flags     []string
			Metadata  map[string][]string
			SubScores []int
		}
		Active bool
	}

	items := []ComplexItem{
		{
			ID:         1,
			Name:       "First Item",
			Tags:       []string{"important", "featured"},
			Categories: []string{"A", "B"},
			Scores:     map[string]int{"math": 90, "science": 85},
			Details: struct {
				Status    string
				Flags     []string
				Metadata  map[string][]string
				SubScores []int
			}{
				Status:    "active",
				Flags:     []string{"new", "premium"},
				Metadata:  map[string][]string{"authors": {"John", "Jane"}},
				SubScores: []int{10, 20, 30},
			},
			Active: true,
		},
		{
			ID:         2,
			Name:       "Second Item",
			Tags:       []string{"archived", "old"},
			Categories: []string{"B", "C"},
			Scores:     map[string]int{"math": 70, "science": 75},
			Details: struct {
				Status    string
				Flags     []string
				Metadata  map[string][]string
				SubScores []int
			}{
				Status:    "inactive",
				Flags:     []string{"deprecated"},
				Metadata:  map[string][]string{"authors": {"Bob"}},
				SubScores: []int{5, 15, 25},
			},
			Active: false,
		},
		{
			ID:         3,
			Name:       "Third Item",
			Tags:       []string{"featured", "popular"},
			Categories: []string{"A", "C"},
			Scores:     map[string]int{"math": 95, "science": 95},
			Details: struct {
				Status    string
				Flags     []string
				Metadata  map[string][]string
				SubScores []int
			}{
				Status:    "active",
				Flags:     []string{"new", "recommended"},
				Metadata:  map[string][]string{"authors": {"John", "Alice"}},
				SubScores: []int{30, 40, 50},
			},
			Active: true,
		},
	}

	tests := []struct {
		name   string
		query  string
		expRes int
	}{
		{
			name:   "Complex query with multiple ANYs and nested fields",
			query:  "(ANY(Tags) = 'featured' OR ANY(Categories) = 'C') AND Details.Status = 'active'",
			expRes: 2, // Items 1 and 3
		},
		{
			name:   "Deep nesting with ANY and multiple conditions",
			query:  "ANY(Details.Flags) = 'new' AND (ANY(Tags) = 'featured' OR ANY(Categories) = 'B')",
			expRes: 2, // Items 1 and 3
		}, {
			name:   "Deeply nested numeric comparison",
			query:  "ANY(Details.SubScores) > '20' AND Active = true",
			expRes: 1, // Item 3
		},
		{
			name:   "Complex multi-level query with negative condition",
			query:  "(ANY(Tags) = ANY('featured', 'popular') AND NOT ANY(Details.Flags) = 'deprecated') OR (ANY(Categories) = 'B' AND Active = false)",
			expRes: 3, // All items
		},
		{
			name:   "Very complex query with multiple ANYs, nesting, and boolean logic",
			query:  "((ANY(Tags) = 'featured' AND ANY(Details.Flags) = 'new') OR (ANY(Categories) = 'C' AND ANY(Details.SubScores) < '10')) AND (Active = true OR Details.Status = 'inactive')",
			expRes: 3, // All items
		},
		{
			name:   "Multiple categories of conditions with ANYs",
			query:  "(ANY(Tags) = 'featured' AND ANY(Categories) = 'A') OR (ANY(Details.Flags) = 'deprecated' AND NOT Active = true)",
			expRes: 3, // All items
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			filtered, err := Parse(tc.query, items)
			if err != nil {
				t.Fatalf("Error parsing complex query %q: %v", tc.query, err)
			}
			if got := len(filtered); got != tc.expRes {
				t.Errorf("Expected %d items, got %d for query %q", tc.expRes, got, tc.query)
				for i, item := range filtered {
					t.Logf("Result %d: %+v", i, item)
				}
			}
		})
	}
}

// Test for interaction between ANY and slices of custom types
func Test_ANY_CustomTypes(t *testing.T) {
	type Tag struct {
		Name  string
		Value string
	}

	type CustomItem struct {
		ID      int
		Tags    []Tag
		Enabled bool
	}

	items := []CustomItem{
		{
			ID: 1,
			Tags: []Tag{
				{Name: "category", Value: "electronics"},
				{Name: "price", Value: "high"},
			},
			Enabled: true,
		},
		{
			ID: 2,
			Tags: []Tag{
				{Name: "category", Value: "clothing"},
				{Name: "price", Value: "medium"},
			},
			Enabled: false,
		},
		{
			ID: 3,
			Tags: []Tag{
				{Name: "category", Value: "electronics"},
				{Name: "price", Value: "medium"},
			},
			Enabled: true,
		},
	}

	tests := []struct {
		name   string
		query  string
		expRes int
	}{
		{
			name:   "ANY with nested struct fields in slice",
			query:  "ANY(Tags.Name) = 'price'",
			expRes: 3, // All items
		},
		{
			name:   "ANY with nested struct field and value",
			query:  "ANY(Tags.Value) = 'high'",
			expRes: 1, // Item 1
		},
		{
			name:   "ANY with complex condition on nested fields",
			query:  "(ANY(Tags.Name) = 'category' AND ANY(Tags.Value) = 'electronics') AND Enabled = true",
			expRes: 2, // Items 1 and 3
		},
		{
			name:   "Combined fields from same nested object",
			query:  "ANY(Tags.Name) = 'price' AND ANY(Tags.Value) = 'medium'",
			expRes: 2, // Items 2 and 3
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			filtered, err := Parse(tc.query, items)
			if err != nil {
				t.Fatalf("Error parsing query %q: %v", tc.query, err)
			}
			if got := len(filtered); got != tc.expRes {
				t.Errorf("Expected %d items, got %d for query %q", tc.expRes, got, tc.query)
				for i, item := range filtered {
					t.Logf("Result %d: %+v", i, item)
				}
			}
		})
	}
}
