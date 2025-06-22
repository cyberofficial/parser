//go:build example
// +build example

package main

import (
	"fmt"
	"log"

	"github.com/zveinn/parser"
)

// Person represents an employee with basic attributes.
type Person struct {
	Name       string
	Age        int
	IsEmployed bool
	Salary     float64
}

func main() {
	// Sample data: a slice of Person structs
	people := []Person{
		{Name: "Alice", Age: 30, IsEmployed: true, Salary: 75000.50},
		{Name: "Bob", Age: 25, IsEmployed: false, Salary: 65000.25},
		{Name: "Charlie", Age: 35, IsEmployed: true, Salary: 85000.75},
		{Name: "Diana", Age: 28, IsEmployed: true, Salary: 72000.00},
	}

	// Define queries to test basic filtering
	queries := []string{
		"AGE > 25 AND isemployed = true", // Case-insensitive field names
		"Name = 'Bob' OR Salary > 80000",
		"Age <= 30", // Test comparison operator
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
			fmt.Printf("- %s (Age: %d, Salary: %.2f, Employed: %v)\n", p.Name, p.Age, p.Salary, p.IsEmployed)
		}
		fmt.Println()
	}

	// Test error handling with an invalid query
	_, err := parser.Parse("InvalidField = 10", people)
	if err != nil {
		fmt.Printf("Expected error for invalid query: %v\n", err)
	}
}
