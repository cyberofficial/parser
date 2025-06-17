package parser

import (
	"fmt"
	"testing"
)

// Test to verify that parentheses are handled correctly
func Test_complex_parentheses_queries(t *testing.T) {
	testCases := []struct {
		name     string
		query    string
		wantType string // Type of the expected expression
		wantErr  bool   // Whether we expect an error
	}{
		{
			name:     "Simple parentheses",
			query:    "Age = 25 OR (Age = 30 AND Address.City = 'Anytown')",
			wantType: "*parser.OrExpression",
			wantErr:  false,
		},
		{
			name:     "Nested parentheses",
			query:    "Age = 25 OR ((Age = 30) AND (Address.City = 'Anytown'))",
			wantType: "*parser.OrExpression",
			wantErr:  false,
		},
		{
			name:     "Triple nested parentheses",
			query:    "(((ID = 1)) OR (ID = 2 AND (Name = 'Alice')))",
			wantType: "*parser.OrExpression",
			wantErr:  false,
		},
		{
			name:     "Missing closing parenthesis",
			query:    "(ID = 1 AND Name = 'Alice'",
			wantType: "",
			wantErr:  true,
		},
		{
			name:     "Extra closing parenthesis",
			query:    "ID = 1 AND Name = 'Alice')",
			wantType: "",
			wantErr:  true,
		},
		{
			name:     "Empty parentheses",
			query:    "()",
			wantType: "*parser.ConjunctionExpression",
			wantErr:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewParser(NewLexer(tc.query))
			expr, err := p.ParseQuery()

			if tc.wantErr && err == nil {
				t.Errorf("Expected error for query %q, got nil", tc.query)
			}

			if !tc.wantErr && err != nil {
				t.Errorf("Unexpected error for query %q: %v", tc.query, err)
			}

			if !tc.wantErr && expr != nil {
				gotType := fmt.Sprintf("%T", expr)
				if gotType != tc.wantType {
					t.Errorf("Wrong expression type for query %q: got %s, want %s", tc.query, gotType, tc.wantType)
				}
			}
		})
	}
}
