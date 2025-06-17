package parser_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/zveinn/parser"
)

// Define benchmark structures
type BenchItem struct {
	ID       int
	Name     string
	Age      int
	IsActive bool
	Score    float64
	Tags     []string
}

func BenchmarkSimpleEquality(b *testing.B) {
	// Create sample data
	data := make([]BenchItem, 1000)
	for i := 0; i < 1000; i++ {
		name := fmt.Sprintf("Item-%d", i)
		data[i] = BenchItem{
			ID:       i,
			Name:     name,
			Age:      20 + (i % 50),
			IsActive: i%2 == 0,
			Score:    float64(i) * 1.5,
			Tags:     []string{"tag1", "tag2"},
		}
	}

	// Run the benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parser.Parse("Name = 'Item-10'", data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSimpleComparison(b *testing.B) {
	// Create sample data
	data := make([]BenchItem, 1000)
	for i := 0; i < 1000; i++ {
		name := fmt.Sprintf("Item-%d", i)
		data[i] = BenchItem{
			ID:       i,
			Name:     name,
			Age:      20 + (i % 50),
			IsActive: i%2 == 0,
			Score:    float64(i) * 1.5,
			Tags:     []string{"tag1", "tag2"},
		}
	}

	// Run the benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parser.Parse("Age > 30", data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkComplexQuery(b *testing.B) {
	// Create sample data
	data := make([]BenchItem, 1000)
	for i := 0; i < 1000; i++ {
		name := fmt.Sprintf("Item-%d", i)
		data[i] = BenchItem{
			ID:       i,
			Name:     name,
			Age:      20 + (i % 50),
			IsActive: i%2 == 0,
			Score:    float64(i) * 1.5,
			Tags:     []string{"tag1", "tag2"},
		}
	}

	// Run the benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parser.Parse("(Age > 30 AND IsActive = true) OR (Age < 25 AND Score > 10)", data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkWithDataSize(b *testing.B) {
	sizes := []int{10, 100, 1000}
	
	for _, size := range sizes {
		b.Run("Size-"+strconv.Itoa(size), func(b *testing.B) {
			// Create sample data of specified size
			data := make([]BenchItem, size)
			for i := 0; i < size; i++ {
				name := fmt.Sprintf("Item-%d", i)
				data[i] = BenchItem{
					ID:       i,
					Name:     name,
					Age:      20 + (i % 50),
					IsActive: i%2 == 0,
					Score:    float64(i) * 1.5,
					Tags:     []string{"tag1", "tag2"},
				}
			}
			
			// Run the benchmark
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := parser.Parse("Age > 30 AND IsActive = true", data)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
