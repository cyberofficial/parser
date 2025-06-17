package parser

import (
	"fmt"
	"testing"
	"time"
)

// Test structure with deeply nested data
type DeepNestedItem struct {
	ID        int
	Name      string
	Level1    *Level1Item
	Tags      []string
	MetaTags  map[string][]string
	Counts    map[string]int
	Enabled   bool
	CreatedAt time.Time
}

type Level1Item struct {
	Name   string
	Level2 *Level2Item
	Tags   []string
}

type Level2Item struct {
	Name   string
	Level3 *Level3Item
	Tags   []string
}

type Level3Item struct {
	Name   string
	Level4 *Level4Item
	Tags   []string
}

type Level4Item struct {
	Name    string
	Values  []int
	Tags    []string
	Options map[string]interface{}
}

func generateDeepNestedTestData(count int) []DeepNestedItem {
	items := make([]DeepNestedItem, count)

	for i := 0; i < count; i++ {
		// Create different item types for variability
		itemType := i % 5

		l4 := &Level4Item{
			Name:   fmt.Sprintf("L4-%d", i),
			Values: []int{i, i * 2, i * 3},
			Tags:   []string{fmt.Sprintf("deep-tag-%d", i%10)},
			Options: map[string]interface{}{
				"priority": i % 3,
				"visible":  i%2 == 0,
				"category": fmt.Sprintf("cat-%d", i%4),
			},
		}

		// Add special tags for some items
		if i%7 == 0 {
			l4.Tags = append(l4.Tags, "needle-in-haystack")
		}

		l3 := &Level3Item{
			Name:   fmt.Sprintf("L3-%d", i),
			Level4: l4,
			Tags:   []string{fmt.Sprintf("l3-tag-%d", i%8)},
		}

		l2 := &Level2Item{
			Name:   fmt.Sprintf("L2-%d", i),
			Level3: l3,
			Tags:   []string{fmt.Sprintf("l2-tag-%d", i%6)},
		}

		l1 := &Level1Item{
			Name:   fmt.Sprintf("L1-%d", i),
			Level2: l2,
			Tags:   []string{fmt.Sprintf("l1-tag-%d", i%4)},
		}

		// Create base item
		items[i] = DeepNestedItem{
			ID:   i,
			Name: fmt.Sprintf("Item-%d-Type-%d", i, itemType),
			Tags: []string{
				fmt.Sprintf("tag-%d", i%10),
				fmt.Sprintf("type-%d", itemType),
			},
			MetaTags: map[string][]string{
				"categories": {
					fmt.Sprintf("category-%d", i%5),
					fmt.Sprintf("segment-%d", i%3),
				},
				"keywords": {
					fmt.Sprintf("kw-%d", i%15),
					fmt.Sprintf("topic-%d", i%7),
				},
			},
			Counts: map[string]int{
				"views":     i * 10,
				"likes":     i * 5,
				"shares":    i * 2,
				"downloads": i % 100,
			},
			Enabled:   i%3 != 0,
			CreatedAt: time.Now().AddDate(0, 0, -i%30),
		}

		// Special cases
		switch itemType {
		case 0:
			// Standard case - complete path
			items[i].Level1 = l1
		case 1:
			// Break at level2 - nil Level3
			l1.Level2.Level3 = nil
			items[i].Level1 = l1
		case 2:
			// Break at level1 - nil Level2
			l1.Level2 = nil
			items[i].Level1 = l1
		case 3:
			// Nil Level1
			items[i].Level1 = nil
		case 4:
			// Complete path with special values
			if i%10 == 0 {
				l4.Tags = append(l4.Tags, "special-value", "performance-test")
			}
			items[i].Level1 = l1
			items[i].Tags = append(items[i].Tags, "root-special")
		}
	}

	return items
}

func Test_ANY_Performance_Deep(t *testing.T) {
	// Generate a larger dataset for performance testing
	items := generateDeepNestedTestData(100)

	tests := []struct {
		name     string
		query    string
		expected int // Expected number of results
	}{
		{
			name:     "Simple_ANY_comparison",
			query:    "ANY(Tags) = 'root-special'",
			expected: 20, // 20% of items have this tag
		},
		{
			name:     "Deeply_nested_ANY_path",
			query:    "ANY(Level1.Level2.Level3.Level4.Tags) = 'needle-in-haystack'",
			expected: 14, // 100/7 items have this tag at level 4
		},
		{
			name:     "Combined_ANY_expressions",
			query:    "ANY(Tags) = 'type-0' AND ANY(MetaTags.categories) = 'category-1'",
			expected: 4, // Type-0 (20%) AND category-1 (20%) = 4%
		},
		{
			name:     "Complex_nested_ANY_with_map_lookup",
			query:    "ANY(Level1.Level2.Level3.Level4.Options.category) = 'cat-2'",
			expected: 25, // 25% of valid items should have cat-2
		},
		{
			name:     "Multi_level_ANY_comparison",
			query:    "ANY(Level1.Tags) = 'l1-tag-1' OR ANY(Level1.Level2.Tags) = 'l2-tag-2'",
			expected: 37, // Complex intersection of path matches
		},
		{
			name:     "Deep_ANY_with_non_existent_values",
			query:    "ANY(Level1.Level2.Level3.Level4.Tags) = 'non-existent'",
			expected: 0,
		},
		{
			name:     "Multiple_chained_ANY",
			query:    "ANY(Tags) = ANY('type-1', 'type-2') AND ANY(MetaTags.keywords) = ANY('kw-1', 'kw-2', 'kw-3')",
			expected: 12, // Combined probability of matches
		},
		{
			name:     "ANY_with_numeric_comparison",
			query:    "ANY(Counts.views) > 500",
			expected: 50, // Items with ID > 50
		},
	}

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

			if len(results) != tt.expected {
				t.Errorf("Expected %d items, got %d for query %q",
					tt.expected, len(results), tt.query)
			}
		})
	}
}

func Test_ANY_Deep_Chained_Fields_Extended(t *testing.T) {
	// Use a smaller dataset for more precise assertions
	items := generateDeepNestedTestData(50)

	// Add specific test items with known values
	specialItem := DeepNestedItem{
		ID:   999,
		Name: "Special Test Item",
		Tags: []string{"special-test"},
		Level1: &Level1Item{
			Name: "Special L1",
			Tags: []string{"l1-special"},
			Level2: &Level2Item{
				Name: "Special L2",
				Tags: []string{"l2-special"},
				Level3: &Level3Item{
					Name: "Special L3",
					Tags: []string{"l3-special"},
					Level4: &Level4Item{
						Name:   "Special L4",
						Tags:   []string{"l4-special", "unique-value"},
						Values: []int{42, 99, 101},
						Options: map[string]interface{}{
							"uniqueKey": "unique-value",
							"nested": map[string]string{
								"deepKey": "deep-value",
							},
						},
					},
				},
			},
		},
	}
	items = append(items, specialItem)

	tests := []struct {
		name        string
		query       string
		expected    int
		expectedIDs []int // Expected IDs of matching items
	}{
		{
			name:        "Find_special_item_by_deep_unique_tag",
			query:       "ANY(Level1.Level2.Level3.Level4.Tags) = 'unique-value'",
			expected:    1,
			expectedIDs: []int{999},
		},
		{
			name:        "Complex_path_with_numeric_filter",
			query:       "ANY(Level1.Level2.Level3.Level4.Values) > 100",
			expected:    1,
			expectedIDs: []int{999},
		},
		{
			name:        "Multi_level_ANY_with_AND",
			query:       "ANY(Level1.Tags) = 'l1-special' AND ANY(Level1.Level2.Tags) = 'l2-special'",
			expected:    1,
			expectedIDs: []int{999},
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

			// Verify expected IDs
			if tt.expectedIDs != nil {
				for i, item := range results {
					found := false
					for _, expectedID := range tt.expectedIDs {
						if item.ID == expectedID {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("Result %d (ID=%d) was not in expected IDs list for query %q",
							i, item.ID, tt.query)
					}
				}
			}
		})
	}
}
