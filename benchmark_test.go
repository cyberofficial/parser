package parser

import (
	"fmt"
	"testing"
)

// Create benchmark test data
type BenchPerson struct {
	Name       string
	Age        int
	IsEmployed bool
	Skills     []string
	Salary     float64
	Department *BenchDepartment
}

type BenchDepartment struct {
	Name     string
	Location string
}

// Benchmark simple queries
func BenchmarkBasicQueries(b *testing.B) {
	people := []BenchPerson{
		{Name: "Alice", Age: 30, IsEmployed: true, Skills: []string{"Go", "Python"}, Salary: 75000.50},
		{Name: "Bob", Age: 25, IsEmployed: false, Skills: []string{"Java", "C++"}, Salary: 65000.25},
		{Name: "Charlie", Age: 35, IsEmployed: true, Skills: []string{"Go", "Rust"}, Salary: 85000.75},
		{Name: "Diana", Age: 28, IsEmployed: true, Skills: []string{"Python", "JavaScript"}, Salary: 72000.00},
		{Name: "Eve", Age: 40, IsEmployed: true, Skills: []string{"C#", ".NET"}, Salary: 90000.00},
	}

	benchmarks := []struct {
		name  string
		query string
	}{
		{"StringEquality", "Name = 'Alice'"},
		{"NumberComparison", "Age > 30"},
		{"BooleanCheck", "IsEmployed = true"},
		{"StringContains", "Name CONTAINS 'li'"},
		{"SimpleAND", "Age > 25 AND IsEmployed = true"},
		{"SimpleOR", "Age = 25 OR Age = 35"},
		{"Parentheses", "(Age > 30 AND IsEmployed = true) OR Name = 'Bob'"},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = Parse(bm.query, people)
			}
		})
	}
}

// BenchmarkComplexQueries benchmarks more advanced query patterns
func BenchmarkComplexQueries(b *testing.B) {
	// Create a more complex data set with nested structures
	people := []BenchPerson{
		{
			Name: "Alice", Age: 30, IsEmployed: true,
			Skills:     []string{"Go", "Python", "SQL"},
			Salary:     75000.50,
			Department: &BenchDepartment{Name: "Engineering", Location: "New York"},
		},
		{
			Name: "Bob", Age: 25, IsEmployed: false,
			Skills:     []string{"Java", "C++", "JavaScript"},
			Salary:     65000.25,
			Department: &BenchDepartment{Name: "Engineering", Location: "Remote"},
		},
		{
			Name: "Charlie", Age: 35, IsEmployed: true,
			Skills:     []string{"Go", "Rust", "Docker"},
			Salary:     85000.75,
			Department: &BenchDepartment{Name: "DevOps", Location: "Seattle"},
		},
		{
			Name: "Diana", Age: 28, IsEmployed: true,
			Skills:     []string{"Python", "JavaScript", "React"},
			Salary:     72000.00,
			Department: &BenchDepartment{Name: "Frontend", Location: "San Francisco"},
		},
		{
			Name: "Eve", Age: 40, IsEmployed: true,
			Skills:     []string{"C#", ".NET", "Azure"},
			Salary:     90000.00,
			Department: &BenchDepartment{Name: "Engineering", Location: "Boston"},
		},
		{
			Name: "Frank", Age: 32, IsEmployed: true,
			Skills:     []string{"Go", "Kubernetes", "AWS"},
			Salary:     80000.00,
			Department: &BenchDepartment{Name: "DevOps", Location: "Seattle"},
		},
		{
			Name: "Grace", Age: 45, IsEmployed: true,
			Skills:     []string{"Java", "Spring", "Hibernate"},
			Salary:     95000.00,
			Department: &BenchDepartment{Name: "Backend", Location: "Chicago"},
		},
		{
			Name: "Henry", Age: 27, IsEmployed: false,
			Skills:     []string{"Python", "Django", "PostgreSQL"},
			Salary:     68000.00,
			Department: &BenchDepartment{Name: "Backend", Location: "Remote"},
		},
		{
			Name: "Ivy", Age: 38, IsEmployed: true,
			Skills:     []string{"JavaScript", "React", "Node.js"},
			Salary:     88000.00,
			Department: &BenchDepartment{Name: "Frontend", Location: "Austin"},
		},
		{
			Name: "Jack", Age: 50, IsEmployed: true,
			Skills:     []string{"Go", "C++", "Rust"},
			Salary:     110000.00,
			Department: nil,
		},
	}

	benchmarks := []struct {
		name  string
		query string
	}{
		// Complex logical combinations
		{"NestedLogic", "(Age > 30 AND Age < 45) AND (Salary > 80000 OR Department.Name = 'DevOps')"},
		{"MixedOperators", "(Age > 35 OR Salary > 85000) AND IsEmployed = true"},
		{"MultipleAND", "Age > 30 AND IsEmployed = true AND Salary > 80000 AND Department.Name = 'Engineering'"},
		{"DeepNesting", "((Age > 25 AND Age < 40) OR (Salary > 90000)) AND (IsEmployed = true OR Skills CONTAINS 'Python')"},

		// Queries with array/slice operations
		{"ArrayContains", "Skills CONTAINS 'Go'"},
		{"ComplexArrayLogic", "(Skills CONTAINS 'Go' AND Salary > 70000) OR (Skills CONTAINS 'React' AND Department.Location = 'San Francisco')"},

		// Nested field access
		{"NestedFields", "Department.Name = 'Engineering' AND Department.Location = 'New York'"},
		{"NestedWithLogic", "Department.Name = 'DevOps' OR (Department.Location = 'Remote' AND Age < 30)"},

		// IS NULL operator
		{"NullCheck", "Department IS NULL"},
		{"NotNullWithLogic", "Department IS NOT NULL AND Age > 40"},

		// ANY operator
		{"AnyOperator", "ANY(Skills) = 'Go'"},
		{"AnyWithMultipleValues", "ANY(Skills) = ANY('Go', 'React', 'AWS')"},
		{"AnyWithLogic", "ANY(Skills) = 'Go' AND Salary > 75000"},

		// NOT operator
		{"NotOperator", "NOT Age < 30"},
		{"ComplexNot", "NOT (Department.Name = 'Engineering' OR Department.Name = 'Frontend')"},

		// Very complex query combining multiple features
		{"SuperComplex", "((Age > 30 AND NOT (Department.Name = 'Engineering')) OR (Salary > 90000 AND Department IS NOT NULL)) AND (Skills CONTAINS 'Go' OR ANY(Skills) = ANY('React', 'AWS', 'Azure'))"},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = Parse(bm.query, people)
			}
		})
	}
}

// BenchmarkDataSizes benchmarks query performance with different dataset sizes
func BenchmarkDataSizes(b *testing.B) {
	// Create test data generator
	generateData := func(count int) []BenchPerson {
		data := make([]BenchPerson, count)
		departments := []*BenchDepartment{
			{Name: "Engineering", Location: "New York"},
			{Name: "Engineering", Location: "Remote"},
			{Name: "DevOps", Location: "Seattle"},
			{Name: "Frontend", Location: "San Francisco"},
			{Name: "Backend", Location: "Chicago"},
			nil,
		}

		skills := [][]string{
			{"Go", "Python", "SQL"},
			{"Java", "C++", "JavaScript"},
			{"Go", "Rust", "Docker"},
			{"Python", "JavaScript", "React"},
			{"C#", ".NET", "Azure"},
		}

		for i := 0; i < count; i++ {
			data[i] = BenchPerson{
				Name:       "Person-" + string(rune(65+(i%26))), // A-Z names
				Age:        20 + (i % 45),                       // Ages 20-64
				IsEmployed: i%3 != 0,                            // 2/3 are employed
				Skills:     skills[i%len(skills)],
				Salary:     60000 + float64(i%50)*1000, // 60k-110k salary
				Department: departments[i%len(departments)],
			}
		}
		return data
	}

	// Define queries of varying complexity
	queries := []struct {
		name  string
		query string
	}{
		{"Simple", "Age > 30"},
		{"Moderate", "Age > 30 AND IsEmployed = true"},
		{"Complex", "(Age > 30 AND IsEmployed = true) OR (Salary > 90000 AND Department IS NOT NULL)"},
		{"VeryComplex", "((Age > 30 OR Salary > 80000) AND Skills CONTAINS 'Go') OR (Department.Name = 'Engineering' AND Department.Location = 'Remote')"},
	}

	// Test with different data sizes
	sizes := []int{10, 100, 1000}
	for _, size := range sizes {
		data := generateData(size)

		for _, q := range queries {
			b.Run(q.name+"-Size-"+fmt.Sprint(size), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					_, _ = Parse(q.query, data)
				}
			})
		}
	}
}

// BenchmarkParserCharacteristics tests specific aspects of parser performance
func BenchmarkParserCharacteristics(b *testing.B) {
	// Create a dataset with 100 items
	data := make([]BenchPerson, 100)
	for i := 0; i < 100; i++ {
		data[i] = BenchPerson{
			Name:       fmt.Sprintf("Person-%d", i),
			Age:        20 + (i % 50),
			IsEmployed: i%2 == 0,
			Skills:     []string{"Skill1", "Skill2", "Skill3"},
			Salary:     50000 + float64(i*1000),
		}
	}

	// Test different parsing characteristics
	benchmarks := []struct {
		name  string
		query string
	}{
		// Test parsing speed versus query complexity
		{"Lexer_SimpleToken", "Age = 30"},
		{"Lexer_ComplexTokens", "Age >= 30 AND IsEmployed = true AND Name CONTAINS 'Person'"},

		// Test operator performance
		{"Op_Equality", "Name = 'Person-50'"},
		{"Op_Comparison", "Age > 50"},
		{"Op_Contains", "Name CONTAINS 'Person'"},

		// Test logical operator performance
		{"Logic_SingleAND", "Age > 30 AND IsEmployed = true"},
		{"Logic_MultipleAND", "Age > 30 AND IsEmployed = true AND Salary > 60000"},
		{"Logic_SingleOR", "Age = 25 OR Age = 35"},
		{"Logic_MultipleOR", "Age = 25 OR Age = 35 OR Age = 45"},
		{"Logic_MixedANDOR", "Age > 30 AND (IsEmployed = true OR Salary > 70000)"},

		// Test parenthesis parsing and nesting depth
		{"Paren_SingleLevel", "(Age > 30 AND IsEmployed = true)"},
		{"Paren_TwoLevels", "((Age > 30) AND (IsEmployed = true))"},
		{"Paren_ThreeLevels", "(((Age > 30) AND (IsEmployed = true)) OR (Salary > 80000))"},

		// Test query selectivity (percentage of records matched)
		{"Select_High", "Age >= 20"},   // Almost all records
		{"Select_Medium", "Age >= 45"}, // About half the records
		{"Select_Low", "Age >= 65"},    // Very few records
		{"Select_None", "Age > 100"},   // No records
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = Parse(bm.query, data)
			}
		})
	}
}

// BenchmarkSpecialOperators focuses on the performance of special operators
func BenchmarkSpecialOperators(b *testing.B) {
	// Create dataset with special cases
	people := []BenchPerson{
		{
			Name:       "Alice",
			Skills:     []string{"Go", "Python", "SQL", "Docker", "Kubernetes", "AWS"}, // Many skills
			Department: &BenchDepartment{Name: "Engineering", Location: "New York"},
		},
		{
			Name:       "Bob",
			Skills:     []string{}, // Empty skills array
			Department: nil,        // Null department
		},
		{
			Name:       "Charlie",
			Skills:     []string{"Go"},                           // Single skill
			Department: &BenchDepartment{Name: "", Location: ""}, // Empty strings
		},
	}

	// Sample of 100 people based on the above templates
	var data []BenchPerson
	for i := 0; i < 100; i++ {
		template := people[i%len(people)]
		person := template
		person.Age = i
		person.IsEmployed = i%3 == 0
		data = append(data, person)
	}

	benchmarks := []struct {
		name  string
		query string
	}{
		// IS NULL operator
		{"IsNull", "Department IS NULL"},
		{"IsNotNull", "Department IS NOT NULL"},

		// NOT operator
		{"Not_Simple", "NOT IsEmployed = true"},
		{"Not_Complex", "NOT (Age > 30 AND IsEmployed = true)"},

		// ANY operator
		{"Any_SingleValue", "ANY(Skills) = 'Go'"},
		{"Any_MultipleValues", "ANY(Skills) = ANY('Go', 'Python', 'Rust')"},
		{"Any_WithContains", "ANY(Skills) CONTAINS 'Go'"},

		// CONTAINS with different string lengths
		{"Contains_Short", "Name CONTAINS 'A'"},
		{"Contains_Medium", "Name CONTAINS 'Ali'"},
		{"Contains_Long", "Name CONTAINS 'Alice'"},

		// Empty arrays and edge cases
		{"Array_Empty", "Skills = ''"},
		{"Array_EmptyCheck", "Skills CONTAINS ''"},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = Parse(bm.query, data)
			}
		})
	}
}

// BenchmarkMemoryPatterns tests memory allocation characteristics
func BenchmarkMemoryPatterns(b *testing.B) {
	// Create a moderate dataset
	data := make([]BenchPerson, 50)
	for i := 0; i < 50; i++ {
		data[i] = BenchPerson{
			Name:       fmt.Sprintf("Person-%d", i),
			Age:        20 + i,
			IsEmployed: i%2 == 0,
			Skills:     []string{"Skill1", "Skill2", "Skill3"},
			Salary:     50000 + float64(i*1000),
			Department: &BenchDepartment{
				Name:     fmt.Sprintf("Dept-%d", i%5),
				Location: fmt.Sprintf("Location-%d", i%3),
			},
		}
	}

	benchmarks := []struct {
		name  string
		query string
	}{
		// Queries with increasing complexity to test memory allocation patterns
		{"Memory_VerySimple", "Age = 30"},                                                                         // Simplest case
		{"Memory_Simple", "Age > 30 AND IsEmployed = true"},                                                       // Simple logical operation
		{"Memory_Medium", "Age > 30 AND IsEmployed = true AND Salary > 60000"},                                    // More conditions
		{"Memory_Complex", "(Age > 30 AND IsEmployed = true) OR (Salary > 70000 AND Department.Name = 'Dept-1')"}, // With parentheses and OR
		{"Memory_VeryComplex", "((Age > 30 AND IsEmployed = true) OR (Salary > 70000)) AND (Department.Name = 'Dept-1' OR Department.Location CONTAINS 'Location')"}, // Nested parentheses

		// Specific operations that might have distinct memory patterns
		{"Memory_Contains", "Name CONTAINS 'Person'"},                                          // String operations
		{"Memory_IsNull", "Department IS NULL"},                                                // NULL checks
		{"Memory_Not", "NOT IsEmployed = true"},                                                // NOT operations
		{"Memory_Any", "ANY(Skills) = 'Skill1'"},                                               // ANY operator
		{"Memory_Nested", "Department.Name = 'Dept-1' AND Department.Location = 'Location-1'"}, // Nested field access
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			b.ReportAllocs() // Explicitly report allocations
			for i := 0; i < b.N; i++ {
				_, _ = Parse(bm.query, data)
			}
		})
	}
}

// BenchmarkQueryReuse tests the performance impact of reusing the same query multiple times
func BenchmarkQueryReuse(b *testing.B) {
	// Generate dataset with 1000 items
	data := make([]BenchPerson, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = BenchPerson{
			Name:       fmt.Sprintf("Person-%d", i),
			Age:        20 + (i % 50),
			IsEmployed: i%2 == 0,
			Skills:     []string{"Skill1", "Skill2", "Skill3"},
			Salary:     50000 + float64(i*1000),
		}
	}

	// Define the query we'll use repeatedly
	query := "Age > 30 AND IsEmployed = true AND Salary > 60000"

	b.Run("FirstRun", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = Parse(query, data)
		}
	})

	// Run the same query multiple times in a loop to see if there's any caching effect
	b.Run("RepeatedRuns", func(b *testing.B) {
		b.StopTimer()
		queriesRun := 0
		b.StartTimer()

		for i := 0; i < b.N; i++ {
			// Run the same query 10 times per iteration
			for j := 0; j < 10; j++ {
				_, _ = Parse(query, data)
				queriesRun++
			}
		}

		// Report the actual number of queries executed
		b.ReportMetric(float64(queriesRun)/float64(b.N), "queries/op")
	})
}
