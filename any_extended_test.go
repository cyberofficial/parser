package parser

import (
	"testing"
)

// More complex test structures for ANY syntax testing
type DataItem struct {
	ID          int
	Name        string
	Properties  map[string]string
	Tags        []string
	Ratings     []float64
	Flags       []bool
	Numbers     []int
	NestedItems []NestedItem
	Active      bool
}

type NestedItem struct {
	Label    string
	Value    int
	SubTags  []string
	Variants map[string][]string
	IsValid  bool
}

func Test_ANY_Extended(t *testing.T) {
	// Create a more complex test dataset
	items := []DataItem{
		{
			ID:      1,
			Name:    "Item One",
			Active:  true,
			Tags:    []string{"important", "featured", "new"},
			Ratings: []float64{4.5, 4.8, 4.2},
			Flags:   []bool{true, false, true},
			Numbers: []int{10, 20, 30},
			Properties: map[string]string{
				"color": "red",
				"size":  "large",
			},
			NestedItems: []NestedItem{
				{
					Label:   "Nested A",
					Value:   100,
					SubTags: []string{"sub1", "sub2"},
					Variants: map[string][]string{
						"colors": {"red", "blue"},
						"sizes":  {"S", "M", "L"},
					},
					IsValid: true,
				},
				{
					Label:   "Nested B",
					Value:   200,
					SubTags: []string{"sub3", "sub4"},
					Variants: map[string][]string{
						"colors": {"green", "yellow"},
						"sizes":  {"XL", "XXL"},
					},
					IsValid: false,
				},
			},
		},
		{
			ID:      2,
			Name:    "Item Two",
			Active:  false,
			Tags:    []string{"archived", "old"},
			Ratings: []float64{3.2, 3.5},
			Flags:   []bool{false, false},
			Numbers: []int{5, 15, 25},
			Properties: map[string]string{
				"color": "blue",
				"size":  "medium",
			},
			NestedItems: []NestedItem{
				{
					Label:   "Nested C",
					Value:   150,
					SubTags: []string{"sub2", "sub5"},
					Variants: map[string][]string{
						"colors": {"black", "white"},
						"sizes":  {"S", "M"},
					},
					IsValid: true,
				},
			},
		},
		{
			ID:      3,
			Name:    "Item Three",
			Active:  true,
			Tags:    []string{"featured", "sale", "limited"},
			Ratings: []float64{5.0, 4.9, 5.0},
			Flags:   []bool{true, true, true},
			Numbers: []int{100, 200, 300},
			Properties: map[string]string{
				"color": "green",
				"size":  "small",
			},
			NestedItems: []NestedItem{
				{
					Label:   "Nested D",
					Value:   500,
					SubTags: []string{"premium", "exclusive"},
					Variants: map[string][]string{
						"colors": {"gold", "silver"},
						"sizes":  {"Custom", "Premium"},
					},
					IsValid: true,
				},
			},
		},
		{
			ID:      4,
			Name:    "Empty Item",
			Active:  false,
			Tags:    []string{},
			Ratings: []float64{},
			Flags:   []bool{},
			Numbers: []int{},
			Properties: map[string]string{
				"status": "empty",
			},
			NestedItems: []NestedItem{},
		},
		{
			ID:         5,
			Name:       "Partial Item",
			Active:     true,
			Tags:       []string{"incomplete"},
			Ratings:    []float64{2.5},
			Flags:      []bool{false},
			Numbers:    []int{-10},
			Properties: nil,
			NestedItems: []NestedItem{
				{
					Label:    "Nested Incomplete",
					Value:    0,
					SubTags:  []string{},
					Variants: nil,
					IsValid:  false,
				},
			},
		},
	}

	tests := []struct {
		name   string
		query  string
		expRes int
	}{
		// Basic ANY tests on different data types
		{
			name:   "ANY with string array - single match",
			query:  "ANY(Tags) = 'featured'",
			expRes: 2, // Items 1 and 3
		},
		{
			name:   "ANY with string array - multiple values",
			query:  "ANY(Tags) = ANY('featured', 'archived')",
			expRes: 3, // Items 1, 2, and 3
		},
		{
			name:   "ANY with numeric array - exact match",
			query:  "ANY(Numbers) = '30'",
			expRes: 1, // Item 1
		},
		{
			name:   "ANY with float array - range match",
			query:  "ANY(Ratings) > '4.5'",
			expRes: 2, // Items 1 and 3
		},
		{
			name:   "ANY with boolean array",
			query:  "ANY(Flags) = 'true' AND Active = true",
			expRes: 2, // Items 1 and 3
		},

		// Complex nested ANY tests
		{
			name:   "ANY with nested array field",
			query:  "ANY(NestedItems.SubTags) = 'sub2'",
			expRes: 2, // Items 1 and 2
		},
		{
			name:   "ANY with deeply nested structure",
			query:  "ANY(NestedItems.SubTags) = 'premium' AND Active = true",
			expRes: 1, // Item 3
		},
		{
			name:   "ANY with nested field and comparison",
			query:  "ANY(NestedItems.Value) > '200'",
			expRes: 1, // Item 3
		},

		// Mixed operations
		{
			name:   "ANY combined with NOT operator",
			query:  "NOT ANY(Tags) = 'featured'",
			expRes: 3, // Items 2, 4, and 5
		},
		{
			name:   "ANY combined with complex boolean logic",
			query:  "(ANY(Tags) = 'featured' OR ANY(Tags) = 'incomplete') AND Active = true",
			expRes: 3, // Items 1, 3, and 5
		},

		// Edge cases
		{
			name:   "ANY on empty array",
			query:  "ANY(Tags) = 'nonexistent' AND ID = 4",
			expRes: 0, // No matches
		},
		{
			name:   "ANY with negative number comparison",
			query:  "ANY(Numbers) < '0'",
			expRes: 1, // Item 5
		}, {
			name:   "ANY with numeric comparison using !=",
			query:  "ANY(Numbers) != '15'",
			expRes: 4, // All items except #5 which only has -10
		},

		// Complex boolean logic
		{
			name:   "ANY with parenthesized boolean logic",
			query:  "(ANY(Tags) = 'featured' AND Active = true) OR (ANY(Numbers) = '5' AND Active = false)",
			expRes: 3, // Items 1, 2, and 3
		}, {
			name:   "ANY with multiple nested conditions",
			query:  "(ANY(Tags) = ANY('featured', 'sale'))",
			expRes: 2, // Items 1 and 3
		},

		// Comparison operators
		{
			name:   "ANY with less-than-or-equal operator",
			query:  "ANY(Ratings) <= '3.5'",
			expRes: 2, // Items 2 and 5
		},
		{
			name:   "ANY with greater-than-or-equal operator",
			query:  "ANY(Ratings) >= '5.0'",
			expRes: 1, // Item 3
		}, {
			name:   "ANY with not-equal operator",
			query:  "ANY(Tags) != 'nonexistent'",
			expRes: 5, // All items have tags that are not 'nonexistent'
		},

		// IS NULL combinations
		{
			name:   "ANY with NULL check on nested items",
			query:  "ANY(Tags) = 'incomplete' AND Properties IS NULL",
			expRes: 1, // Item 5
		},
		{
			name:   "Complex query with NULL checks",
			query:  "(ANY(NestedItems.IsValid) = 'false' OR Properties IS NULL) AND Active = true",
			expRes: 2, // Items 1 and 5
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

// Test specifically for edge cases and error handling
func Test_ANY_EdgeCases(t *testing.T) {
	items := []struct {
		ID          int
		EmptyList   []string
		NilList     []string
		SingleItem  []int
		MixedTypes  []interface{}
		StringValue string
		NumValue    int
		HasData     bool
	}{
		{
			ID:          1,
			EmptyList:   []string{},
			SingleItem:  []int{42},
			MixedTypes:  []interface{}{"string", 123, true},
			StringValue: "test",
			NumValue:    100,
			HasData:     true,
		},
		{
			ID:          2,
			EmptyList:   []string{},
			NilList:     nil,
			SingleItem:  []int{99},
			MixedTypes:  []interface{}{456, "another"},
			StringValue: "sample",
			NumValue:    200,
			HasData:     false,
		},
	}

	tests := []struct {
		name      string
		query     string
		expectErr bool
		expRes    int
	}{
		{
			name:      "ANY on nil list",
			query:     "ANY(NilList) = 'anything'",
			expectErr: false,
			expRes:    0,
		},
		{
			name:      "ANY on empty list",
			query:     "ANY(EmptyList) = 'anything'",
			expectErr: false,
			expRes:    0,
		},
		{
			name:      "ANY with single item list",
			query:     "ANY(SingleItem) = '42'",
			expectErr: false,
			expRes:    1,
		},
		{
			name:      "Malformed ANY query - missing closing paren",
			query:     "ANY(EmptyList = 'test'",
			expectErr: true,
			expRes:    0,
		},
		{
			name:      "Malformed ANY query - no field",
			query:     "ANY() = 'test'",
			expectErr: true,
			expRes:    0,
		},
		{
			name:      "Malformed ANY query - empty value list",
			query:     "ANY(EmptyList) = ANY()",
			expectErr: true,
			expRes:    0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			filtered, err := Parse(tc.query, items)
			if tc.expectErr {
				if err == nil {
					t.Errorf("Expected error for query %q but got none", tc.query)
				}
			} else {
				if err != nil {
					t.Fatalf("Unexpected error for query %q: %v", tc.query, err)
				}
				if got := len(filtered); got != tc.expRes {
					t.Errorf("Expected %d items, got %d for query %q", tc.expRes, got, tc.query)
				}
			}
		})
	}
}

// Test for case sensitivity
func Test_ANY_CaseSensitivity(t *testing.T) {
	items := []struct {
		ID         int
		Tags       []string
		Properties map[string]string
	}{
		{
			ID:   1,
			Tags: []string{"Tag1", "TAG2", "tag3"},
			Properties: map[string]string{
				"Color": "Red",
				"SIZE":  "Large",
			},
		},
		{
			ID:   2,
			Tags: []string{"TAG1", "Tag2", "TAG3"},
			Properties: map[string]string{
				"color": "RED",
				"size":  "large",
			},
		},
	}

	tests := []struct {
		name   string
		query  string
		expRes int
	}{
		{
			name:   "Case-sensitive tag match",
			query:  "ANY(Tags) = 'Tag1'",
			expRes: 1, // Only item 1
		},
		{
			name:   "Case-insensitive property field",
			query:  "ANY(Properties.Color) = 'Red'",
			expRes: 1, // The field name is case-insensitive but values are case-sensitive
		},
		{
			name:   "Case-insensitive ANY with multiple values",
			query:  "ANY(Tags) = ANY('Tag1', 'Tag2')",
			expRes: 2, // Both items should match
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
			}
		})
	}
}
