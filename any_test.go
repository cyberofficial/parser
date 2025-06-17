package parser

import (
	"testing"
)

type Item struct {
	ID       int
	Tags     []string
	Metadata *MetadataInfo
	Scores   []int
	Skills   []string
	Active   bool
}

type MetadataInfo struct {
	Categories []string
	Priority   int
}

func Test_ANY_Syntax(t *testing.T) {
	items := []Item{
		{
			ID:     1,
			Tags:   []string{"python", "code", "tutorial"},
			Scores: []int{90, 85, 92},
			Skills: []string{"coding", "writing"},
			Active: true,
			Metadata: &MetadataInfo{
				Categories: []string{"developer", "beginner"},
				Priority:   1,
			},
		},
		{
			ID:     2,
			Tags:   []string{"javascript", "content", "advanced"},
			Scores: []int{78, 80, 85},
			Skills: []string{"design", "coding"},
			Active: false,
			Metadata: &MetadataInfo{
				Categories: []string{"developer", "advanced"},
				Priority:   2,
			},
		},
		{
			ID:     3,
			Tags:   []string{"content", "article", "beginner"},
			Scores: []int{60, 65, 70},
			Skills: []string{"writing", "research"},
			Active: true,
			Metadata: &MetadataInfo{
				Categories: []string{"content", "beginner"},
				Priority:   3,
			},
		},
		{
			ID:       4,
			Tags:     []string{"python", "article"},
			Scores:   []int{95, 90, 92},
			Skills:   []string{"coding", "research"},
			Active:   false,
			Metadata: nil, // Test null metadata
		},
	}

	tests := []struct {
		name   string
		query  string
		expRes int
	}{
		{
			name:   "ANY equals with multiple values - match both",
			query:  "ANY(Tags) = ANY('python', 'content')",
			expRes: 4, // All items have either python or content
		},
		{
			name:   "ANY not equals with multiple values",
			query:  "ANY(Tags) != ANY('python', 'content')",
			expRes: 4, // All items have at least one tag that's not python or content
		},
		{
			name:   "ANY equals with single value in field",
			query:  "ANY(Metadata.Categories) = 'developer'",
			expRes: 2, // Items 1 and 2 have developer category
		},
		{
			name:   "ANY equals with numeric value",
			query:  "ANY(Scores) = '90'",
			expRes: 2, // Items 1 and 4 have score 90
		},
		{
			name:   "ANY with additional AND condition",
			query:  "ANY(Skills) = 'coding' AND Active = true",
			expRes: 1, // Only item 1 has coding skill and is active
		},
		{
			name:   "ANY with nested field on nil object",
			query:  "ANY(Metadata.Categories) = 'content'",
			expRes: 1, // Only item 3 has content category, item 4 has nil metadata
		},
		{
			name:   "ANY with complex AND-OR expression",
			query:  "(ANY(Tags) = 'python' OR ANY(Tags) = 'article') AND (ANY(Skills) = 'coding' OR ANY(Skills) = 'research')",
			expRes: 3, // Items 1, 3, and 4 match
		},
		{
			name:   "ANY with parentheses and multiple conditions",
			query:  "(ANY(Tags) = ANY('python', 'content')) AND (ANY(Scores) = '90' OR Active = true)",
			expRes: 3, // Items 1, 3, and 4 match
		},
		{
			name:   "ANY with NOT NULL check",
			query:  "ANY(Metadata.Categories) = 'developer' AND Metadata IS NOT NULL",
			expRes: 2, // Items 1 and 2 match
		},
		{
			name:   "Single ANY value comparison",
			query:  "ANY(Tags) = 'python'",
			expRes: 2, // Items 1 and 4 match
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
