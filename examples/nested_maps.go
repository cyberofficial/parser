//go:build example
// +build example

package main

import (
	"fmt"
	"log"

	"github.com/zveinn/parser"
)

// Department represents a department within a company.
type Department struct {
	Name     string
	Location string
}

// Person represents an employee with nested fields and metadata.
type Person struct {
	Name       string
	Age        int
	Skills     []string
	Department *Department
	Metadata   map[string]string
}

func main() {
	// Sample data: a slice of Person structs with nested fields and maps
	people := []Person{
		{
			Name:   "Alice",
			Age:    30,
			Skills: []string{"Go", "Python"},
			Department: &Department{
				Name:     "Engineering",
				Location: "New York",
			},
			Metadata: map[string]string{
				"level":  "senior",
				"remote": "yes",
			},
		},
		{
			Name:   "Bob",
			Age:    25,
			Skills: []string{"Java", "C++"},
			Department: &Department{
				Name:     "Marketing",
				Location: "Remote",
			},
			Metadata: map[string]string{
				"level":  "junior",
				"remote": "yes",
			},
		},
		{
			Name:       "Charlie",
			Age:        35,
			Skills:     []string{"Go", "Rust"},
			Department: nil,
			Metadata: map[string]string{
				"level":  "senior",
				"remote": "no",
			},
		},
	}

	// Define queries to test nested fields and maps
	queries := []string{
		"department.name = 'Engineering'",              // Nested struct access
		"metadata.level = 'senior' AND age > 28",       // Map access with logical operator
		"department IS NULL OR skills CONTAINS 'Rust'", // IS NULL and CONTAINS
	}

	// Run each query and print results
	for _, query := range queries {
		results, err := parser.Parse(query, people)
		if err != nil {
			log.Printf("Error parsing query '%s': %v", query, err)
			continue
		}

		fmt.Printf("Query: %s\n", query)
		fmt.Printf("Matches: %d\n", len(results))
		for _, p := range results {
			deptName := "nil"
			if p.Department != nil {
				deptName = p.Department.Name
			}
			fmt.Printf("- %s (Age: %d, Department: %s, Metadata: %v)\n", p.Name, p.Age, deptName, p.Metadata)
		}
		fmt.Println()
	}
}
