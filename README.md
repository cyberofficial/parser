# Parser

A lightweight, efficient Go library for filtering structured data using a simple query language. This parser allows you to filter collections of structs using SQL-like query expressions without requiring any database setup.

## Installation

```bash
go get github.com/zveinn/parser
```

## Usage

### Importing the package

```go
import "github.com/zveinn/parser"
```

### Basic Usage

The library provides a generic `Parse` function to filter slices of any struct type using query expressions:

```go
// Define your data structure
type Person struct {
    Name       string
    Age        int
    IsEmployed bool
    Skills     []string
    Salary     float64
    Department *Department
}

type Department struct {
    Name     string
    Location string
}

// Sample data
people := []Person{
    {Name: "Alice", Age: 30, IsEmployed: true, Skills: []string{"Go", "Python"}, Salary: 75000.50},
    {Name: "Bob", Age: 25, IsEmployed: false, Skills: []string{"Java", "C++"}, Salary: 65000.25},
    {Name: "Charlie", Age: 35, IsEmployed: true, Skills: []string{"Go", "Rust"}, Salary: 85000.75},
}

// Apply a query to filter the data
results, err := parser.Parse("Age > 25 AND IsEmployed = true", people)
if err != nil {
    log.Fatalf("Error parsing query: %v", err)
}

// Use the filtered results
fmt.Printf("Found %d matching results\n", len(results))
```

### Query Syntax

The query language supports:

#### Comparison Operators
- `=` Equal
- `!=` Not equal
- `>` Greater than
- `<` Less than
- `>=` Greater than or equal
- `<=` Less than or equal
- `CONTAINS` String contains

#### Logical Operators
- `AND` Logical AND
- `OR` Logical OR
- `NOT` Logical NOT

#### Special Operators
- `IS NULL` Check if a field is null/nil
- `ANY` Check if a field matches any value in a list

#### Examples

```
// Basic comparisons
Name = 'Alice'
Age > 30
Salary <= 75000.50
IsEmployed = true
Name CONTAINS 'Al'

// Null checks
Department IS NULL

// Logical operators
Age > 25 AND IsEmployed = true
Name = 'Alice' OR Name = 'Bob'
NOT (Age < 30)

// Nested expressions
(Age > 30 AND Salary > 80000) OR (IsEmployed = false)

// Array access and nested fields
Skills CONTAINS 'Go'
Department.Name = 'Engineering'

// ANY operator
Age = ANY(25, 30, 35)
```

### Advanced Features

- Handles nested struct fields using dot notation (e.g., `Department.Name`)
- Works with array/slice fields
- Supports pointers and nil checks
- Automatic type conversion for numeric comparisons

### Performance Considerations

The parser is designed to be efficient for filtering large datasets:

- No reflection during query compilation, only during evaluation
- Short-circuit evaluation for AND/OR expressions
- Handles nil pointer safety automatically

### Complete Example

```go
package main

import (
	"fmt"
	"log"

	"github.com/zveinn/parser"
)

type Person struct {
	Name       string
	Age        int
	IsEmployed bool
	Skills     []string
	Salary     float64
	Department *Department
	Tags       map[string]string
}

type Department struct {
	Name     string
	Location string
}

func main() {
	// Sample data
	people := []Person{
		{
			Name:       "Alice", 
			Age:        30, 
			IsEmployed: true, 
			Skills:     []string{"Go", "Python"}, 
			Salary:     75000.50,
			Department: &Department{Name: "Engineering", Location: "New York"},
			Tags:       map[string]string{"level": "senior", "remote": "yes"},
		},
		{
			Name:       "Bob", 
			Age:        25, 
			IsEmployed: false, 
			Skills:     []string{"Java", "C++"}, 
			Salary:     65000.25,
			Department: &Department{Name: "Engineering", Location: "Remote"},
			Tags:       map[string]string{"level": "junior", "remote": "yes"},
		},
		{
			Name:       "Charlie", 
			Age:        35, 
			IsEmployed: true, 
			Skills:     []string{"Go", "Rust"}, 
			Salary:     85000.75,
			Department: nil,
			Tags:       map[string]string{"level": "senior", "remote": "no"},
		},
	}

	// Example queries
	queries := []string{
		"Age > 25 AND IsEmployed = true",
		"Skills CONTAINS 'Go' AND Salary > 70000",
		"Department.Location = 'Remote' OR Department IS NULL",
		"Age = ANY(25, 35) OR (Department.Name = 'Engineering' AND Tags.remote = 'yes')",
	}

	for _, query := range queries {
		results, err := parser.Parse(query, people)
		if err != nil {
			log.Fatalf("Error parsing query '%s': %v", query, err)
		}
		
		fmt.Printf("Query: %s\nMatches: %d\n", query, len(results))
		for _, person := range results.([]Person) {
			fmt.Printf("- %s (Age: %d)\n", person.Name, person.Age)
		}
		fmt.Println()
	}
}
```

### Error Handling

The parser provides descriptive error messages for syntax issues:

```go
_, err := parser.Parse("Age >", people)
if err != nil {
    fmt.Println(err) // Will print: "failed to parse query: unexpected EOF"
}

_, err = parser.Parse("InvalidField = 10", people)
if err != nil {
    fmt.Println(err) // Will print something like: "evaluation error: field 'InvalidField' not found"
}
```


