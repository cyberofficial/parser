package parser

import (
	"testing"
	"time"
)

// Complex and deeply nested test structures
type ComplexItem struct {
	ID         int
	Name       string
	Metadata   ComplexMetadata
	Tags       []string
	Categories []Category
	Variants   []Variant
	Stats      map[string]StatValue
	RelatedIDs []int
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Active     bool
}

type ComplexMetadata struct {
	Description string
	Keywords    []string
	Attributes  map[string]interface{}
	Visibility  struct {
		Public  bool
		Groups  []string
		Regions map[string]bool
	}
	Owner struct {
		ID       int
		Username string
		Roles    []string
		Settings map[string]interface{}
	}
}

type Category struct {
	ID       int
	Name     string
	ParentID *int
	Path     []string
	Level    int
	Active   bool
}

type Variant struct {
	ID         int
	SKU        string
	Price      float64
	Stock      int
	Attributes map[string]interface{}
	Options    []Option
	Available  bool
}

type Option struct {
	Name     string
	Value    interface{}
	Metadata map[string]string
}

type StatValue struct {
	Count      int
	Average    float64
	LastUpdate time.Time
	Trend      []float64
}

func generateComplexTestData() []ComplexItem {
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	lastWeek := now.AddDate(0, 0, -7)

	// Create parent category IDs
	parentElectronics := 100
	parentClothing := 200

	items := []ComplexItem{
		{
			ID:   1,
			Name: "Premium Smartphone",
			Metadata: ComplexMetadata{
				Description: "High-end smartphone with advanced features",
				Keywords:    []string{"smartphone", "mobile", "high-end", "premium", "5g"},
				Attributes: map[string]interface{}{
					"brand":           "TechX",
					"model":           "UltraPhone Pro",
					"release_year":    2022,
					"screen_size":     6.7,
					"storage_gb":      256,
					"ram_gb":          12,
					"water_resistant": true,
					"colors":          []string{"black", "silver", "blue"},
				},
				Visibility: struct {
					Public  bool
					Groups  []string
					Regions map[string]bool
				}{
					Public: true,
					Groups: []string{"premium", "featured"},
					Regions: map[string]bool{
						"US":   true,
						"EU":   true,
						"Asia": true,
					},
				},
				Owner: struct {
					ID       int
					Username string
					Roles    []string
					Settings map[string]interface{}
				}{
					ID:       101,
					Username: "tech_admin",
					Roles:    []string{"admin", "editor", "product_manager"},
					Settings: map[string]interface{}{
						"notifications": true,
						"theme":         "dark",
						"language":      "en",
					},
				},
			},
			Tags: []string{"smartphone", "premium", "5g", "new", "featured"},
			Categories: []Category{
				{
					ID:       101,
					Name:     "Smartphones",
					ParentID: &parentElectronics,
					Path:     []string{"Electronics", "Smartphones"},
					Level:    2,
					Active:   true,
				},
				{
					ID:       102,
					Name:     "Premium",
					ParentID: &parentElectronics,
					Path:     []string{"Electronics", "Premium"},
					Level:    2,
					Active:   true,
				},
			},
			Variants: []Variant{
				{
					ID:    1001,
					SKU:   "SP-BLK-256",
					Price: 999.99,
					Stock: 45,
					Attributes: map[string]interface{}{
						"color":   "Black",
						"storage": 256,
						"edition": "Standard",
					},
					Options: []Option{
						{
							Name:  "warranty",
							Value: "2 years",
							Metadata: map[string]string{
								"provider": "TechX Care",
								"type":     "extended",
							},
						},
						{
							Name:  "bundle",
							Value: "charger, earbuds",
							Metadata: map[string]string{
								"value": "$99",
								"type":  "free",
							},
						},
					},
					Available: true,
				},
				{
					ID:    1002,
					SKU:   "SP-SIL-256",
					Price: 999.99,
					Stock: 32,
					Attributes: map[string]interface{}{
						"color":   "Silver",
						"storage": 256,
						"edition": "Standard",
					},
					Options: []Option{
						{
							Name:  "warranty",
							Value: "2 years",
							Metadata: map[string]string{
								"provider": "TechX Care",
								"type":     "extended",
							},
						},
					},
					Available: true,
				},
				{
					ID:    1003,
					SKU:   "SP-BLU-512",
					Price: 1299.99,
					Stock: 15,
					Attributes: map[string]interface{}{
						"color":   "Blue",
						"storage": 512,
						"edition": "Limited",
					},
					Options: []Option{
						{
							Name:  "warranty",
							Value: "3 years",
							Metadata: map[string]string{
								"provider": "TechX Premium Care",
								"type":     "premium",
							},
						},
						{
							Name:  "bundle",
							Value: "wireless charger, premium earbuds, case",
							Metadata: map[string]string{
								"value": "$199",
								"type":  "free",
							},
						},
					},
					Available: true,
				},
			},
			Stats: map[string]StatValue{
				"views": {
					Count:      1250,
					Average:    623.5,
					LastUpdate: yesterday,
					Trend:      []float64{510, 580, 620, 690, 720, 790, 850},
				},
				"ratings": {
					Count:      89,
					Average:    4.7,
					LastUpdate: yesterday,
					Trend:      []float64{4.5, 4.6, 4.7, 4.7, 4.7},
				},
			},
			RelatedIDs: []int{2, 3, 7},
			CreatedAt:  lastWeek,
			UpdatedAt:  yesterday,
			Active:     true,
		},
		{
			ID:   2,
			Name: "Premium T-Shirt",
			Metadata: ComplexMetadata{
				Description: "Premium cotton t-shirt with unique design",
				Keywords:    []string{"clothing", "t-shirt", "premium", "cotton", "fashion"},
				Attributes: map[string]interface{}{
					"brand":        "FashionX",
					"material":     "100% organic cotton",
					"release_year": 2023,
					"gender":       "unisex",
					"sustainable":  true,
					"care":         "machine washable",
					"sizes":        []string{"S", "M", "L", "XL", "XXL"},
					"colors":       []string{"white", "black", "navy", "red"},
				},
				Visibility: struct {
					Public  bool
					Groups  []string
					Regions map[string]bool
				}{
					Public: true,
					Groups: []string{"clothing", "new-arrivals"},
					Regions: map[string]bool{
						"US":   true,
						"EU":   true,
						"Asia": false,
					},
				},
				Owner: struct {
					ID       int
					Username string
					Roles    []string
					Settings map[string]interface{}
				}{
					ID:       202,
					Username: "fashion_admin",
					Roles:    []string{"editor", "fashion_manager"},
					Settings: map[string]interface{}{
						"notifications": true,
						"theme":         "light",
						"language":      "en",
					},
				},
			},
			Tags: []string{"clothing", "t-shirt", "premium", "organic", "sustainable"},
			Categories: []Category{
				{
					ID:       201,
					Name:     "T-Shirts",
					ParentID: &parentClothing,
					Path:     []string{"Clothing", "T-Shirts"},
					Level:    2,
					Active:   true,
				},
				{
					ID:       202,
					Name:     "Premium",
					ParentID: &parentClothing,
					Path:     []string{"Clothing", "Premium"},
					Level:    2,
					Active:   true,
				},
				{
					ID:     203,
					Name:   "Sustainable",
					Path:   []string{"Sustainable"},
					Level:  1,
					Active: true,
				},
			},
			Variants: []Variant{
				{
					ID:    2001,
					SKU:   "TS-BLK-L",
					Price: 29.99,
					Stock: 120,
					Attributes: map[string]interface{}{
						"color": "Black",
						"size":  "L",
					},
					Available: true,
				},
				{
					ID:    2002,
					SKU:   "TS-WHT-M",
					Price: 29.99,
					Stock: 85,
					Attributes: map[string]interface{}{
						"color": "White",
						"size":  "M",
					},
					Available: true,
				},
				{
					ID:    2003,
					SKU:   "TS-RED-XL",
					Price: 29.99,
					Stock: 0,
					Attributes: map[string]interface{}{
						"color": "Red",
						"size":  "XL",
					},
					Available: false,
				},
			},
			Stats: map[string]StatValue{
				"views": {
					Count:      850,
					Average:    325.0,
					LastUpdate: yesterday,
					Trend:      []float64{290, 310, 330, 350, 370, 390, 410},
				},
				"ratings": {
					Count:      42,
					Average:    4.5,
					LastUpdate: yesterday,
					Trend:      []float64{4.2, 4.3, 4.4, 4.5, 4.5},
				},
			},
			RelatedIDs: []int{3, 4, 5},
			CreatedAt:  lastWeek,
			UpdatedAt:  now,
			Active:     true,
		},
	}

	return items
}

func Test_ANY_Complex_Nested_Queries(t *testing.T) {
	items := generateComplexTestData()

	tests := []struct {
		name     string
		query    string
		expected int
	}{
		// Deep field path tests
		{
			name:     "Deep_nested_field_lookup",
			query:    "ANY(Metadata.Keywords) = 'premium'",
			expected: 2, // Both items
		},
		{
			name:     "Multi_level_nested_lookup",
			query:    "ANY(Metadata.Visibility.Groups) = 'featured'",
			expected: 1, // Item 1
		},
		{
			name:     "Deep_nested_map_lookup",
			query:    "ANY(Metadata.Attributes.colors) = 'black'",
			expected: 2, // Both items
		},
		{
			name:     "Extremely_deep_nested_lookup",
			query:    "ANY(Variants.Options.Metadata.type) = 'premium'",
			expected: 1, // Item 1
		},

		// Complex object array lookups
		{
			name:     "Complex_object_array_field",
			query:    "ANY(Categories.Name) = 'Premium'",
			expected: 2, // Both items
		},
		{
			name:     "Complex_object_array_with_condition",
			query:    "ANY(Categories.Path) = 'Sustainable'",
			expected: 1, // Item 2
		},
		{
			name:     "Complex_nested_array_field",
			query:    "ANY(Variants.Attributes.color) = 'Blue'",
			expected: 1, // Item 1
		},
		{
			name:     "Deep_nested_array_in_object_array",
			query:    "ANY(Variants.Options.Metadata.type) = 'extended'",
			expected: 1, // Item 1
		},

		// Map value lookups
		{
			name:     "Map_with_complex_values",
			query:    "ANY(Stats.ratings.Trend) > 4.6",
			expected: 1, // Item 1
		},
		{
			name:     "Complex_map_lookup_with_condition",
			query:    "ANY(Metadata.Visibility.Regions.US) = 'true'",
			expected: 2, // Both items
		},
		{
			name:     "Owner_settings_deep_lookup",
			query:    "ANY(Metadata.Owner.Settings.theme) = 'dark'",
			expected: 1, // Item 1
		},

		// Multiple ANY clauses
		{
			name:     "Multiple_ANY_different_fields",
			query:    "ANY(Tags) = 'premium' AND ANY(Categories.Path) = 'Electronics'",
			expected: 1, // Item 1
		},
		{
			name:     "Multiple_ANY_same_field_different_values",
			query:    "ANY(Tags) = 'premium' AND ANY(Tags) = '5g'",
			expected: 1, // Item 1
		},
		{
			name:     "Multiple_ANY_with_OR",
			query:    "ANY(Tags) = 'smartphone' OR ANY(Tags) = 'clothing'",
			expected: 2, // Both items
		},

		// Complex combined queries
		{
			name:     "Complex_combined_lookups",
			query:    "(ANY(Metadata.Keywords) = 'premium' AND ANY(Categories.Level) > 1) OR (ANY(Stats.views.Count) > 1000 AND Active = true)",
			expected: 2, // Both items
		},
		{
			name:     "Different_value_types_in_query",
			query:    "ANY(Variants.Price) > 1000 AND ANY(Metadata.Attributes.water_resistant) = 'true'",
			expected: 1, // Item 1
		},
		{
			name:     "Multiple_ANY_with_numeric_and_string",
			query:    "ANY(Variants.Attributes.storage) > 256 AND ANY(Metadata.Owner.Roles) = 'admin'",
			expected: 1, // Item 1
		},
		{
			name:     "ANY_with_numeric_array_comparison",
			query:    "ANY(Stats.views.Trend) > 800",
			expected: 1, // Item 1
		},

		// Availability tests
		{
			name:     "Complex_condition_on_variants",
			query:    "ANY(Variants.Available) = 'false' AND ANY(Tags) = 'clothing'",
			expected: 1, // Item 2
		},
		{
			name:     "Stock_check_with_variant",
			query:    "ANY(Variants.Stock) = 0 AND Active = true",
			expected: 1, // Item 2
		},
		{
			name:     "Multiple_complex_conditions",
			query:    "(ANY(Variants.Price) > 1000 OR ANY(Categories.Name) = 'Sustainable') AND ANY(Metadata.Visibility.Groups) = 'featured'",
			expected: 1, // Item 1
		},
	}

	// Run the tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := Parse(tt.query, items)

			if err != nil {
				t.Fatalf("Failed to parse query '%s': %v", tt.query, err)
			}

			if len(results) != tt.expected {
				t.Errorf("Expected %d items, got %d for query %q",
					tt.expected, len(results), tt.query)

				// Print item IDs to help debugging
				var ids []int
				for _, item := range results {
					ids = append(ids, item.ID)
				}
				t.Logf("Matching items: %v", ids)
			}
		})
	}
}

func Test_ANY_Extreme_Edge_Cases(t *testing.T) {
	// Single item for testing extreme edge cases
	items := []struct {
		ID         int
		Name       string
		Empty      struct{}
		NilMap     map[string]string
		EmptyMap   map[string]string
		NilSlice   []string
		EmptySlice []string
		NilPtr     *string
		DeepNest   *****int // Super deep nesting
		Recursive  map[string]interface{}
	}{
		{
			ID:         1,
			Name:       "Edge Case Item",
			Empty:      struct{}{},
			NilMap:     nil,
			EmptyMap:   map[string]string{},
			NilSlice:   nil,
			EmptySlice: []string{},
			NilPtr:     nil,
		},
	}

	// Create a recursive structure (not too deep to avoid stack overflow)
	level3 := map[string]interface{}{"val": "level3"}
	level2 := map[string]interface{}{"next": level3, "val": "level2"}
	level1 := map[string]interface{}{"next": level2, "val": "level1"}
	items[0].Recursive = map[string]interface{}{"next": level1, "val": "root"}

	tests := []struct {
		name        string
		query       string
		expected    int
		expectError bool
	}{
		{
			name:     "Empty_struct_field_lookup",
			query:    "ANY(Empty) = 'anything'",
			expected: 0, // Should return no results
		},
		{
			name:     "Nil_map_lookup",
			query:    "ANY(NilMap.key) = 'value'",
			expected: 0, // Should return no results
		},
		{
			name:     "Empty_map_lookup",
			query:    "ANY(EmptyMap.key) = 'value'",
			expected: 0, // Should return no results
		},
		{
			name:     "Nil_slice_lookup",
			query:    "ANY(NilSlice) = 'anything'",
			expected: 0, // Should return no results
		},
		{
			name:     "Empty_slice_lookup",
			query:    "ANY(EmptySlice) = 'anything'",
			expected: 0, // Should return no results
		},
		{
			name:     "Nil_pointer_lookup",
			query:    "ANY(NilPtr) = 'anything'",
			expected: 0, // Should return no results
		},
		{
			name:     "Recursive_map_first_level",
			query:    "ANY(Recursive.val) = 'root'",
			expected: 1, // Should match
		},
		{
			name:     "Recursive_map_deep_level",
			query:    "ANY(Recursive.next.next.next.val) = 'level3'",
			expected: 1, // Should match
		},
	}

	// Run the tests
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
