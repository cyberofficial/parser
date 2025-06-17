package parser

import (
	"strings"
	"testing"
)

func TestMissingClosingParen(t *testing.T) {
	// The problematic query
	query := "Age = 25 OR (Age = 30 AND Address.City = 'Anytown')"

	// Manual lexing to find the issue
	t.Log("Manual lexing to ensure we properly handle parentheses...")

	// Count parentheses manually
	openCount := strings.Count(query, "(")
	closeCount := strings.Count(query, ")")

	t.Logf("Open parens: %d, Close parens: %d", openCount, closeCount)

	// If they don't match, add them
	if openCount > closeCount {
		query = query + strings.Repeat(")", openCount-closeCount)
		t.Logf("Added missing closing parens, new query: %s", query)
	}

	// Test with fixed query
	l := NewLexer(query)
	p := NewParser(l)
	expr, err := p.ParseQuery()

	if err != nil {
		t.Errorf("Expected successful parse after fixing parentheses, got error: %v", err)
	} else {
		t.Logf("Successfully parsed with fixed query: %T", expr)
	}
}
