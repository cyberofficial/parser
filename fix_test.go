package parser

import (
	"strings"
	"testing"
)

func TestMissingClosingParen(t *testing.T) {
	// The problematic query that was fixed
	query := "Age = 25 OR (Age = 30 AND Address.City = 'Anytown')"

	// Count parentheses to verify they're balanced
	openCount := strings.Count(query, "(")
	closeCount := strings.Count(query, ")")

	// Ensure parentheses are balanced
	if openCount != closeCount {
		t.Errorf("Unbalanced parentheses in test query: %d open, %d close", openCount, closeCount)
	}

	// Parse the query
	l := NewLexer(query)
	p := NewParser(l)
	expr, err := p.ParseQuery()

	if err != nil {
		t.Errorf("Failed to parse balanced query: %v", err)
	}
	
	// Verify we got an expression of the correct type
	if _, ok := expr.(*OrExpression); !ok {
		t.Errorf("Expected OrExpression, got %T", expr)
	}
}
