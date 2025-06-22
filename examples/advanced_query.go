//go:build example
// +build example

package main

import (
	"fmt"
	"log"

	"github.com/zveinn/parser"
)

// Employee represents an employee with advanced attributes.
type Employee struct {
	Name       string
	Salary     float64
	Skills     []string
	Experience int
}

func main() {
	// Sample data: a slice of Employee structs
	employees := []Employee{
		{Name: "Alice", Salary: 75000.50, Skills: []string{"Go", "Python"}, Experience: 5},
		{Name: "Bob", Salary: 65000.25, Skills: []string{"Java", "C++"}, Experience: 3},
		{Name: "Charlie", Salary: 85000.75, Skills: []string{"Go", "Rust"}, Experience: 8},
		{Name: "Diana", Salary: 1.2e5, Skills: []string{"Python", "JavaScript"}, Experience: 6},
	}

	// Define queries to test advanced features
	queries := []string{
		"ANY(skills) = ANY('Go', 'Rust') AND salary > 7.5e4",             // ANY and scientific notation
		"(experience > 5 OR salary >= 100,000) AND NOT (name = 'Alice')", // Complex logic with NOT
		"salary > -1e3 AND skills CONTAINS 'Python'",                     // Negative scientific notation
	}

	// Run each query and print results
	for _, query := range queries {
		results, err := parser.Parse(query, employees)
		if err != nil {
			log.Printf("Error parsing query '%s': %v", query, err)
			continue
		}

		fmt.Printf("Query: %s\n", query)
		fmt.Printf("Matches: %d\n", len(results))
		for _, e := range results {
			fmt.Printf("- %s (Salary: %.2f, Skills: %v, Experience: %d)\n", e.Name, e.Salary, e.Skills, e.Experience)
		}
		fmt.Println()
	}

	// Test error handling with an invalid numeric format
	_, err := parser.Parse("salary > 75abc", employees)
	if err != nil {
		fmt.Printf("Expected error for invalid query: %v\n", err)
	}
}
