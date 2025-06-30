# Parser 🧠

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/zveinn/parser)
![GitHub Issues](https://img.shields.io/github/issues/zveinn/parser)
![GitHub Stars](https://img.shields.io/github/stars/zveinn/parser?style=social)

**Parser** is a lightweight, efficient Go library for filtering slices of structs using a SQL-like query language. It enables in-memory data filtering without a database, ideal for applications like data processing, configuration management, or API response filtering. Built with Go generics, it offers type-safe queries and supports complex expressions, nested fields, and case-insensitive matching.

## 🚀 Features

- **SQL-Like Query Language**: Filter structs with intuitive queries (e.g., `Age > 25 AND Skills CONTAINS 'Go'`).
- **Type-Safe with Generics**: Works with any struct type using Go’s generics.
- **Nested Field Access**: Query nested structs and maps using dot notation (e.g., `Department.Name`).
- **Humanized Values Support**: Parse human-readable values like `10GB`/`10GiB`, `1.5K`, `2TB`/`2TiB`, `1,000` automatically.
- **Rich Operators**: Supports `=`, `!=`, `<`, `>`, `<=`, `>=`, `CONTAINS`, `IS NULL`, `ANY`, `NOT`, `AND`, `OR`.
- **Case-Insensitive Matching**: Field names and keywords (e.g., `AND`, `OR`) are case-insensitive.
- **Efficient Parsing**: Uses an enhanced lexer with support for negative numbers, scientific notation, and comma-separated numbers.
- **Robust Error Handling**: Detailed error messages for syntax and evaluation errors.
- **Zero Dependencies**: Pure Go implementation, no external libraries required.

## 📋 Requirements

- Go 1.24.1 or higher (for generics and latest features)
- No external dependencies

## 📦 Installation

Install the library using Go modules:

```bash
go get github.com/zveinn/parser
```

## 🔧 Usage

### Basic Example

Filter a slice of structs using a query:

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
}

type Department struct {
    Name     string
    Location string
}

func main() {
    people := []Person{
        {Name: "Alice", Age: 30, IsEmployed: true, Skills: []string{"Go", "Python"}, Salary: 75000.50, Department: &Department{Name: "Engineering", Location: "New York"}},
        {Name: "Bob", Age: 25, IsEmployed: false, Skills: []string{"Java", "C++"}, Salary: 65000.25, Department: &Department{Name: "Engineering", Location: "Remote"}},
        {Name: "Charlie", Age: 35, IsEmployed: true, Skills: []string{"Go", "Rust"}, Salary: 85000.75, Department: nil},
    }

    results, err := parser.Parse("Age > 25 AND isemployed = true", people)
    if err != nil {
        log.Fatalf("Error parsing query: %v", err)
    }

    for _, p := range results {
        fmt.Printf("Match: %s (Age: %d)\n", p.Name, p.Age)
    }
}
```

**Output**:
```
Match: Alice (Age: 30)
Match: Charlie (Age: 35)
```

### Query Syntax

The query language supports a variety of operators and expressions:

#### Comparison Operators
| Operator   | Description              | Example                |
|------------|--------------------------|------------------------|
| `=`        | Equal                    | `Name = 'Alice'`       |
| `!=`       | Not equal                | `Age != 30`            |
| `>`        | Greater than             | `Salary > 70,000`      |
| `<`        | Less than                | `Age < 35`             |
| `>=`       | Greater than or equal    | `Salary >= 75000.50`   |
| `<=`       | Less than or equal       | `Age <= 30`            |
| `CONTAINS` | String or slice contains | `Skills CONTAINS 'Go'` |

#### Logical Operators
| Operator | Description | Example                          |
|----------|-------------|----------------------------------|
| `AND`    | Logical AND | `Age > 25 AND IsEmployed = true` |
| `OR`     | Logical OR  | `Name = 'Alice' OR Name = 'Bob'` |
| `NOT`    | Logical NOT | `NOT (Age < 30)`                 |

#### Special Operators
| Operator      | Description               | Example                           |
|---------------|---------------------------|-----------------------------------|
| `IS NULL`     | Check for nil/zero value  | `Department IS NULL`              |
| `IS NOT NULL` | Check for non-nil value   | `Department IS NOT NULL`          |
| `ANY`         | Match any value in a list | `ANY(Skills) = ANY('Go', 'Rust')` |

#### Example Queries
```sql
# Basic filtering
Name = 'Alice'
Salary > 80,000
Skills CONTAINS 'Go'

# Nested fields and maps
Department.Location = 'Remote'
Tags.level = 'senior'

# Complex logic
(Age > 30 AND Salary > 75,000) OR IsEmployed = false
ANY(Skills) = 'Go' AND NOT (Department IS NULL)
```

### Advanced Usage

#### Nested Structs and Maps
Query nested fields or map values using dot notation:

```go
query := "Department.Name = 'Engineering' AND Tags.level = 'senior'"
results, err := parser.Parse(query, people)
```

#### Numeric Formats
The parser supports advanced numeric formats:
- Negative numbers: `Salary > -1000`
- Scientific notation: `Salary > 7.5e4`
- Comma-separated numbers: `Salary > 1,000,000.50`

#### Humanized Values
The parser automatically converts humanized values to their numeric equivalents:

**Byte Sizes (Decimal and Binary Units):**
```sql
# Decimal units (powers of 1000)
Drive.Size > 10GB        # Converts to 10,000,000,000 bytes
Memory > 1.5TB           # Converts to 1,500,000,000,000 bytes  
Storage < 500MB          # Converts to 500,000,000 bytes
Buffer < 100KB           # Converts to 100,000 bytes

# Binary units (powers of 1024)  
Backup > 2.5GiB          # Converts to 2,684,354,560 bytes
Cache > 512MiB           # Converts to 536,870,912 bytes
Temp < 100KiB            # Converts to 102,400 bytes
Archive > 1TiB           # Converts to 1,099,511,627,776 bytes
```

**SI Prefixes:**
```sql
Population > 1.5M        # Converts to 1500000
Count < 5K               # Converts to 5000
Records >= 2.3G          # Converts to 2300000000
```

**Comma-Separated Numbers:**
```sql
Price > 1,000,000        # Converts to 1000000
Users >= 50,000          # Converts to 50000
```

**Example with Real Data:**
```go
type Server struct {
    Name    string
    Memory  int64  // in bytes
    Storage int64  // in bytes
}

servers := []Server{
    {Name: "web1", Memory: 8589934592, Storage: 536870912000},    // 8GB, 500GB
    {Name: "db1", Memory: 34359738368, Storage: 2199023255552},   // 32GB, 2TB
}

// Query using decimal units (powers of 1000)
results, _ := parser.Parse("Memory > 16GB AND Storage < 1TB", servers)

// Query using binary units (powers of 1024) 
results, _ := parser.Parse("Memory > 16GiB AND Storage < 1TiB", servers)

// Both queries work and can be mixed
results, _ := parser.Parse("Memory > 8GB AND Storage > 1TiB", servers)
```

### Performance Considerations
Based on benchmark results:
- **Efficient for Small to Medium Datasets**: Queries on datasets of 10–1000 structs are fast, with simple queries (e.g., `Age > 30`) taking microseconds.
- **Reflection Overhead**: Minimal reflection is used during evaluation, with no reflection during query compilation.
- **Scalability**: Performance scales linearly with dataset size. For very large datasets (>10,000 items), consider batching.
- **Query Complexity**: Complex queries with nested logic or `ANY` operators are slightly slower but optimized with short-circuit evaluation.
- **Memory Usage**: Low memory footprint, with minimal allocations for simple queries (benchmarks show 1–2 allocations per query).

### Error Handling
The parser provides detailed error messages:

```go
_, err := parser.Parse("Age >", people)
if err != nil {
    fmt.Println(err) // Output: "failed to parse query: unexpected EOF"
}

_, err = parser.Parse("InvalidField = 10", people)
if err != nil {
    fmt.Println(err) // Output: "evaluation error: field 'InvalidField' not found"
}
```

## 🛠️ Building and Testing

Clone the repository and build:

```bash
git clone https://github.com/zveinn/parser.git
cd parser
```

Run tests to verify functionality:

```bash
go test -v ./...
```

Run benchmarks to measure performance:

```bash
go test -bench=. ./...
```

## 📚 Documentation

- **API Reference**: Available via [GoDoc](https://pkg.go.dev/github.com/zveinn/parser).
- **Examples**: See the [examples/](examples/) directory for sample queries (create this directory if needed).
- **Source Code Insights**:
  - `parser.go`: Core parsing logic with AST evaluation.
  - `enhanced_lexer.go`: Tokenization with support for advanced numeric formats.
  - `parser_test.go`: Comprehensive test suite for all operators and edge cases.

## 🤝 Contributing

Contributions are welcome! To contribute:

1. Fork the repository.
2. Create a feature branch (`git checkout -b feature/my-feature`).
3. Commit your changes (`git commit -m 'Add my feature'`).
4. Push to the branch (`git push origin feature/my-feature`).
5. Open a Pull Request.

## 📜 License

This project is licensed under the Apache 2.0 License. See [LICENSE](LICENSE) for details.

## 🌟 Acknowledgements

- Built by [zveinn](https://github.com/zveinn).
- Inspired by SQL query engines and libraries like [rql](https://github.com/dvaldivia/rql).
- Thanks to the Go community for feedback and inspiration.

---

⭐ **Star this project** if you find it useful!  
💬 Report issues or suggest features in [Issues](https://github.com/zveinn/parser/issues).
