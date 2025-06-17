package parser

import (
	"testing"
)

// Define example structs
type User struct {
	ID        int
	Name      string
	Email     string
	IsActive  bool
	Age       int
	Balance   float64
	Address   AddressInfo
	Contact   *ContactInfo // Pointer to demonstrate nil handling
	Interests []string
}

type AddressInfo struct {
	Street string
	City   string
	Zip    string
}

type ContactInfo struct {
	Phone string
	Email string
}

type test struct {
	query  string
	expRes int
}

func Test_main(t *testing.T) {
	users := []User{
		{ID: 1, Name: "Alice", Email: "alice@example.com", IsActive: true, Age: 30, Balance: 100.50, Address: AddressInfo{Street: "123 Main St", City: "Anytown", Zip: "12345"}, Contact: &ContactInfo{Phone: "111-222-3333"}},
		{ID: 2, Name: "Bob", Email: "bob@example.com", IsActive: false, Age: 25, Balance: 50.25, Address: AddressInfo{Street: "456 Oak Ave", City: "Anytown", Zip: "12345"}, Contact: &ContactInfo{Phone: ""}}, // Empty phone
		{ID: 3, Name: "Charlie", Email: "charlie@example.com", IsActive: true, Age: 35, Balance: 200.75, Address: AddressInfo{Street: "789 Pine Ln", City: "Otherville", Zip: "67890"}, Contact: nil},          // Nil contact
		{ID: 4, Name: "David", Email: "david@example.com", IsActive: true, Age: 28, Balance: 15.00, Address: AddressInfo{Street: "101 Elm Rd", City: "Anytown", Zip: "12345"}, Contact: &ContactInfo{Phone: "999-888-7777"}},
		{ID: 5, Name: "Eve", Email: "eve@example.com", IsActive: false, Age: 40, Balance: 120.00, Address: AddressInfo{Street: "202 Birch Blvd", City: "Otherville", Zip: "67890"}, Contact: &ContactInfo{Phone: "555-123-4567"}},
	}

	tm := make(map[int]test)
	tm[0] = test{
		query:  `ID = 1`,
		expRes: 1,
	}
	tm[1] = test{
		query:  `Name = 'Bob'`,
		expRes: 1,
	}
	tm[2] = test{
		query:  `Age > 30`,
		expRes: 2,
	}
	tm[3] = test{
		query:  `Balance < 100`,
		expRes: 2,
	}
	tm[4] = test{
		query:  `IsActive = true`,
		expRes: 3,
	}
	tm[5] = test{
		query:  `ID = 1 AND IsActive = true`,
		expRes: 1,
	}
	tm[6] = test{
		query:  `Age > 25 AND Balance > 100`,
		expRes: 3,
	}
	tm[7] = test{
		query:  `Address.City = 'Anytown'`,
		expRes: 3,
	}
	tm[8] = test{
		query:  `Address.Street = '123 Main St' AND Address.Zip = '12345'`,
		expRes: 1,
	}
	tm[9] = test{
		query:  `Contact.Phone != ''`,
		expRes: 3,
	}
	tm[10] = test{
		query:  `Contact.Phone = '111-222-3333' AND ID = 1`,
		expRes: 1,
	}
	tm[11] = test{
		query:  `Contact.Phone = ''`,
		expRes: 1,
	}
	tm[12] = test{
		query:  `Name != 'Charlie' AND Age < 30`,
		expRes: 2,
	}
	tm[13] = test{
		query:  `Age > 20 AND Age < 30`,
		expRes: 2,
	}
	tm[14] = test{
		query:  `Balance > 15.00 AND Balance < 100.00`,
		expRes: 1,
	}
	tm[15] = test{
		query:  `NonExistentField = 'test'`,
		expRes: 0,
	}
	tm[16] = test{
		query:  `(Age = 25 OR Age = 30) AND Address.City = 'Anytown'`,
		expRes: 2, // Bob (25) and Alice (30) both in Anytown
	}
	tm[17] = test{
		query:  `Age = 25 OR Age = 30`,
		expRes: 2, // Bob (25) and Alice (30)
	}
	tm[18] = test{
		query:  `Age = 25 OR (Age = 30 AND Address.City = 'Anytown')`,
		expRes: 2, // Bob (25) or Alice (30, Anytown)
	}

	for _, q := range tm {
		t.Logf("Query: %s\n", q.query)
		filteredUsers, err := Parse(q.query, users)
		if err != nil {
			t.Fatalf("  Error: %v\n", err)
			continue
		}
		if len(filteredUsers) != q.expRes {
			t.Fatalf("expected %d but got %d", q.expRes, len(filteredUsers))
		} else {
			for _, user := range filteredUsers {
				t.Logf("  - %+v\n", user)
			}
		}
		t.Log("---")
	}
}

// Additional struct types and parser tests

type Product struct {
	SKU      string
	Name     string
	Price    float64
	InStock  bool
	Category CategoryInfo
}

type CategoryInfo struct {
	Name string
	ID   int
}

type Order struct {
	OrderID   int
	Product   Product
	Quantity  int
	Total     float64
	Completed bool
}

func Test_parser_product(t *testing.T) {
	products := []Product{
		{SKU: "A1", Name: "Widget", Price: 19.99, InStock: true, Category: CategoryInfo{Name: "Gadgets", ID: 10}},
		{SKU: "B2", Name: "Gizmo", Price: 29.99, InStock: false, Category: CategoryInfo{Name: "Gadgets", ID: 10}},
		{SKU: "C3", Name: "Thingamajig", Price: 9.99, InStock: true, Category: CategoryInfo{Name: "Tools", ID: 20}},
	}

	tests := []struct {
		query  string
		expRes int
	}{
		{"Price > 10", 2},
		{"InStock = true", 2},
		{"Category.Name = 'Gadgets'", 2},
		{"SKU = 'C3' AND InStock = true", 1},
		{"Category.ID = 20", 1},
		{"Name != 'Gizmo'", 2},
	}

	for _, tc := range tests {
		t.Logf("Product Query: %s", tc.query)
		filtered, err := Parse(tc.query, products)
		if err != nil {
			t.Fatalf("  Error: %v", err)
		}
		if len(filtered) != tc.expRes {
			t.Fatalf("expected %d but got %d", tc.expRes, len(filtered))
		}
	}
}

func Test_parser_order(t *testing.T) {
	products := []Product{
		{SKU: "A1", Name: "Widget", Price: 19.99, InStock: true, Category: CategoryInfo{Name: "Gadgets", ID: 10}},
		{SKU: "B2", Name: "Gizmo", Price: 29.99, InStock: false, Category: CategoryInfo{Name: "Gadgets", ID: 10}},
	}
	orders := []Order{
		{OrderID: 100, Product: products[0], Quantity: 2, Total: 39.98, Completed: true},
		{OrderID: 101, Product: products[1], Quantity: 1, Total: 29.99, Completed: false},
		{OrderID: 102, Product: products[0], Quantity: 5, Total: 99.95, Completed: true},
	}

	tests := []struct {
		query  string
		expRes int
	}{
		{"Completed = true", 2},
		{"Product.Name = 'Gizmo'", 1},
		{"Quantity > 1 AND Completed = true", 2},
		{"Product.Category.Name = 'Gadgets' AND Total > 50", 1},
		{"OrderID = 101", 1},
		{"Product.InStock = true", 2},
	}

	for _, tc := range tests {
		t.Logf("Order Query: %s", tc.query)
		filtered, err := Parse(tc.query, orders)
		if err != nil {
			t.Fatalf("  Error: %v", err)
		}
		if len(filtered) != tc.expRes {
			t.Fatalf("expected %d but got %d", tc.expRes, len(filtered))
		}
	}
}

// --- Slice of struct for testing ---
type Tag struct {
	Name  string
	Value string
}

type BlogPost struct {
	Title   string
	Tags    []Tag
	Authors []string
}

func Test_parser_slice_fields(t *testing.T) {
	posts := []BlogPost{
		{Title: "Go Reflection", Tags: []Tag{{Name: "go", Value: "lang"}, {Name: "reflection", Value: "feature"}}, Authors: []string{"Alice", "Bob"}},
		{Title: "Python Tips", Tags: []Tag{{Name: "python", Value: "lang"}}, Authors: []string{"Carol"}},
		{Title: "Music Review", Tags: []Tag{{Name: "music", Value: "art"}}, Authors: []string{"Dave", "Eve"}},
	}
	tests := []struct {
		query  string
		expRes int
	}{
		{"Tags.Name = 'go'", 1},
		{"Tags.Value = 'lang'", 2},
		{"Tags.Name = 'music'", 1},
		{"Authors = 'Alice'", 1},
		{"Authors = 'Eve'", 1},
		{"Tags.Name = 'reflection' AND Authors = 'Bob'", 1},
		{"Tags.Name = 'notfound'", 0},
		{
			"Authors CONTAINS 'Alice'", 1,
		},
		{
			"Authors CONTAINS 'Eve'", 1,
		},
		{
			"Tags.Name CONTAINS 'reflection' AND Authors CONTAINS 'Bob'", 1,
		},
	}
	for _, tc := range tests {
		t.Logf("BlogPost Query: %s", tc.query)
		filtered, err := Parse(tc.query, posts)
		if err != nil {
			t.Fatalf("  Error: %v", err)
		}
		if len(filtered) != tc.expRes {
			t.Fatalf("expected %d but got %d", tc.expRes, len(filtered))
		}
	}
}
