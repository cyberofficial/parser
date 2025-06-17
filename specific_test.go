package parser

import (
	"testing"
)

// Test to specifically check the failure case
func Test_problematic_query(t *testing.T) {
	query := "Age = 25 OR (Age = 30 AND Address.City = 'Anytown')"
	l := NewLexer(query)

	// Let's inspect the tokens
	t.Log("Tokens for query:", query)
	tokens := []Token{}
	for tok := l.NextToken(); tok.Type != EOF; tok = l.NextToken() {
		tokens = append(tokens, tok)
		t.Logf("Token Type: %s, Literal: %s", tok.Type, tok.Literal)
	}

	// Let's check if we have a closing parenthesis
	rparen := false
	for _, tok := range tokens {
		if tok.Type == RPAREN {
			rparen = true
			break
		}
	}
	t.Logf("Has closing parenthesis: %v", rparen)

	// Test the parser directly
	p := NewParser(NewLexer(query))
	expr, err := p.ParseQuery()
	if err != nil {
		t.Logf("Parse error: %v", err)
	} else {
		t.Logf("Successfully parsed expression: %T", expr)
	}

	// Try a more explicit query to see if it works
	query2 := "Age = 25 OR ((Age = 30) AND (Address.City = 'Anytown'))"
	p2 := NewParser(NewLexer(query2))
	expr2, err2 := p2.ParseQuery()
	if err2 != nil {
		t.Logf("Parse error for explicit query: %v", err2)
	} else {
		t.Logf("Successfully parsed explicit query: %T", expr2)
	}
}
