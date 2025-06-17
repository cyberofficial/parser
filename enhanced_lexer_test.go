package parser

import (
	"testing"
)

func TestEnhancedLexerNumericFormats(t *testing.T) {
	tests := []struct {
		input           string
		expectedType    TokenType
		expectedLiteral string
	}{
		{"42", NUMBER, "42"},
		{"-42", NUMBER, "-42"},
		{"3.14", NUMBER, "3.14"},
		{"-3.14", NUMBER, "-3.14"}, {"1e6", NUMBER, "1e6"},
		{"1.5e3", NUMBER, "1.5e3"},
		{"1.5e-3", NUMBER, "1.5e-3"},
		{"1.5e+3", NUMBER, "1.5e+3"},
		// Now handles commas in numbers		{"1,000", NUMBER, "1,000"},
		{"1,000,000.5", NUMBER, "1,000,000.5"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := NewEnhancedLexer(tt.input)
			token := l.NextToken()
			if token.Type != tt.expectedType {
				t.Errorf("Expected token type %q, got %q", tt.expectedType, token.Type)
			}
			if token.Literal != tt.expectedLiteral {
				t.Errorf("Expected token literal %q, got %q", tt.expectedLiteral, token.Literal)
			}
		})
	}
}

func TestLexerInvalidNumericToken(t *testing.T) {
	input := "Age > 25abc"
	l := NewEnhancedLexer(input)

	tokens := []Token{}
	for i := 0; i < 10; i++ { // Get first 10 tokens or until EOF
		token := l.NextToken()
		tokens = append(tokens, token)
		if token.Type == EOF {
			break
		}
	}

	t.Logf("Tokens for '%s':", input)
	for i, token := range tokens {
		t.Logf(" %d: Type=%s, Literal=%s", i, token.Type, token.Literal)
	}
}
