# Parser Benchmarks

This file contains comprehensive benchmarks for the parser implementation.

## Benchmarks Overview

The benchmarks are organized into several categories to test different aspects of the parser:

### 1. Basic Queries
Tests simple query operations like string equality, numeric comparison, boolean checks, and basic logical operations.

### 2. Complex Queries
Tests more advanced query patterns including:
- Nested logical operators
- Complex combinations of AND/OR
- Array/slice operations
- Nested field access
- Special operators (IS NULL, ANY, NOT)

### 3. Performance with Different Data Sizes
Tests how query performance scales with different dataset sizes (10, 100, 1000 records).

### 4. Parser Characteristics
Tests specific aspects of parser performance:
- Token lexing complexity
- Different operator types
- Logical operator behavior
- Parenthesis nesting depth
- Query selectivity

### 5. Special Operators
Focused benchmarks for special operators like IS NULL, NOT, and ANY.

### 6. Memory Allocation Patterns
Tests memory allocation characteristics for queries of different complexity.

### 7. Query Reuse
Tests the performance impact of running the same query multiple times.

## Running the Benchmarks

To run all benchmarks:
```
go test -bench=. -benchmem
```

To run a specific benchmark suite:
```
go test -bench=BenchmarkComplexQueries -benchmem
```

To run a specific benchmark:
```
go test -bench=BenchmarkBasicQueries/StringEquality -benchmem
```

To limit the number of iterations (useful for long-running benchmarks):
```
go test -bench=BenchmarkComplexQueries -benchmem -benchtime=10x
```

## Interpreting Results

The benchmark results include:
- Operations per second: Higher is better
- ns/op: Time per operation in nanoseconds (lower is better)
- B/op: Bytes allocated per operation (lower is better)
- allocs/op: Number of allocations per operation (lower is better)

## Benchmark Analysis

### Parser Performance Characteristics

1. **Lexing and Parsing Complexity**:
   - Simple queries are processed much faster than complex ones
   - Each additional token adds overhead to lexing time
   - Logical operators and parentheses increase parsing complexity

2. **Memory Usage Patterns**:
   - Memory allocation increases with query complexity
   - Special operators like ANY tend to use more memory
   - Deep nesting of expressions has higher memory overhead

3. **Query Complexity Impact**:
   - AND/OR operations have different performance characteristics
   - Deeply nested expressions are significantly slower
   - Parentheses add parsing overhead but clarify evaluation precedence

4. **Data Volume Sensitivity**:
   - Performance scales roughly linearly with dataset size
   - Highly selective queries (matching fewer records) can be faster
   - Queries on large arrays or nested structures have higher overhead

### Areas for Potential Optimization

1. Memory allocation reductions for repeated query patterns
2. Optimizing common comparison operations
3. Better handling of array/slice operations
4. More efficient parsing of complex logical expressions
