//go:build example
// +build example

package main

import (
	"fmt"

	"github.com/zveinn/parser"
)

type User struct {
	Name    string
	City    string
	Address string
	Phone   string
	Email   string
}

func main() {
	users := []User{
		{Name: "Alice", City: "New York", Address: "123 Main St", Phone: "(212) 555-0101", Email: "alice@example.com"},
		{Name: "bob", City: "Los Angeles", Address: "456 Oak Ave", Phone: "(213) 555-0102", Email: "bob@example.com"},
		{Name: "Charlie", City: "chicago", Address: "789 Pine Ln", Phone: "(312) 555-0103", Email: "charlie@example.com"},
		{Name: "Kyle", City: "Houston", Address: "123 Main St", Phone: "(713) 555-0104", Email: "kyle@example.com"},
		{Name: "Diana", City: "New York", Address: "123 Main St", Phone: "(212) 555-0101", Email: "diana@example.com"},
		{Name: "Ethan", City: "chicago", Address: "789 Pine Ln", Phone: "(312) 555-0105", Email: "ethan@example.com"},
		{Name: "Fiona", City: "Miami", Address: "321 Palm Rd", Phone: "(305) 555-0106", Email: "FIONA@example.com"},
		{Name: "George", City: "New York", Address: "123 Main St", Phone: "(917) 555-0107", Email: "george@example.com"},
		{Name: "Hannah", City: "Los Angeles", Address: "456 Oak Ave", Phone: "(323) 555-0108", Email: "HANNAH@example.com"},
	}

	// --- Test: Default Case-Insensitive Search ---
	fmt.Println("--- Test: Default Case-Insensitive Search ---")
	results, _ := parser.Parse("Name = 'alice'", users)
	fmt.Printf("Found %d user with Name = 'alice': %v\n", len(results), results)

	results, _ = parser.Parse("Name = 'ALICE'", users)
	fmt.Printf("Found %d user with Name = 'ALICE': %v\n", len(results), results)

	// --- Test: UPPER Function ---
	fmt.Println("\n--- Test: UPPER Function ---")
	results, _ = parser.Parse("UPPER(Name) = 'BOB'", users)
	fmt.Printf("Found %d user with UPPER(Name) = 'BOB': %v\n", len(results), results)

	results, _ = parser.Parse("UPPER(Name) = 'alice'", users)
	fmt.Printf("Found %d user with UPPER(Name) = 'alice': %v\n", len(results), results)

	// --- Test: LOWER Function ---
	fmt.Println("\n--- Test: LOWER Function ---")
	results, _ = parser.Parse("LOWER(Name) = 'charlie'", users)
	fmt.Printf("Found %d user with LOWER(Name) = 'charlie': %v\n", len(results), results)

	results, _ = parser.Parse("LOWER(Email) = 'fiona@example.com'", users)
	fmt.Printf("Found %d user with LOWER(Email) = 'fiona@example.com': %v\n", len(results), results)

	// --- Test: EXACT Function for Case-Sensitive Matching ---
	fmt.Println("\n--- Test: EXACT Function for Case-Sensitive Matching ---")
	results, _ = parser.Parse("EXACT(Name) = 'Kyle'", users)
	fmt.Printf("Found %d user with EXACT(Name) = 'Kyle': %v\n", len(results), results)

	results, _ = parser.Parse("EXACT(Name) = 'kyle'", users)
	fmt.Printf("Found %d user with EXACT(Name) = 'kyle': %v\n", len(results), results)

	results, _ = parser.Parse("EXACT(Email) = 'HANNAH@example.com'", users)
	fmt.Printf("Found %d user with EXACT(Email) = 'HANNAH@example.com': %v\n", len(results), results)

	// --- Test: Mixed Operations with AND, OR, NOT ---
	fmt.Println("\n--- Test: Mixed Operations with AND, OR, NOT ---")
	results, _ = parser.Parse("LOWER(City) = 'new york' AND Address CONTAINS 'Main'", users)
	fmt.Printf("Found %d users in New York on Main St: %v\n", len(results), results)

	results, _ = parser.Parse("(LOWER(Name) = 'alice' OR LOWER(Name) = 'bob') AND Phone CONTAINS '555'", users)
	fmt.Printf("Found %d users named Alice or Bob with '555' in phone: %v\n", len(results), results)

	results, _ = parser.Parse("LOWER(City) = 'chicago' AND NOT Address CONTAINS '123'", users)
	fmt.Printf("Found %d users in Chicago not living at an address with '123': %v\n", len(results), results)

	results, _ = parser.Parse("EXACT(Name) = 'Fiona' OR (LOWER(City) = 'los angeles' AND Email CONTAINS 'hannah')", users)
	fmt.Printf("Found %d users who are Fiona or live in LA and have 'hannah' in their email: %v\n", len(results), results)

	// --- Test: Searching by Area Code and Shared Information ---
	fmt.Println("\n--- Test: Searching by Area Code and Shared Information ---")
	results, _ = parser.Parse("Phone CONTAINS '(212)'", users)
	fmt.Printf("Found %d users in the (212) area code: %v\n", len(results), results)

	results, _ = parser.Parse("Address = '123 Main St' AND Phone != '(713) 555-0104'", users)
	fmt.Printf("Found %d users at '123 Main St' not named Kyle: %v\n", len(results), results)

	results, _ = parser.Parse("LOWER(City) = 'los angeles' OR Phone CONTAINS '(312)'", users)
	fmt.Printf("Found %d users in LA or with a Chicago area code: %v\n", len(results), results)

	// --- Test: Exclusionary Searches for Specificity ---
	fmt.Println("\n--- Test: Exclusionary Searches for Specificity ---")
	// Find people at '123 Main St' but exclude Kyle.
	results, _ = parser.Parse("Address = '123 Main St' AND NOT EXACT(Name) = 'Kyle'", users)
	fmt.Printf("Found %d users at '123 Main St' who are not Kyle: %v\n", len(results), results)

	// Find people in New York, but only those with a '(917)' area code.
	results, _ = parser.Parse("LOWER(City) = 'new york' AND Phone CONTAINS '(917)'", users)
	fmt.Printf("Found %d users in New York with a (917) area code: %v\n", len(results), results)

	// Find people who share a phone number with Alice, but are not Alice.
	results, _ = parser.Parse("Phone = '(212) 555-0101' AND Name != 'Alice'", users)
	fmt.Printf("Found %d users with Alice's phone number who are not Alice: %v\n", len(results), results)

	// Find people living at '456 Oak Ave' but exclude Bob and anyone with 'HANNAH' in their email.
	results, _ = parser.Parse("Address = '456 Oak Ave' AND NOT (LOWER(Name) = 'bob' OR Email CONTAINS 'HANNAH')", users)
	fmt.Printf("Found %d users at '456 Oak Ave' who are not Bob or Hannah: %v\n", len(results), results)
}
